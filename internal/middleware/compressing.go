package middleware

import (
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

func WithCompressing(handler http.Handler) http.Handler {
	compressFunc := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("accept-encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("content-encoding", "gzip")
		gzipWriter, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gzipWriter.Close()
		handler.ServeHTTP(gzipResponseWriter{ResponseWriter: w, Writer: gzipWriter}, r)
	}
	return http.HandlerFunc(compressFunc)
}
