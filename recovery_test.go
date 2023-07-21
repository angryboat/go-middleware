package middleware_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/angryboat/go-middleware"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	tLogger := new(bytes.Buffer)

	recovery := middleware.Recovery(tLogger)

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("without panic", func(t *testing.T) {
		rr := httptest.NewRecorder()

		defer tLogger.Reset()

		recovery(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusNotFound)
		})).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("given a string", func(t *testing.T) {
		rr := httptest.NewRecorder()

		defer tLogger.Reset()

		recovery(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			panic("testing recovery")
		})).ServeHTTP(rr, req)

		assert.Greater(t, tLogger.Len(), 0)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("given an error", func(t *testing.T) {
		rr := httptest.NewRecorder()

		defer tLogger.Reset()

		recovery(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			panic(errors.New("testing recovery"))
		})).ServeHTTP(rr, req)

		assert.Greater(t, tLogger.Len(), 0)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("given an int", func(t *testing.T) {
		rr := httptest.NewRecorder()

		defer tLogger.Reset()

		recovery(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			panic(1)
		})).ServeHTTP(rr, req)

		assert.Greater(t, tLogger.Len(), 0)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
