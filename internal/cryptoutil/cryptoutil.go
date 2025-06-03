package cryptoutil

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func Encrypte(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {

	return rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)

}

func Decrypte(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {

	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)

}

func getPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("неверный формат публичного ключа: %w", err)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге сертификата: %w", err)
	}

	// Получение публичного ключа
	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("ошибка при получении ключа из сертификата: %w", err)
	}

	return pubKey, nil

}

func EncrypteBody(data []byte, path string) ([]byte, error) {
	publicKey, err := getPublicKey(path)
	if err != nil {
		return nil, err
	}

	return Encrypte(publicKey, data)
}
