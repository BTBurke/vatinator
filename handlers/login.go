package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/BTBurke/vatinator"
	"github.com/pkg/errors"
)

func LoginHandler(account vatinator.AccountService, session vatinator.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Email    string
			Password string
		}{}
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		if err := dec.Decode(&req); err != nil {
			handleError(w, http.StatusBadRequest, errors.New("bad request"))
			return
		}

		id, err := account.CheckPassword(req.Email, req.Password)
		if err != nil {
			handleError(w, http.StatusUnauthorized, errors.New("bad email or password"))
			return
		}

		if err := session.New(w, r, id); err != nil {
			handleError(w, http.StatusUnauthorized, errors.Wrap(err, "failed to set cookie"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
