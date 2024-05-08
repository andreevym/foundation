package mock

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/anoideaopen/foundation/core"
	"github.com/anoideaopen/foundation/proto"
	pb "github.com/golang/protobuf/proto" //nolint:staticcheck
)

func (w *Wallet) BatcherSignedInvoke(ch string, fn string, args ...string) ([]byte, error) {
	response, _, err := w.BatcherSignedInvokeWithTxEventReturned(ch, fn, args...)
	if err != nil {
		return nil, err
	}

	return response.GetResult(), nil
}

func (w *Wallet) BatcherSignedInvokeWithTxEventReturned(
	ch string,
	fn string,
	args ...string,
) (
	*proto.BatcherRequestResponse,
	*proto.BatcherRequestEvent,
	error,
) {
	err := w.verifyIncoming(ch, fn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to verify incoming args: %w", err)
	}

	// setup creator
	cert, err := hex.DecodeString(batchRobotCert)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode hex string batchRobotCert: %w", err)
	}
	w.ledger.stubs[ch].SetCreator(cert)

	// sign argument and use output args with signature for invoke chaincode 'batcherBatchExecute'
	argsWithSign, _ := w.sign(fn, ch, args...)

	r := core.BatcherRequest{
		BatcherRequestID:   strconv.FormatInt(rand.Int63(), 10),
		Chaincode:          ch,
		Channel:            ch,
		Method:             fn,
		Args:               argsWithSign,
		BatcherRequestType: core.TxBatcherRequestType,
	}

	requests := core.BatcherBatchExecuteRequestDTO{Requests: []core.BatcherRequest{r}}
	bytes, err := json.Marshal(requests)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal requests BatcherBatchExecuteRequestDTO: %w", err)
	}

	// do invoke chaincode
	peerResponse, err := w.ledger.doInvokeWithPeerResponse(ch, txIDGen(), core.BatcherBatchExecute, string(bytes))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to invoke method %s: %w", core.BatcherBatchExecute, err)
	}

	if peerResponse.GetStatus() != http.StatusOK {
		return nil, nil, fmt.Errorf("failed to invoke method %s, status: '%v', message: '%s'", core.BatcherBatchExecute, peerResponse.GetStatus(), peerResponse.GetMessage())
	}

	var batchResponse proto.BatcherBatchResponse
	err = json.Unmarshal(peerResponse.GetPayload(), &batchResponse)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal BatcherBatchExecuteResponseDTO: %w", err)
	}

	if len(batchResponse.GetRequestResponses()) != 1 {
		return nil, nil, fmt.Errorf("failed to handle response, current response len is %d", len(batchResponse.GetRequestResponses()))
	}

	txEvent, err := w.getBatcherTxEventFromChannelByRequestID(ch, r.BatcherRequestID)
	if err != nil {
		return nil, nil, err
	}

	response := batchResponse.GetRequestResponses()[0]
	if responseErr := response.GetError(); responseErr != nil {
		return nil, nil, errors.New(responseErr.GetError())
	}

	return response, txEvent, nil
}

func (w *Wallet) getBatcherTxEventFromChannelByRequestID(
	channel string,
	requestID string,
) (
	*proto.BatcherRequestEvent,
	error,
) {
	e := <-w.ledger.stubs[channel].ChaincodeEventsChannel
	if e.GetEventName() == core.BatcherBatchExecuteEvent {
		events := &proto.BatcherBatchEvent{}
		err := pb.Unmarshal(e.GetPayload(), events)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal BatcherBatchEvent: %w", err)
		}
		for _, ev := range events.GetEvents() {
			if ev.GetBatcherRequestId() == requestID {
				return ev, nil
			}
		}
	}
	return nil,
		fmt.Errorf(
			"failed to find event %s for request %s",
			core.BatcherBatchExecuteEvent,
			requestID,
		)
}
