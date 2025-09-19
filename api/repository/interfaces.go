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

// QuizRepositoryInterface defines the contract for quiz data operations
type QuizRepositoryInterface interface {
	// Quiz read operations
	GetQuizByID(id uuid.UUID) (*models.Quiz, error)
	GetQuizByIDWithQuestions(id uuid.UUID) (*models.Quiz, error)
	GetQuizByIDWithStatistics(id uuid.UUID) (*models.Quiz, error)
	GetQuizByIDWithAll(id uuid.UUID) (*models.Quiz, error)

	// Quiz list operations with filtering and pagination
	GetQuizzes(filters *models.QuizFilters) ([]models.Quiz, int, error)
	GetFeaturedQuizzes(limit int) ([]models.Quiz, error)
	GetQuizzesByCategory(category string, limit int) ([]models.Quiz, error)
	GetQuizzesByUser(userID uuid.UUID, offset, limit int) ([]models.Quiz, int, error)

	// Question read operations
	GetQuestionsByQuizID(quizID uuid.UUID) ([]models.Question, error)

	// Quiz statistics read operations
	GetQuizStatistics(quizID uuid.UUID) (*models.QuizStatistics, error)

	// Quiz attempt operations
	CreateQuizAttempt(attempt *models.QuizAttempt) error
	UpdateQuizAttempt(attempt *models.QuizAttempt) error
	GetQuizAttempt(id uuid.UUID) (*models.QuizAttempt, error)
	GetUserQuizAttempts(userID, quizID uuid.UUID) ([]models.QuizAttempt, error)
	GetActiveQuizAttempt(userID, quizID uuid.UUID) (*models.QuizAttempt, error)
	DeleteActiveQuizAttempt(userID, quizID uuid.UUID) error
}

// Repository aggregates all repository interfaces
type Repository struct {
	User UserRepositoryInterface
	Quiz QuizRepositoryInterface
}

// NewRepository creates a new repository instance
func NewRepository() *Repository {
	return &Repository{
		User: NewUserRepository(),
		Quiz: NewQuizRepository(),
	}
}