package svc

import (
	"fmt"
	"strings"
	"time"

	"github.com/BTBurke/vatinator/db"
	"github.com/vmihailenco/msgpack/v5"
)

type Batch struct {
	StartID     string
	NumReceipts int
	VAT         int
	Total       int
	Closed      int64
}

func (b *Batch) MarshalBinary() ([]byte, error) {
	return msgpack.Marshal(b)
}

func (b *Batch) UnmarshalBinary(data []byte) error {
	return msgpack.Unmarshal(data, b)
}

func (b *Batch) Type() byte {
	return db.Batch
}

func (b *Batch) TTL() time.Duration {
	// default - batches never expire
	return 0
}

type BatchKey struct {
	AccountID string
	BatchID   string
}

func (b *BatchKey) MarshalBinary() ([]byte, error) {
	if len(b.AccountID) == 0 || len(b.BatchID) == 0 {
		return nil, fmt.Errorf("batch key error: acct: %s batch: %s", b.AccountID, b.BatchID)
	}
	return []byte(fmt.Sprintf("a/%s/b/%s", b.AccountID, b.BatchID)), nil
}

func (b *BatchKey) UnmarshalBinary(data []byte) error {
	key := splitKey(data)
	acctID, ok := key["a"]
	if !ok {
		return fmt.Errorf("batch missing account ID: %s", string(data))
	}
	batchID, ok := key["b"]
	if !ok {
		return fmt.Errorf("batch missing batch ID: %s", string(data))
	}

	b.AccountID = acctID
	b.BatchID = batchID

	return nil
}

// splits key string like a/b/c/d and assembles pairwise map a=b, c=d
func splitKey(b []byte) map[string]string {
	out := make(map[string]string)

	elements := strings.Split(string(b), "/")
	if len(elements)%2 != 0 {
		// TODO: replace panic here, but this should not happen, compile time guarantee
		panic(fmt.Sprintf("unallowed key: %s", string(b)))
	}

	for i := 0; i < len(elements); i += 2 {
		if i+1 > len(elements)-1 {
			continue
		}
		out[elements[i]] = elements[i+1]
	}
	return out
}

var _ db.Key = &BatchKey{}
var _ db.Entity = &Batch{}
