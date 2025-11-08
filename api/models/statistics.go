package models

import (
	"time"

	"github.com/google/uuid"
)

// UserStatistics represents comprehensive user statistics
type UserStatistics struct {
	UserID                uuid.UUID             `json:"user_id"`
	TotalAttempts         int                   `json:"total_attempts"`
	CompletedQuizzes      int                   `json:"total_quizzes_completed"`
	CompletionRate        float64               `json:"completion_rate"`
	AverageScore          float64               `json:"average_score"`
	TotalPoints           int                   `json:"total_points"`
	CurrentStreak         int                   `json:"current_streak"`
	BestStreak            int                   `json:"best_streak"`
	AverageCompletionTime int                   `json:"average_completion_time"` // in seconds
	CategoryPerformance   []CategoryPerformance `json:"category_performance"`
	RecentActivity        []RecentActivityItem  `json:"recent_activity"`
	LastUpdated           time.Time             `json:"last_updated"`
	QuizzesByDifficulty   map[string]int        `json:"quizzes_by_difficulty"`
	ScoreDistribution     ScoreDistribution     `json:"score_distribution"`
	MonthlyProgress       []MonthlyProgressItem `json:"monthly_progress"`
}

// CategoryPerformance represents performance metrics for a specific category
type CategoryPerformance struct {
	CategoryID       string     `json:"category_id"`
	CategoryName     string     `json:"category_name"`
	QuizzesCompleted int        `json:"quizzes_completed"`
	AverageScore     float64    `json:"average_score"`
	TotalAttempts    int        `json:"total_attempts"`
	BestScore        float64    `json:"best_score"`
	LastAttempt      *time.Time `json:"last_attempt,omitempty"`
}

// RecentActivityItem represents a recent quiz activity
type RecentActivityItem struct {
	QuizID      uuid.UUID `json:"quiz_id"`
	QuizTitle   string    `json:"quiz_title"`
	Score       float64   `json:"score"`
	Category    string    `json:"category"`
	CompletedAt time.Time `json:"completed_at"`
	TimeSpent   int       `json:"time_spent"` // in seconds
	Difficulty  string    `json:"difficulty"`
}

// ScoreDistribution represents the distribution of scores
type ScoreDistribution struct {
	Range0to20   int `json:"range_0_to_20"`
	Range21to40  int `json:"range_21_to_40"`
	Range41to60  int `json:"range_41_to_60"`
	Range61to80  int `json:"range_61_to_80"`
	Range81to100 int `json:"range_81_to_100"`
}

// MonthlyProgressItem represents progress for a specific month
type MonthlyProgressItem struct {
	Month            string  `json:"month"`
	QuizzesCompleted int     `json:"quizzes_completed"`
	AverageScore     float64 `json:"average_score"`
	TotalPoints      int     `json:"total_points"`
}

// UserStatisticsResponse represents the API response for user statistics
type UserStatisticsResponse struct {
	Statistics UserStatistics `json:"statistics"`
	Message    string         `json:"message,omitempty"`
}
