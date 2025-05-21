package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {

	// создаём шаблон сертификата
	cert := &x509.Certificate{
		// уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// информация о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"metrics"},
			Country:      []string{"RU"},
		},
		// разрешение на использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// кодируем сертификат в формате PEM, который будет
	// использоваться в качестве публичного ключа?
	var certPEM bytes.Buffer
	pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	// кодируем приватный ключ в формате PEM
	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Сохранение приватного ключа в файл
	privateKeyFile, err := os.Create("private_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer privateKeyFile.Close()

	privateKeyFile.Write(privateKeyPEM.Bytes())

	// Сохранение публичного ключа в файл
	certFile, err := os.Create("cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer certFile.Close()

	certFile.Write(certPEM.Bytes())

	// кодируем публичный ключ в формате PEM
	var publicKeyPEM bytes.Buffer
	pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	// Сохранение приватного ключа в файл
	publicKeyFile, err := os.Create("public_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer publicKeyFile.Close()

	privateKeyFile.Write(publicKeyPEM.Bytes())

}
