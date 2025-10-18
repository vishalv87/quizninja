package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(format string, args ...interface{}) *ValidationError {
	return &ValidationError{
		Field:   "",
		Message: fmt.Sprintf(format, args...),
	}
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// ForbiddenError represents a forbidden access error
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{
		Message: message,
	}
}

// APIError represents a structured API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ErrorResponse sends a structured error response
func ErrorResponse(c *gin.Context, statusCode int, message string, details ...string) {
	error := APIError{
		Code:    statusCode,
		Message: message,
	}

	if len(details) > 0 {
		error.Details = details[0]
	}

	c.JSON(statusCode, gin.H{"error": error})
}

// HandleError handles different types of errors and sends appropriate responses
func HandleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *ValidationError:
		ErrorResponse(c, http.StatusBadRequest, "Validation error", e.Message)
	case *NotFoundError:
		ErrorResponse(c, http.StatusNotFound, "Resource not found", e.Error())
	case *ForbiddenError:
		ErrorResponse(c, http.StatusForbidden, "Access forbidden", e.Message)
	default:
		ErrorResponse(c, http.StatusInternalServerError, "Internal server error", e.Error())
	}
}

// SuccessResponse sends a structured success response
func SuccessResponse(c *gin.Context, data interface{}, message ...string) {
	response := gin.H{"data": data}

	if len(message) > 0 {
		response["message"] = message[0]
	}

	c.JSON(http.StatusOK, response)
}

// CreatedResponse sends a structured success response for created resources
func CreatedResponse(c *gin.Context, data interface{}, message ...string) {
	response := gin.H{"data": data}

	if len(message) > 0 {
		response["message"] = message[0]
	} else {
		response["message"] = "Resource created successfully"
	}

	c.JSON(http.StatusCreated, response)
}
