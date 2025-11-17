package ew11

import (
	"encoding/hex"
	"strings"
)

func EncodeBCD(n int) byte {
	tens := n / 10
	ones := n % 10
	return byte((tens << 4) | ones)
}

func DecodeBCD(b byte) int {
	tens := (b >> 4) & 0x0F
	ones := b & 0x0F
	return int(tens*10 + ones)
}

func Ptr[T any](v T) *T {
	return &v
}

func PrettyHex(data []byte) string {
	var sb strings.Builder
	for _, b := range data {
		h := hex.EncodeToString([]byte{b})
		sb.WriteString(h)
		sb.WriteString(" ")
	}
	str := strings.TrimSpace(sb.String())
	return strings.ToUpper(str)
}
