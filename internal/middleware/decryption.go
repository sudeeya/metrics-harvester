package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
)

// WithDecryption provides middleware that decrypts the request body using the provided private key.
func WithDecryption(privateKey *rsa.PrivateKey, handler http.Handler) http.Handler {
	decryptFunc := func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		decryptedBody, err := privateKey.Decrypt(rand.Reader, body, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(decryptFunc)
}
