package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"quizninja-api/database"
	"quizninja-api/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ChallengesRepository struct {
	db *sql.DB
}

// NewChallengesRepository creates a new challenges repository instance
func NewChallengesRepository() ChallengesRepositoryInterface {
	return &ChallengesRepository{
		db: database.DB,
	}
}

// CreateChallenge creates a new challenge
func (r *ChallengesRepository) CreateChallenge(challenge *models.Challenge) error {
	log.Printf("===== CreateChallenge Repository =====")
	log.Printf("DEBUG: Input challengerID=%s, challengeeID=%s, quizID=%s",
		challenge.ChallengerID, challenge.ChallengeeID, challenge.QuizID)
	log.Printf("DEBUG: Message=%v, ExpiresAt=%v, IsGroupChallenge=%v, IsTestData=%v",
		challenge.Message, challenge.ExpiresAt, challenge.IsGroupChallenge, challenge.IsTestData)

	query := `
		INSERT INTO challenges (
			challenger_id, challengee_id, quiz_id, message, expires_at,
			is_group_challenge, participant_ids, participant_scores, is_test_data
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, status, created_at, updated_at
	`
	log.Printf("DEBUG: SQL query prepared: %s", query)

	var participantScores interface{}
	if challenge.ParticipantScores != nil {
		participantScoresJSON, err := json.Marshal(challenge.ParticipantScores)
		if err != nil {
			log.Printf("DEBUG: Failed to marshal participant scores: %v", err)
			return fmt.Errorf("failed to marshal participant scores: %w", err)
		}
		participantScores = string(participantScoresJSON)
	} else {
		participantScores = nil
	}

	log.Printf("DEBUG: About to execute INSERT with params: [$1=%s, $2=%s, $3=%s, $4=%v, $5=%v, $6=%v, $7=%v, $8=%v, $9=%v]",
		challenge.ChallengerID, challenge.ChallengeeID, challenge.QuizID, challenge.Message,
		challenge.ExpiresAt, challenge.IsGroupChallenge, challenge.ParticipantIDs, participantScores, challenge.IsTestData)

	err := r.db.QueryRow(
		query,
		challenge.ChallengerID,
		challenge.ChallengeeID,
		challenge.QuizID,
		challenge.Message,
		challenge.ExpiresAt,
		challenge.IsGroupChallenge,
		pq.Array(challenge.ParticipantIDs),
		participantScores,
		challenge.IsTestData,
	).Scan(
		&challenge.ID,
		&challenge.Status,
		&challenge.CreatedAt,
		&challenge.UpdatedAt,
	)

	if err != nil {
		log.Printf("DEBUG: Database error occurred: %v", err)
		log.Printf("DEBUG: Error type: %T", err)
		return fmt.Errorf("failed to create challenge: %w", err)
	}

	log.Printf("DEBUG: Challenge created successfully with ID=%s, Status=%s", challenge.ID, challenge.Status)
	return nil
}

// GetChallengeByID retrieves a challenge by ID
func (r *ChallengesRepository) GetChallengeByID(id uuid.UUID) (*models.Challenge, error) {
	log.Printf("GetChallengeByID called: id=%s", id)

	query := `
		SELECT id, challenger_id, challengee_id, quiz_id, status,
			   challenger_score, challengee_score, message, expires_at,
			   is_group_challenge, participant_ids, participant_scores,
			   created_at, updated_at, is_test_data
		FROM challenges
		WHERE id = $1
	`

	var challenge models.Challenge
	var participantIDs pq.StringArray
	var participantScoresJSON sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&challenge.ID,
		&challenge.ChallengerID,
		&challenge.ChallengeeID,
		&challenge.QuizID,
		&challenge.Status,
		&challenge.ChallengerScore,
		&challenge.ChallengeeScore,
		&challenge.Message,
		&challenge.ExpiresAt,
		&challenge.IsGroupChallenge,
		&participantIDs,
		&participantScoresJSON,
		&challenge.CreatedAt,
		&challenge.UpdatedAt,
		&challenge.IsTestData,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("challenge not found")
		}
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}

	// Convert string array to UUID array
	for _, idStr := range participantIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			challenge.ParticipantIDs = append(challenge.ParticipantIDs, id)
		}
	}

	// Parse participant scores JSON
	if participantScoresJSON.Valid {
		err := json.Unmarshal([]byte(participantScoresJSON.String), &challenge.ParticipantScores)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal participant scores: %w", err)
		}
	}

	return &challenge, nil
}

