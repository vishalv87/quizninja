package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"quizninja-api/database"
	"quizninja-api/models"
	"strings"

	"github.com/google/uuid"
)

type QuizSessionRepository struct {
	db *sql.DB
}

func NewQuizSessionRepository() *QuizSessionRepository {
	return &QuizSessionRepository{
		db: database.DB,
	}
}

// CreateSession creates a new quiz session
func (r *QuizSessionRepository) CreateSession(session *models.QuizSession) error {
	query := `
		INSERT INTO quiz_sessions (id, attempt_id, user_id, quiz_id, current_question_index,
		                          current_answers, session_state, time_remaining, time_spent_so_far,
		                          last_activity_at, created_at, updated_at, is_test_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	// Convert current answers to JSON
	answersJSON, err := json.Marshal(session.CurrentAnswers)
	if err != nil {
		return fmt.Errorf("failed to marshal current answers: %w", err)
	}

	_, err = r.db.Exec(query,
		session.ID, session.AttemptID, session.UserID, session.QuizID,
		session.CurrentQuestionIndex, answersJSON, session.SessionState,
		session.TimeRemaining, session.TimeSpentSoFar, session.LastActivityAt,
		session.CreatedAt, session.UpdatedAt, session.IsTestData)
	if err != nil {
		return fmt.Errorf("failed to create quiz session: %w", err)
	}

	return nil
}

// UpdateSession updates an existing quiz session
func (r *QuizSessionRepository) UpdateSession(session *models.QuizSession) error {
	query := `
		UPDATE quiz_sessions
		SET current_question_index = $3, current_answers = $4, session_state = $5,
		    time_remaining = $6, time_spent_so_far = $7, last_activity_at = $8,
		    paused_at = $9, updated_at = $10
		WHERE id = $1 AND user_id = $2`

	// Convert current answers to JSON
	answersJSON, err := json.Marshal(session.CurrentAnswers)
	if err != nil {
		return fmt.Errorf("failed to marshal current answers: %w", err)
	}

	result, err := r.db.Exec(query,
		session.ID, session.UserID, session.CurrentQuestionIndex, answersJSON,
		session.SessionState, session.TimeRemaining, session.TimeSpentSoFar,
		session.LastActivityAt, session.PausedAt, session.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update quiz session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("quiz session not found or unauthorized")
	}

	return nil
}

// GetSessionByID retrieves a quiz session by its ID
func (r *QuizSessionRepository) GetSessionByID(id uuid.UUID) (*models.QuizSession, error) {
	query := `
		SELECT id, attempt_id, user_id, quiz_id, current_question_index,
		       current_answers, session_state, time_remaining, time_spent_so_far,
		       last_activity_at, paused_at, created_at, updated_at
		FROM quiz_sessions
		WHERE id = $1`

	return r.scanQuizSession(r.db.QueryRow(query, id))
}

// GetSessionByAttemptID retrieves a quiz session by attempt ID
func (r *QuizSessionRepository) GetSessionByAttemptID(attemptID uuid.UUID) (*models.QuizSession, error) {
	query := `
		SELECT id, attempt_id, user_id, quiz_id, current_question_index,
		       current_answers, session_state, time_remaining, time_spent_so_far,
		       last_activity_at, paused_at, created_at, updated_at
		FROM quiz_sessions
		WHERE attempt_id = $1`

	return r.scanQuizSession(r.db.QueryRow(query, attemptID))
}

// DeleteSession deletes a quiz session
func (r *QuizSessionRepository) DeleteSession(id uuid.UUID) error {
	query := `DELETE FROM quiz_sessions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete quiz session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("quiz session not found")
	}

	return nil
}

