package models

import (
	"time"

	"github.com/google/uuid"
)

// Quiz represents a quiz
type Quiz struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	Title         string          `json:"title" db:"title"`
	Description   string          `json:"description" db:"description"`
	Category      string          `json:"category" db:"category_id"`
	Difficulty    string          `json:"difficulty" db:"difficulty"`
	TimeLimit     int             `json:"time_limit" db:"time_limit_minutes"` // in minutes
	QuestionCount int             `json:"question_count" db:"total_questions"`
	Points        int             `json:"points" db:"points"`
	IsFeatured    bool            `json:"is_featured" db:"is_featured"`
	IsPublic      bool            `json:"is_public" db:"is_public"`
	CreatedBy     uuid.UUID       `json:"created_by" db:"created_by"`
	Tags          StringArray     `json:"tags" db:"tags"`
	ThumbnailURL  *string         `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
	Questions     []Question      `json:"questions,omitempty"`
	Statistics    *QuizStatistics `json:"statistics,omitempty"`
	AverageRating *float64        `json:"average_rating,omitempty"`
	TotalRatings  *int            `json:"total_ratings,omitempty"`
}

// Question represents a quiz question
type Question struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	QuizID        uuid.UUID   `json:"quiz_id" db:"quiz_id"`
	QuestionText  string      `json:"question_text" db:"question_text"`
	QuestionType  string      `json:"question_type" db:"question_type"` // multiple_choice, true_false, short_answer
	Options       StringArray `json:"options,omitempty" db:"options"`
	CorrectAnswer string      `json:"correct_answer" db:"correct_answer"`
	Explanation   *string     `json:"explanation,omitempty" db:"explanation"`
	Points        int         `json:"points" db:"points"`
	Order         int         `json:"order" db:"order_index"`
	ImageURL      *string     `json:"image_url,omitempty" db:"image_url"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}

// QuizStatistics represents statistics for a quiz
type QuizStatistics struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	QuizID            uuid.UUID  `json:"quiz_id" db:"quiz_id"`
	TotalAttempts     int        `json:"total_attempts" db:"total_attempts"`
	CompletedAttempts int        `json:"completed_attempts" db:"completed_attempts"`
	AverageScore      float64    `json:"average_score" db:"average_score"`
	AverageTime       int        `json:"average_time" db:"average_time"` // in seconds
	HighestScore      float64    `json:"highest_score" db:"highest_score"`
	LowestScore       float64    `json:"lowest_score" db:"lowest_score"`
	PopularityScore   int        `json:"popularity_score" db:"popularity_score"`
	CompletionRate    float64    `json:"completion_rate"` // Computed: (completed_attempts / total_attempts) * 100
	LastAttemptAt     *time.Time `json:"last_attempt_at,omitempty" db:"last_attempt_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// Quiz DTOs

// CreateQuizRequest represents the request to create a new quiz
type CreateQuizRequest struct {
	Title        string                  `json:"title" binding:"required,min=1,max=200"`
	Description  string                  `json:"description" binding:"required,min=1,max=1000"`
	Category     string                  `json:"category" binding:"required"`
	Difficulty   string                  `json:"difficulty" binding:"required,oneof=Easy Medium Hard"`
	TimeLimit    int                     `json:"time_limit" binding:"required,min=30,max=3600"`
	IsFeatured   bool                    `json:"is_featured"`
	IsPublic     bool                    `json:"is_public"`
	Tags         []string                `json:"tags"`
	ThumbnailURL *string                 `json:"thumbnail_url,omitempty"`
	Questions    []CreateQuestionRequest `json:"questions" binding:"required,dive"`
}

// CreateQuestionRequest represents the request to create a new question
type CreateQuestionRequest struct {
	QuestionText  string   `json:"question_text" binding:"required,min=1,max=500"`
	QuestionType  string   `json:"question_type" binding:"required,oneof=multiple_choice true_false short_answer"`
	Options       []string `json:"options,omitempty"`
	CorrectAnswer string   `json:"correct_answer" binding:"required"`
	Explanation   *string  `json:"explanation,omitempty"`
	Points        int      `json:"points" binding:"required,min=1,max=100"`
	Order         int      `json:"order" binding:"required,min=1"`
	ImageURL      *string  `json:"image_url,omitempty"`
}

// UpdateQuizRequest represents the request to update a quiz
type UpdateQuizRequest struct {
	Title        *string  `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Description  *string  `json:"description,omitempty" binding:"omitempty,min=1,max=1000"`
	Category     *string  `json:"category,omitempty"`
	Difficulty   *string  `json:"difficulty,omitempty" binding:"omitempty,oneof=Easy Medium Hard"`
	TimeLimit    *int     `json:"time_limit,omitempty" binding:"omitempty,min=30,max=3600"`
	IsFeatured   *bool    `json:"is_featured,omitempty"`
	IsPublic     *bool    `json:"is_public,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	ThumbnailURL *string  `json:"thumbnail_url,omitempty"`
}

// QuizListResponse represents the response for quiz list
type QuizListResponse struct {
	Quizzes    []QuizSummary `json:"quizzes"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// QuizSummary represents a summary of a quiz
type QuizSummary struct {
	ID            uuid.UUID              `json:"id"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Category      string                 `json:"category"`
	Difficulty    string                 `json:"difficulty"`
	TimeLimit     int                    `json:"time_limit"`
	QuestionCount int                    `json:"question_count"`
	Points        int                    `json:"points"`
	IsFeatured    bool                   `json:"is_featured"`
	Tags          []string               `json:"tags"`
	ThumbnailURL  *string                `json:"thumbnail_url,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	Statistics    *QuizStatisticsSummary `json:"statistics,omitempty"`
}

// QuizStatisticsSummary represents a summary of quiz statistics
type QuizStatisticsSummary struct {
	TotalAttempts int     `json:"total_attempts"`
	AverageScore  float64 `json:"average_score"`
	AverageTime   int     `json:"average_time"`
}

// QuizDetailResponse represents the response for quiz detail
type QuizDetailResponse struct {
	Quiz Quiz `json:"quiz"`
}

// QuizFilters represents filters for quiz queries
type QuizFilters struct {
	Category   string `form:"category"`   // Comma-separated category IDs for multiple category filtering
	Difficulty string `form:"difficulty" binding:"omitempty,oneof=beginner intermediate advanced"`
	Featured   *bool  `form:"featured"`
	Tags       string `form:"tags"`
	Search     string `form:"search"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=10" binding:"min=1,max=100"`
}
