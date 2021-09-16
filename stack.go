package middleware

import (
	"net/http"
)

// Stack returns the middleware functions setup as a single MiddlewareFunc. The
// middlewares will be called in the order they were added.
func Stack(in ...MiddlewareFunc) MiddlewareFunc {
	if len(in) == 0 {
		panic(`must provide at least one middleware`)
	}

	return func(next http.Handler) http.Handler {
		for i := len(in) - 1; i >= 0; i-- {
			next = in[i](next)
		}

		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(rw, r)
		})
	}
}
