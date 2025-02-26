package authsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func VerifySig(receivedHash string, data []byte, key []byte) bool {
	calculatedHash := CalculateHash(data, key)
	return receivedHash == calculatedHash
}

func CalculateHash(data []byte, key []byte) string {
	hash := hmac.New(sha256.New, key)
	hash.Write([]byte(data))
	calculatedHash := hash.Sum(nil)
	return hex.EncodeToString(calculatedHash)
}
