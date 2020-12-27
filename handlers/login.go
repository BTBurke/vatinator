package handlers

import "net/http"

func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: check password
		w.WriteHeader(http.StatusOK)
	}
}
