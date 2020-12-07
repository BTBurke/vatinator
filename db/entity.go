package db

import (
	"encoding"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v2"
)

// Model flags that are stored as metadata along with the binary representation to help with
// decoding
const (
	Unknown byte = 1 << iota
	Account
	User
	Batch
	Receipt
	Export
	Image
)

// EntityTypeError is returned when the value at key does not match the type of the expected entity that it is
// supposed to be marshaled into
var EntityTypeError = fmt.Errorf("entity metadata did not match receiver type")

// Entity is an interface implemented by models in the database
type Entity interface {
	Type() byte
	TTL() time.Duration
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// Key is an enterface for keys that know how to marshal and unmarshal themselves from/to []byte
type Key interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// Set will set a value in the database with an associated entity type and
// TTL
func Set(txn *badger.Txn, key Key, e Entity) error {
	val, err := e.MarshalBinary()
	if err != nil {
		return err
	}
	k, err := key.MarshalBinary()
	if err != nil {
		return err
	}

	entry := badger.NewEntry(k, val).WithMeta(e.Type())
	if e.TTL() > 0 {
		entry.WithTTL(e.TTL())
	}
	return txn.SetEntry(entry)
}

// Get retreives a value from the database and unmarshals into the entity.  It can return
// an EntityTypeError if the types do not match.
func Get(txn *badger.Txn, key Key, e Entity) error {
	k, err := key.MarshalBinary()
	if err != nil {
		return err
	}

	item, err := txn.Get(k)
	if err != nil {
		return err
	}
	return FromItem(item, e)
}

// Del deletes a key and value from the database
func Del(txn *badger.Txn, key Key) error {
	k, err := key.MarshalBinary()
	if err != nil {
		return err
	}
	return txn.Delete(k)
}

// FromItem takes a badger entry, copies it, then unmarshals it into e.  It can return
// an EntityTypeError if the entity types do not match.  Other errors are possible but not typed.
func FromItem(item *badger.Item, e Entity) error {
	b, err := item.ValueCopy(nil)
	if err != nil {
		return err
	}
	if item.UserMeta() != e.Type() {
		return EntityTypeError
	}
	return e.UnmarshalBinary(b)
}
