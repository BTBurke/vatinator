package vatinator

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var SessionNotValid error = errors.New("session expired or doesnt exist")
var defaultOptions = sessions.Options{
	Domain: ".vatinator.com",
	Path:   "/",
	// 45 days
	MaxAge: 60 * 60 * 24 * 45,
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
	sessionPath := filepath.Join(path, "sessions")
	if err := os.MkdirAll(sessionPath, 0755); err != nil {
		return nil, err
	}

	store := sessions.NewFilesystemStore(sessionPath, keys...)
	return &sessionService{store}, nil
}

type sessionService struct {
	store *sessions.FilesystemStore
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

// returns sessions keys, this data structure is fucked.  It should be a keypair
// with (sign, encrypt) so pass (key, nil) when only signing key should be used.
// TODO: remove the patch for encrypted keys that were issued in Feb 2020 because of bug
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

	// if it's been more than a 3 months, rotate keys and persist
	if time.Since(resp[0].Created) > 90*24*time.Hour {
		key := securecookie.GenerateRandomKey(32)
		if err := insertKey(db, key); err != nil {
			return nil, err
		}
		out := [][]byte{key, nil}
		for _, r := range resp {
			out = append(out, r.Key, nil)
		}
		// TODO: remove patch
		if len(resp) >= 2 {
			out = patch(out, resp[len(resp)-2].Key, resp[len(resp)-1].Key)
		}
		return out, nil
	}

	// all good, return keys
	var out [][]byte
	for _, r := range resp {
		out = append(out, r.Key, nil)
	}
	// TODO: remove patch
	if len(resp) >= 2 {
		out = patch(out, resp[len(resp)-2].Key, resp[len(resp)-1].Key)
	}
	return out, nil
}

// TODO: remove the patch after first two keys are long expired.
// This is related to issue where some cookies were inadvertently
// encrypted due to a bug in month one and two
func patch(keys [][]byte, key0 []byte, key1 []byte) [][]byte {
	return append(keys, key0, key1, key1, key0)
}

func insertKey(db *DB, key []byte) error {
	q := "INSERT INTO keys (key) VALUES ($1);"
	if _, err := db.Exec(q, key); err != nil {
		return err
	}
	return nil
}

type devSession struct{}

func NewDevSessionService() SessionService {
	return devSession{}
}

func (devSession) New(w http.ResponseWriter, r *http.Request, id AccountID) error {
	return nil
}

func (devSession) Get(w http.ResponseWriter, r *http.Request) (AccountID, error) {
	return AccountID(1), nil
}
