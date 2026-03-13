package models

import (
	"time"

	"github.com/google/uuid"
)

// ========================================
// Attempt Validation Models
// ========================================

// ValidateAttemptRequest is the request body for attempt validation
type ValidateAttemptRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	QuizID uuid.UUID `json:"quiz_id" binding:"required"`
}

// ValidateAttemptResponse is the response from attempt validation
type ValidateAttemptResponse struct {
	Valid     bool         `json:"valid"`
	Attempt   *AttemptInfo `json:"attempt,omitempty"`
	ErrorCode string       `json:"error_code,omitempty"`
	ErrorMsg  string       `json:"error_message,omitempty"`
}

// AttemptInfo contains basic attempt information
type AttemptInfo struct {
	ID          uuid.UUID  `json:"id"`
	QuizID      uuid.UUID  `json:"quiz_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Status      string     `json:"status"`
	IsCompleted bool       `json:"is_completed"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// ========================================
// Quiz Questions Models
// ========================================

// QuestionWithAnswer represents a question with its correct answer (for internal use only)
type QuestionWithAnswer struct {
	ID            uuid.UUID `json:"id"`
	QuizID        uuid.UUID `json:"quiz_id"`
	QuestionText  string    `json:"question_text"`
	QuestionType  string    `json:"question_type"`
	Options       []string  `json:"options"`
	CorrectAnswer string    `json:"correct_answer"`
	Points        int       `json:"points"`
	OrderIndex    int       `json:"order_index"`
}

// GetQuestionsResponse is the response containing quiz questions with answers
type GetQuestionsResponse struct {
	QuizID    uuid.UUID            `json:"quiz_id"`
	Questions []QuestionWithAnswer `json:"questions"`
}

// ========================================
// Scoring Models
// ========================================

// CalculateScoreRequest is the request body for score calculation
type CalculateScoreRequest struct {
	Questions []QuestionForScoring `json:"questions" binding:"required"`
	Answers   []SubmittedAnswer    `json:"answers" binding:"required"`
}

// QuestionForScoring contains question data needed for scoring
type QuestionForScoring struct {
	ID            uuid.UUID `json:"id"`
	QuestionType  string    `json:"question_type"`
	Options       []string  `json:"options"`
	CorrectAnswer string    `json:"correct_answer"`
}

// SubmittedAnswer represents a user's submitted answer
type SubmittedAnswer struct {
	QuestionID          uuid.UUID `json:"question_id"`
	SelectedOptionIndex *int      `json:"selected_option_index,omitempty"`
	SelectedOption      *string   `json:"selected_option,omitempty"`
	TextAnswer          string    `json:"text_answer,omitempty"`
}

// CalculateScoreResponse is the response from score calculation
type CalculateScoreResponse struct {
	TotalQuestions   int               `json:"total_questions"`
	CorrectAnswers   int               `json:"correct_answers"`
	Score            float64           `json:"score"`
	PercentageScore  float64           `json:"percentage_score"`
	Passed           bool              `json:"passed"`
	ValidatedAnswers []ValidatedAnswer `json:"validated_answers"`
}

// ValidatedAnswer represents a validated and scored answer
type ValidatedAnswer struct {
	QuestionID     uuid.UUID `json:"question_id"`
	SelectedOption int       `json:"selected_option"`
	TextAnswer     string    `json:"text_answer"`
	IsCorrect      bool      `json:"is_correct"`
	PointsEarned   int       `json:"points_earned"`
}

// ========================================
// Update Attempt Models
// ========================================

// UpdateAttemptRequest is the request body for updating an attempt
type UpdateAttemptRequest struct {
	UserID          uuid.UUID         `json:"user_id" binding:"required"`
	Answers         []ValidatedAnswer `json:"answers"`
	Score           float64           `json:"score"`
	TotalPoints     int               `json:"total_points"`
	TimeSpent       int               `json:"time_spent"`
	PercentageScore float64           `json:"percentage_score"`
	Passed          bool              `json:"passed"`
	Status          string            `json:"status"`
}

// UpdateAttemptResponse is the response from updating an attempt
type UpdateAttemptResponse struct {
	Success     bool       `json:"success"`
	AttemptID   uuid.UUID  `json:"attempt_id"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// ========================================
// User Statistics Models
// ========================================

// UpdateStatisticsRequest is the request body for updating user statistics
type UpdateStatisticsRequest struct {
	QuizID          uuid.UUID `json:"quiz_id,omitempty"`
	Score           float64   `json:"score"`
	PercentageScore float64   `json:"percentage_score"`
	TimeSpent       int       `json:"time_spent,omitempty"`
}

// UpdateStatisticsResponse is the response from updating statistics
type UpdateStatisticsResponse struct {
	Success                 bool    `json:"success"`
	TotalQuizzesCompleted   int     `json:"total_quizzes_completed"`
	AverageScore            float64 `json:"average_score"`
	CurrentStreak           int     `json:"current_streak"`
	BestStreak              int     `json:"best_streak"`
}

// ========================================
// Achievement Models
// ========================================

// CheckAchievementsRequest is the request body for checking achievements
type CheckAchievementsRequest struct {
	Trigger string                 `json:"trigger" binding:"required"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// CheckAchievementsResponse is the response from checking achievements
type CheckAchievementsResponse struct {
	NewAchievements []AchievementNotification `json:"new_achievements"`
	TotalChecked    int                       `json:"total_checked"`
	TotalUnlocked   int                       `json:"total_unlocked"`
}

// AchievementNotification represents an achievement notification
type AchievementNotification struct {
	AchievementID uuid.UUID `json:"achievement_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Icon          string    `json:"icon"`
	Color         string    `json:"color"`
	PointsAwarded int       `json:"points_awarded"`
	IsRare        bool      `json:"is_rare"`
}