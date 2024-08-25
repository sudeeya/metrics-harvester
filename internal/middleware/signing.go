package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
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
	w.Header().Set("HashSHA256", string(h.Sum(nil)))
	return w.Write(b)
}

func WithSigning(key []byte, handler http.Handler) http.Handler {
	signFunc := func(w http.ResponseWriter, r *http.Request) {
		if len(key) == 0 {
			handler.ServeHTTP(w, r)
			return
		}
		expected := []byte(r.Header.Get("HashSHA256"))
		var body []byte
		if _, err := r.Body.Read(body); err != nil {
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
		handler.ServeHTTP(hmacResponseWriter{ResponseWriter: w, key: key}, r)
	}
	return http.HandlerFunc(signFunc)
}
