//go:build go1.18
// +build go1.18

// Fuzzing test(s) are gated behind build tags in order to avoid generating errors
// when running on versions of Go that do not support fuzzing. This was added in
// Go 1.18.

package ber

import (
	"os"
	"testing"
	"time"
)

func FuzzDecodePacket(f *testing.F) {
	// Seed the fuzz corpus with the test cases defined in suite_test.go
	for _, tc := range testCases {
		file := tc.File

		dataIn, err := os.ReadFile(file)
		if err != nil {
			f.Fatalf("failed to load file %s into fuzz corpus: %v", file, err)
			continue
		}
		f.Add(dataIn)
	}

	// Seed the fuzz corpus with data known to cause panics in the past
	f.Add([]byte{0x09, 0x02, 0x85, 0x30})
	f.Add([]byte{0x09, 0x01, 0xcf})

	// Set a limit on the length decoded in readPacket() since the call to
	// make([]byte, length) can allocate up to MaxPacketLengthBytes which is
	// currently 2 GB. This can cause memory related crashes when fuzzing in
	// parallel or on memory constrained devices.
	MaxPacketLengthBytes = 65536
	f.Fuzz(func(t *testing.T, data []byte) {
		stime := time.Now()
		p, err := DecodePacketErr(data)

		if e := time.Since(stime); e > (time.Millisecond * 500) {
			t.Fatalf("DecodePacketErr took too long: %s", e)
		}

		if p == nil && err == nil {
			t.Fatalf("DecodePacketErr returned a nil packet and no error")
		}
	})
}
