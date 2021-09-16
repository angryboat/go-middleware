/*
Package middleware provides http.Handler functions for common server middleware
tasks.

Basic Usage

  // Create a logger middleware using the default request logger and recovery
  stack := middleware.Stack(
		middleware.Logger(middleware.DefaultRequestLogger),
		middleware.Recovery(log.Default().Writer()),
  )

  // Wrap the fooHandler in the logger middleware.
  http.Handle("/foo", stack(fooHandler))
*/
package middleware

import "net/http"

// MiddlewareFunc defines a standard middleware pattern. It takes a http.Handler
// as an argument and returns a new http.Handler. The returned handler will be
// responsible for processing the request and calling the supplied handler to
// continue the chain.
type MiddlewareFunc func(http.Handler) http.Handler
