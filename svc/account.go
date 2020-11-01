package svc

import (
	"fmt"

	"github.com/dgraph-io/badger/v2"
)

type User struct {
	Email     string
	Password  string
	AccountID string
}

type Account struct {
	ID    string
	Users []User
}

type AccountService interface {
	Create(a *Account) error
	Get(accountID string) (*Account, error)
	Update(a *Account) error
	Delete(accountID string) error
}

type a struct {
	db *badger.DB
}

func (a a) Delete(id string) error {
	key := []byte(fmt.Sprintf("a/%s", id))
	return a.db.DropPrefix(key)
}
