package decryptmiddleware

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"metrics/internal/config"
	"metrics/internal/cryptoutil"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
)

type DecryptMiddleware struct {
	config *config.Config
	logger *zap.Logger
}

func NewDecrypteMW(cfg *config.Config, logger *zap.Logger) *DecryptMiddleware {
	return &DecryptMiddleware{
		config: cfg,
		logger: logger,
	}
}

func (dmw *DecryptMiddleware) Decrypte(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if strings.Contains(r.Header.Get("Content-Encrypted"), "true") {

			privateKey, err := getPrivateKey(dmw.config.PrivateKeyPath)
			if err != nil {
				http.Error(w, "Ошибка загрузки приватного ключа", http.StatusBadRequest)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			encryptedData, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			decrypted, err := cryptoutil.Decrypte(privateKey, encryptedData)
			if err != nil {
				http.Error(w, "Ошибка дешифровки", http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(decrypted))

		}

		h.ServeHTTP(ow, r)
	}
}

func getPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("неверный формат приватного ключа")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)

}
