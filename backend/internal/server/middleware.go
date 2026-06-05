package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	corsAllowMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	corsAllowHeaders = "Accept, Content-Type, Authorization"
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

func cors(allowedOrigins []string) gin.HandlerFunc {
	origins := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		origins[origin] = struct{}{}
	}

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if _, ok := origins[origin]; ok {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Access-Control-Allow-Methods", corsAllowMethods)
			ctx.Header("Access-Control-Allow-Headers", corsAllowHeaders)
			ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
			ctx.Header("Vary", "Origin")
		}

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
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
