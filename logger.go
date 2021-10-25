package middleware

import (
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
			w := NewResponseWriter()

			startTime := time.Now()
			next.ServeHTTP(w, r)
			deltaTime := time.Since(startTime)

			logRequest := r.Clone(context.Background())

			defer logger.Log(w.StatusCode, deltaTime, logRequest)

			w.Apply(rw)
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
