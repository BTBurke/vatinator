package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/BTBurke/vatinator"
	"github.com/pkg/errors"
)

// fake in-memory database
var data map[string]string

func init() {
	data = map[string]string{
		"first_name":    "Bryan",
		"last_name":     "Burke",
		"full_name":     "Bryan Burke",
		"diplomatic_id": "B99991111",
		"embassy":       "US Embassy",
		"address":       "Kentmanni 20",
		"bank":          "SWEDBANK, HABEZX, Liviiala 8, 15040 Tallinn, EE220000099999",
		"bank_name":     "SWEDBANK, HABEZX, Liviiala 8, 15040 Tallinn",
		"account":       "EE220000099999",
	}
}

func GetAccountHandler(account vatinator.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := GetAccountID(r)
		if err != nil {
			handleError(w, http.StatusForbidden, errors.Wrap(err, "no account information present"))
			return
		}

		fd, err := account.GetFormData(id)
		if err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "database error"))
		}
		if fd == nil {
			// form data has not been set yet, return 204, no body to trigger form fill
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(fd)
	}
}

func UpdateAccountHandler(account vatinator.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := GetAccountID(r)
		if err != nil {
			handleError(w, http.StatusForbidden, errors.Wrap(err, "no account information present"))
			return
		}

		b := new(bytes.Buffer)
		defer r.Body.Close()
		if _, err := io.Copy(b, r.Body); err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "error processing account update"))
			return
		}

		if err := account.UpdateFormData(id, b.Bytes()); err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "error saving form data"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type accountRequest struct {
	Email    string
	Password string
}

func CreateAccountHandler(account vatinator.AccountService, session vatinator.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var m accountRequest
		if err := dec.Decode(&m); err != nil {
			handleError(w, http.StatusBadRequest, errors.New("could not decode request"))
			return
		}

		if len(m.Email) == 0 || len(m.Password) == 0 {
			handleError(w, http.StatusBadRequest, errors.New("email and password must be provided"))
			return
		}

		id, err := account.Create(m.Email, m.Password)
		if err != nil {
			log.Printf("error creating account: %v", err)
			handleError(w, http.StatusInternalServerError, errors.New("error creating account"))
			return
		}

		if err := session.New(w, r, id); err != nil {
			log.Printf("error creating session: %v", err)
			handleError(w, http.StatusInternalServerError, errors.New("error creating session"))
			return
		}
		for k, v := range w.Header() {
			log.Printf("%s: %s", k, v)
		}

		w.WriteHeader(http.StatusOK)
	}
}
