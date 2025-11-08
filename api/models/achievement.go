package models

import (
	"time"

	"github.com/google/uuid"
)

// Achievement represents an achievement that users can unlock
type Achievement struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Key          string    `json:"key" db:"key"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	Icon         *string   `json:"icon,omitempty" db:"icon"`
	Color        string    `json:"color" db:"color"`
	PointsReward int       `json:"points_reward" db:"points_reward"`
	Category     string    `json:"category" db:"category"`
	IsRare       bool      `json:"is_rare" db:"is_rare"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserAchievement represents a user's unlocked achievement
type UserAchievement struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	UserID        uuid.UUID    `json:"user_id" db:"user_id"`
	AchievementID uuid.UUID    `json:"achievement_id" db:"achievement_id"`
	UnlockedAt    time.Time    `json:"unlocked_at" db:"unlocked_at"`
	PointsAwarded int          `json:"points_awarded" db:"points_awarded"`
	Achievement   *Achievement `json:"achievement,omitempty"`
}

// Achievement DTOs

// AchievementListResponse represents the response for user achievements
type AchievementListResponse struct {
	Achievements []UserAchievement `json:"achievements"`
	Total        int               `json:"total"`
}

// AchievementNotification represents a notification for a newly unlocked achievement
type AchievementNotification struct {
	AchievementID uuid.UUID `json:"achievement_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Icon          *string   `json:"icon,omitempty"`
	Color         string    `json:"color"`
	PointsAwarded int       `json:"points_awarded"`
	IsRare        bool      `json:"is_rare"`
}

// AchievementProgress represents progress towards an achievement
type AchievementProgress struct {
	AchievementID uuid.UUID `json:"achievement_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Icon          *string   `json:"icon,omitempty"`
	Color         string    `json:"color"`
	Category      string    `json:"category"`
	IsRare        bool      `json:"is_rare"`
	CurrentValue  int       `json:"current_value"`
	TargetValue   int       `json:"target_value"`
	Progress      float64   `json:"progress"` // Percentage (0-100)
	IsUnlocked    bool      `json:"is_unlocked"`
}

// AchievementProgressResponse represents the response for achievement progress
type AchievementProgressResponse struct {
	Progress []AchievementProgress `json:"progress"`
	Total    int                   `json:"total"`
}