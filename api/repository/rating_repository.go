package repository

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/google/uuid"
	"quizninja-api/database"
	"quizninja-api/models"
)

type RatingRepository struct {
	db *sql.DB
}

func NewRatingRepository() *RatingRepository {
	return &RatingRepository{
		db: database.DB,
	}
}

// CreateRating creates a new quiz rating
func (r *RatingRepository) CreateRating(rating *models.QuizRating) error {
	query := `
		INSERT INTO quiz_ratings (id, quiz_id, user_id, rating, review, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING created_at`

	err := r.db.QueryRow(
		query,
		rating.ID,
		rating.QuizID,
		rating.UserID,
		rating.Rating,
		rating.Review,
	).Scan(&rating.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create rating: %w", err)
	}

	return nil
}

// GetRatingByID retrieves a rating by ID
func (r *RatingRepository) GetRatingByID(id uuid.UUID) (*models.QuizRating, error) {
	query := `
		SELECT id, quiz_id, user_id, rating, review, created_at
		FROM quiz_ratings
		WHERE id = $1`

	var rating models.QuizRating
	err := r.db.QueryRow(query, id).Scan(
		&rating.ID,
		&rating.QuizID,
		&rating.UserID,
		&rating.Rating,
		&rating.Review,
		&rating.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("rating not found")
		}
		return nil, fmt.Errorf("failed to get rating: %w", err)
	}

	return &rating, nil
}

// GetUserRating retrieves a specific user's rating for a quiz
func (r *RatingRepository) GetUserRating(userID, quizID uuid.UUID) (*models.QuizRating, error) {
	query := `
		SELECT id, quiz_id, user_id, rating, review, created_at
		FROM quiz_ratings
		WHERE user_id = $1 AND quiz_id = $2`

	var rating models.QuizRating
	err := r.db.QueryRow(query, userID, quizID).Scan(
		&rating.ID,
		&rating.QuizID,
		&rating.UserID,
		&rating.Rating,
		&rating.Review,
		&rating.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No rating found is not an error
		}
		return nil, fmt.Errorf("failed to get user rating: %w", err)
	}

	return &rating, nil
}

// GetQuizRatings retrieves all ratings for a quiz with pagination
func (r *RatingRepository) GetQuizRatings(quizID uuid.UUID, page, pageSize int) ([]models.RatingResponse, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM quiz_ratings WHERE quiz_id = $1`
	err := r.db.QueryRow(countQuery, quizID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count ratings: %w", err)
	}

	// Get ratings with user info
	offset := (page - 1) * pageSize
	query := `
		SELECT r.id, r.quiz_id, r.user_id, u.username, r.rating, r.review, r.created_at
		FROM quiz_ratings r
		JOIN users u ON r.user_id = u.id
		WHERE r.quiz_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, quizID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get ratings: %w", err)
	}
	defer rows.Close()

	ratings := []models.RatingResponse{}
	for rows.Next() {
		var rating models.RatingResponse
		err := rows.Scan(
			&rating.ID,
			&rating.QuizID,
			&rating.UserID,
			&rating.UserName,
			&rating.Rating,
			&rating.Review,
			&rating.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan rating: %w", err)
		}
		ratings = append(ratings, rating)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating ratings: %w", err)
	}

	return ratings, total, nil
}

// GetAverageRating calculates the average rating for a quiz
func (r *RatingRepository) GetAverageRating(quizID uuid.UUID) (float64, int, error) {
	query := `
		SELECT COALESCE(AVG(rating), 0), COUNT(*)
		FROM quiz_ratings
		WHERE quiz_id = $1`

	var avgRating float64
	var count int
	err := r.db.QueryRow(query, quizID).Scan(&avgRating, &count)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	// Round to 1 decimal place
	avgRating = math.Round(avgRating*10) / 10

	return avgRating, count, nil
}

// UpdateRating updates an existing rating
func (r *RatingRepository) UpdateRating(rating *models.QuizRating) error {
	query := `
		UPDATE quiz_ratings
		SET rating = $1, review = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		rating.Rating,
		rating.Review,
		rating.ID,
	).Scan(&rating.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("rating not found")
		}
		return fmt.Errorf("failed to update rating: %w", err)
	}

	return nil
}

// DeleteRating deletes a rating
func (r *RatingRepository) DeleteRating(id uuid.UUID) error {
	query := `DELETE FROM quiz_ratings WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete rating: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("rating not found")
	}

	return nil
}
