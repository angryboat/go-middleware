package middleware

import (
	"bytes"
	"net/http"
	"strings"
)

type ResponseWriter struct {
	StatusCode int
	Headers    http.Header
	BodyBuffer *bytes.Buffer

	hasWrittenStatus bool
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		StatusCode: http.StatusOK,
		Headers:    make(http.Header),
		BodyBuffer: bytes.NewBuffer([]byte{}),
	}
}

func (r *ResponseWriter) WriteHeader(code int) {
	r.hasWrittenStatus = true
	r.StatusCode = code
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	return r.BodyBuffer.Write(b)
}

func (r *ResponseWriter) Header() http.Header {
	return r.Headers
}

func (r *ResponseWriter) Apply(w http.ResponseWriter) {
	for key, value := range r.Headers {
		w.Header().Set(key, strings.Join(value, ","))
	}

	if r.hasWrittenStatus {
		w.WriteHeader(r.StatusCode)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Write(r.BodyBuffer.Bytes())
}

var _ http.ResponseWriter = new(ResponseWriter)
