package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"net/http"
)

// WithDecryption provides middleware that decrypts the request body using the provided symmetric key.
func WithDecryption(key *[]byte, handler http.Handler) http.Handler {
	decryptFunc := func(w http.ResponseWriter, r *http.Request) {
		if len(*key) == 0 {
			handler.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(body) == 0 {
			handler.ServeHTTP(w, r)
			return
		}

		block, err := aes.NewCipher(*key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		nonceSize := gcm.NonceSize()
		nonce := body[:nonceSize]
		body = body[nonceSize:]

		decryptedBody, err := gcm.Open(nil, nonce, body, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(decryptFunc)
}
