package handlers

import (
	"context"
	"net/http"

	"github.com/BTBurke/vatinator"
	"github.com/pkg/errors"
)

type Middleware func(http.Handler) http.Handler

type accountCtx struct{}

var accountCtxKey accountCtx = accountCtx{}

// SessionMiddleware gets the current session and populates the account ID to the request context.  Responds
// with 403 forbidden if no session exists
func SessionMiddleware(session vatinator.SessionService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := session.Get(w, r)
			if err != nil {
				handleError(w, http.StatusForbidden, errors.New("forbidden"))
			}
			ctx := context.WithValue(r.Context(), accountCtxKey, id)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func GetAccountID(r *http.Request) (vatinator.AccountID, error) {
	val := r.Context().Value(accountCtxKey)
	var id vatinator.AccountID
	id, ok := val.(vatinator.AccountID)
	if !ok {
		return vatinator.AccountID(-1), errors.New("no account id")
	}
	return id, nil
}
