package handlers

import (
	"net/http"

	"github.com/michaeljs1990/sqlitestore"
)

type Middleware func(http.Handler) http.Handler

type accountCtx struct{}

var accountCtxKey accountCtx = accountCtx{}

// SessionMiddleware gets the current session and populates the account ID to the request context.  Responds
// with 403 forbidden if no session exists
func SessionMiddleware(store *sqlitestore.SqliteStore) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: add session handling
		})
	}
}
