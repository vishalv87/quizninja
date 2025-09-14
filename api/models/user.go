package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	Email                 string     `json:"email" db:"email"`
	PasswordHash          string     `json:"-" db:"password_hash"`
	Name                  string     `json:"name" db:"name"`
	Age                   *int       `json:"age,omitempty" db:"age"`
	Level                 string     `json:"level" db:"level"`
	TotalPoints           int        `json:"total_points" db:"total_points"`
	CurrentStreak         int        `json:"current_streak" db:"current_streak"`
	BestStreak            int        `json:"best_streak" db:"best_streak"`
	TotalQuizzesCompleted int        `json:"total_quizzes_completed" db:"total_quizzes_completed"`
	AverageScore          float64    `json:"average_score" db:"average_score"`
	IsOnline              bool       `json:"is_online" db:"is_online"`
	LastActive            time.Time  `json:"last_active" db:"last_active"`
	AvatarURL             *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
	Preferences           *UserPreferences `json:"preferences,omitempty"`
}

type UserPreferences struct {
	ID                      uuid.UUID    `json:"id" db:"id"`
	UserID                  uuid.UUID    `json:"user_id" db:"user_id"`
	SelectedInterests       StringArray  `json:"selected_interests" db:"selected_interests"`
	DifficultyPreference    string       `json:"difficulty_preference" db:"difficulty_preference"`
	NotificationsEnabled    bool         `json:"notifications_enabled" db:"notifications_enabled"`
	NotificationFrequency   string       `json:"notification_frequency" db:"notification_frequency"`
	CreatedAt               time.Time    `json:"created_at" db:"created_at"`
}

// StringArray handles PostgreSQL array type
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return pq.Array(a).Value()
}

func (a *StringArray) Scan(value any) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}
	return pq.Array(a).Scan(value)
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type RegisterRequest struct {
	Email       string                     `json:"email" binding:"required,email"`
	Password    string                     `json:"password" binding:"required,min=6"`
	Name        string                     `json:"name" binding:"required"`
	Age         *int                       `json:"age,omitempty"`
	Preferences *UserPreferencesRequest    `json:"preferences,omitempty"`
}

type UserPreferencesRequest struct {
	SelectedInterests       []string `json:"selected_interests"`
	DifficultyPreference    string   `json:"difficulty_preference" binding:"omitempty,oneof=Easy Medium Hard"`
	NotificationsEnabled    bool     `json:"notifications_enabled"`
	NotificationFrequency   string   `json:"notification_frequency" binding:"omitempty,oneof=Daily Weekly Never"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Interest represents quiz topic interests
type Interest struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"display_name" db:"display_name"`
	Description string `json:"description" db:"description"`
	IconURL     string `json:"icon_url" db:"icon_url"`
	IsActive    bool   `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// DifficultyLevel represents quiz difficulty levels
type DifficultyLevel struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"display_name" db:"display_name"`
	Description string `json:"description" db:"description"`
	Order       int    `json:"order" db:"order"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// NotificationFrequency represents notification frequency options
type NotificationFrequency struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"display_name" db:"display_name"`
	Description string `json:"description" db:"description"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// AppSettings represents application configuration settings
type AppSettings struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Key              string    `json:"key" db:"key"`
	Value            string    `json:"value" db:"value"`
	Description      string    `json:"description" db:"description"`
	Category         string    `json:"category" db:"category"`
	IsPublic         bool      `json:"is_public" db:"is_public"`
	LastModifiedBy   uuid.UUID `json:"last_modified_by" db:"last_modified_by"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Category represents quiz categories for API response
type Category struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Description string      `json:"description"`
	IconURL     string      `json:"icon_url"`
	IsActive    bool        `json:"is_active"`
	Interests   []Interest  `json:"interests,omitempty"`
}

// API Request/Response types
type UpdatePreferencesRequest struct {
	SelectedInterests       []string `json:"selected_interests" binding:"dive,required"`
	DifficultyPreference    string   `json:"difficulty_preference" binding:"required,oneof=Easy Medium Hard"`
	NotificationsEnabled    bool     `json:"notifications_enabled"`
	NotificationFrequency   string   `json:"notification_frequency" binding:"required,oneof=Daily Weekly Never"`
}

type AppSettingsResponse struct {
	Settings map[string]any `json:"settings"`
	Version  string         `json:"version"`
	LastUpdated time.Time   `json:"last_updated"`
}

// Quiz models

type Quiz struct {
	ID                uuid.UUID    `json:"id" db:"id"`
	Title             string       `json:"title" db:"title"`
	Description       string       `json:"description" db:"description"`
	Category          string       `json:"category" db:"category"`
	Difficulty        string       `json:"difficulty" db:"difficulty"`
	TimeLimit         int          `json:"time_limit" db:"time_limit"` // in seconds
	QuestionCount     int          `json:"question_count" db:"question_count"`
	IsFeatured        bool         `json:"is_featured" db:"is_featured"`
	IsPublic          bool         `json:"is_public" db:"is_public"`
	CreatedBy         uuid.UUID    `json:"created_by" db:"created_by"`
	Tags              StringArray  `json:"tags" db:"tags"`
	ThumbnailURL      *string      `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	CreatedAt         time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at" db:"updated_at"`
	Questions         []Question   `json:"questions,omitempty"`
	Statistics        *QuizStatistics `json:"statistics,omitempty"`
}

