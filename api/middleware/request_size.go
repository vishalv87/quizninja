package middleware

import (
	"net/http"

	"quizninja-api/config"
	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
)

var requestSizeConfig *config.Config

// InitRequestSizeLimits initializes the request size limit configuration
func InitRequestSizeLimits(cfg *config.Config) {
	requestSizeConfig = cfg
}

// DefaultRequestSizeLimit applies the default request size limit
func DefaultRequestSizeLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if requestSizeConfig == nil || !requestSizeConfig.RequestSizeLimitEnabled {
			c.Next()
			return
		}

		// Apply size limit using MaxBytesReader
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, requestSizeConfig.RequestSizeDefault)

		c.Next()

		// Check if the request was too large
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			utils.ErrorResponse(c, http.StatusRequestEntityTooLarge,
				"Request body too large. Maximum allowed size is 10MB")
			c.Abort()
		}
	}
}

// AuthRequestSizeLimit applies stricter size limits for authentication endpoints
func AuthRequestSizeLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if requestSizeConfig == nil || !requestSizeConfig.RequestSizeLimitEnabled {
			c.Next()
			return
		}

		// Apply auth-specific size limit
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, requestSizeConfig.RequestSizeAuth)

		c.Next()

		// Check if the request was too large
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			utils.ErrorResponse(c, http.StatusRequestEntityTooLarge,
				"Request body too large. Maximum allowed size for auth endpoints is 1MB")
			c.Abort()
		}
	}
}

// WriteRequestSizeLimit applies size limits for write operations
func WriteRequestSizeLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if requestSizeConfig == nil || !requestSizeConfig.RequestSizeLimitEnabled {
			c.Next()
			return
		}

		// Apply write-specific size limit
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, requestSizeConfig.RequestSizeWrite)

		c.Next()

		// Check if the request was too large
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			utils.ErrorResponse(c, http.StatusRequestEntityTooLarge,
				"Request body too large. Maximum allowed size for write operations is 5MB")
			c.Abort()
		}
	}
}

// RequestSizeLimitWithCustomSize creates a middleware with a custom size limit
func RequestSizeLimitWithCustomSize(maxBytes int64, errorMessage string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if requestSizeConfig == nil || !requestSizeConfig.RequestSizeLimitEnabled {
			c.Next()
			return
		}

		// Apply custom size limit
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

		c.Next()

		// Check if the request was too large
		if c.Writer.Status() == http.StatusRequestEntityTooLarge {
			utils.ErrorResponse(c, http.StatusRequestEntityTooLarge, errorMessage)
			c.Abort()
		}
	}
}