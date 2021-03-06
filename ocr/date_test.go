package ocr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDates(t *testing.T) {
	tt := []struct {
		name string
		in   string
		out  string
	}{
		{"dot delimited full year", "09.12.2020", "09/12/2020"},
		{"dot delimited short year", "09.12.20", "09/12/2020"},
		{"spaces and extra info", "date: 09.12.20", "09/12/2020"},
		{"slash delimited full year", "09/12/2020", "09/12/2020"},
		{"slash delimited short year", "09/12/20", "09/12/2020"},
		{"not delimited full year", "09122020", "09/12/2020"},
		{"not delimited short year", "091220", "09/12/2020"},
		{"not delimited full year 21", "09122021", "09/12/2021"},
		{"spaces inside year", "09.12. 2020", "09/12/2020"},
		{"hyphen delimited", "09-12-2020", "09/12/2020"},
		{"pathological from ocr fail", "29.12, 2020", "29/12/2020"},
		// TODO: only handles dates through 2021
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := extractDate([]string{tc.in})
			assert.Equal(t, tc.out, got)
		})
	}
}

func TestReverseDates(t *testing.T) {
	tt := []struct {
		name string
		in   string
		out  string
	}{
		{"reversed date", "2020-12-29", "29/12/2020"},
		{"reversed date short year", "20-12-29", "29/12/2020"},
		// TODO: only handles dates through 2021
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := extractDateReversed([]string{tc.in})
			assert.Equal(t, tc.out, got)
		})
	}
}
