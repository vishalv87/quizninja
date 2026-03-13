package models

import (
	"time"

	"github.com/google/uuid"
)

// UserQuizFavorite represents a user's favorite quiz
type UserQuizFavorite struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	UserID      uuid.UUID    `json:"user_id" db:"user_id"`
	QuizID      uuid.UUID    `json:"quiz_id" db:"quiz_id"`
	FavoritedAt time.Time    `json:"favorited_at" db:"favorited_at"`
	Quiz        *QuizSummary `json:"quiz,omitempty"`
}

// FavoritesListResponse represents the response for user favorites
type FavoritesListResponse struct {
	Favorites  []UserQuizFavorite `json:"favorites"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// AddFavoriteRequest represents the request to add a quiz to favorites
type AddFavoriteRequest struct {
	QuizID uuid.UUID `json:"quiz_id" binding:"required"`
}