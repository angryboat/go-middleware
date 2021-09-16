package middleware

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"
)

// Logger performs the next handler then calls RequestLogger.Log with the
// response status code, runtime of the sub-handler, and a clone of the original
// http.Request.
func Logger(logger RequestLogger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			w := &responseWriter{
				Code:       0,
				HeaderMap:  make(http.Header),
				BodyBuffer: new(bytes.Buffer),
			}

			startTime := time.Now()
			next.ServeHTTP(w, r)
			deltaTime := time.Since(startTime)

			logRequest := r.Clone(context.Background())

			defer logger.Log(w.Code, deltaTime, logRequest)

			rw.WriteHeader(w.Code)
			for key, value := range w.HeaderMap {
				for _, header := range value {
					rw.Header().Add(key, header)
				}
			}
			rw.Write(w.BodyBuffer.Bytes())
		})
	}
}

// RequestLogger is used to generate a log message for a request.
type RequestLogger interface {
	// Log received the HTTP status code, duration of the sub-handler, and a copy of the http.Request
	Log(int, time.Duration, *http.Request)
}

type dRequestLogger struct{}

func (l *dRequestLogger) Log(status int, runtime time.Duration, request *http.Request) {
	log.Default().Printf("%s %s completed %d %s in %s", request.Method, request.URL.EscapedPath(), status, http.StatusText(status), runtime)
}

// DefaultRequestLogger provides a simple log message written to log.Default()
var DefaultRequestLogger RequestLogger = new(dRequestLogger)

type responseWriter struct {
	Code       int
	HeaderMap  http.Header
	BodyBuffer *bytes.Buffer
}

var _ http.ResponseWriter = new(responseWriter)

func (r *responseWriter) WriteHeader(status int) {
	r.Code = status
}

func (r *responseWriter) Write(b []byte) (int, error) {
	if r.Code == 0 {
		r.Code = http.StatusOK
	}
	return r.BodyBuffer.Write(b)
}

func (r *responseWriter) Header() http.Header {
	return r.HeaderMap
}
