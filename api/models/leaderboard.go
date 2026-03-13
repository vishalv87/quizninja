package models

import (
	"time"

	"github.com/google/uuid"
)

// LeaderboardEntry represents a user's position in the leaderboard
type LeaderboardEntry struct {
	UserID           uuid.UUID      `json:"user_id" db:"user_id"`
	Name             string         `json:"name" db:"name"`
	Avatar           *string        `json:"avatar" db:"avatar_url"`
	Rank             int            `json:"rank"`
	Points           int            `json:"points" db:"total_points"`
	QuizzesCompleted int            `json:"quizzes_completed" db:"total_quizzes_completed"`
	AverageScore     float64        `json:"average_score" db:"average_score"`
	CurrentStreak    int            `json:"current_streak" db:"current_streak"`
	Level            string         `json:"level" db:"level"`
	IsCurrentUser    bool           `json:"is_current_user"`
	IsFriend         bool           `json:"is_friend"`
	LastActive       time.Time      `json:"last_active" db:"last_active"`
	Achievements     []string       `json:"achievements"`
	CategoryPoints   map[string]int `json:"category_points"`
}

// LeaderboardFilters represents filters for leaderboard queries
type LeaderboardFilters struct {
	Period      string `form:"period,default=alltime" binding:"omitempty,oneof=today week month alltime"`
	FriendsOnly bool   `form:"friends_only,default=false"`
	Limit       int    `form:"limit,default=50" binding:"min=1,max=100"`
	Offset      int    `form:"offset,default=0" binding:"min=0"`
}

// LeaderboardResponse represents the API response for leaderboard
type LeaderboardResponse struct {
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	UserRank    *UserRankInfo      `json:"user_rank,omitempty"`
	Total       int                `json:"total"`
	Period      string             `json:"period"`
	FriendsOnly bool               `json:"friends_only"`
}

// UserRankInfo represents the current user's rank information
type UserRankInfo struct {
	Rank                 int       `json:"rank"`
	TotalUsers           int       `json:"total_users"`
	User                 UserInfo  `json:"user"`
	TotalPoints          int       `json:"total_points"`
	QuizzesCompleted     int       `json:"quizzes_completed"`
	AchievementsUnlocked int       `json:"achievements_unlocked"`
	PointsToNext         int       `json:"points_to_next"`
	RankChange           int       `json:"rank_change"` // +/- change from previous period
}

// UserInfo represents user profile information in rank data
type UserInfo struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}
