package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"quizninja-api/database"
	"quizninja-api/models"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type QuizRepository struct {
	db *sql.DB
}

func NewQuizRepository() *QuizRepository {
	return &QuizRepository{
		db: database.DB,
	}
}

// GetQuizByID retrieves a quiz by its ID (basic info only)
func (r *QuizRepository) GetQuizByID(id uuid.UUID) (*models.Quiz, error) {
	query := `
		SELECT id, title, description, category_id, difficulty, time_limit_minutes, total_questions, points,
		       is_featured, is_public, created_by, tags, thumbnail_url, created_at, updated_at
		FROM quizzes
		WHERE id = $1`

	var quiz models.Quiz
	var tags pq.StringArray

	err := r.db.QueryRow(query, id).Scan(
		&quiz.ID, &quiz.Title, &quiz.Description, &quiz.Category, &quiz.Difficulty,
		&quiz.TimeLimit, &quiz.QuestionCount, &quiz.Points, &quiz.IsFeatured, &quiz.IsPublic,
		&quiz.CreatedBy, &tags, &quiz.ThumbnailURL, &quiz.CreatedAt, &quiz.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("quiz not found")
		}
		return nil, fmt.Errorf("failed to get quiz: %w", err)
	}

	quiz.Tags = models.StringArray(tags)
	return &quiz, nil
}

// GetQuizByIDWithQuestions retrieves a quiz with its questions
func (r *QuizRepository) GetQuizByIDWithQuestions(id uuid.UUID) (*models.Quiz, error) {
	quiz, err := r.GetQuizByID(id)
	if err != nil {
		return nil, err
	}

	questions, err := r.GetQuestionsByQuizID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	quiz.Questions = questions
	return quiz, nil
}

