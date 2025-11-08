package models

import "time"

// OnboardingCompleteRequest represents the request to complete onboarding
type OnboardingCompleteRequest struct {
	SelectedCategories    []string `json:"selected_categories" binding:"required,dive,required"`
	DifficultyPreference  string   `json:"difficulty_preference" binding:"required,oneof=Easy Medium Hard"`
	NotificationsEnabled  bool     `json:"notifications_enabled"`
	NotificationFrequency string   `json:"notification_frequency" binding:"required,oneof=Daily Weekly Never"`
}

// OnboardingStatusResponse represents the response for onboarding status
type OnboardingStatusResponse struct {
	IsCompleted bool             `json:"is_completed"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
	Preferences *UserPreferences `json:"preferences,omitempty"`
}