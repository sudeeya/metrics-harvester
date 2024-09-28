package agent

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func generateSymmetricKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Reader.Read(key); err != nil {
		log.Fatal(err)
	}
	return key
}

func extractPublicKey(file string) *rsa.PublicKey {
	pemData, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		log.Fatalf("PEM file contains %s, not RSA PUBLIC KEY", block.Type)
	}
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	return publicKey
}
