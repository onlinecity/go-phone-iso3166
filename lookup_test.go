package phoneiso3166

import (
	"bytes"
	"strconv"
	"testing"
)

func TestCountry(t *testing.T) {
	tests := []struct {
		msisdn uint64
		expect string
	}{
		{45, "DK"},
		{4566118311, "DK"},
		{38640118311, "SI"},
		{38340118311, "XK"},
		{37740118311, "MC"},
		{1204, "CA"},
		{12024561111, "US"}, // White house comment line
		{14412921234, "BM"}, // Bermuda city hall
		{0, ""},
	}
	for _, tt := range tests {
		m := tt.msisdn
		name := strconv.FormatUint(tt.msisdn, 10)
		ex := tt.expect
		t.Run(name, func(t *testing.T) {
			res := E164.Lookup(m)
			if res != ex {
				t.Errorf(
					"lookup(%d) returned %q, but expected %q",
					m, res, ex,
				)
			}
		})
	}
}

const testNum uint64 = 4566118311
const testString = "4566118311"

func TestPrealloc(t *testing.T) {
	var buf = make([]byte, 0, 16)
	{
		b := strconv.AppendUint(buf, testNum, 10)
		if !bytes.Equal(b, []byte(testString)) {
			t.Error("pre alloc fail")
			t.Log(string(b))
		}
	}
	{
		b := strconv.AppendUint(buf, 4512345678, 10)
		if !bytes.Equal(b, []byte("4512345678")) {
			t.Error("pre alloc fail")
			t.Log(string(b))
		}
	}
}

func BenchmarkE164Lookup(b *testing.B) {
	for n := 0; n < b.N; n++ {
		E164.Lookup(testNum)
	}
}

// This option is interesting if you plan to lookup a lot of numbers
func BenchmarkE164LookupPrealloc(b *testing.B) {
	var buf = make([]byte, 0, 16)
	for n := 0; n < b.N; n++ {
		E164.LookupByteString(
			strconv.AppendUint(buf, testNum, 10),
		)
	}
}

func BenchmarkE164LookupBuffer(b *testing.B) {
	var buf bytes.Buffer
	for n := 0; n < b.N; n++ {
		buf.Reset()
		E164.LookupByteString(
			strconv.AppendUint(buf.Bytes(), testNum, 10),
		)
	}
}

func BenchmarkE164LookupString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		E164.LookupString(testString)
	}
}

func BenchmarkE164LookupBytes(b *testing.B) {
	m := []byte(testString)
	for n := 0; n < b.N; n++ {
		E164.LookupByteString(m)
	}
}

func BenchmarkE164LookupNoExist(b *testing.B) {
	for n := 0; n < b.N; n++ {
		E164.Lookup(19820000000000)
	}
}

func TestOperator(t *testing.T) {
	tests := []struct {
		mcc    uint16
		mnc    uint16
		expect string
	}{
		{238, 0, "DK"},
		{340, 1, "GP"},
		{340, 12, "MQ"},
		{0, 0, ""},
	}
	for _, tt := range tests {
		mcc := tt.mcc
		mnc := tt.mnc
		ex := tt.expect
		m := strconv.FormatUint(uint64(tt.mcc), 10)
		n := strconv.FormatUint(uint64(tt.mnc), 10)
		t.Run(m+n, func(t *testing.T) {
			res := E212.Lookup(mcc, mnc)
			if res != ex {
				t.Errorf(
					"lookup(%d, %d) returned %q, but expected %q",
					mcc, mnc, res, ex,
				)
			}
		})
	}
}

func BenchmarkE212Lookup(b *testing.B) {
	for n := 0; n < b.N; n++ {
		E212.Lookup(340, 12)
	}
}

func TestOperatorName(t *testing.T) {
	tests := []struct {
		mcc    uint16
		mnc    uint16
		expect string
	}{
		{238, 1, "TDC A/S"},
		{238, 0, ""},
		{0, 0, ""},
	}
	for _, tt := range tests {
		mcc := tt.mcc
		mnc := tt.mnc
		ex := tt.expect
		m := strconv.FormatUint(uint64(tt.mcc), 10)
		n := strconv.FormatUint(uint64(tt.mnc), 10)
		t.Run(m+n, func(t *testing.T) {
			res := NetworkName(mcc, mnc)
			if res != ex {
				t.Errorf(
					"lookup(%d, %d) returned %q, but expected %q",
					mcc, mnc, res, ex,
				)
			}
		})
	}
}
