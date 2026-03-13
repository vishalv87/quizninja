package models

import (
	"time"

	"github.com/google/uuid"
)

// QuizRating represents a user's rating and review for a quiz
type QuizRating struct {
	ID        uuid.UUID `json:"id" db:"id"`
	QuizID    uuid.UUID `json:"quiz_id" db:"quiz_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Rating    int       `json:"rating" db:"rating"` // 1-5 stars
	Review    *string   `json:"review,omitempty" db:"review"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Rating DTOs

// CreateRatingRequest represents the request to create a new rating
type CreateRatingRequest struct {
	Rating int     `json:"rating" binding:"required,min=1,max=5"`
	Review *string `json:"review,omitempty" binding:"omitempty,max=1000"`
}

// UpdateRatingRequest represents the request to update a rating
type UpdateRatingRequest struct {
	Rating *int    `json:"rating,omitempty" binding:"omitempty,min=1,max=5"`
	Review *string `json:"review,omitempty" binding:"omitempty,max=1000"`
}

// RatingResponse represents a rating with additional user info
type RatingResponse struct {
	ID        uuid.UUID `json:"id"`
	QuizID    uuid.UUID `json:"quiz_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	Rating    int       `json:"rating"`
	Review    *string   `json:"review,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// RatingListResponse represents the response for rating list
type RatingListResponse struct {
	Ratings       []RatingResponse `json:"ratings"`
	AverageRating float64          `json:"average_rating"`
	TotalRatings  int              `json:"total_ratings"`
	Total         int              `json:"total"`
	Page          int              `json:"page"`
	PageSize      int              `json:"page_size"`
	TotalPages    int              `json:"total_pages"`
}

// AverageRatingResponse represents the average rating for a quiz
type AverageRatingResponse struct {
	QuizID        uuid.UUID `json:"quiz_id"`
	AverageRating float64   `json:"average_rating"`
	TotalRatings  int       `json:"total_ratings"`
}
