package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"

	"github.com/sudeeya/metrics-harvester/internal/keygen"
)

func main() {
	cfg, err := keygen.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, keygen.KeySize)
	if err != nil {
		log.Fatal(err)
	}

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	var publicKeyPEM bytes.Buffer
	pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})

	keygen.SaveKeys(
		privateKeyPEM.Bytes(), publicKeyPEM.Bytes(),
		cfg.PrivateKeyPath, cfg.PublicKeyPath)
}
