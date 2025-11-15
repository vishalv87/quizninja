package models

import (
	"time"

	"github.com/google/uuid"
)

// UserProfileResponse represents a user's profile with privacy-aware data
// Used when viewing another user's profile
type UserProfileResponse struct {
	// Basic profile info (always visible based on privacy settings)
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"` // Same as ID, for frontend compatibility
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"` // Same as Name, for frontend compatibility
	Email     string    `json:"email"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Bio       *string   `json:"bio,omitempty"` // Currently not in DB schema, but added for future
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Stats (null if user has hidden stats in privacy settings)
	Stats *UserStats `json:"stats,omitempty"`

	// Privacy settings (subset visible to others)
	Preferences *UserProfilePreferences `json:"preferences,omitempty"`

	// Friendship status
	IsFriend            bool   `json:"is_friend"`
	FriendRequestStatus string `json:"friend_request_status"` // "none", "pending_sent", "pending_received", "friends"
}

// UserStats represents user statistics
type UserStats struct {
	UserID                uuid.UUID `json:"user_id"`
	TotalQuizzesTaken     int       `json:"total_quizzes_taken"`
	TotalQuizzesCompleted int       `json:"total_quizzes_completed"`
	TotalPoints           int       `json:"total_points"`
	AverageScore          float64   `json:"average_score"`
	TotalTimeSpentMinutes int       `json:"total_time_spent_minutes"`
	CurrentStreak         int       `json:"current_streak"`
	LongestStreak         int       `json:"longest_streak"`
	AchievementsUnlocked  int       `json:"achievements_unlocked"`
	ChallengesWon         int       `json:"challenges_won"`
	ChallengesLost        int       `json:"challenges_lost"`
	Rank                  int       `json:"rank"`
}

// UserProfilePreferences represents privacy settings visible to others
type UserProfilePreferences struct {
	ProfileVisibility  bool `json:"profile_visibility"`  // If false, profile is private
	ShowAchievements   bool `json:"show_achievements"`   // Not in current schema, default true
	ShowStats          bool `json:"show_stats"`          // Not in current schema, default true
	AllowFriendRequest bool `json:"allow_friend_request"` // From user_preferences
}
