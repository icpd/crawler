package randkey

import (
	"math/rand"
	"time"
)

var (
	charSet               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charSet[seededRand.Intn(len(charSet))]
	}
	return string(b)
}
