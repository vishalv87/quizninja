package middleware

import (
	"time"

	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Logger provides structured logging for HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Generate request ID for correlation
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Prepare log fields
		fields := logrus.Fields{
			"request_id":  requestID,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"latency_ms":  latency.Milliseconds(),
			"ip":          c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
		}

		// Add user ID if authenticated
		if userID, exists := c.Get("user_id"); exists {
			fields["user_id"] = userID
		}

		// Add query params if present
		if len(c.Request.URL.RawQuery) > 0 {
			fields["query"] = c.Request.URL.RawQuery
		}

		// Add error if present
		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.String()
		}

		// Log based on status code
		statusCode := c.Writer.Status()
		message := "HTTP request completed"

		if statusCode >= 500 {
			utils.WithFields(fields).Error(message)
		} else if statusCode >= 400 {
			utils.WithFields(fields).Warn(message)
		} else {
			utils.WithFields(fields).Info(message)
		}
	}
}

// ErrorHandler handles errors with structured logging
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			fields := logrus.Fields{
				"error": err.Err.Error(),
				"type":  err.Type,
			}

			// Add request context
			if requestID := c.GetString("request_id"); requestID != "" {
				fields["request_id"] = requestID
			}
			if userID, exists := c.Get("user_id"); exists {
				fields["user_id"] = userID
			}

			utils.WithFields(fields).Error("Request error occurred")

			// Send appropriate response based on error type
			switch err.Type {
			case gin.ErrorTypeBind:
				if !c.Writer.Written() {
					c.JSON(400, gin.H{"error": "Invalid request format"})
				}
			case gin.ErrorTypePublic:
				if !c.Writer.Written() {
					c.JSON(500, gin.H{"error": err.Error()})
				}
			default:
				if !c.Writer.Written() {
					c.JSON(500, gin.H{"error": "Internal server error"})
				}
			}
		}
	}
}
