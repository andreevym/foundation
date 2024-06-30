package core

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/anoideaopen/foundation/core/cachestub"
	"github.com/anoideaopen/foundation/core/contract"
	"github.com/anoideaopen/foundation/core/logger"
	"github.com/anoideaopen/foundation/core/telemetry"
	"github.com/anoideaopen/foundation/core/types"
	"github.com/anoideaopen/foundation/proto"
	pb "github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const ExecuteTasksEvent = "executeTasks"

var ErrTasksNotFound = errors.New("no tasks found")

type ValidatedTx struct {
	task          *proto.Task
	senderAddress *proto.Address
	method        contract.Method
	args          []string
	txCacheStub   *cachestub.TxCacheStub
	nonce         uint64
	err           error
	sender        *types.Sender
}

// TaskExecutor handles the execution of a group of tasks.
type TaskExecutor struct {
	sync.Mutex

	BatchCacheStub *cachestub.BatchCacheStub
	Chaincode      *Chaincode
	SKI            string
	TracingHandler *telemetry.TracingHandler
}

// NewTaskExecutor initializes a new TaskExecutor.
func NewTaskExecutor(stub shim.ChaincodeStubInterface, cc *Chaincode, tracingHandler *telemetry.TracingHandler) *TaskExecutor {
	return &TaskExecutor{
		BatchCacheStub: cachestub.NewBatchCacheStub(stub),
		Chaincode:      cc,
		TracingHandler: tracingHandler,
	}
}

// TasksExecutorHandler executes multiple sub-transactions (tasks) within a single transaction in Hyperledger Fabric,
// using cached state between tasks to solve the MVCC problem. Each request in the arguments contains its own set of
// arguments for the respective chaincode method calls.
func TasksExecutorHandler(
	traceCtx telemetry.TraceContext,
	stub shim.ChaincodeStubInterface,
	args []string,
	cc *Chaincode,
) ([]byte, error) {
	tracingHandler := cc.contract.TracingHandler()
	traceCtx, span := tracingHandler.StartNewSpan(traceCtx, ExecuteTasks)
	defer span.End()

	log := logger.Logger()
	txID := stub.GetTxID()
	span.SetAttributes(attribute.String("tx_id", txID))
	start := time.Now()
	defer func() {
		log.Infof("tasks executor: tx id: %s, elapsed: %s", txID, time.Since(start))
	}()

	if len(args) != 1 {
		err := fmt.Errorf("failed to validate args for transaction %s: expected exactly 1 argument, received %d", txID, len(args))
		return nil, handleTasksError(span, err)
	}

	var executeTaskRequest proto.ExecuteTasksRequest
	if err := pb.Unmarshal([]byte(args[0]), &executeTaskRequest); err != nil {
		err = fmt.Errorf("failed to unmarshal argument to ExecuteTasksRequest for transaction %s, argument: %s", txID, args[0])
		return nil, handleTasksError(span, err)
	}

	log.Warningf("tasks executor: tx id: %s, txs: %d", txID, len(executeTaskRequest.GetTasks()))

	if len(executeTaskRequest.GetTasks()) == 0 {
		err := fmt.Errorf("failed to validate argument: no tasks found in ExecuteTasksRequest for transaction %s: %w", txID, ErrTasksNotFound)
		return nil, handleTasksError(span, err)
	}

	executor := NewTaskExecutor(stub, cc, tracingHandler)

	response, event, err := executor.ExecuteTasks(traceCtx, executeTaskRequest.GetTasks())
	if err != nil {
		return nil, handleTasksError(span, fmt.Errorf("failed to handle task for transaction %s: %w", txID, err))
	}

	eventData, err := pb.Marshal(event)
	if err != nil {
		return nil, handleTasksError(span, fmt.Errorf("failed to marshal event for transaction %s: %w", txID, err))
	}

	err = stub.SetEvent(ExecuteTasksEvent, eventData)
	if err != nil {
		return nil, handleTasksError(span, fmt.Errorf("failed to set event for transaction %s: %w", txID, err))
	}

	data, err := pb.Marshal(response)
	if err != nil {
		return nil, handleTasksError(span, fmt.Errorf("failed to marshal response for transaction %s: %w", txID, err))
	}

	return data, nil
}

