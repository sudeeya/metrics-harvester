package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var typesToCompress = []string{
	"application/json",
	"text/html",
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithCompressing(handler http.Handler) http.Handler {
	compressFunc := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("accept-encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}
		if contentType := w.Header().Get("content-type"); !IsCompressible(contentType) {
			handler.ServeHTTP(w, r)
			return
		}
		gzipWriter, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gzipWriter.Close()
		w.Header().Set("content-encoding", "gzip")
		handler.ServeHTTP(gzipResponseWriter{ResponseWriter: w, Writer: gzipWriter}, r)
	}
	return http.HandlerFunc(compressFunc)
}

func IsCompressible(contentType string) bool {
	for _, ct := range typesToCompress {
		if strings.Contains(contentType, ct) {
			return true
		}
	}
	return false
}
