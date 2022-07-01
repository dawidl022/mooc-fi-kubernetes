// code from https://blog.questionable.services/article/guide-logging-middleware-go/

package server

import (
	"io"
	"log"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			body, _ := io.ReadAll(r.Body)
			logger.Printf("%d %s %s %s", wrapped.status, r.Method, r.URL.EscapedPath(), body)
		})
	}
}
