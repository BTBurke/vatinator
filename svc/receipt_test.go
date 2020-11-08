package svc

import (
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/shamaton/msgpack"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Id   string
	Msg  string
	Date int64
}

func TestMsgpack(t *testing.T) {
	t1 := &testStruct{"test", "test2", time.Now().Unix()}
	b, err := msgpack.Encode(t1)
	assert.NoError(t, err)

	t2 := &testStruct{}
	err2 := msgpack.Decode(b, t2)
	assert.NoError(t, err2)
	assert.Equal(t, t1, t2)

}

func TestRKMarshal(t *testing.T) {
	r := &ReceiptKey{
		AccountID: "test",
		ReceiptID: "test2",
	}

	b, err := r.MarshalBinary()
	assert.NoError(t, err)
	r2 := &ReceiptKey{}
	err2 := r2.UnmarshalBinary(b)
	assert.NoError(t, err2)
	assert.Equal(t, r, r2)

}

func TestRMarshal(t *testing.T) {
	r := &Receipt{
		ID:    "test",
		Total: 100,
	}

	b, err := msgpack.Encode(r)
	assert.NoError(t, err)
	r2 := &Receipt{}
	err2 := msgpack.Decode(b, r2)
	assert.NoError(t, err2)
	assert.Equal(t, r, r2)
}

func TestRMashalInterface(t *testing.T) {
	r := &Receipt{
		ID:                xid.New().String(),
		Total:             100,
		VAT:               20,
		Vendor:            "Me AS",
		Date:              "04/09/2020",
		CurrencyPrecision: Digit2,
		Reviewed:          time.Now().Unix(),
	}

	b, err := r.MarshalBinary()
	assert.NoError(t, err)
	r2 := &Receipt{}

	err2 := r2.UnmarshalBinary(b)
	assert.NoError(t, err2)
	assert.Equal(t, r, r2)
}

func TestCurrencyConv(t *testing.T) {
	tt := []struct {
		name string
		in   int
		p    Precision
		out  string
	}{
		{name: "1dp2", in: 5, p: Digit2, out: "0.05"},
		{name: "1dp3", in: 5, p: Digit3, out: "0.005"},
		{name: "2dp2", in: 15, p: Digit2, out: "0.15"},
		{name: "2dp3", in: 15, p: Digit3, out: "0.015"},
		{name: "3dp2", in: 105, p: Digit2, out: "1.05"},
		{name: "3dp3", in: 105, p: Digit3, out: "0.105"},
		{name: "4dp2", in: 1055, p: Digit2, out: "10.55"},
		{name: "4dp3", in: 1055, p: Digit3, out: "1.055"},
		{name: "5dp2", in: 11115, p: Digit2, out: "111.15"},
		{name: "5dp3", in: 11115, p: Digit3, out: "11.115"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.p {
			case Digit2:
				assert.Equal(t, tc.out, currency2ToString(tc.in))
			default:
				assert.Equal(t, tc.out, currency3ToString(tc.in))
			}
		})
	}
}