// PauseSession pauses a quiz session
func (r *QuizSessionRepository) PauseSession(attemptID uuid.UUID, pauseData *models.PauseSessionRequest) error {
	query := `
		UPDATE quiz_sessions
		SET current_question_index = $2, current_answers = $3, session_state = 'paused',
		    time_spent_so_far = $4, time_remaining = $5, paused_at = CURRENT_TIMESTAMP,
		    last_activity_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE attempt_id = $1 AND session_state = 'active'`

	// Convert current answers to JSON
	answersJSON, err := json.Marshal(pauseData.CurrentAnswers)
	if err != nil {
		return fmt.Errorf("failed to marshal current answers: %w", err)
	}

	result, err := r.db.Exec(query, attemptID, pauseData.CurrentQuestionIndex,
		answersJSON, pauseData.TimeSpentSoFar, pauseData.TimeRemaining)
	if err != nil {
		return fmt.Errorf("failed to pause quiz session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("active quiz session not found for attempt")
	}

	// Also update the quiz attempt status
	updateAttemptQuery := `
		UPDATE quiz_attempts
		SET status = 'paused', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = 'started'`

	_, err = r.db.Exec(updateAttemptQuery, attemptID)
	if err != nil {
		return fmt.Errorf("failed to update quiz attempt status: %w", err)
	}

	return nil
}

// ResumeSession resumes a paused quiz session
func (r *QuizSessionRepository) ResumeSession(attemptID uuid.UUID) error {
	query := `
		UPDATE quiz_sessions
		SET session_state = 'active', paused_at = NULL,
		    last_activity_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE attempt_id = $1 AND session_state = 'paused'`

	result, err := r.db.Exec(query, attemptID)
	if err != nil {
		return fmt.Errorf("failed to resume quiz session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("paused quiz session not found for attempt")
	}

	// Also update the quiz attempt status
	updateAttemptQuery := `
		UPDATE quiz_attempts
		SET status = 'started', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = 'paused'`

	_, err = r.db.Exec(updateAttemptQuery, attemptID)
	if err != nil {
		return fmt.Errorf("failed to update quiz attempt status: %w", err)
	}

	return nil
}

// AbandonSession marks a session as abandoned
func (r *QuizSessionRepository) AbandonSession(attemptID uuid.UUID) error {
	query := `
		UPDATE quiz_sessions
		SET session_state = 'abandoned', last_activity_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE attempt_id = $1 AND session_state IN ('active', 'paused')`

	result, err := r.db.Exec(query, attemptID)
	if err != nil {
		return fmt.Errorf("failed to abandon quiz session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("active/paused quiz session not found for attempt")
	}

	// Also update the quiz attempt status
	updateAttemptQuery := `
		UPDATE quiz_attempts
		SET status = 'abandoned', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status IN ('started', 'paused')`

	_, err = r.db.Exec(updateAttemptQuery, attemptID)
	if err != nil {
		return fmt.Errorf("failed to update quiz attempt status: %w", err)
	}

	return nil
}

// CompleteSession marks a session as completed
func (r *QuizSessionRepository) CompleteSession(attemptID uuid.UUID) error {
	query := `
		UPDATE quiz_sessions
		SET session_state = 'completed', last_activity_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE attempt_id = $1 AND session_state IN ('active', 'paused')`

	result, err := r.db.Exec(query, attemptID)
	if err != nil {
		return fmt.Errorf("failed to complete quiz session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("active/paused quiz session not found for attempt")
	}

	return nil
}

// GetActiveSession retrieves an active or paused session for a user and quiz
func (r *QuizSessionRepository) GetActiveSession(userID, quizID uuid.UUID) (*models.QuizSession, error) {
	query := `
		SELECT id, attempt_id, user_id, quiz_id, current_question_index,
		       current_answers, session_state, time_remaining, time_spent_so_far,
		       last_activity_at, paused_at, created_at, updated_at
		FROM quiz_sessions
		WHERE user_id = $1 AND quiz_id = $2 AND session_state IN ('active', 'paused')
		ORDER BY last_activity_at DESC
		LIMIT 1`

	return r.scanQuizSession(r.db.QueryRow(query, userID, quizID))
}

// GetUserActiveSessions retrieves all active/paused sessions for a user with filtering
func (r *QuizSessionRepository) GetUserActiveSessions(userID uuid.UUID, filters *models.SessionFilters) ([]models.QuizSessionWithDetails, int, error) {
	// Build WHERE conditions
	conditions := []string{"qs.user_id = $1"}
	args := []interface{}{userID}
	argCount := 1

	// Add filters
	if filters.SessionState != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("qs.session_state = $%d", argCount))
		args = append(args, filters.SessionState)
	} else {
		// Default to active and paused sessions only
		conditions = append(conditions, "qs.session_state IN ('active', 'paused')")
	}

	if filters.QuizID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("qs.quiz_id = $%d", argCount))
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
		conditions = append(conditions, fmt.Sprintf("qs.created_at >= $%d", argCount))
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("qs.created_at <= $%d", argCount))
		args = append(args, *filters.EndDate)
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// Build ORDER BY clause
	orderClause := fmt.Sprintf("ORDER BY qs.%s %s", filters.SortBy, strings.ToUpper(filters.SortOrder))

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM quiz_sessions qs
		JOIN quizzes q ON qs.quiz_id = q.id
		%s`, whereClause)

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user active sessions: %w", err)
	}

	// Calculate pagination
	offset := (filters.Page - 1) * filters.PageSize

	// Main query
	query := fmt.Sprintf(`
		SELECT
			qs.id, qs.attempt_id, qs.user_id, qs.quiz_id, qs.current_question_index,
			qs.current_answers, qs.session_state, qs.time_remaining, qs.time_spent_so_far,
			qs.last_activity_at, qs.paused_at, qs.created_at, qs.updated_at, qs.is_test_data,
			q.title, q.category_id, q.difficulty, q.thumbnail_url, q.total_questions,
			q.time_limit_minutes, q.is_test_data
		FROM quiz_sessions qs
		JOIN quizzes q ON qs.quiz_id = q.id
		%s
		%s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderClause, argCount+1, argCount+2)

	args = append(args, filters.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user active sessions: %w", err)
	}
	defer rows.Close()

	sessions := make([]models.QuizSessionWithDetails, 0)
	for rows.Next() {
		var session models.QuizSessionWithDetails
		var answersJSON []byte

		err := rows.Scan(
			&session.ID, &session.AttemptID, &session.UserID, &session.QuizID,
			&session.CurrentQuestionIndex, &answersJSON, &session.SessionState,
			&session.TimeRemaining, &session.TimeSpentSoFar, &session.LastActivityAt,
			&session.PausedAt, &session.CreatedAt, &session.UpdatedAt, &session.IsTestData,
			&session.QuizTitle, &session.QuizCategory, &session.QuizDifficulty,
			&session.QuizThumbnail, &session.TotalQuestions, &session.OriginalTimeLimit,
			&session.QuizIsTestData,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan session row: %w", err)
		}

		// Unmarshal current answers JSON
		if len(answersJSON) > 0 {
			err = json.Unmarshal(answersJSON, &session.CurrentAnswers)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal current answers: %w", err)
			}
		}

		// Convert time limit from minutes to seconds
		session.OriginalTimeLimit = session.OriginalTimeLimit * 60

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating session rows: %w", err)
	}

	return sessions, total, nil
}

