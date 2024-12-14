package utils

import (
	"time"

	"golang.org/x/exp/rand"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

const numberCharset = "0123456789"
const upperCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lowerCharset = "abcdefghijklmnopqrstuvwxyz"
const charset = upperCharset + lowerCharset + numberCharset

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}

func RandomNumberString(length int) string {
	return StringWithCharset(length, numberCharset)
}
