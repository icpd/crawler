package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

func Encrypt(key, plainInputText string) (string, error) {
	plainText := padding([]byte(plainInputText))
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return string(cipherText), nil
}

func Decrypt(key, plainInputText string) (string, error) {
	outputText := []byte(plainInputText)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(outputText) < aes.BlockSize {
		return "", errors.New("cipher text must be longer than blocksize")
	} else if len(outputText)%aes.BlockSize != 0 {
		return "", errors.New("cipher text must be multiple of blocksize(128bit)")
	}
	iv := outputText[:aes.BlockSize]
	outputText = outputText[aes.BlockSize:]
	plainText := make([]byte, len(outputText))

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(plainText, outputText)

	return string(plainText), nil
}

func padding(b []byte) []byte {
	size := aes.BlockSize - (len(b) % aes.BlockSize)
	pad := bytes.Repeat([]byte{byte(size)}, size)
	return append(b, pad...)
}
