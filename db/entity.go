package db

import (
	"encoding"
	"github.com/dgraph-io/badger/v2"
	"time"
)

const (
	Unknown byte = 1 << iota
	Account
	User
	Batch
	Receipt
)

type Entity interface {
	Type() byte
	TTL() time.Duration
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func Set(txn *badger.Txn, key []byte, e Entity) error {
	val, err := e.MarshalBinary()
	if err != nil {
		return err
	}
	entry := badger.NewEntry(key, val).WithMeta(e.Type())
	if e.TTL() > 0 {
		entry.WithTTL(e.TTL())
	}
	return txn.SetEntry(entry)
}

func Get(txn *badger.Txn, key []byte, e Entity) error {
	item, err := txn.Get(key)
	if err != nil {
		return err
	}
	b, err := item.ValueCopy(nil)
	if err != nil {
		return err
	}
	if err := e.UnmarshalBinary(b); err != nil {
		return err
	}
	return nil
}
