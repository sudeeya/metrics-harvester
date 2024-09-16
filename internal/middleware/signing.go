package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

type hmacResponseWriter struct {
	http.ResponseWriter
	key []byte
}

func (w hmacResponseWriter) Write(b []byte) (int, error) {
	h := hmac.New(sha256.New, w.key)
	if _, err := h.Write(b); err != nil {
		return 0, err
	}
	w.Header().Set("HashSHA256", hex.EncodeToString(h.Sum(nil)))
	return w.ResponseWriter.Write(b)
}

// WithLogging provides middleware that handles signing of HTTP requests and responses.
// If the key is empty does nothing.
// Otherwise, it checks the HashSHA256 header of the request.
// If the signature is incorrect, the response status code is 400 (Bad Request).
func WithSigning(key []byte, handler http.Handler) http.Handler {
	signFunc := func(w http.ResponseWriter, r *http.Request) {
		if len(key) == 0 {
			handler.ServeHTTP(w, r)
			return
		}

		hexHash := r.Header.Get("HashSHA256")
		if hexHash == "" {
			handler.ServeHTTP(w, r)
			return
		}
		expected, err := hex.DecodeString(hexHash)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h := hmac.New(sha256.New, key)
		if _, err := h.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		actual := h.Sum(nil)

		if !hmac.Equal(expected, actual) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		handler.ServeHTTP(hmacResponseWriter{ResponseWriter: w, key: key}, r)
	}
	return http.HandlerFunc(signFunc)
}
