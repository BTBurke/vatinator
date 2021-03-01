package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BTBurke/vatinator"
	"github.com/pkg/errors"
)

func RequestPasswordReset(account vatinator.AccountService, token vatinator.TokenService, email vatinator.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Email string
		}{}

		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		if err := dec.Decode(&req); err != nil || len(req.Email) == 0 {
			handleError(w, http.StatusBadRequest, errors.Wrap(err, "malformed password reset request"))
			return
		}

		id, err := account.GetAccountID(req.Email)
		if err != nil || id == vatinator.AccountID(-1) {
			// return 200 even if request is bogus to not leak private info
			w.WriteHeader(http.StatusOK)
			return
		}

		encToken, err := token.NewPasswordReset(req.Email)
		if err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to create password reset token"))
			return
		}

		if err := email.SendPasswordResetEmail(req.Email, vatinator.EmailData{
			Link: fmt.Sprintf("https://vatinator.com/reset/%s", encToken),
		}); err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to send password reset email"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DoPasswordReset(account vatinator.AccountService, token vatinator.TokenService, session vatinator.SessionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Password string
			Token    string
		}{}

		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		if err := dec.Decode(&req); err != nil || len(req.Password) == 0 || len(req.Token) == 0 {
			handleError(w, http.StatusBadRequest, errors.Wrap(err, "malformed password reset request"))
			return
		}

		email, err := token.CheckPasswordReset(req.Token)
		if err != nil {
			handleError(w, http.StatusUnauthorized, errors.Wrap(err, "password reset token not valid"))
			return
		}

		if err := account.SetPassword(email, req.Password); err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to set new password"))
			return
		}
		// this returns an error but making best shot at deleting cookie
		_ = session.Del(w, r, 0)

		w.WriteHeader(http.StatusOK)
	}
}
