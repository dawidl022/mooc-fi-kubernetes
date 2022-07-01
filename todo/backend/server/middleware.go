// code from https://blog.questionable.services/article/guide-logging-middleware-go/

package server

import (
	"bytes"
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

type readCloser struct {
	io.Reader
	io.Closer
}

func NewTeeCloser(rc io.ReadCloser, w io.Writer) io.ReadCloser {
	tee := io.TeeReader(rc, w)
	return readCloser{tee, rc}
}

func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := wrapResponseWriter(w)

			var bodyBuf bytes.Buffer
			tee := NewTeeCloser(r.Body, &bodyBuf)
			r.Body = tee

			next.ServeHTTP(wrapped, r)
			if wrapped.status == 0 {
				wrapped.status = 200
			}

			body, _ := io.ReadAll(&bodyBuf)
			logger.Printf("%d %s %s %s", wrapped.status, r.Method, r.URL.EscapedPath(), string(body))
		})
	}
}
