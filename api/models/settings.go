package models

import (
	"time"

	"github.com/google/uuid"
)

// Category represents quiz topic categories
type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Description string    `json:"description" db:"description"`
	IconURL     string    `json:"icon_url" db:"icon_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	QuizCount   int       `json:"quiz_count" db:"quiz_count"`
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
	ID             uuid.UUID `json:"id" db:"id"`
	Key            string    `json:"key" db:"key"`
	Value          string    `json:"value" db:"value"`
	Description    string    `json:"description" db:"description"`
	Category       string    `json:"category" db:"category"`
	IsPublic       bool      `json:"is_public" db:"is_public"`
	LastModifiedBy uuid.UUID `json:"last_modified_by" db:"last_modified_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CategoryGroup represents quiz category groups for API response
type CategoryGroup struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	DisplayName string     `json:"display_name"`
	Description string     `json:"description"`
	IconURL     string     `json:"icon_url"`
	IsActive    bool       `json:"is_active"`
	Categories  []Category `json:"categories,omitempty"`
}

// UpdatePreferencesRequest represents the request to update user preferences
type UpdatePreferencesRequest struct {
	SelectedCategories    []string               `json:"selected_categories" binding:"dive,required"`
	DifficultyPreference  string                 `json:"difficulty_preference" binding:"required,oneof=Easy Medium Hard"`
	NotificationsEnabled  bool                   `json:"notifications_enabled"`
	NotificationFrequency string                 `json:"notification_frequency" binding:"required,oneof=Daily Weekly Never"`
	ProfileVisibility     *bool                  `json:"profile_visibility,omitempty"`
	ShowOnlineStatus      *bool                  `json:"show_online_status,omitempty"`
	AllowFriendRequests   *bool                  `json:"allow_friend_requests,omitempty"`
	ShareActivityStatus   *bool                  `json:"share_activity_status,omitempty"`
	NotificationTypes     map[string]interface{} `json:"notification_types,omitempty"`
}

type AppSettingsResponse struct {
	Settings    map[string]any `json:"settings"`
	Version     string         `json:"version"`
	LastUpdated time.Time      `json:"last_updated"`
}
