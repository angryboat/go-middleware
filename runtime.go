package middleware

import (
	"net/http"
	"time"
)

// ResponseRuntime will include a request duration header response value.
func ResponseRuntime(header string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := NewResponseWriter()
			startTime := time.Now()
			next.ServeHTTP(ww, r)
			ww.Header().Set(header, time.Since(startTime).String())
			ww.Apply(w)
		})
	}
}
