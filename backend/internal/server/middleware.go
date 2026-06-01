package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func securityHeaders(production bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		if production {
			ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		ctx.Next()
	}
}

func limitBody(limit int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.ContentLength > limit {
			ctx.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return
		}

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, limit)
		ctx.Next()
	}
}

type rateLimiter struct {
	limit   int
	window  time.Duration
	clock   func() time.Time
	mutex   sync.Mutex
	clients map[string]rateLimitBucket
}

type rateLimitBucket struct {
	count   int
	resetAt time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:   limit,
		window:  window,
		clock:   time.Now,
		clients: make(map[string]rateLimitBucket),
	}
}

func (limiter *rateLimiter) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !limiter.allow(ctx.ClientIP()) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}

		ctx.Next()
	}
}

func (limiter *rateLimiter) allow(key string) bool {
	now := limiter.clock()

	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	bucket := limiter.clients[key]
	if bucket.resetAt.IsZero() || !now.Before(bucket.resetAt) {
		limiter.clients[key] = rateLimitBucket{
			count:   1,
			resetAt: now.Add(limiter.window),
		}
		return true
	}

	if bucket.count >= limiter.limit {
		return false
	}

	bucket.count++
	limiter.clients[key] = bucket
	return true
}