// GetChallengeWithDetails retrieves a challenge with full user and quiz details
func (r *ChallengesRepository) GetChallengeWithDetails(id uuid.UUID) (*models.ChallengeWithDetails, error) {
	log.Printf("GetChallengeWithDetails called: id=%s", id)

	query := `
		SELECT c.id, c.challenger_id, c.challengee_id, c.quiz_id, c.status,
			   c.challenger_score, c.challengee_score, c.message, c.expires_at,
			   c.is_group_challenge, c.participant_ids, c.participant_scores,
			   c.created_at, c.updated_at, c.is_test_data,
			   u1.name, COALESCE(u1.avatar_url, ''),
			   u2.name, COALESCE(u2.avatar_url, ''),
			   q.title, q.category_id
		FROM challenges c
		JOIN users u1 ON c.challenger_id = u1.id
		JOIN users u2 ON c.challengee_id = u2.id
		JOIN quizzes q ON c.quiz_id = q.id
		WHERE c.id = $1
	`

	var challenge models.ChallengeWithDetails
	var participantIDs pq.StringArray
	var participantScoresJSON sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&challenge.ID,
		&challenge.ChallengerID,
		&challenge.ChallengeeID,
		&challenge.QuizID,
		&challenge.Status,
		&challenge.ChallengerScore,
		&challenge.ChallengeeScore,
		&challenge.Message,
		&challenge.ExpiresAt,
		&challenge.IsGroupChallenge,
		&participantIDs,
		&participantScoresJSON,
		&challenge.CreatedAt,
		&challenge.UpdatedAt,
		&challenge.IsTestData,
		&challenge.ChallengerName,
		&challenge.ChallengerAvatar,
		&challenge.ChallengeeName,
		&challenge.ChallengeeAvatar,
		&challenge.QuizTitle,
		&challenge.QuizCategory,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("challenge not found")
		}
		return nil, fmt.Errorf("failed to get challenge with details: %w", err)
	}

	// Convert string array to UUID array
	for _, idStr := range participantIDs {
		if id, err := uuid.Parse(idStr); err == nil {
			challenge.ParticipantIDs = append(challenge.ParticipantIDs, id)
		}
	}

	// Parse participant scores JSON
	if participantScoresJSON.Valid {
		err := json.Unmarshal([]byte(participantScoresJSON.String), &challenge.ParticipantScores)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal participant scores: %w", err)
		}
	}

	return &challenge, nil
}

