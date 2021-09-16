package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/angryboat/go-middleware"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	tFormatter := new(testLogFormatter)
	logger := middleware.Logger(tFormatter)

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	t.Run("logs a response with delta-time", func(t *testing.T) {
		defer tFormatter.Reset()

		logger(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			http.Redirect(rw, r, "/foo-bar", http.StatusTemporaryRedirect)
		})).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		assert.Equal(t, 1, tFormatter.Count())
	})
}
