package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lunancy1992/jianghu-server/internal/pkg/response"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    burst,
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	l, exists := rl.limiters[key]
	if !exists {
		l = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = l
	}
	return l
}

// RateLimit creates a per-IP rate limiting middleware.
func RateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	rl := NewRateLimiter(r, burst)
	return func(c *gin.Context) {
		key := c.ClientIP()

		// If user is authenticated, use user ID instead
		if uid, exists := c.Get("user_id"); exists {
			key = "user:" + toString(uid)
		}

		if !rl.getLimiter(key).Allow() {
			response.Error(c, http.StatusTooManyRequests, 4001, "rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case int64:
		return string(rune(val))
	default:
		return "unknown"
	}
}

// SMSRateLimit is stricter rate limiting for SMS endpoints.
func SMSRateLimit() gin.HandlerFunc {
	rl := NewRateLimiter(rate.Limit(0.1), 3) // 1 per 10 seconds, burst 3
	return func(c *gin.Context) {
		key := "sms:" + c.ClientIP()
		if !rl.getLimiter(key).Allow() {
			response.Error(c, http.StatusTooManyRequests, 4001, "too many SMS requests")
			c.Abort()
			return
		}
		c.Next()
	}
}
