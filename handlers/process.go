package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type processRequest struct {
	BatchID string `json:"batch_id"`
	Date    string `json:"date"`
}

func ProcessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		pr := new(processRequest)
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(pr); err != nil {
			handleError(w, http.StatusBadRequest, errors.Wrap(err, "failed to read process request"))
			return
		}
		fmt.Printf("got process request: %v", pr)
		time.Sleep(5 * time.Second)

		w.WriteHeader(http.StatusOK)
	}
}