// GetQuizByIDWithStatistics retrieves a quiz with its statistics
func (r *QuizRepository) GetQuizByIDWithStatistics(id uuid.UUID) (*models.Quiz, error) {
	quiz, err := r.GetQuizByID(id)
	if err != nil {
		return nil, err
	}

	stats, err := r.GetQuizStatistics(id)
	if err != nil && err.Error() != "quiz statistics not found" {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	quiz.Statistics = stats
	return quiz, nil
}

// GetQuizByIDWithAll retrieves a quiz with questions and statistics
func (r *QuizRepository) GetQuizByIDWithAll(id uuid.UUID) (*models.Quiz, error) {
	quiz, err := r.GetQuizByIDWithQuestions(id)
	if err != nil {
		return nil, err
	}

	stats, err := r.GetQuizStatistics(id)
	if err != nil && err.Error() != "quiz statistics not found" {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	quiz.Statistics = stats
	return quiz, nil
}

// GetQuizzes retrieves quizzes with filtering and pagination
func (r *QuizRepository) GetQuizzes(filters *models.QuizFilters) ([]models.Quiz, int, error) {
	// Build WHERE clause
	whereClause := "WHERE is_public = true"
	args := []interface{}{}
	argIndex := 1

	if filters.Category != "" {
		// Support multiple categories separated by comma
		categories := strings.Split(filters.Category, ",")
		for i := range categories {
			categories[i] = strings.TrimSpace(categories[i])
		}

		if len(categories) == 1 {
			// Single category
			whereClause += fmt.Sprintf(" AND category_id = $%d", argIndex)
			args = append(args, categories[0])
			argIndex++
		} else {
			// Multiple categories - use IN clause
			whereClause += fmt.Sprintf(" AND category_id = ANY($%d)", argIndex)
			args = append(args, pq.Array(categories))
			argIndex++
		}
	}

	if filters.Difficulty != "" {
		whereClause += fmt.Sprintf(" AND difficulty = $%d", argIndex)
		args = append(args, filters.Difficulty)
		argIndex++
	}

	if filters.Featured != nil {
		whereClause += fmt.Sprintf(" AND is_featured = $%d", argIndex)
		args = append(args, *filters.Featured)
		argIndex++
	}

	if filters.Tags != "" {
		tags := strings.Split(filters.Tags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		whereClause += fmt.Sprintf(" AND tags && $%d", argIndex)
		args = append(args, pq.Array(tags))
		argIndex++
	}

	if filters.Search != "" {
		searchTerm := "%" + strings.ToLower(filters.Search) + "%"
		whereClause += fmt.Sprintf(" AND (LOWER(title) LIKE $%d OR LOWER(description) LIKE $%d)", argIndex, argIndex+1)
		args = append(args, searchTerm, searchTerm)
		argIndex += 2
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM quizzes %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count quizzes: %w", err)
	}

	// Get paginated results
	offset := (filters.Page - 1) * filters.PageSize
	query := fmt.Sprintf(`
		SELECT q.id, q.title, q.description, q.category_id, q.difficulty, q.time_limit_minutes,
		       q.total_questions, q.points, q.is_featured, q.is_public, q.created_by, q.tags,
		       q.thumbnail_url, q.created_at, q.updated_at,
		       qs.total_attempts, qs.average_score, qs.average_time_seconds
		FROM quizzes q
		LEFT JOIN quiz_statistics qs ON q.id = qs.quiz_id
		%s
		ORDER BY q.is_featured DESC, q.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, filters.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query quizzes: %w", err)
	}
	defer rows.Close()

	var quizzes []models.Quiz
	for rows.Next() {
		var quiz models.Quiz
		var tags pq.StringArray
		var totalAttempts, averageTime sql.NullInt32
		var averageScore sql.NullFloat64

		err := rows.Scan(
			&quiz.ID, &quiz.Title, &quiz.Description, &quiz.Category, &quiz.Difficulty,
			&quiz.TimeLimit, &quiz.QuestionCount, &quiz.Points, &quiz.IsFeatured, &quiz.IsPublic,
			&quiz.CreatedBy, &tags, &quiz.ThumbnailURL, &quiz.CreatedAt, &quiz.UpdatedAt,
			&totalAttempts, &averageScore, &averageTime,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan quiz: %w", err)
		}

		quiz.Tags = models.StringArray(tags)

		// Add statistics if available
		if totalAttempts.Valid {
			quiz.Statistics = &models.QuizStatistics{
				TotalAttempts: int(totalAttempts.Int32),
				AverageScore:  averageScore.Float64,
				AverageTime:   int(averageTime.Int32),
			}
		}

		quizzes = append(quizzes, quiz)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	return quizzes, total, nil
}

// GetFeaturedQuizzes retrieves featured quizzes
func (r *QuizRepository) GetFeaturedQuizzes(limit int) ([]models.Quiz, error) {
	filters := &models.QuizFilters{
		Featured: &[]bool{true}[0],
		Page:     1,
		PageSize: limit,
	}

	quizzes, _, err := r.GetQuizzes(filters)
	return quizzes, err
}

// GetQuizzesByCategory retrieves quizzes by category
func (r *QuizRepository) GetQuizzesByCategory(category string, limit int) ([]models.Quiz, error) {
	filters := &models.QuizFilters{
		Category: category,
		Page:     1,
		PageSize: limit,
	}

	quizzes, _, err := r.GetQuizzes(filters)
	return quizzes, err
}

// GetQuizzesByUser retrieves quizzes created by a specific user
func (r *QuizRepository) GetQuizzesByUser(userID uuid.UUID, offset, limit int) ([]models.Quiz, int, error) {
	// Count total records
	countQuery := "SELECT COUNT(*) FROM quizzes WHERE created_by = $1"
	var total int
	err := r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user quizzes: %w", err)
	}

	// Get paginated results
	query := `
		SELECT q.id, q.title, q.description, q.category_id, q.difficulty, q.time_limit_minutes,
		       q.total_questions, q.points, q.is_featured, q.is_public, q.created_by, q.tags,
		       q.thumbnail_url, q.created_at, q.updated_at,
		       qs.total_attempts, qs.average_score, qs.average_time_seconds
		FROM quizzes q
		LEFT JOIN quiz_statistics qs ON q.id = qs.quiz_id
		WHERE q.created_by = $1
		ORDER BY q.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query user quizzes: %w", err)
	}
	defer rows.Close()

	var quizzes []models.Quiz
	for rows.Next() {
		var quiz models.Quiz
		var tags pq.StringArray
		var totalAttempts, averageTime sql.NullInt32
		var averageScore sql.NullFloat64

		err := rows.Scan(
			&quiz.ID, &quiz.Title, &quiz.Description, &quiz.Category, &quiz.Difficulty,
			&quiz.TimeLimit, &quiz.QuestionCount, &quiz.Points, &quiz.IsFeatured, &quiz.IsPublic,
			&quiz.CreatedBy, &tags, &quiz.ThumbnailURL, &quiz.CreatedAt, &quiz.UpdatedAt,
			&totalAttempts, &averageScore, &averageTime,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan quiz: %w", err)
		}

		quiz.Tags = models.StringArray(tags)

		// Add statistics if available
		if totalAttempts.Valid {
			quiz.Statistics = &models.QuizStatistics{
				TotalAttempts: int(totalAttempts.Int32),
				AverageScore:  averageScore.Float64,
				AverageTime:   int(averageTime.Int32),
			}
		}

		quizzes = append(quizzes, quiz)
	}

	return quizzes, total, nil
}

// GetCompletedQuizzesByUser retrieves quizzes that a user has completed
func (r *QuizRepository) GetCompletedQuizzesByUser(userID uuid.UUID, offset, limit int) ([]models.Quiz, int, error) {
	// Count total unique quizzes completed by user
	countQuery := `
		SELECT COUNT(DISTINCT qa.quiz_id)
		FROM quiz_attempts qa
		WHERE qa.user_id = $1 AND qa.is_completed = true`
	var total int
	err := r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count completed quizzes: %w", err)
	}

	// Get paginated results with quiz details
	// Using a subquery to get the most recent completion date for each quiz first,
	// then paginate, and finally join to get full quiz details
	query := `
		WITH user_completed_quizzes AS (
			SELECT DISTINCT qa.quiz_id, MAX(qa.completed_at) as last_completed
			FROM quiz_attempts qa
			WHERE qa.user_id = $1 AND qa.is_completed = true
			GROUP BY qa.quiz_id
			ORDER BY last_completed DESC
			LIMIT $2 OFFSET $3
		)
		SELECT q.id, q.title, q.description, q.category_id, q.difficulty, q.time_limit_minutes,
		       q.total_questions, q.points, q.is_featured, q.is_public, q.created_by, q.tags,
		       q.thumbnail_url, q.created_at, q.updated_at,
		       qs.total_attempts, qs.average_score, qs.average_time_seconds,
		       ucq.last_completed
		FROM user_completed_quizzes ucq
		INNER JOIN quizzes q ON ucq.quiz_id = q.id
		LEFT JOIN quiz_statistics qs ON q.id = qs.quiz_id
		ORDER BY ucq.last_completed DESC`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query completed quizzes: %w", err)
	}
	defer rows.Close()

	var quizzes []models.Quiz
	for rows.Next() {
		var quiz models.Quiz
		var tags pq.StringArray
		var totalAttempts, averageTime sql.NullInt32
		var averageScore sql.NullFloat64
		var completedAt sql.NullTime

		err := rows.Scan(
			&quiz.ID, &quiz.Title, &quiz.Description, &quiz.Category, &quiz.Difficulty,
			&quiz.TimeLimit, &quiz.QuestionCount, &quiz.Points, &quiz.IsFeatured, &quiz.IsPublic,
			&quiz.CreatedBy, &tags, &quiz.ThumbnailURL, &quiz.CreatedAt, &quiz.UpdatedAt,
			&totalAttempts, &averageScore, &averageTime, &completedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan completed quiz: %w", err)
		}

		quiz.Tags = models.StringArray(tags)

		// Add statistics if available
		if totalAttempts.Valid {
			quiz.Statistics = &models.QuizStatistics{
				TotalAttempts: int(totalAttempts.Int32),
				AverageScore:  averageScore.Float64,
				AverageTime:   int(averageTime.Int32),
			}
		}

		quizzes = append(quizzes, quiz)
	}

	return quizzes, total, nil
}

// GetQuestionsByQuizID retrieves all questions for a quiz
func (r *QuizRepository) GetQuestionsByQuizID(quizID uuid.UUID) ([]models.Question, error) {
	query := `
		SELECT id, quiz_id, question_text, question_type, options, correct_answer,
		       explanation, order_index, created_at
		FROM questions
		WHERE quiz_id = $1
		ORDER BY order_index ASC`

	rows, err := r.db.Query(query, quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to query questions: %w", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var question models.Question
		var options pq.StringArray

		err := rows.Scan(
			&question.ID, &question.QuizID, &question.QuestionText, &question.QuestionType,
			&options, &question.CorrectAnswer, &question.Explanation, &question.Order, &question.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}

		question.Options = models.StringArray(options)

		// Set default values for fields not in database
		question.Points = 1                     // Default points per question
		question.ImageURL = nil                 // No image URL in current schema
		question.UpdatedAt = question.CreatedAt // Use created_at as default for updated_at

		questions = append(questions, question)
	}

	return questions, nil
}

// GetQuizStatistics retrieves quiz statistics
func (r *QuizRepository) GetQuizStatistics(quizID uuid.UUID) (*models.QuizStatistics, error) {
	query := `
		SELECT quiz_id, total_attempts, total_completions, average_score,
		       average_time_seconds, difficulty_rating, popularity_score, updated_at
		FROM quiz_statistics
		WHERE quiz_id = $1`

	var stats models.QuizStatistics
	var difficultyRating float64 // temporary variable for difficulty_rating which isn't in model
	err := r.db.QueryRow(query, quizID).Scan(
		&stats.QuizID, &stats.TotalAttempts, &stats.CompletedAttempts,
		&stats.AverageScore, &stats.AverageTime, &difficultyRating, &stats.PopularityScore,
		&stats.UpdatedAt,
	)

	// Set default values for fields not retrieved from database
	stats.ID = stats.QuizID           // Use QuizID as ID since it's the primary key
	stats.CreatedAt = stats.UpdatedAt // Use updated_at as created_at for now
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("quiz statistics not found")
		}
		return nil, fmt.Errorf("failed to get quiz statistics: %w", err)
	}

	return &stats, nil
}

