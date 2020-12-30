package vatinator

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BTBurke/sqlitestore"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var SessionNotValid error = errors.New("session expired or doesnt exist")
var defaultOptions = sessions.Options{
	Path: "/",
	// 2 weeks
	MaxAge: 60 * 60 * 24 * 14,
}

type SessionService interface {
	New(w http.ResponseWriter, r *http.Request, id AccountID) error
	Get(w http.ResponseWriter, r *http.Request) (AccountID, error)
}

type Key []byte

func NewSessionService(path string, keys ...[]byte) (SessionService, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("must supply key for session service")
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	store, err := sqlitestore.NewStore(db, keys...)
	if err != nil {
		return nil, err
	}
	return &sessionService{store}, nil
}

type sessionService struct {
	store *sqlitestore.Store
}

func (s *sessionService) New(w http.ResponseWriter, r *http.Request, id AccountID) error {
	sess, err := s.store.New(r, "vat")
	if err != nil {
		log.Printf("got error higher")
		return err
	}
	sess.Options = &defaultOptions
	sess.Values["account_id"] = int64(id)
	if err := sess.Save(r, w); err != nil {
		log.Printf("got error here")
		return err
	}
	return nil
}

func (s *sessionService) Get(w http.ResponseWriter, r *http.Request) (AccountID, error) {
	sess, err := s.store.New(r, "vat")
	if err != nil || sess.IsNew {
		return AccountID(-1), SessionNotValid
	}
	sess.Options = &defaultOptions
	val := sess.Values["account_id"]
	id, ok := val.(int64)
	if !ok {
		return AccountID(-1), SessionNotValid
	}
	// if session update fails, still return id
	if err := sess.Save(r, w); err != nil {
		return AccountID(id), err
	}

	return AccountID(id), nil
}

func GetSessionKeys(db *DB) ([][]byte, error) {
	q := "SELECT * FROM keys ORDER BY created DESC;"
	resp := []struct {
		ID      int
		Key     []byte
		Created time.Time
	}{}
	if err := db.Select(&resp, q); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// if no existing keys, create one and store it in DB
	if len(resp) == 0 {
		key := securecookie.GenerateRandomKey(32)
		if err := insertKey(db, key); err != nil {
			return nil, err
		}
		return [][]byte{key}, nil
	}

	// if it's been more than a month, rotate keys and persist
	if time.Since(resp[0].Created) > 30*24*time.Hour {
		key := securecookie.GenerateRandomKey(32)
		if err := insertKey(db, key); err != nil {
			return nil, err
		}
		out := [][]byte{key}
		for _, r := range resp {
			out = append(out, r.Key)
		}
		return out, nil
	}

	// all good, return keys
	var out [][]byte
	for _, r := range resp {
		out = append(out, r.Key)
	}
	return out, nil
}

func insertKey(db *DB, key []byte) error {
	q := "INSERT INTO keys (key) VALUES ($1);"
	if _, err := db.Exec(q, key); err != nil {
		return err
	}
	return nil
}
