package models

import (
	"time"

	"github.com/google/uuid"
)

// Challenge represents a quiz challenge between users
type Challenge struct {
	ID                uuid.UUID              `json:"id" db:"id"`
	ChallengerID      uuid.UUID              `json:"challenger_id" db:"challenger_id"`
	ChallengeeID      uuid.UUID              `json:"challengee_id" db:"challengee_id"`
	QuizID            uuid.UUID              `json:"quiz_id" db:"quiz_id"`
	Status            string                 `json:"status" db:"status"`
	ChallengerScore   *float64               `json:"challenger_score,omitempty" db:"challenger_score"`
	ChallengeeScore   *float64               `json:"challengee_score,omitempty" db:"challengee_score"`
	Message           *string                `json:"message,omitempty" db:"message"`
	ExpiresAt         *time.Time             `json:"expires_at,omitempty" db:"expires_at"`
	IsGroupChallenge  bool                   `json:"is_group_challenge" db:"is_group_challenge"`
	ParticipantIDs    []uuid.UUID            `json:"participant_ids,omitempty" db:"participant_ids"`
	ParticipantScores map[string]interface{} `json:"participant_scores,omitempty" db:"participant_scores"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`

	// Asynchronous challenge tracking fields
	ChallengerAttemptID   *uuid.UUID `json:"challenger_attempt_id,omitempty" db:"challenger_attempt_id"`
	ChallengeeAttemptID   *uuid.UUID `json:"challengee_attempt_id,omitempty" db:"challengee_attempt_id"`
	ChallengerCompletedAt *time.Time `json:"challenger_completed_at,omitempty" db:"challenger_completed_at"`
	ChallengeeCompletedAt *time.Time `json:"challengee_completed_at,omitempty" db:"challengee_completed_at"`

	// Relationships
	Challenger *User        `json:"challenger,omitempty"`
	Challengee *User        `json:"challengee,omitempty"`
	Quiz       *QuizSummary `json:"quiz,omitempty"`
}

// ChallengeWithDetails represents a challenge with full user and quiz details
type ChallengeWithDetails struct {
	Challenge
	ChallengerName   string `json:"challenger_name"`
	ChallengerAvatar string `json:"challenger_avatar"`
	ChallengeeName   string `json:"challengee_name"`
	ChallengeeAvatar string `json:"challengee_avatar"`
	QuizTitle        string `json:"quiz_title"`
	QuizCategory     string `json:"quiz_category"`
}

// Challenge DTOs

// CreateChallengeRequest represents the request to create a new challenge
type CreateChallengeRequest struct {
	ChallengeeUserID uuid.UUID   `json:"challengee_user_id" binding:"required"`
	QuizID           uuid.UUID   `json:"quiz_id" binding:"required"`
	Message          *string     `json:"message,omitempty"`
	ExpiresAt        *time.Time  `json:"expires_at,omitempty"`
	IsGroupChallenge bool        `json:"is_group_challenge"`
	ParticipantIDs   []uuid.UUID `json:"participant_ids,omitempty"`
}

// AcceptChallengeRequest represents the request to accept a challenge
type AcceptChallengeRequest struct {
	// No additional fields needed - just the challenge ID from URL
}

// UpdateChallengeScoreRequest represents the request to update challenge scores
type UpdateChallengeScoreRequest struct {
	UserScore float64 `json:"user_score" binding:"required,min=0,max=100"`
}

// ChallengeFilters represents filters for challenge queries
type ChallengeFilters struct {
	Status    string     `form:"status" binding:"omitempty,oneof=pending accepted completed expired declined"`
	QuizID    *uuid.UUID `form:"quiz_id"`
	UserType  string     `form:"user_type" binding:"omitempty,oneof=challenger challengee all"`
	StartDate *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"end_date" time_format:"2006-01-02"`
	Page      int        `form:"page,default=1" binding:"min=1"`
	PageSize  int        `form:"page_size,default=10" binding:"min=1,max=100"`
	SortBy    string     `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at expires_at status"`
	SortOrder string     `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// ChallengeListResponse represents the response for challenge list
type ChallengeListResponse struct {
	Challenges []ChallengeWithDetails `json:"challenges"`
	Total      int                    `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

// ChallengeDetailResponse represents the response for a single challenge detail
type ChallengeDetailResponse struct {
	Challenge ChallengeWithDetails `json:"challenge"`
}

// ChallengeStatsResponse represents challenge statistics for a user
type ChallengeStatsResponse struct {
	TotalChallenges     int     `json:"total_challenges"`
	PendingChallenges   int     `json:"pending_challenges"`
	ActiveChallenges    int     `json:"active_challenges"`
	CompletedChallenges int     `json:"completed_challenges"`
	WonChallenges       int     `json:"won_challenges"`
	LostChallenges      int     `json:"lost_challenges"`
	WinRate             float64 `json:"win_rate"`
	AverageScore        float64 `json:"average_score"`
}