package svc

import (
	"fmt"
	"time"
)

// ShortDateToTime returns a time.Time from a short date in the form dd/mm/yyyy or dd/mm/yy
func ShortDateToTime(s string) (time.Time, error) {
	layouts := []string{
		"02/01/06",
		"02/01/2006",
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("failed to parse time string %s: unknown format", s)
}

// TimeToTimeToShortDate returns the time formatted in dd/mm/yyyy
func TimeToShortDate(t time.Time) string {
	return t.Format("02/01/2006")
}
