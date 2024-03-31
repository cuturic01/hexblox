package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hexblox/internal/config"
	"strings"
)

func GenerateHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateKeyPair() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(config.Curve, rand.Reader)
	if err != nil {
		panic(err)
	}
	return key
}

func EncodeKey(key *ecdsa.PublicKey) string {
	uncompressedBytes := elliptic.MarshalCompressed(key.Curve, key.X, key.Y)
	return hex.EncodeToString(uncompressedBytes)
}

func DecodeKey(hexEncoded string) (*ecdsa.PublicKey, error) {
	uncompressedBytes, err := hex.DecodeString(hexEncoded)
	if err != nil {
		return nil, err
	}

	x, y := elliptic.UnmarshalCompressed(config.Curve, uncompressedBytes)
	if x == nil || y == nil {
		return nil, fmt.Errorf("invalid uncompressed public key")
	}

	publicKey := &ecdsa.PublicKey{
		Curve: config.Curve,
		X:     x,
		Y:     y,
	}
	return publicKey, nil
}

func VerifySignature(publicKeyHex string, signature string, hash string) bool {
	publicKey, err := DecodeKey(publicKeyHex)
	if err != nil {
		panic(err)
	}

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		panic(err)
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		panic(err)
	}

	return ecdsa.VerifyASN1(publicKey, hashBytes, signatureBytes)
}

func IndentString(input string, indent string) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.ReplaceAll(strings.Join(lines, "\n")+string(input[len(input)-1]), "\n\n  ", "\n")
}
