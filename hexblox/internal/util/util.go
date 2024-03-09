package util

import (
	"crypto/rand"
	"crypto/rsa"
)

func GenerateKeyPair() *rsa.PrivateKey {
	bitSize := 4096
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}
	return key
}