// Quiz attempt operations
func (r *QuizRepository) CreateQuizAttempt(attempt *models.QuizAttempt) error {
	query := `
		INSERT INTO quiz_attempts (id, quiz_id, user_id, answers, score, total_points, time_spent,
		                         percentage_score, passed, status, is_completed, started_at,
		                         completed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`

	// Convert answers to JSON for storage
	answersJSON, err := json.Marshal(attempt.Answers)
	if err != nil {
		return fmt.Errorf("failed to marshal answers: %w", err)
	}

	_, err = r.db.Exec(query, attempt.ID, attempt.QuizID, attempt.UserID,
		answersJSON, attempt.Score, attempt.TotalPoints, attempt.TimeSpent,
		attempt.PercentageScore, attempt.Passed, attempt.Status, attempt.IsCompleted,
		attempt.StartedAt, attempt.CompletedAt, attempt.CreatedAt, attempt.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create quiz attempt: %w", err)
	}

	return nil
}

func (r *QuizRepository) UpdateQuizAttempt(attempt *models.QuizAttempt) error {
	query := `
		UPDATE quiz_attempts
		SET answers = $3, score = $4, total_points = $5, time_spent = $6,
		    percentage_score = $7, passed = $8, status = $9, is_completed = $10,
		    completed_at = $11, updated_at = $12
		WHERE id = $1 AND user_id = $2`

	// Convert answers to JSON for storage
	answersJSON, err := json.Marshal(attempt.Answers)
	if err != nil {
		return fmt.Errorf("failed to marshal answers: %w", err)
	}

	result, err := r.db.Exec(query, attempt.ID, attempt.UserID, answersJSON,
		attempt.Score, attempt.TotalPoints, attempt.TimeSpent, attempt.PercentageScore,
		attempt.Passed, attempt.Status, attempt.IsCompleted, attempt.CompletedAt, attempt.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update quiz attempt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("quiz attempt not found")
	}

	return nil
}

