package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// WithCompressing provides middleware that handles
// gzip compression and decompression for HTTP requests and responses.
// If the Content-Encoding header of request contains "gzip", the request body is decompressed.
// If the Accept-Encoding header of request contains "gzip", the response body is compressed.
func WithCompressing(handler http.Handler) http.Handler {
	compressFunc := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			rawBody, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			body, err := io.ReadAll(rawBody)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gzipWriter, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer gzipWriter.Close()
		handler.ServeHTTP(gzipResponseWriter{ResponseWriter: w, Writer: gzipWriter}, r)
	}
	return http.HandlerFunc(compressFunc)
}
