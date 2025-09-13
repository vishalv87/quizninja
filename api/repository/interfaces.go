package repository

import (
	"quizninja-api/models"

	"github.com/google/uuid"
)

// UserRepositoryInterface defines the contract for user data operations
type UserRepositoryInterface interface {
	// User CRUD operations
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uuid.UUID) error

	// User preferences operations
	CreateUserPreferences(preferences *models.UserPreferences) error
	GetUserPreferences(userID uuid.UUID) (*models.UserPreferences, error)
	UpdateUserPreferences(preferences *models.UserPreferences) error
	DeleteUserPreferences(userID uuid.UUID) error

	// User with preferences operations
	GetUserWithPreferences(userID uuid.UUID) (*models.User, error)

	// Refresh token operations
	SaveRefreshToken(refreshToken *models.RefreshToken) error
	GetRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteUserRefreshTokens(userID uuid.UUID) error

	// User status operations
	UpdateUserOnlineStatus(userID uuid.UUID, isOnline bool) error
	UpdateUserLastActive(userID uuid.UUID) error
}

// Repository aggregates all repository interfaces
type Repository struct {
	User UserRepositoryInterface
}

// NewRepository creates a new repository instance
func NewRepository() *Repository {
	return &Repository{
		User: NewUserRepository(),
	}
}