type Question struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	QuizID        uuid.UUID    `json:"quiz_id" db:"quiz_id"`
	QuestionText  string       `json:"question_text" db:"question_text"`
	QuestionType  string       `json:"question_type" db:"question_type"` // multiple_choice, true_false, short_answer
	Options       StringArray  `json:"options,omitempty" db:"options"`
	CorrectAnswer string       `json:"correct_answer" db:"correct_answer"`
	Explanation   *string      `json:"explanation,omitempty" db:"explanation"`
	Points        int          `json:"points" db:"points"`
	Order         int          `json:"order" db:"order"`
	ImageURL      *string      `json:"image_url,omitempty" db:"image_url"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

type QuizStatistics struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	QuizID             uuid.UUID `json:"quiz_id" db:"quiz_id"`
	TotalAttempts      int       `json:"total_attempts" db:"total_attempts"`
	CompletedAttempts  int       `json:"completed_attempts" db:"completed_attempts"`
	AverageScore       float64   `json:"average_score" db:"average_score"`
	AverageTime        int       `json:"average_time" db:"average_time"` // in seconds
	HighestScore       float64   `json:"highest_score" db:"highest_score"`
	LowestScore        float64   `json:"lowest_score" db:"lowest_score"`
	LastAttemptAt      *time.Time `json:"last_attempt_at,omitempty" db:"last_attempt_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// Quiz DTOs

type CreateQuizRequest struct {
	Title         string                   `json:"title" binding:"required,min=1,max=200"`
	Description   string                   `json:"description" binding:"required,min=1,max=1000"`
	Category      string                   `json:"category" binding:"required"`
	Difficulty    string                   `json:"difficulty" binding:"required,oneof=Easy Medium Hard"`
	TimeLimit     int                      `json:"time_limit" binding:"required,min=30,max=3600"`
	IsFeatured    bool                     `json:"is_featured"`
	IsPublic      bool                     `json:"is_public"`
	Tags          []string                 `json:"tags"`
	ThumbnailURL  *string                  `json:"thumbnail_url,omitempty"`
	Questions     []CreateQuestionRequest  `json:"questions" binding:"required,dive"`
}

type CreateQuestionRequest struct {
	QuestionText  string   `json:"question_text" binding:"required,min=1,max=500"`
	QuestionType  string   `json:"question_type" binding:"required,oneof=multiple_choice true_false short_answer"`
	Options       []string `json:"options,omitempty"`
	CorrectAnswer string   `json:"correct_answer" binding:"required"`
	Explanation   *string  `json:"explanation,omitempty"`
	Points        int      `json:"points" binding:"required,min=1,max=100"`
	Order         int      `json:"order" binding:"required,min=1"`
	ImageURL      *string  `json:"image_url,omitempty"`
}

type UpdateQuizRequest struct {
	Title         *string  `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Description   *string  `json:"description,omitempty" binding:"omitempty,min=1,max=1000"`
	Category      *string  `json:"category,omitempty"`
	Difficulty    *string  `json:"difficulty,omitempty" binding:"omitempty,oneof=Easy Medium Hard"`
	TimeLimit     *int     `json:"time_limit,omitempty" binding:"omitempty,min=30,max=3600"`
	IsFeatured    *bool    `json:"is_featured,omitempty"`
	IsPublic      *bool    `json:"is_public,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	ThumbnailURL  *string  `json:"thumbnail_url,omitempty"`
}

type QuizListResponse struct {
	Quizzes    []QuizSummary `json:"quizzes"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

type QuizSummary struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Difficulty    string    `json:"difficulty"`
	TimeLimit     int       `json:"time_limit"`
	QuestionCount int       `json:"question_count"`
	IsFeatured    bool      `json:"is_featured"`
	Tags          []string  `json:"tags"`
	ThumbnailURL  *string   `json:"thumbnail_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Statistics    *QuizStatisticsSummary `json:"statistics,omitempty"`
}

type QuizStatisticsSummary struct {
	TotalAttempts int     `json:"total_attempts"`
	AverageScore  float64 `json:"average_score"`
	AverageTime   int     `json:"average_time"`
}

type QuizDetailResponse struct {
	Quiz Quiz `json:"quiz"`
}

type QuizFilters struct {
	Category   string `form:"category"`
	Difficulty string `form:"difficulty" binding:"omitempty,oneof=Easy Medium Hard"`
	Featured   *bool  `form:"featured"`
	Tags       string `form:"tags"`
	Search     string `form:"search"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=10" binding:"min=1,max=100"`
}

type QuizAttempt struct {
	ID          uuid.UUID `json:"id" db:"id"`
	QuizID      uuid.UUID `json:"quiz_id" db:"quiz_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Score       float64   `json:"score" db:"score"`
	TotalPoints int       `json:"total_points" db:"total_points"`
	TimeSpent   int       `json:"time_spent" db:"time_spent"` // in seconds
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
	StartedAt   time.Time `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}