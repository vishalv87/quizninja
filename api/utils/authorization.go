package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthorizationError represents an authorization failure
type AuthorizationError struct {
	Message string
	Code    int
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

// GetUserIDFromContext safely extracts user ID from context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, &AuthorizationError{
			Message: "User not authenticated",
			Code:    http.StatusUnauthorized,
		}
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil, &AuthorizationError{
			Message: "Invalid user ID in context",
			Code:    http.StatusUnauthorized,
		}
	}

	return userID, nil
}

// RequireOwnership verifies that the authenticated user owns the resource
func RequireOwnership(c *gin.Context, resourceOwnerID uuid.UUID, resourceType string) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	if userID != resourceOwnerID {
		return &AuthorizationError{
			Message: "You don't have permission to access this " + resourceType,
			Code:    http.StatusForbidden,
		}
	}

	return nil
}

// RequireAnyOwnership verifies user owns at least one of the provided IDs
func RequireAnyOwnership(c *gin.Context, resourceOwnerIDs []uuid.UUID, resourceType string) error {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	for _, ownerID := range resourceOwnerIDs {
		if userID == ownerID {
			return nil
		}
	}

	return &AuthorizationError{
		Message: "You don't have permission to access this " + resourceType,
		Code:    http.StatusForbidden,
	}
}

// HandleAuthError sends appropriate error response for authorization errors
func HandleAuthError(c *gin.Context, err error) bool {
	if authErr, ok := err.(*AuthorizationError); ok {
		ErrorResponse(c, authErr.Code, authErr.Message)
		return true
	}
	return false
}