// GetSessionWithDetails retrieves a session with full quiz details
func (r *QuizSessionRepository) GetSessionWithDetails(sessionID uuid.UUID) (*models.QuizSessionWithDetails, error) {
	query := `
		SELECT
			qs.id, qs.attempt_id, qs.user_id, qs.quiz_id, qs.current_question_index,
			qs.current_answers, qs.session_state, qs.time_remaining, qs.time_spent_so_far,
			qs.last_activity_at, qs.paused_at, qs.created_at, qs.updated_at,
			q.title, q.category_id, q.difficulty, q.thumbnail_url, q.total_questions,
			q.time_limit_minutes, q.is_test_data
		FROM quiz_sessions qs
		JOIN quizzes q ON qs.quiz_id = q.id
		WHERE qs.id = $1`

	var session models.QuizSessionWithDetails
	var answersJSON []byte

	err := r.db.QueryRow(query, sessionID).Scan(
		&session.ID, &session.AttemptID, &session.UserID, &session.QuizID,
		&session.CurrentQuestionIndex, &answersJSON, &session.SessionState,
		&session.TimeRemaining, &session.TimeSpentSoFar, &session.LastActivityAt,
		&session.PausedAt, &session.CreatedAt, &session.UpdatedAt,
		&session.QuizTitle, &session.QuizCategory, &session.QuizDifficulty,
		&session.QuizThumbnail, &session.TotalQuestions, &session.OriginalTimeLimit,
		&session.QuizIsTestData,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session with details: %w", err)
	}

	// Unmarshal current answers JSON
	if len(answersJSON) > 0 {
		err = json.Unmarshal(answersJSON, &session.CurrentAnswers)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal current answers: %w", err)
		}
	}

	// Convert time limit from minutes to seconds
	session.OriginalTimeLimit = session.OriginalTimeLimit * 60

	return &session, nil
}

