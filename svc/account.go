package svc

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type AccountService interface {
	Delete(id string) error
}

type a struct {
	db *badger.DB
}

func (a a) Delete(id string) error {
	idB, err := xid.FromString(id)
	if err != nil {
		return errors.Wrap(err, "failed to delete account")
	}
	return a.db.DropPrefix(append([]byte("a"), idB.Bytes()...))
}

var _ AccountService = a{}