// UpdateChallenge updates a challenge
func (r *ChallengesRepository) UpdateChallenge(challenge *models.Challenge) error {
	log.Printf("UpdateChallenge called: id=%s", challenge.ID)

	query := `
		UPDATE challenges
		SET challenger_score = $2, challengee_score = $3, status = $4,
			message = $5, expires_at = $6, participant_scores = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		challenge.ID,
		challenge.ChallengerScore,
		challenge.ChallengeeScore,
		challenge.Status,
		challenge.Message,
		challenge.ExpiresAt,
		challenge.ParticipantScores,
	)

	if err != nil {
		return fmt.Errorf("failed to update challenge: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("challenge not found")
	}

	return nil
}

// UpdateChallengeStatus updates only the status of a challenge
func (r *ChallengesRepository) UpdateChallengeStatus(challengeID uuid.UUID, status string) error {
	log.Printf("UpdateChallengeStatus called: id=%s, status=%s", challengeID, status)

	query := `
		UPDATE challenges
		SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(query, challengeID, status)
	if err != nil {
		return fmt.Errorf("failed to update challenge status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("challenge not found")
	}

	return nil
}

// UpdateChallengeScore updates the score for a specific user in a challenge
func (r *ChallengesRepository) UpdateChallengeScore(challengeID uuid.UUID, userID uuid.UUID, score float64) error {
	log.Printf("UpdateChallengeScore called: challengeID=%s, userID=%s, score=%f", challengeID, userID, score)

	// First get the challenge to determine which score to update
	challenge, err := r.GetChallengeByID(challengeID)
	if err != nil {
		return err
	}

	var query string
	if challenge.ChallengerID == userID {
		query = `
			UPDATE challenges
			SET challenger_score = $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
	} else if challenge.ChallengeeID == userID {
		query = `
			UPDATE challenges
			SET challengee_score = $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
	} else {
		return fmt.Errorf("user is not part of this challenge")
	}

	result, err := r.db.Exec(query, challengeID, score)
	if err != nil {
		return fmt.Errorf("failed to update challenge score: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("challenge not found")
	}

	// Check if both scores are now set and update status to completed
	updatedChallenge, err := r.GetChallengeByID(challengeID)
	if err != nil {
		return err
	}

	if updatedChallenge.ChallengerScore != nil && updatedChallenge.ChallengeeScore != nil {
		return r.UpdateChallengeStatus(challengeID, "completed")
	}

	return nil
}

// DeleteChallenge deletes a challenge
func (r *ChallengesRepository) DeleteChallenge(id uuid.UUID) error {
	log.Printf("DeleteChallenge called: id=%s", id)

	query := `DELETE FROM challenges WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete challenge: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("challenge not found")
	}

	return nil
}

// GetUserChallenges retrieves challenges for a user with filtering and pagination
func (r *ChallengesRepository) GetUserChallenges(userID uuid.UUID, filters *models.ChallengeFilters) ([]models.ChallengeWithDetails, int, error) {
	log.Printf("GetUserChallenges called: userID=%s", userID)

	// Build WHERE clause
	whereConditions := []string{"(c.challenger_id = $1 OR c.challengee_id = $1)"}
	args := []interface{}{userID}
	argIndex := 2

	if filters.Status != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("c.status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.QuizID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("c.quiz_id = $%d", argIndex))
		args = append(args, *filters.QuizID)
		argIndex++
	}

	if filters.UserType != "" && filters.UserType != "all" {
		if filters.UserType == "challenger" {
			whereConditions = append(whereConditions, fmt.Sprintf("c.challenger_id = $%d", argIndex))
		} else if filters.UserType == "challenged" {
			whereConditions = append(whereConditions, fmt.Sprintf("c.challengee_id = $%d", argIndex))
		}
		args = append(args, userID)
		argIndex++
	}

	if filters.StartDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("c.created_at >= $%d", argIndex))
		args = append(args, *filters.StartDate)
		argIndex++
	}

	if filters.EndDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("c.created_at <= $%d", argIndex))
		args = append(args, *filters.EndDate)
		argIndex++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM challenges c
		WHERE %s
	`, whereClause)

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count challenges: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "c.created_at DESC"
	if filters.SortBy != "" {
		direction := "DESC"
		if filters.SortOrder == "asc" {
			direction = "ASC"
		}
		orderBy = fmt.Sprintf("c.%s %s", filters.SortBy, direction)
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT c.id, c.challenger_id, c.challengee_id, c.quiz_id, c.status,
			   c.challenger_score, c.challengee_score, c.message, c.expires_at,
			   c.is_group_challenge, c.participant_ids, c.participant_scores,
			   c.created_at, c.updated_at, c.is_test_data,
			   u1.name, COALESCE(u1.avatar_url, ''),
			   u2.name, COALESCE(u2.avatar_url, ''),
			   q.title, q.category_id
		FROM challenges c
		JOIN users u1 ON c.challenger_id = u1.id
		JOIN users u2 ON c.challengee_id = u2.id
		JOIN quizzes q ON c.quiz_id = q.id
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argIndex, argIndex+1)

	offset := (filters.Page - 1) * filters.PageSize
	args = append(args, filters.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user challenges: %w", err)
	}
	defer rows.Close()

	var challenges []models.ChallengeWithDetails
	for rows.Next() {
		var challenge models.ChallengeWithDetails
		var participantIDs pq.StringArray
		var participantScoresJSON sql.NullString

		err := rows.Scan(
			&challenge.ID,
			&challenge.ChallengerID,
			&challenge.ChallengeeID,
			&challenge.QuizID,
			&challenge.Status,
			&challenge.ChallengerScore,
			&challenge.ChallengeeScore,
			&challenge.Message,
			&challenge.ExpiresAt,
			&challenge.IsGroupChallenge,
			&participantIDs,
			&participantScoresJSON,
			&challenge.CreatedAt,
			&challenge.UpdatedAt,
			&challenge.IsTestData,
			&challenge.ChallengerName,
			&challenge.ChallengerAvatar,
			&challenge.ChallengeeName,
			&challenge.ChallengeeAvatar,
			&challenge.QuizTitle,
			&challenge.QuizCategory,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan challenge row: %w", err)
		}

		// Convert string array to UUID array
		for _, idStr := range participantIDs {
			if id, err := uuid.Parse(idStr); err == nil {
				challenge.ParticipantIDs = append(challenge.ParticipantIDs, id)
			}
		}

		// Parse participant scores JSON
		if participantScoresJSON.Valid {
			err := json.Unmarshal([]byte(participantScoresJSON.String), &challenge.ParticipantScores)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal participant scores: %w", err)
			}
		}

		challenges = append(challenges, challenge)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating challenge rows: %w", err)
	}

	return challenges, total, nil
}

