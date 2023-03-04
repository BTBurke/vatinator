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
		// 2023 update
		{"dot delimited full year", "09.12.2023", "09/12/2023"},
		{"dot delimited short year", "09.12.23", "09/12/2023"},
		{"spaces and extra info", "date: 09.12.23", "09/12/2023"},
		{"slash delimited full year", "09/12/2023", "09/12/2023"},
		{"slash delimited short year", "09/12/23", "09/12/2023"},
		{"not delimited full year", "09122023", "09/12/2023"},
		{"not delimited short year", "091223", "09/12/2023"},
		{"not delimited full year 22", "09122023", "09/12/2023"},
		{"spaces inside year", "09.12. 2023", "09/12/2023"},
		{"hyphen delimited", "09-12-2023", "09/12/2023"},
		{"pathological from ocr fail", "29.12, 2023", "29/12/2023"},
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