// GetUserCompletedAttempts retrieves all completed attempts for a quiz by a user
func (r *QuizRepository) GetUserCompletedAttempts(userID, quizID uuid.UUID) ([]models.QuizAttempt, error) {
	query := `
		SELECT id, quiz_id, user_id, answers, score, total_points, time_spent,
		       percentage_score, passed, status, is_completed,
		       started_at, completed_at, created_at, updated_at
		FROM quiz_attempts
		WHERE user_id = $1 AND quiz_id = $2 AND is_completed = true
		ORDER BY completed_at DESC`

	rows, err := r.db.Query(query, userID, quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed quiz attempts: %w", err)
	}
	defer rows.Close()

	var attempts []models.QuizAttempt
	for rows.Next() {
		var attempt models.QuizAttempt
		var answersJSON []byte

		err := rows.Scan(
			&attempt.ID, &attempt.QuizID, &attempt.UserID, &answersJSON,
			&attempt.Score, &attempt.TotalPoints, &attempt.TimeSpent,
			&attempt.PercentageScore, &attempt.Passed, &attempt.Status,
			&attempt.IsCompleted, &attempt.StartedAt, &attempt.CompletedAt,
			&attempt.CreatedAt, &attempt.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quiz attempt: %w", err)
		}

		// Unmarshal answers JSON
		if len(answersJSON) > 0 {
			err = json.Unmarshal(answersJSON, &attempt.Answers)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal answers: %w", err)
			}
		}

		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

// HasUserCompletedQuiz checks if a user has completed a quiz (enforces one-attempt-per-quiz policy)
func (r *QuizRepository) HasUserCompletedQuiz(userID, quizID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM quiz_attempts
			WHERE user_id = $1 AND quiz_id = $2 AND is_completed = true
		)`

	var hasCompleted bool
	err := r.db.QueryRow(query, userID, quizID).Scan(&hasCompleted)
	if err != nil {
		return false, fmt.Errorf("failed to check quiz completion status: %w", err)
	}

	return hasCompleted, nil
}

func (r *QuizRepository) GetQuizAttempt(id uuid.UUID) (*models.QuizAttempt, error) {
	query := `
		SELECT id, quiz_id, user_id, answers, score, total_points, time_spent,
		       percentage_score, passed, status, is_completed,
		       started_at, completed_at, created_at, updated_at
		FROM quiz_attempts
		WHERE id = $1`

	var attempt models.QuizAttempt
	var answersJSON []byte
	err := r.db.QueryRow(query, id).Scan(
		&attempt.ID, &attempt.QuizID, &attempt.UserID, &answersJSON,
		&attempt.Score, &attempt.TotalPoints, &attempt.TimeSpent,
		&attempt.PercentageScore, &attempt.Passed, &attempt.Status, &attempt.IsCompleted,
		&attempt.StartedAt, &attempt.CompletedAt, &attempt.CreatedAt, &attempt.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("quiz attempt not found")
		}
		return nil, fmt.Errorf("failed to get quiz attempt: %w", err)
	}

	// Unmarshal answers JSON
	if len(answersJSON) > 0 {
		err = json.Unmarshal(answersJSON, &attempt.Answers)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal answers: %w", err)
		}
	}

	return &attempt, nil
}

func (r *QuizRepository) GetUserQuizAttempts(userID, quizID uuid.UUID) ([]models.QuizAttempt, error) {
	query := `
		SELECT id, quiz_id, user_id, answers, score, total_points, time_spent,
		       percentage_score, passed, status, is_completed,
		       started_at, completed_at, created_at, updated_at
		FROM quiz_attempts
		WHERE user_id = $1 AND quiz_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID, quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to query quiz attempts: %w", err)
	}
	defer rows.Close()

	var attempts []models.QuizAttempt
	for rows.Next() {
		var attempt models.QuizAttempt
		var answersJSON []byte
		err := rows.Scan(
			&attempt.ID, &attempt.QuizID, &attempt.UserID, &answersJSON,
			&attempt.Score, &attempt.TotalPoints, &attempt.TimeSpent,
			&attempt.PercentageScore, &attempt.Passed, &attempt.Status, &attempt.IsCompleted,
			&attempt.StartedAt, &attempt.CompletedAt, &attempt.CreatedAt, &attempt.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quiz attempt: %w", err)
		}

		// Unmarshal answers JSON
		if answersJSON != nil {
			err = json.Unmarshal(answersJSON, &attempt.Answers)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal answers: %w", err)
			}
		}

		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

