package middleware

import (
	"fmt"
	"net/http"
	"time"

	"quizninja-api/config"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var (
	// Global rate limiters
	globalLimiter *limiter.Limiter
	authLimiter   *limiter.Limiter
	strictLimiter *limiter.Limiter
	userLimiter   *limiter.Limiter
)

// InitRateLimiters initializes rate limiters with config values
func InitRateLimiters(cfg *config.Config) {
	// Global: configurable requests per minute per IP
	globalStore := memory.NewStore()
	globalRate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  cfg.RateLimitGlobal,
	}
	globalLimiter = limiter.New(globalStore, globalRate)

	// Auth endpoints: configurable requests per minute per IP (prevent brute force)
	authStore := memory.NewStore()
	authRate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  cfg.RateLimitAuth,
	}
	authLimiter = limiter.New(authStore, authRate)

	// Strict endpoints (create/write operations): configurable requests per minute per IP
	strictStore := memory.NewStore()
	strictRate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  cfg.RateLimitWrite,
	}
	strictLimiter = limiter.New(strictStore, strictRate)

	// Per-user rate limiter: configurable requests per minute per user
	userStore := memory.NewStore()
	userRate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  cfg.RateLimitPerUser,
	}
	userLimiter = limiter.New(userStore, userRate)
}

// GlobalRateLimit applies global rate limiting
func GlobalRateLimit() gin.HandlerFunc {
	if globalLimiter == nil {
		// Return no-op middleware if not initialized
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return mgin.NewMiddleware(globalLimiter)
}

// AuthRateLimit applies strict rate limiting for auth endpoints
func AuthRateLimit() gin.HandlerFunc {
	if authLimiter == nil {
		// Return no-op middleware if not initialized
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return mgin.NewMiddleware(authLimiter)
}

// StrictRateLimit applies rate limiting for write operations
func StrictRateLimit() gin.HandlerFunc {
	if strictLimiter == nil {
		// Return no-op middleware if not initialized
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return mgin.NewMiddleware(strictLimiter)
}

// PerUserRateLimit limits requests per authenticated user
func PerUserRateLimit() gin.HandlerFunc {
	if userLimiter == nil {
		// Return no-op middleware if not initialized
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			// No user ID, skip user-specific rate limiting
			c.Next()
			return
		}

		// Use user ID as key
		key := fmt.Sprintf("user:%v", userID)

		context, err := userLimiter.Get(c.Request.Context(), key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limit check failed",
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", context.Reset))

		if context.Reached {
			c.Header("Retry-After", fmt.Sprintf("%d", context.Reset))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests from this user",
				"retry_after": context.Reset,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
