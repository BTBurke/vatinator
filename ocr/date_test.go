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
		// 2026 update
		{"dot delimited full year", "09.12.2026", "09/12/2026"},
		{"dot delimited short year", "09.12.26", "09/12/2026"},
		{"spaces and extra info", "date: 09.12.26", "09/12/2026"},
		{"slash delimited full year", "09/12/2026", "09/12/2026"},
		{"slash delimited short year", "09/12/26", "09/12/2026"},
		{"not delimited full year", "09122026", "09/12/2026"},
		{"not delimited short year", "091226", "09/12/2026"},
		{"not delimited full year 22", "09122026", "09/12/2026"},
		{"spaces inside year", "09.12. 2026", "09/12/2026"},
		{"hyphen delimited", "09-12-2026", "09/12/2026"},
		{"pathological from ocr fail", "29.12, 2026", "29/12/2026"},
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
		{"reversed date", "2026-12-29", "29/12/2026"},
		{"reversed date short year", "26-12-29", "29/12/2026"},
		{"hyphen reversed", "2026-09-12", "12/09/2026"},
		// TODO: only handles dates through 2022
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := extractDateReversed([]string{tc.in})
			assert.Equal(t, tc.out, got)
		})
	}
}