// GetActiveQuizAttempt retrieves an active (incomplete) quiz attempt for a user and quiz
func (r *QuizRepository) GetActiveQuizAttempt(userID, quizID uuid.UUID) (*models.QuizAttempt, error) {
	query := `
		SELECT id, quiz_id, user_id, answers, score, total_points, time_spent,
		       percentage_score, passed, status, is_completed,
		       started_at, completed_at, created_at, updated_at
		FROM quiz_attempts
		WHERE user_id = $1 AND quiz_id = $2 AND is_completed = false
		ORDER BY created_at DESC
		LIMIT 1`

	var attempt models.QuizAttempt
	var answersJSON []byte
	err := r.db.QueryRow(query, userID, quizID).Scan(
		&attempt.ID, &attempt.QuizID, &attempt.UserID, &answersJSON,
		&attempt.Score, &attempt.TotalPoints, &attempt.TimeSpent,
		&attempt.PercentageScore, &attempt.Passed, &attempt.Status, &attempt.IsCompleted,
		&attempt.StartedAt, &attempt.CompletedAt, &attempt.CreatedAt, &attempt.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No active attempt found
		}
		return nil, fmt.Errorf("failed to get active quiz attempt: %w", err)
	}

	// Unmarshal answers JSON
	if answersJSON != nil {
		err = json.Unmarshal(answersJSON, &attempt.Answers)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal answers: %w", err)
		}
	}

	return &attempt, nil
}

