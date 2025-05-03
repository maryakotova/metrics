// Пакет authsign реализует механизм подписи передаваемых данных по алгоритму SHA256. Для этого используется hash.
// от всего тела запроса, которвый разещается в HTTP-заголовке HashSHA256.
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
