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