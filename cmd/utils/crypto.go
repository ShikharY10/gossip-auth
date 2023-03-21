package utils

import (
	crand "crypto/rand"

	"encoding/base64"
	"io"
	"math/big"
)

func Encode(data []byte) string {
	hb := base64.StdEncoding.EncodeToString([]byte(data))
	return hb
}

// Decoding the base string to array of bytes
func Decode(data string) []byte {
	hb, _ := base64.StdEncoding.DecodeString(data)
	return hb
}

func GenerateRandomId() string {
	r32, _ := crand.Int(crand.Reader, big.NewInt(999999999999999))
	p32, _ := crand.Prime(crand.Reader, 5)
	s := r32.String() + p32.String()
	return s
}

func GenerateOTP(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(crand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
