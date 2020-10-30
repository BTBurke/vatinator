package svc

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/vmihailenco/msgpack/v5"
)

type BatchService interface {
	CreateBatch(acctID string, startID string) (*Batch, error)
	GetBatch(acctID string, id string) (*Batch, error)
	CloseBatch(acctID string, id string) error
}

type Batch struct {
	StartID string
	Closed  int64
}

type b struct {
	db *badger.DB
}

func (b b) CreateBatch(acctID string, startID string) (*Batch, error) {
	batch := &Batch{StartID: startID}
	bEnc, err := msgpack.Marshal(batch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal batch to msgpack")
	}

	if err := b.db.Update(func(txn *badger.Txn) error {
		key, err := createBatchKey(acctID)
		if err != nil {
			return errors.Wrap(err, "failed to create batch key")
		}
		if err := txn.Set(key, bEnc); err != nil {
			return errors.Wrap(err, "failed to persist batch")
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to create batch")
	}
	return batch, nil
}

func createBatchKey(acctID string) ([]byte, error) {
	id := xid.New()
	aID, err := xid.FromString(acctID)
	if err != nil {
		return nil, err
	}
	return join([]byte("/a/"), aID.Bytes(), []byte("/b/"), id.Bytes()), nil
}

func join(b ...[]byte) []byte {
	var out []byte
	for _, b1 := range b {
		out = append(out, b1...)
	}
	return out
}