// SaveSessionProgress saves the current progress of a quiz session
func (r *QuizSessionRepository) SaveSessionProgress(sessionID uuid.UUID, updateData *models.UpdateQuizSessionRequest) error {
	query := `
		UPDATE quiz_sessions
		SET current_question_index = $2, current_answers = $3, time_spent_so_far = $4,
		    time_remaining = $5, last_activity_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND session_state IN ('active', 'paused')`

	// Convert current answers to JSON
	answersJSON, err := json.Marshal(updateData.CurrentAnswers)
	if err != nil {
		return fmt.Errorf("failed to marshal current answers: %w", err)
	}

	result, err := r.db.Exec(query, sessionID, updateData.CurrentQuestionIndex,
		answersJSON, updateData.TimeSpentSoFar, updateData.TimeRemaining)
	if err != nil {
		return fmt.Errorf("failed to save session progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("active/paused session not found")
	}

	return nil
}

// UpdateSessionActivity updates the last activity timestamp for a session
func (r *QuizSessionRepository) UpdateSessionActivity(sessionID uuid.UUID) error {
	query := `
		UPDATE quiz_sessions
		SET last_activity_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND session_state IN ('active', 'paused')`

	result, err := r.db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session activity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("active/paused session not found")
	}

	return nil
}

// CleanupExpiredSessions cleans up sessions that have been inactive for too long
func (r *QuizSessionRepository) CleanupExpiredSessions() (int, error) {
	// Use the database function we created in the migration
	query := `SELECT cleanup_expired_quiz_sessions()`

	var cleanedCount int
	err := r.db.QueryRow(query).Scan(&cleanedCount)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	return cleanedCount, nil
}

// HasActiveSession checks if a user has an active session for a quiz
func (r *QuizSessionRepository) HasActiveSession(userID, quizID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM quiz_sessions
			WHERE user_id = $1 AND quiz_id = $2 AND session_state IN ('active', 'paused')
		)`

	var exists bool
	err := r.db.QueryRow(query, userID, quizID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check active session: %w", err)
	}

	return exists, nil
}

// CanResumeSession checks if a session can be resumed by a user
func (r *QuizSessionRepository) CanResumeSession(attemptID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM quiz_sessions
			WHERE attempt_id = $1 AND user_id = $2 AND session_state = 'paused'
			AND last_activity_at > CURRENT_TIMESTAMP - INTERVAL '24 hours'
		)`

	var canResume bool
	err := r.db.QueryRow(query, attemptID, userID).Scan(&canResume)
	if err != nil {
		return false, fmt.Errorf("failed to check resume session: %w", err)
	}

	return canResume, nil
}

// scanQuizSession is a helper function to scan a quiz session from a database row
func (r *QuizSessionRepository) scanQuizSession(row *sql.Row) (*models.QuizSession, error) {
	var session models.QuizSession
	var answersJSON []byte

	err := row.Scan(
		&session.ID, &session.AttemptID, &session.UserID, &session.QuizID,
		&session.CurrentQuestionIndex, &answersJSON, &session.SessionState,
		&session.TimeRemaining, &session.TimeSpentSoFar, &session.LastActivityAt,
		&session.PausedAt, &session.CreatedAt, &session.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("quiz session not found")
		}
		return nil, fmt.Errorf("failed to scan quiz session: %w", err)
	}

	// Unmarshal current answers JSON
	if len(answersJSON) > 0 {
		err = json.Unmarshal(answersJSON, &session.CurrentAnswers)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal current answers: %w", err)
		}
	}

	return &session, nil
}