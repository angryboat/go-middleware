package middleware

import (
	"context"
	"log"
	"log/slog"
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

			logRequest := r.Clone(r.Context())

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

func NewStructuredRequestLogger(l SLogger) RequestLogger {
	return &sRequestLogger{l}
}

type sRequestLogger struct {
	l SLogger
}

func (l *sRequestLogger) Log(status int, runtime time.Duration, request *http.Request) {
	l.l.InfoContext(request.Context(), "Request Completed",
		slog.Group("request",
			slog.String("method", request.Method),
			slog.String("path", request.URL.EscapedPath()),
		),
		slog.Int("status", status),
		slog.String("runtime", runtime.String()),
	)
}

type SLogger interface {
	DebugContext(context.Context, string, ...any)
	InfoContext(context.Context, string, ...any)
	ErrorContext(context.Context, string, ...any)
}
