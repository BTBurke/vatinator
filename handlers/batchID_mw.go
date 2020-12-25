package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

type batchCtx struct{}

var batchCtxKey batchCtx = batchCtx{}

// BatchIDMiddleware extracts the URL parameter `batch_id` and returns bad request if missing
func BatchIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		batchID := chi.URLParam(r, "batch_id")

		ctx := context.WithValue(r.Context(), batchCtxKey, batchID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getBatchID(r *http.Request) string {
	return r.Context().Value(batchCtxKey).(string)
}
