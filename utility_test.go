package middleware_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/angryboat/go-middleware"
	"github.com/stretchr/testify/assert"
)

func TestRemoteIP(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("given no addr", func(t *testing.T) {
		ips := middleware.RemoteIP(req.Clone(context.Background()))

		assert.Equal(t, []string{`127.0.0.1`}, ips)
	})

	t.Run("given x-forwarded-for", func(t *testing.T) {
		r := req.Clone(context.Background())
		r.Header.Add("X-Forwarded-For", "123.0.0.1,456.0.0.1")

		ips := middleware.RemoteIP(r)

		assert.Equal(t, []string{`123.0.0.1`, `456.0.0.1`}, ips)
	})

	t.Run("given x-original-forwarded-for", func(t *testing.T) {
		r := req.Clone(context.Background())
		r.Header.Add("X-Original-Forwarded-For", "789.0.0.1,123.0.0.1")

		ips := middleware.RemoteIP(r)

		assert.Equal(t, []string{`789.0.0.1`, `123.0.0.1`}, ips)
	})
}
