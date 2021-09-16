package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/angryboat/go-middleware"
	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	t.Run("given multiple middlewares", func(t *testing.T) {
		tLogger := new(bytes.Buffer)
		tFormatter := new(testLogFormatter)

		stack := middleware.Stack(
			middleware.Logger(tFormatter),
			middleware.Recovery(tLogger),
		)

		req, err := http.NewRequest("GET", "/stack-check", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		t.Run("and a panic", func(t *testing.T) {
			defer tLogger.Reset()
			defer tFormatter.Reset()

			stack(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				panic("testing recovery")
			})).ServeHTTP(rr, req)

			assert.Equal(t, http.StatusInternalServerError, rr.Code)
			assert.Equal(t, 1, tFormatter.Count())
			assert.Greater(t, tLogger.Len(), 0)
		})
	})
}
