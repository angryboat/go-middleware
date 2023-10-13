package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery provides a panic handler and logs panic details to the passed
// io.Writer. Log message will be JSON printed to a single line with keys
// "panic" and "stackTrace". The stackTrace key is a Base64 encoded stack trace.
func Recovery(w io.Writer) MiddlewareFunc {
	log := slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	return RecoveryWithLogger(log)
}

func RecoveryWithLogger(logger SLogger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rec := NewResponseWriter()

			defer func() {
				panicErr := recover()

				if panicErr != nil {
					var err error
					switch t := panicErr.(type) {
					case string:
						err = errors.New(t)
					case error:
						err = t
					default:
						err = fmt.Errorf("unknown recovery: (%T) %v", panicErr, panicErr)
					}

					errorContext, _ := json.Marshal(map[string]interface{}{
						"panic":      err.Error(),
						"stackTrace": debug.Stack(),
					})
					errorContext = append(errorContext, '\n')

					logger.ErrorContext(r.Context(), string(errorContext))

					http.Error(rw, err.Error(), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rec, r)

			rec.Apply(rw)
		})
	}
}
