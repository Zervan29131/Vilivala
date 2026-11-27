package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type windowCounter struct {
	count   int
	resetAt time.Time
}

// RateLimiter enforces a fixed-window request limit per key (e.g. client IP).
type RateLimiter struct {
	mu          sync.Mutex
	limit       int
	window      time.Duration
	entries     map[string]*windowCounter
	lastCleanup time.Time
	now         func() time.Time
}

// NewRateLimiter creates a fixed-window limiter allowing `limit` requests per `window`.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:       limit,
		window:      window,
		entries:     make(map[string]*windowCounter),
		lastCleanup: time.Now(),
		now:         time.Now,
	}
}

// Allow reports whether a request from key is within the configured rate limit.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := rl.now()
	if now.Sub(rl.lastCleanup) >= rl.window {
		for k, v := range rl.entries {
			if !now.Before(v.resetAt) {
				delete(rl.entries, k)
			}
		}
		rl.lastCleanup = now
	}

	entry, ok := rl.entries[key]
	if !ok || !now.Before(entry.resetAt) {
		rl.entries[key] = &windowCounter{count: 1, resetAt: now.Add(rl.window)}
		return true
	}
	if entry.count >= rl.limit {
		return false
	}
	entry.count++
	return true
}

// Middleware wraps the limiter as a Gin handler keyed by client IP.
func (rl *RateLimiter) Middleware(message string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.Allow(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": message,
				"data":    nil,
			})
			return
		}
		c.Next()
	}
}
