package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/BTBurke/vatinator"
	"github.com/pkg/errors"
)

type processRequest struct {
	BatchID string `json:"batch_id"`
	Date    string `json:"date"`
}

func ProcessHandler(process vatinator.ProcessService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := GetAccountID(r)
		if err != nil {
			handleError(w, http.StatusBadRequest, errors.Wrap(err, "unknown account"))
			return
		}
		defer r.Body.Close()

		pr := new(processRequest)
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(pr); err != nil {
			handleError(w, http.StatusBadRequest, errors.Wrap(err, "failed to read process request"))
			return
		}
		t, err := time.Parse("January 2006", pr.Date)
		if err != nil {
			handleError(w, http.StatusBadRequest, errors.Wrapf(err, "unable to parse date: %s", pr.Date))
			return
		}

		if err := process.Do(id, pr.BatchID, t.Format("January"), t.Year()); err != nil {
			handleError(w, http.StatusInternalServerError, errors.Wrap(err, "failed to queue batch for processing"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
