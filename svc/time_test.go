package svc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeConversion(t *testing.T) {
	tt := []struct {
		name string
		in   string
		out  time.Time
		back string
	}{
		{name: "short year", in: "01/06/20", out: time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC), back: "01/06/2020"},
		{name: "long year", in: "01/06/2020", out: time.Date(2020, 06, 01, 0, 0, 0, 0, time.UTC), back: "01/06/2020"},
	}

	for _, tc := range tt {
		t1, err := ShortDateToTime(tc.in)
		assert.NoError(t, err)
		assert.Equal(t, t1, tc.out)

		t2 := TimeToShortDate(t1)
		assert.Equal(t, t2, tc.back)
	}
}
