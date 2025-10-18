package utils

import (
	"os"
	"runtime"
	"strings"

	"quizninja-api/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// InitLogger initializes the application logger with configuration
func InitLogger(cfg *config.Config) {
	Log = logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
		Log.Warnf("Invalid log level '%s', defaulting to INFO", cfg.LogLevel)
	}
	Log.SetLevel(level)

	// Set log format
	if cfg.LogFormat == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set output
	switch cfg.LogOutput {
	case "file":
		file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			Log.SetOutput(file)
		} else {
			Log.Warn("Failed to log to file, using default stderr")
		}
	case "both":
		// Log to both file and stdout (for now, just log to stdout)
		// Note: For production, you might want to use io.MultiWriter here
		Log.SetOutput(os.Stdout)
	default:
		Log.SetOutput(os.Stdout)
	}

	// Enable caller reporting for better debugging
	Log.SetReportCaller(true)
}

// SanitizeFields removes sensitive information from log fields
func SanitizeFields(fields logrus.Fields) logrus.Fields {
	sanitized := make(logrus.Fields)
	sensitiveKeys := []string{
		"password", "token", "secret", "api_key", "apikey",
		"authorization", "auth", "credential", "private_key",
	}

	for key, value := range fields {
		keyLower := strings.ToLower(key)
		isSensitive := false

		for _, sensitiveKey := range sensitiveKeys {
			if strings.Contains(keyLower, sensitiveKey) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			sanitized[key] = "[REDACTED]"
		} else {
			sanitized[key] = value
		}
	}

	return sanitized
}

// WithContext adds common context from Gin context to log fields
func WithContext(c *gin.Context) *logrus.Entry {
	fields := logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}

	// Add request ID if exists
	if requestID := c.GetString("request_id"); requestID != "" {
		fields["request_id"] = requestID
	}

	// Add user ID if exists
	if userID, exists := c.Get("user_id"); exists {
		fields["user_id"] = userID
	}

	return Log.WithFields(SanitizeFields(fields))
}

// WithFields creates a log entry with custom fields (sanitized)
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(SanitizeFields(fields))
}

// WithField creates a log entry with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return WithFields(logrus.Fields{key: value})
}

// GetCaller returns the file and line number of the caller
func GetCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	// Get just the filename, not the full path
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return file + ":" + string(rune(line+'0'))
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

// LogError logs an error with context
func LogError(err error, message string, fields logrus.Fields) {
	if err == nil {
		return
	}

	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["error"] = err.Error()

	Log.WithFields(SanitizeFields(fields)).Error(message)
}

// LogRequest logs HTTP request details
func LogRequest(c *gin.Context, message string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	// Add standard request fields
	fields["method"] = c.Request.Method
	fields["path"] = c.Request.URL.Path
	fields["ip"] = c.ClientIP()
	fields["user_agent"] = c.Request.UserAgent()

	// Add request ID if exists
	if requestID := c.GetString("request_id"); requestID != "" {
		fields["request_id"] = requestID
	}

	// Add user ID if exists
	if userID, exists := c.Get("user_id"); exists {
		fields["user_id"] = userID
	}

	Log.WithFields(SanitizeFields(fields)).Info(message)
}
