package models

// RegisterRequest for frontend-initiated Supabase auth
type RegisterRequest struct {
	// Supabase user data (received after frontend auth)
	SupabaseUserID string                  `json:"supabase_user_id" binding:"required"`
	Email          string                  `json:"email" binding:"required,email"`
	Name           string                  `json:"name" binding:"required"`
	Preferences    *UserPreferencesRequest `json:"preferences,omitempty"`

	// Optional: additional user metadata
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type UserPreferencesRequest struct {
	SelectedCategories    []string `json:"selected_categories"`
	DifficultyPreference  string   `json:"difficulty_preference" binding:"omitempty,oneof=Easy Medium Hard"`
	NotificationsEnabled  bool     `json:"notifications_enabled"`
	NotificationFrequency string   `json:"notification_frequency" binding:"omitempty,oneof=Daily Weekly Never"`
}

// LoginRequest for frontend-initiated Supabase auth (user sync)
type LoginRequest struct {
	// Supabase user data (received after frontend auth)
	SupabaseUserID string  `json:"supabase_user_id" binding:"required"`
	Email          string  `json:"email" binding:"required,email"`
	Name           string  `json:"name" binding:"required"`
	AvatarURL      *string `json:"avatar_url,omitempty"`
}

// AuthResponse for frontend-initiated auth (returns user profile only)
type AuthResponse struct {
	User    User   `json:"user"`
	Message string `json:"message,omitempty"`
}

type UpdateProfileRequest struct {
	Name      *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}
