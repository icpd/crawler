package randkey

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomKey() string {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(key)
}