// ExecuteTasks processes a group of tasks, returning a group response and event.
func (e *TaskExecutor) ExecuteTasks(
	traceCtx telemetry.TraceContext,
	tasks []*proto.Task,
) (
	*proto.BatchResponse,
	*proto.BatchEvent,
	error,
) {
	traceCtx, span := e.TracingHandler.StartNewSpan(traceCtx, "TaskExecutor.ExecuteTasks")
	defer span.End()

	batchResponse := &proto.BatchResponse{}
	batchEvent := &proto.BatchEvent{}

	validatedTxsMap := make(map[string]ValidatedTx)
	wg := &sync.WaitGroup{}

	// init router once
	e.Chaincode.Router()

	m := &sync.Mutex{}
	for _, task := range tasks {
		wg.Add(1)
		go func(t *proto.Task) {
			defer wg.Done()

			validatedTx := ValidatedTx{
				task:        t,
				txCacheStub: e.BatchCacheStub.NewTxCacheStub(t.GetId()),
			}

			span.AddEvent("parsing chaincode method")
			method, err := e.Chaincode.Method(task.GetMethod())
			if err != nil {
				err = fmt.Errorf("failed to parse chaincode method '%s' for task %s: %w", task.GetMethod(), task.GetId(), err)
				span.SetStatus(codes.Error, err.Error())
				txResponse, txEvent := handleTaskError(span, validatedTx.task, err)
				batchResponse.TxResponses = append(batchResponse.TxResponses, txResponse)
				batchEvent.Events = append(batchEvent.Events, txEvent)
				return
			}

			validatedTx.method = method

			m.Lock()
			defer m.Unlock()
			validatedTxsMap[t.GetId()] = validatedTx
		}(task)
	}
	wg.Wait()

	for id, validatedTx := range validatedTxsMap {
		span.AddEvent("validating and extracting invocation context")
		senderAddress, args, nonce, err := e.Chaincode.validateAndExtractInvocationContext(e.BatchCacheStub, validatedTx.method, validatedTx.task.GetArgs())
		if err != nil {
			err = fmt.Errorf("failed to validate and extract invocation context for task %s: %w", validatedTx.task.GetId(), err)
			span.SetStatus(codes.Error, err.Error())
			txResponse, txEvent := handleTaskError(span, validatedTx.task, err)
			batchResponse.TxResponses = append(batchResponse.TxResponses, txResponse)
			batchEvent.Events = append(batchEvent.Events, txEvent)
			break
		}

		validatedTx.nonce = nonce
		validatedTx.senderAddress = senderAddress
		validatedTx.args = args
		validatedTxsMap[id] = validatedTx
	}

	checkResult := make(chan *ValidatedTx, len(validatedTxsMap))
	defer close(checkResult)
	for _, validatedTx := range validatedTxsMap {
		wg.Add(1)
		go func(validatedTx ValidatedTx) {
			defer wg.Done()

			span.AddEvent("validating authorization")
			if !validatedTx.method.RequiresAuth || validatedTx.senderAddress == nil {
				err := fmt.Errorf("failed to validate authorization for task %s: sender address is missing", validatedTx.task.GetId())
				span.SetStatus(codes.Error, err.Error())
				validatedTx.err = err
				checkResult <- &validatedTx
				return
			}
			argsToValidate := append([]string{validatedTx.senderAddress.AddrString()}, validatedTx.args...)

			span.AddEvent("validating arguments")
			if err := e.Chaincode.Router().Check(validatedTx.method.MethodName, argsToValidate...); err != nil {
				err = fmt.Errorf("failed to validate arguments for task %s: %w", validatedTx.task.GetId(), err)
				span.SetStatus(codes.Error, err.Error())
				validatedTx.err = err
				checkResult <- &validatedTx
				return
			}

			span.AddEvent("validating nonce")
			validatedTx.sender = types.NewSenderFromAddr((*types.Address)(validatedTx.senderAddress))
			validatedTx.args = validatedTx.args[:validatedTx.method.NumArgs-1]
			checkResult <- &validatedTx
		}(validatedTx)
	}
	wg.Wait()
exit:
	for {
		select {
		case validatedTx := <-checkResult:
			if validatedTx.err != nil {
				txResponse, txEvent := handleTaskError(span, validatedTx.task, validatedTx.err)
				batchResponse.TxResponses = append(batchResponse.TxResponses, txResponse)
				batchEvent.Events = append(batchEvent.Events, txEvent)
				delete(validatedTxsMap, validatedTx.task.GetId())
			} else {
				validatedTxsMap[validatedTx.task.Id] = *validatedTx
			}
		default:
			break exit
		}
	}

	for id, validatedTx := range validatedTxsMap {
		err := checkNonce(e.BatchCacheStub, validatedTx.sender, validatedTx.nonce)
		if err != nil {
			err = fmt.Errorf("failed to validate nonce for task %s, nonce %d: %w", validatedTx.task.GetId(), validatedTx.nonce, err)
			span.SetStatus(codes.Error, err.Error())

			txResponse, txEvent := handleTaskError(span, validatedTx.task, err)
			batchResponse.TxResponses = append(batchResponse.TxResponses, txResponse)
			batchEvent.Events = append(batchEvent.Events, txEvent)
			delete(validatedTxsMap, id)
		}
	}

	for _, validatedTx := range validatedTxsMap {
		txResponse, txEvent := e.ExecuteTask(traceCtx, validatedTx)
		batchResponse.TxResponses = append(batchResponse.TxResponses, txResponse)
		batchEvent.Events = append(batchEvent.Events, txEvent)
	}

	if err := e.BatchCacheStub.Commit(); err != nil {
		return nil, nil, fmt.Errorf("failed to commit changes using BatchCacheStub: %w", err)
	}

	return batchResponse, batchEvent, nil
}