// DeleteActiveQuizAttempt deletes any active (incomplete) quiz attempts for a user and quiz
func (r *QuizRepository) DeleteActiveQuizAttempt(userID, quizID uuid.UUID) error {
	query := `
		DELETE FROM quiz_attempts
		WHERE user_id = $1 AND quiz_id = $2 AND is_completed = false`

	result, err := r.db.Exec(query, userID, quizID)
	if err != nil {
		return fmt.Errorf("failed to delete active quiz attempt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// Log how many attempts were deleted (for debugging)
	if rowsAffected > 0 {
		fmt.Printf("Deleted %d active quiz attempt(s) for user %s and quiz %s\n", rowsAffected, userID, quizID)
	}

	return nil
}

// GetUserAttempts retrieves user's quiz attempts with filtering and pagination
func (r *QuizRepository) GetUserAttempts(userID uuid.UUID, filters *models.AttemptFilters) ([]models.QuizAttemptWithDetails, int, error) {
	fmt.Printf("[QuizRepository] GetUserAttempts called for user: %s with filters: %+v\n", userID, filters)

	// Build WHERE conditions
	conditions := []string{"qa.user_id = $1"}
	args := []interface{}{userID}
	argCount := 1

	// Add filters
	if filters.QuizID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("qa.quiz_id = $%d", argCount))
		args = append(args, *filters.QuizID)
	}

	if filters.Category != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("q.category_id = $%d", argCount))
		args = append(args, filters.Category)
	}

	if filters.Difficulty != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("q.difficulty = $%d", argCount))
		args = append(args, filters.Difficulty)
	}

	if filters.StartDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("qa.completed_at >= $%d", argCount))
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("qa.completed_at <= $%d", argCount))
		args = append(args, *filters.EndDate)
	}

	if filters.MinScore != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("qa.score >= $%d", argCount))
		args = append(args, *filters.MinScore)
	}

	if filters.MaxScore != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("qa.score <= $%d", argCount))
		args = append(args, *filters.MaxScore)
	}

	// Only include completed attempts
	conditions = append(conditions, "qa.is_completed = true")

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// Build ORDER BY clause
	orderClause := fmt.Sprintf("ORDER BY qa.%s %s", filters.SortBy, strings.ToUpper(filters.SortOrder))

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM quiz_attempts qa
		JOIN quizzes q ON qa.quiz_id = q.id
		%s`, whereClause)

	fmt.Printf("[QuizRepository] Executing count query: %s with args: %v\n", countQuery, args)

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		fmt.Printf("[QuizRepository] Error executing count query: %v\n", err)
		return nil, 0, fmt.Errorf("failed to count user attempts: %w", err)
	}

	fmt.Printf("[QuizRepository] Count query returned total: %d\n", total)

	// Calculate pagination
	offset := (filters.Page - 1) * filters.PageSize

	// Main query
	query := fmt.Sprintf(`
		SELECT
			qa.id, qa.quiz_id, qa.user_id, qa.answers, qa.score, qa.total_points, qa.time_spent,
			qa.percentage_score, qa.passed, qa.status, qa.is_completed, qa.started_at,
			qa.completed_at, qa.created_at, qa.updated_at,
			q.id, q.title, q.description, q.category_id, q.difficulty, q.time_limit_minutes,
			q.total_questions, q.points, q.is_featured, q.tags, q.thumbnail_url, q.created_at
		FROM quiz_attempts qa
		JOIN quizzes q ON qa.quiz_id = q.id
		%s
		%s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderClause, argCount+1, argCount+2)

	args = append(args, filters.PageSize, offset)

	fmt.Printf("[QuizRepository] Executing main query: %s\n", query)
	fmt.Printf("[QuizRepository] Query args: %v\n", args)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		fmt.Printf("[QuizRepository] Error executing main query: %v\n", err)
		return nil, 0, fmt.Errorf("failed to get user attempts: %w", err)
	}
	defer rows.Close()

	attempts := make([]models.QuizAttemptWithDetails, 0)
	for rows.Next() {
		var attempt models.QuizAttemptWithDetails
		var tags pq.StringArray
		var answersJSON []byte

		err := rows.Scan(
			&attempt.ID, &attempt.QuizID, &attempt.UserID, &answersJSON, &attempt.Score,
			&attempt.TotalPoints, &attempt.TimeSpent, &attempt.PercentageScore,
			&attempt.Passed, &attempt.Status, &attempt.IsCompleted,
			&attempt.StartedAt, &attempt.CompletedAt, &attempt.CreatedAt, &attempt.UpdatedAt,
			&attempt.Quiz.ID, &attempt.Quiz.Title, &attempt.Quiz.Description,
			&attempt.Quiz.Category, &attempt.Quiz.Difficulty, &attempt.Quiz.TimeLimit,
			&attempt.Quiz.QuestionCount, &attempt.Quiz.Points, &attempt.Quiz.IsFeatured, &tags,
			&attempt.Quiz.ThumbnailURL, &attempt.Quiz.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan attempt row: %w", err)
		}

		// Unmarshal answers JSON
		if len(answersJSON) > 0 {
			err = json.Unmarshal(answersJSON, &attempt.Answers)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal answers: %w", err)
			}
		}

		// Convert tags
		attempt.Quiz.Tags = []string(tags)
		// Convert time limit from minutes to seconds for JSON response
		attempt.Quiz.TimeLimit = attempt.Quiz.TimeLimit * 60

		attempts = append(attempts, attempt)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("[QuizRepository] Error iterating rows: %v\n", err)
		return nil, 0, fmt.Errorf("error iterating attempt rows: %w", err)
	}

	fmt.Printf("[QuizRepository] Successfully retrieved %d attempts for user %s\n", len(attempts), userID)
	return attempts, total, nil
}

