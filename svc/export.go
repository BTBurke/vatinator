package svc

import (
	"fmt"
	"time"

	"github.com/BTBurke/vatinator/db"
	"github.com/shamaton/msgpack"
)

type Export struct {
	ID      string
	ZipFile []byte
}

func (e *Export) TTL() time.Duration {
	return 0
}

func (e *Export) Type() byte {
	return db.Export
}

func (e *Export) MarshalBinary() ([]byte, error) {
	return msgpack.Encode(e)
}

func (e *Export) UnmarshalBinary(data []byte) error {
	return msgpack.Decode(data, e)
}

type ExportKey struct {
	AccountID string
	BatchID   string
}

func (e *ExportKey) MarshalBinary() ([]byte, error) {
	return []byte(fmt.Sprintf("a/%s/e/%s", e.AccountID, e.BatchID)), nil
}

func (e *ExportKey) UnmarshalBinary(data []byte) error {
	key := splitKey(data)
	acctID, ok := key["a"]
	if !ok {
		return fmt.Errorf("batch missing account ID: %s", string(data))
	}
	batchID, ok := key["e"]
	if !ok {
		return fmt.Errorf("batch missing export batch ID: %s", string(data))
	}

	e.AccountID = acctID
	e.BatchID = batchID

	return nil
}

var _ db.Entity = &Export{}
var _ db.Key = &ExportKey{}
