package models

import (
	"time"

	"github.com/google/uuid"
)

// QuizAttempt represents a user's attempt at a quiz
type QuizAttempt struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	QuizID          uuid.UUID       `json:"quiz_id" db:"quiz_id"`
	UserID          uuid.UUID       `json:"user_id" db:"user_id"`
	Answers         []AttemptAnswer `json:"answers" db:"answers"`
	Score           float64         `json:"score" db:"score"`
	TotalPoints     int             `json:"total_points" db:"total_points"`
	TimeSpent       int             `json:"time_spent" db:"time_spent"` // in seconds
	PercentageScore float64         `json:"percentage_score" db:"percentage_score"`
	Passed          bool            `json:"passed" db:"passed"`
	Status          string          `json:"status" db:"status"` // started, completed, abandoned
	IsCompleted     bool            `json:"is_completed" db:"is_completed"`
	StartedAt       time.Time       `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`

	// Challenge tracking fields
	ChallengeID        *uuid.UUID `json:"challenge_id,omitempty" db:"challenge_id"`
	IsChallengeAttempt bool       `json:"is_challenge_attempt" db:"is_challenge_attempt"`
}

// AttemptAnswer represents a user's answer to a quiz question
type AttemptAnswer struct {
	QuestionID     uuid.UUID `json:"question_id" validate:"required"`
	SelectedOption int       `json:"selected_option"` // For multiple choice
	TextAnswer     string    `json:"text_answer"`     // For text answers
	IsCorrect      bool      `json:"is_correct"`
	PointsEarned   int       `json:"points_earned"`
}

// UpdateAttemptRequest represents the request body for updating/completing a quiz attempt
type UpdateAttemptRequest struct {
	Answers   []AttemptAnswer `json:"answers" validate:"required"`
	TimeSpent int             `json:"time_spent" validate:"min=1"`
	Status    string          `json:"status" validate:"required,oneof=completed abandoned"`
}

// AttemptFilters represents filters for user attempt history
type AttemptFilters struct {
	QuizID     *string    `form:"quiz_id"`
	Category   string     `form:"category"`
	Difficulty string     `form:"difficulty" binding:"omitempty,oneof=Easy Medium Hard"`
	StartDate  *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate    *time.Time `form:"end_date" time_format:"2006-01-02"`
	MinScore   *float64   `form:"min_score" binding:"omitempty,min=0,max=100"`
	MaxScore   *float64   `form:"max_score" binding:"omitempty,min=0,max=100"`
	Page       int        `form:"page,default=1" binding:"min=1"`
	PageSize   int        `form:"page_size,default=10" binding:"min=1,max=100"`
	SortBy     string     `form:"sort_by,default=completed_at" binding:"omitempty,oneof=completed_at score time_spent"`
	SortOrder  string     `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// QuizAttemptWithDetails represents a quiz attempt with quiz details for API responses
type QuizAttemptWithDetails struct {
	QuizAttempt
	Quiz QuizSummary `json:"quiz"`
}

// AttemptHistoryResponse represents the API response for user attempt history
type AttemptHistoryResponse struct {
	Attempts   []QuizAttemptWithDetails `json:"attempts"`
	Total      int                      `json:"total"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPages int                      `json:"total_pages"`
}

// AttemptDetailResponse represents the API response for a single attempt detail
type AttemptDetailResponse struct {
	Attempt QuizAttemptWithDetails `json:"attempt"`
}
