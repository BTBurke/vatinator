package vatinator

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"

	ss "github.com/elithrar/simple-scrypt"
)

var nothing = AccountID(-1)

var LoginFailed error = errors.New("login failed")

type AccountService interface {
	Create(email, password string) (AccountID, error)
	CheckPassword(email, password string) (AccountID, error)
	UpdateFormData(id AccountID, fd []byte) error
	GetFormData(id AccountID) ([]byte, error)
	GetFormAndEmailData(id AccountID) (string, []byte, error)
}

type accountService struct {
	db *DB
}

func NewAccountService(db *DB) AccountService {
	return &accountService{db}
}

func (a accountService) Create(email, password string) (AccountID, error) {
	q := "INSERT INTO accounts (email, password) VALUES (LOWER($1), $2);"

	hash, err := ss.GenerateFromPassword([]byte(password), ss.DefaultParams)
	if err != nil {
		return nothing, err
	}

	resp, err := a.db.Exec(q, email, hash)
	if err != nil {
		return nothing, err
	}
	id, err := resp.LastInsertId()
	if err != nil {
		return nothing, err
	}

	return AccountID(id), nil
}

func (a accountService) CheckPassword(email, password string) (AccountID, error) {
	q := "SELECT id, password FROM accounts WHERE email = LOWER($1);"
	resp := struct {
		ID       AccountID
		Password []byte
	}{}
	if err := a.db.Get(&resp, q, email); err != nil {
		return nothing, err
	}

	if err := ss.CompareHashAndPassword(resp.Password, []byte(password)); err != nil {
		return nothing, LoginFailed
	}
	return resp.ID, nil
}

func (a accountService) UpdateFormData(id AccountID, fd []byte) error {
	q := "UPDATE accounts SET form_data = $1 WHERE id = $2;"
	if _, err := a.db.Exec(q, fd, id); err != nil {
		return err
	}
	return nil
}

func (a accountService) GetFormData(id AccountID) ([]byte, error) {
	q := "SELECT form_data FROM accounts WHERE id = $1;"
	var b []byte
	if err := a.db.Get(&b, q, id); err != nil {
		return nil, err
	}
	return b, nil
}

func (a accountService) GetFormAndEmailData(id AccountID) (string, []byte, error) {
	q := "SELECT email, form_data FROM accounts WHERE id = $1;"
	resp := struct {
		Email    string
		FormData []byte `db:"form_data"`
	}{}
	if err := a.db.Get(&resp, q, id); err != nil {
		return "", nil, err
	}
	return resp.Email, resp.FormData, nil
}

// AccountID is the ID in the database
type AccountID int64

// Fulfills stringer interface
func (a AccountID) String() string {
	return strconv.Itoa(int(a))
}

// Valuer interface for SQL
func (a AccountID) Value() (driver.Value, error) {
	// value needs to be a base driver.Value type
	return int64(a), nil
}

// Scanner interface for SQL
func (a *AccountID) Scan(value interface{}) error {
	// if value is nil, error
	if value == nil {
		return fmt.Errorf("accountID is nil")
	}
	if v, ok := value.(int64); ok {
		*a = AccountID(v)
		return nil
	}
	// otherwise, return an error
	return errors.New("failed to scan account ID")
}
