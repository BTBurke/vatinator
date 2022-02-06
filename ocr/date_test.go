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
		// 2022 update
		{"dot delimited full year", "09.12.2022", "09/12/2022"},
		{"dot delimited short year", "09.12.22", "09/12/2022"},
		{"spaces and extra info", "date: 09.12.22", "09/12/2022"},
		{"slash delimited full year", "09/12/2022", "09/12/2022"},
		{"slash delimited short year", "09/12/22", "09/12/2022"},
		{"not delimited full year", "09122022", "09/12/2022"},
		{"not delimited short year", "091222", "09/12/2022"},
		{"not delimited full year 22", "09122022", "09/12/2022"},
		{"spaces inside year", "09.12. 2022", "09/12/2022"},
		{"hyphen delimited", "09-12-2022", "09/12/2022"},
		{"pathological from ocr fail", "29.12, 2022", "29/12/2022"},
		// TODO: only handles dates through 2022
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
		{"reversed date", "2021-12-29", "29/12/2021"},
		{"reversed date short year", "21-12-29", "29/12/2021"},
		// TODO: only handles dates through 2022
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := extractDateReversed([]string{tc.in})
			assert.Equal(t, tc.out, got)
		})
	}
}
