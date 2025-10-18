package utils

import (
	"html"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Validation error messages
const (
	ErrInvalidEmail         = "Invalid email format"
	ErrInvalidName          = "Name must be 2-100 characters and contain only letters, numbers, spaces, and common punctuation"
	ErrInvalidURL           = "Invalid URL format"
	ErrMessageTooLong       = "Message exceeds maximum length of 500 characters"
	ErrSuspiciousContent    = "Content contains potentially malicious patterns"
	ErrInvalidStringLength  = "String length is invalid"
)

var (
	// Email validation regex (RFC 5322 compliant)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// Name validation: letters, numbers, spaces, and common punctuation
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,']+$`)

	// SQL injection patterns
	sqlInjectionPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(union\s+select|insert\s+into|delete\s+from|drop\s+table|update\s+.+\s+set)`),
		regexp.MustCompile(`(?i)(exec\s*\(|execute\s*\(|script\s*>)`),
		regexp.MustCompile(`[';]--`),
		regexp.MustCompile(`\bor\b\s+[0-9]+\s*=\s*[0-9]+`),
		regexp.MustCompile(`\band\b\s+[0-9]+\s*=\s*[0-9]+`),
	}

	// XSS patterns
	xssPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
		regexp.MustCompile(`(?i)<iframe[^>]*>.*?</iframe>`),
		regexp.MustCompile(`(?i)javascript:`),
		regexp.MustCompile(`(?i)on(load|error|click|mouse\w+)\s*=`),
		regexp.MustCompile(`(?i)<embed[^>]*>`),
		regexp.MustCompile(`(?i)<object[^>]*>`),
	}
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}

	if len(email) > 254 {
		return &ValidationError{Field: "email", Message: "Email exceeds maximum length of 254 characters"}
	}

	if !emailRegex.MatchString(email) {
		return &ValidationError{Field: "email", Message: ErrInvalidEmail}
	}

	return nil
}

// ValidateName validates name/username format and length
func ValidateName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return &ValidationError{Field: "name", Message: "Name is required"}
	}

	length := utf8.RuneCountInString(name)
	if length < 2 || length > 100 {
		return &ValidationError{Field: "name", Message: "Name must be between 2 and 100 characters"}
	}

	if !nameRegex.MatchString(name) {
		return &ValidationError{Field: "name", Message: ErrInvalidName}
	}

	return nil
}

// ValidateURL validates URL format
func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return nil // Empty URLs are allowed for optional fields
	}

	urlStr = strings.TrimSpace(urlStr)

	if len(urlStr) > 2048 {
		return &ValidationError{Field: "url", Message: "URL exceeds maximum length of 2048 characters"}
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return &ValidationError{Field: "url", Message: ErrInvalidURL}
	}

	// Ensure it's an HTTP or HTTPS URL
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return &ValidationError{Field: "url", Message: "URL must use HTTP or HTTPS scheme"}
	}

	if parsedURL.Host == "" {
		return &ValidationError{Field: "url", Message: ErrInvalidURL}
	}

	return nil
}

// ValidateMessage validates message content
func ValidateMessage(message string) error {
	if message == "" {
		return nil // Empty messages are allowed for optional fields
	}

	message = strings.TrimSpace(message)

	length := utf8.RuneCountInString(message)
	if length > 500 {
		return &ValidationError{Field: "message", Message: ErrMessageTooLong}
	}

	// Check for XSS patterns
	if ContainsXSS(message) {
		return &ValidationError{Field: "message", Message: ErrSuspiciousContent}
	}

	// Check for SQL injection patterns
	if ContainsSQLInjection(message) {
		return &ValidationError{Field: "message", Message: ErrSuspiciousContent}
	}

	return nil
}

// ValidateStringLength validates string length within min and max bounds
func ValidateStringLength(value string, min, max int, fieldName string) error {
	length := utf8.RuneCountInString(value)

	if length < min || length > max {
		return &ValidationError{
			Field:   fieldName,
			Message: fieldName + " must be between " + string(rune(min+'0')) + " and " + string(rune(max+'0')) + " characters",
		}
	}

	return nil
}

// SanitizeString removes potentially dangerous characters and trims whitespace
func SanitizeString(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove control characters except newline and tab
	var sanitized strings.Builder
	for _, r := range input {
		if r == '\n' || r == '\t' || r >= 32 {
			sanitized.WriteRune(r)
		}
	}

	return sanitized.String()
}

// SanitizeHTML encodes HTML special characters to prevent XSS
func SanitizeHTML(input string) string {
	// First sanitize as string
	input = SanitizeString(input)

	// HTML encode special characters
	input = html.EscapeString(input)

	return input
}

// StripHTML removes all HTML tags from input
func StripHTML(input string) string {
	// Remove all HTML tags
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	input = htmlTagRegex.ReplaceAllString(input, "")

	// Decode HTML entities
	input = html.UnescapeString(input)

	// Sanitize the result
	return SanitizeString(input)
}

// ContainsXSS checks if string contains potential XSS patterns
func ContainsXSS(input string) bool {
	lowercaseInput := strings.ToLower(input)

	for _, pattern := range xssPatterns {
		if pattern.MatchString(lowercaseInput) {
			return true
		}
	}

	return false
}

// ContainsSQLInjection checks if string contains potential SQL injection patterns
func ContainsSQLInjection(input string) bool {
	for _, pattern := range sqlInjectionPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}

	return false
}

// SanitizeEmail sanitizes and normalizes email address
func SanitizeEmail(email string) string {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	return SanitizeString(email)
}

// SanitizeName sanitizes name while preserving valid characters
func SanitizeName(name string) string {
	name = strings.TrimSpace(name)

	// Remove consecutive spaces
	spaceRegex := regexp.MustCompile(`\s+`)
	name = spaceRegex.ReplaceAllString(name, " ")

	return SanitizeString(name)
}