// GetAttemptWithDetails retrieves a single attempt with quiz details
func (r *QuizRepository) GetAttemptWithDetails(attemptID uuid.UUID) (*models.QuizAttemptWithDetails, error) {
	fmt.Printf("[QuizRepository] GetAttemptWithDetails called for attempt: %s\n", attemptID)

	query := `
		SELECT
			qa.id, qa.quiz_id, qa.user_id, qa.answers, qa.score, qa.total_points, qa.time_spent,
			qa.percentage_score, qa.passed, qa.status, qa.is_completed, qa.started_at,
			qa.completed_at, qa.created_at, qa.updated_at,
			q.id, q.title, q.description, q.category_id, q.difficulty, q.time_limit_minutes,
			q.total_questions, q.points, q.is_featured, q.tags, q.thumbnail_url, q.created_at
		FROM quiz_attempts qa
		JOIN quizzes q ON qa.quiz_id = q.id
		WHERE qa.id = $1`

	fmt.Printf("[QuizRepository] Executing single attempt query for attempt: %s\n", attemptID)

	var attempt models.QuizAttemptWithDetails
	var tags pq.StringArray
	var answersJSON []byte

	err := r.db.QueryRow(query, attemptID).Scan(
		&attempt.ID, &attempt.QuizID, &attempt.UserID, &answersJSON, &attempt.Score,
		&attempt.TotalPoints, &attempt.TimeSpent, &attempt.PercentageScore,
		&attempt.Passed, &attempt.Status, &attempt.IsCompleted,
		&attempt.StartedAt, &attempt.CompletedAt, &attempt.CreatedAt, &attempt.UpdatedAt,
		&attempt.Quiz.ID, &attempt.Quiz.Title, &attempt.Quiz.Description,
		&attempt.Quiz.Category, &attempt.Quiz.Difficulty, &attempt.Quiz.TimeLimit,
		&attempt.Quiz.QuestionCount, &attempt.Quiz.Points, &attempt.Quiz.IsFeatured, &tags,
		&attempt.Quiz.ThumbnailURL, &attempt.Quiz.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("[QuizRepository] Attempt not found: %s\n", attemptID)
			return nil, fmt.Errorf("attempt not found")
		}
		fmt.Printf("[QuizRepository] Error scanning attempt details: %v\n", err)
		return nil, fmt.Errorf("failed to get attempt with details: %w", err)
	}

	// Unmarshal answers JSON
	if answersJSON != nil {
		err = json.Unmarshal(answersJSON, &attempt.Answers)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal answers: %w", err)
		}
	}

	// Convert tags
	attempt.Quiz.Tags = []string(tags)
	// Convert time limit from minutes to seconds for JSON response
	attempt.Quiz.TimeLimit = attempt.Quiz.TimeLimit * 60

	fmt.Printf("[QuizRepository] Successfully retrieved attempt details: user=%s, quiz='%s', score=%.2f\n",
		attempt.UserID, attempt.Quiz.Title, attempt.Score)

	return &attempt, nil
}

// GetUserQuizAttemptWithDetails retrieves THE user's attempt for a specific quiz with quiz details
// Since one-attempt-per-quiz policy is enforced, this returns a single result
func (r *QuizRepository) GetUserQuizAttemptWithDetails(userID, quizID uuid.UUID) (*models.QuizAttemptWithDetails, error) {
	query := `
		SELECT
			qa.id, qa.quiz_id, qa.user_id, qa.answers, qa.score, qa.total_points, qa.time_spent,
			qa.percentage_score, qa.passed, qa.status, qa.is_completed, qa.started_at,
			qa.completed_at, qa.created_at, qa.updated_at,
			q.id, q.title, q.description, q.category_id, q.difficulty, q.time_limit_minutes,
			q.total_questions, q.points, q.is_featured, q.tags, q.thumbnail_url, q.created_at
		FROM quiz_attempts qa
		JOIN quizzes q ON qa.quiz_id = q.id
		WHERE qa.user_id = $1 AND qa.quiz_id = $2 AND qa.is_completed = true
		LIMIT 1`

	var attempt models.QuizAttemptWithDetails
	var tags pq.StringArray
	var answersJSON []byte

	err := r.db.QueryRow(query, userID, quizID).Scan(
		&attempt.ID, &attempt.QuizID, &attempt.UserID, &answersJSON, &attempt.Score,
		&attempt.TotalPoints, &attempt.TimeSpent, &attempt.PercentageScore,
		&attempt.Passed, &attempt.Status, &attempt.IsCompleted,
		&attempt.StartedAt, &attempt.CompletedAt, &attempt.CreatedAt, &attempt.UpdatedAt,
		&attempt.Quiz.ID, &attempt.Quiz.Title, &attempt.Quiz.Description,
		&attempt.Quiz.Category, &attempt.Quiz.Difficulty, &attempt.Quiz.TimeLimit,
		&attempt.Quiz.QuestionCount, &attempt.Quiz.Points, &attempt.Quiz.IsFeatured, &tags,
		&attempt.Quiz.ThumbnailURL, &attempt.Quiz.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no attempt found for this quiz")
		}
		return nil, fmt.Errorf("failed to get user quiz attempt: %w", err)
	}

	// Unmarshal answers JSON
	if answersJSON != nil {
		err = json.Unmarshal(answersJSON, &attempt.Answers)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal answers: %w", err)
		}
	}

	// Convert tags and time limit
	attempt.Quiz.Tags = []string(tags)
	attempt.Quiz.TimeLimit = attempt.Quiz.TimeLimit * 60 // Convert minutes to seconds

	return &attempt, nil
}

