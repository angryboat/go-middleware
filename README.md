# Go Middleware

[![Testing](https://github.com/angryboat/go-middleware/actions/workflows/testing.yml/badge.svg)](https://github.com/angryboat/go-middleware/actions/workflows/testing.yml)

## Links

- [Documentation](https://pkg.go.dev/github.com/angryboat/go-middleware)

## Basic Usage

```go
// Create a logger middleware using the default request logger and recovery
stack := middleware.Stack(
  middleware.Logger(middleware.DefaultRequestLogger),
  middleware.Recovery(log.Default().Writer()),
)

// Wrap the fooHandler in the logger middleware.
http.Handle("/foo", stack(fooHandler))
```
