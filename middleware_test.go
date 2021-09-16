package middleware_test

import (
	"net/http"
	"sync"
	"time"

	"github.com/angryboat/go-middleware"
)

type testLogFormatter struct {
	mu    sync.Mutex
	count int
}

func (t *testLogFormatter) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.count
}

func (t *testLogFormatter) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.count = 0
}

func (t *testLogFormatter) Log(_ int, _ time.Duration, _ *http.Request) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.count += 1
}

var _ middleware.RequestLogger = new(testLogFormatter)
