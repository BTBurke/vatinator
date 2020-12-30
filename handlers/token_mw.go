package handlers

import (
	"fmt"
	"net/http"

	"github.com/BTBurke/vatinator"
)

func TokenMiddleware(token vatinator.TokenService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			encToken := r.URL.Query().Get("token")

			if len(path) == 0 || len(encToken) == 0 {
				handleError(w, http.StatusUnauthorized, fmt.Errorf("token failed validation"))
				return
			}

			if err := token.CheckPath(encToken, path); err != nil {
				handleError(w, http.StatusUnauthorized, fmt.Errorf("token failed validation"))
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
