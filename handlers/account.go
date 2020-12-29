package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

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

// TODO: add db in constructor
func GetAccountHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: get account from context

		// TODO: look up account in DB

		// TODO: return account data

		b, err := json.Marshal(data)
		if err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "error serializing account data"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func UpdateAccountHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := new(bytes.Buffer)
		defer r.Body.Close()
		if _, err := io.Copy(b, r.Body); err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "error processing account update"))
			return
		}

		if err := json.Unmarshal(b.Bytes(), &data); err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "error deserializing account update"))
		}
		w.WriteHeader(http.StatusOK)
	}
}

type accountRequest struct {
	Email string
	Password string
}

func CreateAccountHandler(account vat.AccountService, session vat.SessionService) http.HandlerFunc {
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

		id, err := account.Create(email, password string)
		if err != nil {
			handleError(w, http.StatusInternalServerError, errors.New("error creating account"))
			return
		}
		
		if err := session.New(w, r, id); err != nil {
			handleError(w, http.StatusInternalServerError, errors.New("error creating session"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