// GetPendingChallenges retrieves pending challenges for a user
func (r *ChallengesRepository) GetPendingChallenges(userID uuid.UUID) ([]models.ChallengeWithDetails, error) {
	filters := &models.ChallengeFilters{
		Status:   "pending",
		Page:     1,
		PageSize: 100, // Get all pending challenges
	}

	challenges, _, err := r.GetUserChallenges(userID, filters)
	return challenges, err
}

// GetActiveChallenges retrieves active challenges for a user
func (r *ChallengesRepository) GetActiveChallenges(userID uuid.UUID) ([]models.ChallengeWithDetails, error) {
	filters := &models.ChallengeFilters{
		Status:   "accepted",
		Page:     1,
		PageSize: 100, // Get all active challenges
	}

	challenges, _, err := r.GetUserChallenges(userID, filters)
	return challenges, err
}

// GetCompletedChallenges retrieves completed challenges for a user
func (r *ChallengesRepository) GetCompletedChallenges(userID uuid.UUID) ([]models.ChallengeWithDetails, error) {
	filters := &models.ChallengeFilters{
		Status:   "completed",
		Page:     1,
		PageSize: 100, // Get all completed challenges
	}

	challenges, _, err := r.GetUserChallenges(userID, filters)
	return challenges, err
}

// AcceptChallenge accepts a challenge
func (r *ChallengesRepository) AcceptChallenge(challengeID uuid.UUID, userID uuid.UUID) error {
	log.Printf("AcceptChallenge called: challengeID=%s, userID=%s", challengeID, userID)

	// Verify user is the challenged user and challenge is pending
	query := `
		UPDATE challenges
		SET status = 'accepted', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND challengee_id = $2 AND status = 'pending'
	`

	result, err := r.db.Exec(query, challengeID, userID)
	if err != nil {
		return fmt.Errorf("failed to accept challenge: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("challenge not found or cannot be accepted")
	}

	return nil
}

// DeclineChallenge declines a challenge
func (r *ChallengesRepository) DeclineChallenge(challengeID uuid.UUID, userID uuid.UUID) error {
	log.Printf("DeclineChallenge called: challengeID=%s, userID=%s", challengeID, userID)

	// Verify user is the challenged user and challenge is pending
	query := `
		UPDATE challenges
		SET status = 'declined', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND challengee_id = $2 AND status = 'pending'
	`

	result, err := r.db.Exec(query, challengeID, userID)
	if err != nil {
		return fmt.Errorf("failed to decline challenge: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("challenge not found or cannot be declined")
	}

	return nil
}

// CompleteChallenge marks a challenge as completed
func (r *ChallengesRepository) CompleteChallenge(challengeID uuid.UUID) error {
	return r.UpdateChallengeStatus(challengeID, "completed")
}

// GetChallengeStats retrieves challenge statistics for a user
func (r *ChallengesRepository) GetChallengeStats(userID uuid.UUID) (*models.ChallengeStatsResponse, error) {
	log.Printf("GetChallengeStats called: userID=%s", userID)

	query := `
		SELECT
			COUNT(*) as total_challenges,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_challenges,
			COUNT(CASE WHEN status = 'accepted' THEN 1 END) as active_challenges,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_challenges,
			COUNT(CASE
				WHEN status = 'completed' AND (
					(challenger_id = $1 AND challenger_score > challengee_score) OR
					(challengee_id = $1 AND challengee_score > challenger_score)
				) THEN 1
			END) as won_challenges,
			COUNT(CASE
				WHEN status = 'completed' AND (
					(challenger_id = $1 AND challenger_score < challengee_score) OR
					(challengee_id = $1 AND challengee_score < challenger_score)
				) THEN 1
			END) as lost_challenges,
			AVG(CASE
				WHEN challenger_id = $1 THEN challenger_score
				WHEN challengee_id = $1 THEN challengee_score
			END) as average_score
		FROM challenges
		WHERE challenger_id = $1 OR challengee_id = $1
	`

	var stats models.ChallengeStatsResponse
	var avgScore sql.NullFloat64

	err := r.db.QueryRow(query, userID).Scan(
		&stats.TotalChallenges,
		&stats.PendingChallenges,
		&stats.ActiveChallenges,
		&stats.CompletedChallenges,
		&stats.WonChallenges,
		&stats.LostChallenges,
		&avgScore,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get challenge stats: %w", err)
	}

	if avgScore.Valid {
		stats.AverageScore = avgScore.Float64
	}

	// Calculate win rate
	if stats.CompletedChallenges > 0 {
		stats.WinRate = float64(stats.WonChallenges) / float64(stats.CompletedChallenges) * 100
	}

	return &stats, nil
}

// CanUserChallenge checks if a user can challenge another user
func (r *ChallengesRepository) CanUserChallenge(challengerID, challengedID uuid.UUID) (bool, error) {
	// Check if users are friends (assuming friendship is required for challenges)
	query := `
		SELECT COUNT(*)
		FROM friendships
		WHERE (user1_id = $1 AND user2_id = $2) OR (user1_id = $2 AND user2_id = $1)
	`

	var count int
	err := r.db.QueryRow(query, challengerID, challengedID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check friendship: %w", err)
	}

	return count > 0, nil
}

// HasPendingChallenge checks if there's already a pending challenge between users for a quiz
func (r *ChallengesRepository) HasPendingChallenge(challengerID, challengedID uuid.UUID, quizID uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM challenges
		WHERE challenger_id = $1 AND challengee_id = $2 AND quiz_id = $3 AND status = 'pending'
	`

	var count int
	err := r.db.QueryRow(query, challengerID, challengedID, quizID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check pending challenge: %w", err)
	}

	return count > 0, nil
}

// ExpireChallenges marks expired challenges as expired
func (r *ChallengesRepository) ExpireChallenges() error {
	log.Println("ExpireChallenges called")

	query := `
		UPDATE challenges
		SET status = 'expired', updated_at = CURRENT_TIMESTAMP
		WHERE expires_at < CURRENT_TIMESTAMP
		AND status IN ('pending', 'accepted')
	`

	result, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to expire challenges: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	log.Printf("Expired %d challenges", rowsAffected)
	return nil
}

// LinkAttemptToChallenge links a quiz attempt to a challenge
// This is called when a user starts taking a quiz for a challenge
func (r *ChallengesRepository) LinkAttemptToChallenge(challengeID uuid.UUID, attemptID uuid.UUID, userID uuid.UUID) error {
	log.Printf("LinkAttemptToChallenge called: challengeID=%s, attemptID=%s, userID=%s", challengeID, attemptID, userID)

	// First get the challenge to determine which attempt field to update
	challenge, err := r.GetChallengeByID(challengeID)
	if err != nil {
		return err
	}

	// Determine which attempt ID field to update based on user role
	var query string
	if challenge.ChallengerID == userID {
		query = `
			UPDATE challenges
			SET challenger_attempt_id = $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
	} else if challenge.ChallengeeID == userID {
		query = `
			UPDATE challenges
			SET challengee_attempt_id = $2, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
	} else {
		return fmt.Errorf("user is not part of this challenge")
	}

	// Execute the update
	result, err := r.db.Exec(query, challengeID, attemptID)
	if err != nil {
		return fmt.Errorf("failed to link attempt to challenge: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("challenge not found")
	}

	// Also update the quiz_attempt record to mark it as a challenge attempt
	updateAttemptQuery := `
		UPDATE quiz_attempts
		SET challenge_id = $2, is_challenge_attempt = true, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err = r.db.Exec(updateAttemptQuery, attemptID, challengeID)
	if err != nil {
		return fmt.Errorf("failed to update quiz attempt: %w", err)
	}

	log.Printf("Successfully linked attempt %s to challenge %s", attemptID, challengeID)
	return nil
}

// CompleteChallengeAttempt marks a user's challenge attempt as complete and updates the challenge status
// This method handles the asynchronous nature of challenges - it updates the appropriate completion timestamp
// and transitions the challenge status based on whether one or both users have completed
func (r *ChallengesRepository) CompleteChallengeAttempt(challengeID uuid.UUID, userID uuid.UUID, score float64) error {
	log.Printf("CompleteChallengeAttempt called: challengeID=%s, userID=%s, score=%f", challengeID, userID, score)

	// First get the challenge to determine which user is completing
	challenge, err := r.GetChallengeByID(challengeID)
	if err != nil {
		return err
	}

	// Determine query based on user role and check opponent's completion status
	var query string
	var newStatus string

	if challenge.ChallengerID == userID {
		// Challenger is completing
		query = `
			UPDATE challenges
			SET challenger_score = $2,
			    challenger_completed_at = CURRENT_TIMESTAMP,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
			RETURNING challengee_completed_at
		`
	} else if challenge.ChallengeeID == userID {
		// Challengee is completing
		query = `
			UPDATE challenges
			SET challengee_score = $2,
			    challengee_completed_at = CURRENT_TIMESTAMP,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
			RETURNING challenger_completed_at
		`
	} else {
		return fmt.Errorf("user is not part of this challenge")
	}

	// Execute the update and get opponent's completion status
	var opponentCompletedAt *time.Time
	err = r.db.QueryRow(query, challengeID, score).Scan(&opponentCompletedAt)
	if err != nil {
		return fmt.Errorf("failed to update challenge completion: %w", err)
	}

	// Determine the new status based on completion state
	if opponentCompletedAt != nil {
		// Both users have completed
		newStatus = "completed"
	} else {
		// Only this user has completed
		if challenge.ChallengerID == userID {
			newStatus = "challenger_completed"
		} else {
			newStatus = "challengee_completed"
		}
	}

	// Update the challenge status
	err = r.UpdateChallengeStatus(challengeID, newStatus)
	if err != nil {
		return fmt.Errorf("failed to update challenge status: %w", err)
	}

	log.Printf("Successfully completed challenge attempt. New status: %s", newStatus)
	return nil
}
