package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// InternalAuthMiddleware validates internal API calls using a shared secret
// This ensures only internal services can call internal endpoints
func InternalAuthMiddleware() gin.HandlerFunc {
	internalSecret := os.Getenv("INTERNAL_API_SECRET")

	return func(c *gin.Context) {
		// Check for internal API key header
		providedSecret := c.GetHeader("X-Internal-API-Key")

		if providedSecret == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "MISSING_INTERNAL_KEY",
					"message": "Internal API key is required",
				},
			})
			c.Abort()
			return
		}

		if internalSecret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_CONFIG_ERROR",
					"message": "Internal API secret not configured",
				},
			})
			c.Abort()
			return
		}

		if providedSecret != internalSecret {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "INVALID_INTERNAL_KEY",
					"message": "Invalid internal API key",
				},
			})
			c.Abort()
			return
		}

		// Copy request ID from incoming request for tracing
		requestID := c.GetHeader("X-Request-ID")
		if requestID != "" {
			c.Set("request_id", requestID)
		}

		c.Next()
	}
}