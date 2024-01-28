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
		// 2024 update
		{"dot delimited full year", "09.12.2024", "09/12/2024"},
		{"dot delimited short year", "09.12.24", "09/12/2024"},
		{"spaces and extra info", "date: 09.12.24", "09/12/2024"},
		{"slash delimited full year", "09/12/2024", "09/12/2024"},
		{"slash delimited short year", "09/12/24", "09/12/2024"},
		{"not delimited full year", "09122024", "09/12/2024"},
		{"not delimited short year", "091224", "09/12/2024"},
		{"not delimited full year 22", "09122024", "09/12/2024"},
		{"spaces inside year", "09.12. 2024", "09/12/2024"},
		{"hyphen delimited", "09-12-2024", "09/12/2024"},
		{"pathological from ocr fail", "29.12, 2024", "29/12/2024"},
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
		{"reversed date", "2023-12-29", "29/12/2023"},
		{"reversed date short year", "23-12-29", "29/12/2023"},
		{"hyphen reversed", "2023-09-12", "12/09/2023"},
		// TODO: only handles dates through 2022
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := extractDateReversed([]string{tc.in})
			assert.Equal(t, tc.out, got)
		})
	}
}
