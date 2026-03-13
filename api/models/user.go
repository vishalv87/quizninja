package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	Email                 string    `json:"email" db:"email"`
	PasswordHash          string    `json:"-" db:"password_hash"`
	Name                  string    `json:"name" db:"name"`
	Level                 string    `json:"level" db:"level"`
	TotalPoints           int       `json:"total_points" db:"total_points"`
	CurrentStreak         int       `json:"current_streak" db:"current_streak"`
	BestStreak            int       `json:"best_streak" db:"best_streak"`
	TotalQuizzesCompleted int       `json:"total_quizzes_completed" db:"total_quizzes_completed"`
	AverageScore          float64   `json:"average_score" db:"average_score"`
	IsOnline              bool      `json:"is_online" db:"is_online"`
	LastActive            time.Time `json:"last_active" db:"last_active"`
	AvatarURL             *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`

	// Supabase auth integration fields
	AuthMethod     string     `json:"auth_method" db:"auth_method"`           // "supabase" or "jwt"
	SupabaseID     *string    `json:"supabase_id,omitempty" db:"supabase_id"` // Supabase user ID
	LastAuthMethod string     `json:"last_auth_method" db:"last_auth_method"` // Track last successful auth
	MigratedAt     *time.Time `json:"migrated_at,omitempty" db:"migrated_at"` // When user was migrated between auth systems

	Preferences *UserPreferences `json:"preferences,omitempty"`
}

type UserPreferences struct {
	ID                    uuid.UUID              `json:"id" db:"id"`
	UserID                uuid.UUID              `json:"user_id" db:"user_id"`
	SelectedCategories    StringArray            `json:"selected_categories" db:"selected_categories"`
	DifficultyPreference  string                 `json:"difficulty_preference" db:"difficulty_preference"`
	NotificationsEnabled  bool                   `json:"notifications_enabled" db:"notifications_enabled"`
	NotificationFrequency string                 `json:"notification_frequency" db:"notification_frequency"`
	ProfileVisibility     bool                   `json:"profile_visibility" db:"profile_visibility"`
	ShowOnlineStatus      bool                   `json:"show_online_status" db:"show_online_status"`
	AllowFriendRequests   bool                   `json:"allow_friend_requests" db:"allow_friend_requests"`
	ShareActivityStatus   bool                   `json:"share_activity_status" db:"share_activity_status"`
	NotificationTypes     map[string]interface{} `json:"notification_types" db:"notification_types"`
	OnboardingCompletedAt *time.Time             `json:"onboarding_completed_at,omitempty" db:"onboarding_completed_at"`
	CreatedAt             time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at" db:"updated_at"`
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