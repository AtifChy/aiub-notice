package notice

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func Test_httpGetWithRetry(t *testing.T) {
	origSleep := sleep
	sleep = func(time.Duration) {}
	defer func() { sleep = origSleep }()

	tests := []struct {
		name     string
		setup    func() (*httptest.Server, int)
		validate func(t *testing.T, resp *http.Response, err error)
	}{
		{
			name: "successful request on first try",
			setup: func() (*httptest.Server, int) {
				s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("OK"))
				}))
				return s, 3
			},
			validate: func(t *testing.T, resp *http.Response, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				defer func() { _ = resp.Body.Close() }()
				b, _ := io.ReadAll(resp.Body)
				if string(b) != "OK" {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "successful request after retries",
			setup: func() (*httptest.Server, int) {
				var attempt int32
				failureCount := int32(2)
				s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					n := atomic.AddInt32(&attempt, 1)
					if n <= failureCount {
						hj, ok := w.(http.Hijacker)
						if ok {
							conn, _, _ := hj.Hijack()
							_ = conn.Close()
						}
						return
					}
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("OK"))
				}))
				return s, 5
			},
			validate: func(t *testing.T, resp *http.Response, err error) {
				if err != nil {
					t.Fatalf("expected success after retries, got error: %v", err)
				}
				defer func() { _ = resp.Body.Close() }()
			},
		},
		{
			name: "all retries fail",
			setup: func() (*httptest.Server, int) {
				s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					hj, ok := w.(http.Hijacker)
					if ok {
						conn, _, _ := hj.Hijack()
						_ = conn.Close()
					}
				}))
				return s, 3
			},
			validate: func(t *testing.T, resp *http.Response, err error) {
				if err == nil {
					t.Fatalf("expected error after all retries fail, got success")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, maxRetries := tt.setup()
			defer server.Close()
			resp, err := httpGetWithRetry(server.URL, maxRetries)
			tt.validate(t, resp, err)
		})
	}
}
