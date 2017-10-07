// +build gofuzz

package ber

import (
	"bytes"
)

func Fuzz(data []byte) int {
	rd := bytes.NewReader(data)
	if _, err := ReadPacket(rd); err != nil {
		return 0
	}
	return 1
}