// CreateOrUpdateQuizStatistics creates or updates quiz statistics when an attempt is completed
func (r *QuizRepository) CreateOrUpdateQuizStatistics(quizID uuid.UUID, score float64, timeSpent int) error {
	// First, try to get existing statistics
	existingStats, err := r.GetQuizStatistics(quizID)

	if err != nil && err.Error() != "quiz statistics not found" {
		return fmt.Errorf("failed to check existing quiz statistics: %w", err)
	}

	if existingStats == nil {
		// Create new statistics record
		query := `
			INSERT INTO quiz_statistics (quiz_id, total_attempts, total_completions, average_score,
			                          average_time_seconds, difficulty_rating, popularity_score, updated_at)
			VALUES ($1, 1, 1, $2, $3, 0.0, 1, CURRENT_TIMESTAMP)`

		_, err = r.db.Exec(query, quizID, score, timeSpent)
		if err != nil {
			return fmt.Errorf("failed to create quiz statistics: %w", err)
		}
	} else {
		// Update existing statistics
		newTotalAttempts := existingStats.TotalAttempts + 1
		newTotalCompletions := existingStats.CompletedAttempts + 1

		// Calculate new averages
		newAverageScore := ((existingStats.AverageScore * float64(existingStats.CompletedAttempts)) + score) / float64(newTotalCompletions)
		newAverageTime := ((existingStats.AverageTime * existingStats.CompletedAttempts) + timeSpent) / newTotalCompletions

		// Update popularity score (simple increment for now)
		newPopularityScore := existingStats.PopularityScore + 1

		query := `
			UPDATE quiz_statistics
			SET total_attempts = $2, total_completions = $3, average_score = $4,
			    average_time_seconds = $5, popularity_score = $6, updated_at = CURRENT_TIMESTAMP
			WHERE quiz_id = $1`

		_, err = r.db.Exec(query, quizID, newTotalAttempts, newTotalCompletions,
			newAverageScore, newAverageTime, newPopularityScore)
		if err != nil {
			return fmt.Errorf("failed to update quiz statistics: %w", err)
		}
	}

	return nil
}

// AddFavorite adds a quiz to user's favorites
func (r *QuizRepository) AddFavorite(userID, quizID uuid.UUID) error {
	query := `
		INSERT INTO user_quiz_favorites (id, user_id, quiz_id, favorited_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, quiz_id) DO NOTHING`

	_, err := r.db.Exec(query, uuid.New(), userID, quizID)
	if err != nil {
		return fmt.Errorf("failed to add favorite: %w", err)
	}

	return nil
}

// RemoveFavorite removes a quiz from user's favorites
func (r *QuizRepository) RemoveFavorite(userID, quizID uuid.UUID) error {
	query := `DELETE FROM user_quiz_favorites WHERE user_id = $1 AND quiz_id = $2`

	result, err := r.db.Exec(query, userID, quizID)
	if err != nil {
		return fmt.Errorf("failed to remove favorite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("favorite not found")
	}

	return nil
}

// GetUserFavorites retrieves user's favorite quizzes with pagination
func (r *QuizRepository) GetUserFavorites(userID uuid.UUID, page, pageSize int) ([]models.UserQuizFavorite, int, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get total count
	countQuery := `SELECT COUNT(*) FROM user_quiz_favorites WHERE user_id = $1`
	var total int
	err := r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get favorites count: %w", err)
	}

	// Get favorites with quiz details
	query := `
		SELECT
			f.id, f.user_id, f.quiz_id, f.favorited_at,
			q.id, q.title, q.description, q.category_id, q.difficulty,
			q.time_limit_minutes, q.total_questions, q.points, q.is_featured,
			q.tags, q.thumbnail_url, q.created_at
		FROM user_quiz_favorites f
		JOIN quizzes q ON f.quiz_id = q.id
		WHERE f.user_id = $1
		ORDER BY f.favorited_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get favorites: %w", err)
	}
	defer rows.Close()

	favorites := make([]models.UserQuizFavorite, 0)
	for rows.Next() {
		var favorite models.UserQuizFavorite
		var quiz models.QuizSummary
		var tags pq.StringArray

		err := rows.Scan(
			&favorite.ID, &favorite.UserID, &favorite.QuizID, &favorite.FavoritedAt,
			&quiz.ID, &quiz.Title, &quiz.Description, &quiz.Category, &quiz.Difficulty,
			&quiz.TimeLimit, &quiz.QuestionCount, &quiz.Points, &quiz.IsFeatured,
			&tags, &quiz.ThumbnailURL, &quiz.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan favorite: %w", err)
		}

		// Convert tags and time limit
		quiz.Tags = []string(tags)
		quiz.TimeLimit = quiz.TimeLimit * 60 // Convert minutes to seconds

		favorite.Quiz = &quiz
		favorites = append(favorites, favorite)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate favorites: %w", err)
	}

	return favorites, total, nil
}

// IsFavorite checks if a quiz is in user's favorites
func (r *QuizRepository) IsFavorite(userID, quizID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_quiz_favorites WHERE user_id = $1 AND quiz_id = $2)`

	var exists bool
	err := r.db.QueryRow(query, userID, quizID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check favorite status: %w", err)
	}

	return exists, nil
}
