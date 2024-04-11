package cachestub

import (
	"sort"

	"github.com/anoideaopen/foundation/core/types"
	"github.com/anoideaopen/foundation/core/types/big"
	"github.com/anoideaopen/foundation/proto"
)

type TxCacheStub struct {
	*BatchCacheStub
	txID       string
	txCache    map[string]*proto.WriteElement
	events     map[string][]byte
	Accounting []*proto.AccountingRecord
}

func (bs *BatchCacheStub) NewTxCacheStub(txID string) *TxCacheStub {
	return &TxCacheStub{
		BatchCacheStub: bs,
		txID:           txID,
		txCache:        make(map[string]*proto.WriteElement),
		events:         make(map[string][]byte),
	}
}

// GetTxID returns TxCacheStub transaction ID
func (bts *TxCacheStub) GetTxID() string {
	return bts.txID
}

// GetState returns state from TxCacheStub cache or, if absent, from batchState cache
func (bts *TxCacheStub) GetState(key string) ([]byte, error) {
	existsElement, ok := bts.txCache[key]
	if ok {
		return existsElement.Value, nil
	}
	return bts.BatchCacheStub.GetState(key)
}

// PutState puts state to the TxCacheStub's cache
func (bts *TxCacheStub) PutState(key string, value []byte) error {
	bts.txCache[key] = &proto.WriteElement{Value: value}
	return nil
}

// SetEvent sets payload to a TxCacheStub events
func (bts *TxCacheStub) SetEvent(name string, payload []byte) error {
	bts.events[name] = payload
	return nil
}

func (bts *TxCacheStub) AddAccountingRecord(token string, from *types.Address, to *types.Address, amount *big.Int, reason string) {
	bts.Accounting = append(bts.Accounting, &proto.AccountingRecord{
		Token:     token,
		Sender:    from.Bytes(),
		Recipient: to.Bytes(),
		Amount:    amount.Bytes(),
		Reason:    reason,
	})
}

// Commit puts state from a TxCacheStub cache to the BatchCacheStub cache
func (bts *TxCacheStub) Commit() ([]*proto.WriteElement, []*proto.Event) {
	writeKeys := make([]string, 0, len(bts.txCache))
	for k, v := range bts.txCache {
		bts.batchCache[k] = v
		writeKeys = append(writeKeys, k)
	}
	sort.Strings(writeKeys)
	writes := make([]*proto.WriteElement, 0, len(writeKeys))
	for _, k := range writeKeys {
		writes = append(writes, &proto.WriteElement{
			Key:       k,
			Value:     bts.txCache[k].Value,
			IsDeleted: bts.txCache[k].IsDeleted,
		})
	}

	eventKeys := make([]string, 0, len(bts.events))
	for k := range bts.events {
		eventKeys = append(eventKeys, k)
	}
	sort.Strings(eventKeys)
	events := make([]*proto.Event, 0, len(eventKeys))
	for _, k := range eventKeys {
		events = append(events, &proto.Event{
			Name:  k,
			Value: bts.events[k],
		})
	}
	return writes, events
}

// DelState marks state in TxCacheStub as deleted
func (bts *TxCacheStub) DelState(key string) error {
	bts.txCache[key] = &proto.WriteElement{Key: key, IsDeleted: true}
	return nil
}