package ber

import (
	"math"
	"testing"
)

var negativeZero = math.Copysign(0, -1)

func TestRealEncoding(t *testing.T) {
	for _, value := range []float64{
		0.15625,
		-0.15625,
		math.Inf(1),
		math.Inf(-1),
		math.NaN(),
		negativeZero,
		0.0,
	} {
		enc := encodeFloat(value)
		dec, err := ParseReal(enc)
		if err != nil {
			t.Errorf("Failed to decode %f (%v): %s", value, enc, err)
		}
		if dec != value {
			if !(math.IsNaN(dec) && math.IsNaN(value)) {
				t.Errorf("decoded value != orig: %f <=> %f", value, dec)
			}
		}
	}
}

func TestRealBinaryDecodingTC10(t *testing.T) {
	// This is the content of tests/tc10.ber. The orignal test suite would emit a
	// "Needlessly long format" warning which we don't care about.
	dec, err := DecodePacketErr([]byte{0x09, 0x07, 0x83, 0x04, 0xff, 0xff, 0xff, 0xfb, 0x05})
	var expected float64 = 0.156250
	if err != nil {
		t.Errorf("Failed to decode: %s", err)
	}
	result := dec.Value.(float64)
	if result != expected {
		t.Errorf("invalid value parsed in tc10: %f <=> %f", result, expected)
	}
}
