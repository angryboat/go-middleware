package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/angryboat/go-middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseRuntime(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	require.NoError(t, err)

	runtime := middleware.ResponseRuntime("X-Runtime")

	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	t.Run("responds with an X-Runtime header", func(t *testing.T) {
		rr := httptest.NewRecorder()

		runtime(app).ServeHTTP(rr, req)

		assert.NotEmpty(t, rr.Header().Get("X-Runtime"))
	})
}
