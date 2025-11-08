package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// QuizSession represents an active or paused quiz session for continuation functionality
type QuizSession struct {
	ID                   uuid.UUID       `json:"id" db:"id"`
	AttemptID            uuid.UUID       `json:"attempt_id" db:"attempt_id"`
	UserID               uuid.UUID       `json:"user_id" db:"user_id"`
	QuizID               uuid.UUID       `json:"quiz_id" db:"quiz_id"`
	CurrentQuestionIndex int             `json:"current_question_index" db:"current_question_index"`
	CurrentAnswers       []AttemptAnswer `json:"current_answers" db:"current_answers"`
	SessionState         string          `json:"session_state" db:"session_state"`             // active, paused, completed, abandoned
	TimeRemaining        *int            `json:"time_remaining,omitempty" db:"time_remaining"` // seconds
	TimeSpentSoFar       int             `json:"time_spent_so_far" db:"time_spent_so_far"`     // seconds
	LastActivityAt       time.Time       `json:"last_activity_at" db:"last_activity_at"`
	PausedAt             *time.Time      `json:"paused_at,omitempty" db:"paused_at"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
	Quiz                 *QuizSummary    `json:"quiz,omitempty"`
	Attempt              *QuizAttempt    `json:"attempt,omitempty"`
}

// IsActive returns true if the session is currently active
func (qs *QuizSession) IsActive() bool {
	return qs.SessionState == "active"
}

// IsPaused returns true if the session is currently paused
func (qs *QuizSession) IsPaused() bool {
	return qs.SessionState == "paused"
}

// CanResume returns true if the session can be resumed
func (qs *QuizSession) CanResume() bool {
	return qs.SessionState == "paused" || qs.SessionState == "active"
}

// GetProgress returns the completion progress as a percentage (0-100)
func (qs *QuizSession) GetProgress() float64 {
	if qs.Quiz == nil || qs.Quiz.QuestionCount == 0 {
		return 0.0
	}
	return (float64(qs.CurrentQuestionIndex) / float64(qs.Quiz.QuestionCount)) * 100
}

// GetQuestionsRemaining returns the number of questions left to complete
func (qs *QuizSession) GetQuestionsRemaining() int {
	if qs.Quiz == nil {
		return 0
	}
	return qs.Quiz.QuestionCount - qs.CurrentQuestionIndex
}

// GetTimeRemainingFormatted returns the remaining time in a human-readable format
func (qs *QuizSession) GetTimeRemainingFormatted() string {
	if qs.TimeRemaining == nil {
		return "No time limit"
	}

	minutes := *qs.TimeRemaining / 60
	seconds := *qs.TimeRemaining % 60

	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// HasExpired returns true if the session has been inactive for too long
func (qs *QuizSession) HasExpired() bool {
	// Sessions no longer expire automatically - users can resume anytime
	return false
}

// QuizSessionWithDetails represents a quiz session with full quiz and attempt details
type QuizSessionWithDetails struct {
	QuizSession
	QuizTitle         string  `json:"quiz_title"`
	QuizCategory      string  `json:"quiz_category"`
	QuizDifficulty    string  `json:"quiz_difficulty"`
	QuizThumbnail     *string `json:"quiz_thumbnail,omitempty"`
	TotalQuestions    int     `json:"total_questions"`
	OriginalTimeLimit int     `json:"original_time_limit"` // in seconds
}

// Quiz Session DTOs

// CreateQuizSessionRequest represents the request to create a new quiz session
type CreateQuizSessionRequest struct {
	AttemptID            uuid.UUID `json:"attempt_id" binding:"required"`
	CurrentQuestionIndex int       `json:"current_question_index" binding:"min=0"`
	TimeRemaining        *int      `json:"time_remaining,omitempty" binding:"omitempty,min=0"`
}

// UpdateQuizSessionRequest represents the request to update a quiz session
type UpdateQuizSessionRequest struct {
	CurrentQuestionIndex int             `json:"current_question_index" binding:"min=0"`
	CurrentAnswers       []AttemptAnswer `json:"current_answers"`
	TimeSpentSoFar       int             `json:"time_spent_so_far" binding:"min=0"`
	TimeRemaining        *int            `json:"time_remaining,omitempty" binding:"omitempty,min=0"`
}

// PauseSessionRequest represents the request to pause a quiz session
type PauseSessionRequest struct {
	CurrentQuestionIndex int             `json:"current_question_index" binding:"min=0"`
	CurrentAnswers       []AttemptAnswer `json:"current_answers"`
	TimeSpentSoFar       int             `json:"time_spent_so_far" binding:"min=0"`
	TimeRemaining        *int            `json:"time_remaining,omitempty" binding:"omitempty,min=0"`
}

// ResumeSessionRequest represents the request to resume a quiz session
type ResumeSessionRequest struct {
	// No additional fields needed - session state is loaded from database
}

// ActiveSessionsResponse represents the response for user's active/paused sessions
type ActiveSessionsResponse struct {
	Sessions    []QuizSessionWithDetails `json:"sessions"`
	Total       int                      `json:"total"`
	ActiveCount int                      `json:"active_count"`
	PausedCount int                      `json:"paused_count"`
}

// SessionFilters represents filters for quiz session queries
type SessionFilters struct {
	SessionState string     `form:"session_state" binding:"omitempty,oneof=active paused completed abandoned"`
	QuizID       *uuid.UUID `form:"quiz_id"`
	Category     string     `form:"category"`
	Difficulty   string     `form:"difficulty" binding:"omitempty,oneof=Easy Medium Hard"`
	StartDate    *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate      *time.Time `form:"end_date" time_format:"2006-01-02"`
	Page         int        `form:"page,default=1" binding:"min=1"`
	PageSize     int        `form:"page_size,default=10" binding:"min=1,max=100"`
	SortBy       string     `form:"sort_by,default=last_activity_at" binding:"omitempty,oneof=last_activity_at created_at"`
	SortOrder    string     `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// SessionActionResponse represents the response for session actions (pause/resume/abandon)
type SessionActionResponse struct {
	SessionID     uuid.UUID `json:"session_id"`
	Action        string    `json:"action"` // paused, resumed, abandoned
	SessionState  string    `json:"session_state"`
	Message       string    `json:"message"`
	TimeRemaining *int      `json:"time_remaining,omitempty"`
	Progress      float64   `json:"progress"`
}
