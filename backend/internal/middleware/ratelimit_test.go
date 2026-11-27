package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiterAllowAndReset(t *testing.T) {
	now := time.Date(2026, 11, 27, 12, 0, 0, 0, time.UTC)
	rl := NewRateLimiter(2, time.Minute)
	rl.now = func() time.Time { return now }

	if !rl.Allow("1.2.3.4") {
		t.Fatal("first request should be allowed")
	}
	if !rl.Allow("1.2.3.4") {
		t.Fatal("second request should be allowed")
	}
	if rl.Allow("1.2.3.4") {
		t.Fatal("third request within window should be blocked")
	}
	if !rl.Allow("5.6.7.8") {
		t.Fatal("different IP should have its own bucket")
	}

	now = now.Add(time.Minute + time.Second)
	if !rl.Allow("1.2.3.4") {
		t.Fatal("request after window expiry should be allowed")
	}
}

func TestRateLimiterMiddlewareReturns429(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rl := NewRateLimiter(1, time.Minute)

	r := gin.New()
	r.POST("/login", rl.Middleware("too many requests"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req1 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req1.RemoteAddr = "10.0.0.1:12345"
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w1.Code)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/login", nil)
	req2.RemoteAddr = "10.0.0.1:12345"
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w2.Code)
	}
}