// ExecuteTask processes an individual task, returning a transaction response and event.
func (e *TaskExecutor) ExecuteTask(
	traceCtx telemetry.TraceContext,
	validatedTx ValidatedTx,
) (*proto.TxResponse, *proto.BatchTxEvent) {
	traceCtx, span := e.TracingHandler.StartNewSpan(traceCtx, "TaskExecutor.ExecuteTasks")
	defer span.End()
	log := logger.Logger()
	start := time.Now()
	span.SetAttributes(
		attribute.String("task_method", validatedTx.task.GetMethod()),
		attribute.StringSlice("task_args", validatedTx.task.GetArgs()),
		attribute.String("task_id", validatedTx.task.GetId()),
	)
	defer func() {
		log.Infof("task method %s task %s elapsed: %s", validatedTx.task.GetMethod(), validatedTx.task.GetId(), time.Since(start))
	}()

	span.AddEvent("calling method")
	response, err := e.Chaincode.InvokeContractMethod(traceCtx, validatedTx.txCacheStub, validatedTx.method, validatedTx.senderAddress, validatedTx.args)
	if err != nil {
		return handleTaskError(span, validatedTx.task, err)
	}

	span.AddEvent("commit")
	writes, events := validatedTx.txCacheStub.Commit()

	sort.Slice(validatedTx.txCacheStub.Accounting, func(i, j int) bool {
		return strings.Compare(validatedTx.txCacheStub.Accounting[i].String(), validatedTx.txCacheStub.Accounting[j].String()) < 0
	})

	span.SetStatus(codes.Ok, "")
	return &proto.TxResponse{
			Id:     []byte(validatedTx.task.GetId()),
			Method: validatedTx.task.GetMethod(),
			Writes: writes,
		},
		&proto.BatchTxEvent{
			Id:         []byte(validatedTx.task.GetId()),
			Method:     validatedTx.task.GetMethod(),
			Accounting: validatedTx.txCacheStub.Accounting,
			Events:     events,
			Result:     response,
		}
}

func handleTasksError(span trace.Span, err error) error {
	logger.Logger().Error(err)
	span.SetStatus(codes.Error, err.Error())
	return err
}

func handleTaskError(span trace.Span, task *proto.Task, err error) (*proto.TxResponse, *proto.BatchTxEvent) {
	logger.Logger().Errorf("%s: %s: %s", task.GetMethod(), task.GetId(), err)
	span.SetStatus(codes.Error, err.Error())

	ee := proto.ResponseError{Error: err.Error()}
	return &proto.TxResponse{
			Id:     []byte(task.GetId()),
			Method: task.GetMethod(),
			Error:  &ee,
		}, &proto.BatchTxEvent{
			Id:     []byte(task.GetId()),
			Method: task.GetMethod(),
			Error:  &ee,
		}
}
