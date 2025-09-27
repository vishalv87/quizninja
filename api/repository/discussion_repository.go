package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"quizninja-api/database"
	"quizninja-api/models"

	"github.com/google/uuid"
)

type DiscussionRepository struct {
	db *sql.DB
}

// DiscussionRepositoryInterface defines the contract for discussion data operations
type DiscussionRepositoryInterface interface {
	// Discussion CRUD operations
	CreateDiscussion(discussion *models.Discussion) error
	GetDiscussionByID(id uuid.UUID) (*models.Discussion, error)
	GetDiscussionWithDetails(id uuid.UUID, userID *uuid.UUID) (*models.DiscussionWithDetails, error)
	UpdateDiscussion(discussion *models.Discussion) error
	DeleteDiscussion(id uuid.UUID, userID uuid.UUID) error

	// Discussion list operations
	GetDiscussions(filters *models.DiscussionFilters, userID *uuid.UUID) ([]models.Discussion, int, error)
	GetDiscussionsByQuiz(quizID uuid.UUID, userID *uuid.UUID, limit, offset int) ([]models.Discussion, int, error)
	GetDiscussionsByQuestion(questionID uuid.UUID, userID *uuid.UUID, limit, offset int) ([]models.Discussion, int, error)

	// Discussion reply operations
	CreateDiscussionReply(reply *models.DiscussionReply) error
	GetDiscussionReplies(discussionID uuid.UUID, userID *uuid.UUID, limit, offset int) ([]models.DiscussionReply, int, error)
	GetReplyByID(id uuid.UUID) (*models.DiscussionReply, error)
	UpdateDiscussionReply(reply *models.DiscussionReply) error
	DeleteDiscussionReply(id uuid.UUID, userID uuid.UUID) error

	// Like operations
	LikeDiscussion(discussionID uuid.UUID, userID uuid.UUID) error
	UnlikeDiscussion(discussionID uuid.UUID, userID uuid.UUID) error
	IsDiscussionLikedByUser(discussionID uuid.UUID, userID uuid.UUID) (bool, error)

	LikeDiscussionReply(replyID uuid.UUID, userID uuid.UUID) error
	UnlikeDiscussionReply(replyID uuid.UUID, userID uuid.UUID) error
	IsReplyLikedByUser(replyID uuid.UUID, userID uuid.UUID) (bool, error)

	// Count operations
	UpdateDiscussionLikesCount(discussionID uuid.UUID) error
	UpdateDiscussionRepliesCount(discussionID uuid.UUID) error
	UpdateReplyLikesCount(replyID uuid.UUID) error

	// Statistics operations
	GetDiscussionStats(userID *uuid.UUID) (*models.DiscussionStatsResponse, error)
}

// NewDiscussionRepository creates a new discussion repository instance
func NewDiscussionRepository() DiscussionRepositoryInterface {
	return &DiscussionRepository{
		db: database.DB,
	}
}

// CreateDiscussion creates a new discussion
func (r *DiscussionRepository) CreateDiscussion(discussion *models.Discussion) error {
	log.Printf("CreateDiscussion called: quizID=%s, userID=%s", discussion.QuizID, discussion.UserID)

	query := `
		INSERT INTO discussions (quiz_id, question_id, user_id, content, type, is_test_data)
		VALUES ($1, $2, $3, $4, $5, true)
		RETURNING id, created_at, updated_at, is_test_data
	`

	err := r.db.QueryRow(query,
		discussion.QuizID,
		discussion.QuestionID,
		discussion.UserID,
		discussion.Content,
		discussion.Type,
	).Scan(
		&discussion.ID,
		&discussion.CreatedAt,
		&discussion.UpdatedAt,
		&discussion.IsTestData,
	)

	if err != nil {
		return fmt.Errorf("failed to create discussion: %w", err)
	}

	return nil
}

// GetDiscussionByID retrieves a discussion by ID
func (r *DiscussionRepository) GetDiscussionByID(id uuid.UUID) (*models.Discussion, error) {
	log.Printf("GetDiscussionByID called: id=%s", id)

	query := `
		SELECT id, quiz_id, question_id, user_id, content, likes_count, replies_count, type, created_at, updated_at, is_test_data
		FROM discussions
		WHERE id = $1
	`

	var discussion models.Discussion
	err := r.db.QueryRow(query, id).Scan(
		&discussion.ID,
		&discussion.QuizID,
		&discussion.QuestionID,
		&discussion.UserID,
		&discussion.Content,
		&discussion.LikesCount,
		&discussion.RepliesCount,
		&discussion.Type,
		&discussion.CreatedAt,
		&discussion.UpdatedAt,
		&discussion.IsTestData,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("discussion not found")
		}
		return nil, fmt.Errorf("failed to get discussion: %w", err)
	}

	return &discussion, nil
}

// GetDiscussionWithDetails retrieves a discussion with full user and quiz details
func (r *DiscussionRepository) GetDiscussionWithDetails(id uuid.UUID, userID *uuid.UUID) (*models.DiscussionWithDetails, error) {
	log.Printf("GetDiscussionWithDetails called: id=%s, userID=%v", id, userID)

	query := `
		SELECT
			d.id, d.quiz_id, d.question_id, d.user_id, d.content, d.likes_count, d.replies_count,
			d.type, d.created_at, d.updated_at, d.is_test_data,
			u.name as user_name, u.avatar_url as user_avatar, u.is_test_data as user_is_test_data,
			q.title as quiz_title, q.category_id as quiz_category, q.is_test_data as quiz_is_test_data,
			qs.question_text, qs.is_test_data as question_is_test_data
		FROM discussions d
		JOIN users u ON d.user_id = u.id
		JOIN quizzes q ON d.quiz_id = q.id
		LEFT JOIN questions qs ON d.question_id = qs.id
		WHERE d.id = $1
	`

	var discussion models.DiscussionWithDetails
	var userAvatar sql.NullString
	var questionText sql.NullString
	var userIsTestData bool
	var quizIsTestData bool
	var questionIsTestData sql.NullBool

	err := r.db.QueryRow(query, id).Scan(
		&discussion.ID,
		&discussion.QuizID,
		&discussion.QuestionID,
		&discussion.UserID,
		&discussion.Content,
		&discussion.LikesCount,
		&discussion.RepliesCount,
		&discussion.Type,
		&discussion.CreatedAt,
		&discussion.UpdatedAt,
		&discussion.IsTestData,
		&discussion.UserName,
		&userAvatar,
		&userIsTestData,
		&discussion.QuizTitle,
		&discussion.QuizCategory,
		&quizIsTestData,
		&questionText,
		&questionIsTestData,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("discussion not found")
		}
		return nil, fmt.Errorf("failed to get discussion with details: %w", err)
	}

	if userAvatar.Valid {
		discussion.UserAvatar = userAvatar.String
	}
	if questionText.Valid {
		discussion.QuestionText = questionText.String
	}

	// Check if user has liked this discussion
	if userID != nil {
		isLiked, err := r.IsDiscussionLikedByUser(id, *userID)
		if err != nil {
			log.Printf("Warning: failed to check if discussion is liked by user: %v", err)
		}
		discussion.IsLikedByUser = isLiked
	}

	return &discussion, nil
}

// UpdateDiscussion updates a discussion
func (r *DiscussionRepository) UpdateDiscussion(discussion *models.Discussion) error {
	log.Printf("UpdateDiscussion called: id=%s", discussion.ID)

	query := `
		UPDATE discussions
		SET content = $1, type = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(query,
		discussion.Content,
		discussion.Type,
		discussion.ID,
		discussion.UserID,
	).Scan(&discussion.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("discussion not found or unauthorized")
		}
		return fmt.Errorf("failed to update discussion: %w", err)
	}

	return nil
}

// DeleteDiscussion deletes a discussion
func (r *DiscussionRepository) DeleteDiscussion(id uuid.UUID, userID uuid.UUID) error {
	log.Printf("DeleteDiscussion called: id=%s, userID=%s", id, userID)

	query := `DELETE FROM discussions WHERE id = $1 AND user_id = $2`

	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete discussion: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("discussion not found or unauthorized")
	}

	return nil
}

// GetDiscussions retrieves discussions with filters
func (r *DiscussionRepository) GetDiscussions(filters *models.DiscussionFilters, userID *uuid.UUID) ([]models.Discussion, int, error) {
	log.Printf("GetDiscussions called with filters: %+v", filters)

	var whereConditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE clause based on filters
	if filters.QuizID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("d.quiz_id = $%d", argIndex))
		args = append(args, *filters.QuizID)
		argIndex++
	}

	if filters.QuestionID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("d.question_id = $%d", argIndex))
		args = append(args, *filters.QuestionID)
		argIndex++
	}

	if filters.UserID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("d.user_id = $%d", argIndex))
		args = append(args, *filters.UserID)
		argIndex++
	}

	if filters.Type != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("d.type = $%d", argIndex))
		args = append(args, filters.Type)
		argIndex++
	}

	if filters.Search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("d.content ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total rows
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM discussions d
		%s
	`, whereClause)

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count discussions: %w", err)
	}

	// Build sort clause
	sortClause := fmt.Sprintf("ORDER BY d.%s %s", filters.SortBy, strings.ToUpper(filters.SortOrder))

	// Build main query
	mainQuery := fmt.Sprintf(`
		SELECT
			d.id, d.quiz_id, d.question_id, d.user_id, d.content, d.likes_count, d.replies_count,
			d.type, d.created_at, d.updated_at, d.is_test_data,
			u.name as user_name, u.avatar_url as user_avatar, u.is_test_data as user_is_test_data
		FROM discussions d
		JOIN users u ON d.user_id = u.id
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortClause, argIndex, argIndex+1)

	// Add pagination args
	args = append(args, filters.PageSize, (filters.Page-1)*filters.PageSize)

	rows, err := r.db.Query(mainQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query discussions: %w", err)
	}
	defer rows.Close()

	discussions := make([]models.Discussion, 0)
	for rows.Next() {
		var discussion models.Discussion
		var user models.User
		var avatarURL sql.NullString

		err := rows.Scan(
			&discussion.ID,
			&discussion.QuizID,
			&discussion.QuestionID,
			&discussion.UserID,
			&discussion.Content,
			&discussion.LikesCount,
			&discussion.RepliesCount,
			&discussion.Type,
			&discussion.CreatedAt,
			&discussion.UpdatedAt,
			&discussion.IsTestData,
			&user.Name,
			&avatarURL,
			&user.IsTestData,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan discussion: %w", err)
		}

		user.ID = discussion.UserID
		if avatarURL.Valid {
			user.AvatarURL = &avatarURL.String
		}
		discussion.User = &user

		// Check if user has liked this discussion
		if userID != nil {
			isLiked, err := r.IsDiscussionLikedByUser(discussion.ID, *userID)
			if err != nil {
				log.Printf("Warning: failed to check if discussion is liked by user: %v", err)
			}
			discussion.IsLikedByUser = isLiked
		}

		discussions = append(discussions, discussion)
	}

	return discussions, total, nil
}

// GetDiscussionsByQuiz retrieves discussions for a specific quiz
func (r *DiscussionRepository) GetDiscussionsByQuiz(quizID uuid.UUID, userID *uuid.UUID, limit, offset int) ([]models.Discussion, int, error) {
	filters := &models.DiscussionFilters{
		QuizID:    &quizID,
		Page:      (offset / limit) + 1,
		PageSize:  limit,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	return r.GetDiscussions(filters, userID)
}

// GetDiscussionsByQuestion retrieves discussions for a specific question
func (r *DiscussionRepository) GetDiscussionsByQuestion(questionID uuid.UUID, userID *uuid.UUID, limit, offset int) ([]models.Discussion, int, error) {
	filters := &models.DiscussionFilters{
		QuestionID: &questionID,
		Page:       (offset / limit) + 1,
		PageSize:   limit,
		SortBy:     "created_at",
		SortOrder:  "desc",
	}
	return r.GetDiscussions(filters, userID)
}

// CreateDiscussionReply creates a new reply to a discussion
func (r *DiscussionRepository) CreateDiscussionReply(reply *models.DiscussionReply) error {
	log.Printf("CreateDiscussionReply called: discussionID=%s, userID=%s", reply.DiscussionID, reply.UserID)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert the reply
	query := `
		INSERT INTO discussion_replies (discussion_id, user_id, content, is_test_data)
		VALUES ($1, $2, $3, true)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRow(query,
		reply.DiscussionID,
		reply.UserID,
		reply.Content,
	).Scan(
		&reply.ID,
		&reply.CreatedAt,
		&reply.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create discussion reply: %w", err)
	}

	// Update replies count in discussion
	err = r.updateDiscussionRepliesCountTx(tx, reply.DiscussionID)
	if err != nil {
		return fmt.Errorf("failed to update discussion replies count: %w", err)
	}

	return tx.Commit()
}

// GetDiscussionReplies retrieves replies for a discussion
func (r *DiscussionRepository) GetDiscussionReplies(discussionID uuid.UUID, userID *uuid.UUID, limit, offset int) ([]models.DiscussionReply, int, error) {
	log.Printf("GetDiscussionReplies called: discussionID=%s", discussionID)

	// Count total replies
	countQuery := `SELECT COUNT(*) FROM discussion_replies WHERE discussion_id = $1`
	var total int
	err := r.db.QueryRow(countQuery, discussionID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count discussion replies: %w", err)
	}

	// Get replies with user details
	query := `
		SELECT
			r.id, r.discussion_id, r.user_id, r.content, r.likes_count, r.created_at, r.updated_at, r.is_test_data,
			u.name as user_name, u.avatar_url as user_avatar, u.is_test_data as user_is_test_data
		FROM discussion_replies r
		JOIN users u ON r.user_id = u.id
		WHERE r.discussion_id = $1
		ORDER BY r.created_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, discussionID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query discussion replies: %w", err)
	}
	defer rows.Close()

	var replies []models.DiscussionReply
	for rows.Next() {
		var reply models.DiscussionReply
		var user models.User
		var avatarURL sql.NullString

		err := rows.Scan(
			&reply.ID,
			&reply.DiscussionID,
			&reply.UserID,
			&reply.Content,
			&reply.LikesCount,
			&reply.CreatedAt,
			&reply.UpdatedAt,
			&reply.IsTestData,
			&user.Name,
			&avatarURL,
			&user.IsTestData,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan discussion reply: %w", err)
		}

		user.ID = reply.UserID
		if avatarURL.Valid {
			user.AvatarURL = &avatarURL.String
		}
		reply.User = &user

		// Check if user has liked this reply
		if userID != nil {
			isLiked, err := r.IsReplyLikedByUser(reply.ID, *userID)
			if err != nil {
				log.Printf("Warning: failed to check if reply is liked by user: %v", err)
			}
			reply.IsLikedByUser = isLiked
		}

		replies = append(replies, reply)
	}

	return replies, total, nil
}

// GetReplyByID retrieves a reply by ID
func (r *DiscussionRepository) GetReplyByID(id uuid.UUID) (*models.DiscussionReply, error) {
	log.Printf("GetReplyByID called: id=%s", id)

	query := `
		SELECT id, discussion_id, user_id, content, likes_count, created_at, updated_at, is_test_data
		FROM discussion_replies
		WHERE id = $1
	`

	var reply models.DiscussionReply
	err := r.db.QueryRow(query, id).Scan(
		&reply.ID,
		&reply.DiscussionID,
		&reply.UserID,
		&reply.Content,
		&reply.LikesCount,
		&reply.CreatedAt,
		&reply.UpdatedAt,
		&reply.IsTestData,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reply not found")
		}
		return nil, fmt.Errorf("failed to get reply: %w", err)
	}

	return &reply, nil
}

// UpdateDiscussionReply updates a reply
func (r *DiscussionRepository) UpdateDiscussionReply(reply *models.DiscussionReply) error {
	log.Printf("UpdateDiscussionReply called: id=%s", reply.ID)

	query := `
		UPDATE discussion_replies
		SET content = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND user_id = $3
		RETURNING updated_at
	`

	err := r.db.QueryRow(query,
		reply.Content,
		reply.ID,
		reply.UserID,
	).Scan(&reply.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("reply not found or unauthorized")
		}
		return fmt.Errorf("failed to update reply: %w", err)
	}

	return nil
}

// DeleteDiscussionReply deletes a reply
func (r *DiscussionRepository) DeleteDiscussionReply(id uuid.UUID, userID uuid.UUID) error {
	log.Printf("DeleteDiscussionReply called: id=%s, userID=%s", id, userID)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the discussion ID first
	var discussionID uuid.UUID
	err = tx.QueryRow("SELECT discussion_id FROM discussion_replies WHERE id = $1 AND user_id = $2", id, userID).Scan(&discussionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("reply not found or unauthorized")
		}
		return fmt.Errorf("failed to get discussion ID: %w", err)
	}

	// Delete the reply
	result, err := tx.Exec("DELETE FROM discussion_replies WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete reply: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reply not found or unauthorized")
	}

	// Update replies count in discussion
	err = r.updateDiscussionRepliesCountTx(tx, discussionID)
	if err != nil {
		return fmt.Errorf("failed to update discussion replies count: %w", err)
	}

	return tx.Commit()
}

// LikeDiscussion likes a discussion
func (r *DiscussionRepository) LikeDiscussion(discussionID uuid.UUID, userID uuid.UUID) error {
	log.Printf("LikeDiscussion called: discussionID=%s, userID=%s", discussionID, userID)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert like (ignore if already exists)
	query := `
		INSERT INTO discussion_likes (discussion_id, user_id, is_test_data)
		VALUES ($1, $2, true)
		ON CONFLICT (discussion_id, user_id) DO NOTHING
	`

	_, err = tx.Exec(query, discussionID, userID)
	if err != nil {
		return fmt.Errorf("failed to like discussion: %w", err)
	}

	// Update likes count
	err = r.updateDiscussionLikesCountTx(tx, discussionID)
	if err != nil {
		return fmt.Errorf("failed to update discussion likes count: %w", err)
	}

	return tx.Commit()
}

// UnlikeDiscussion unlikes a discussion
func (r *DiscussionRepository) UnlikeDiscussion(discussionID uuid.UUID, userID uuid.UUID) error {
	log.Printf("UnlikeDiscussion called: discussionID=%s, userID=%s", discussionID, userID)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete like
	query := `DELETE FROM discussion_likes WHERE discussion_id = $1 AND user_id = $2`

	_, err = tx.Exec(query, discussionID, userID)
	if err != nil {
		return fmt.Errorf("failed to unlike discussion: %w", err)
	}

	// Update likes count
	err = r.updateDiscussionLikesCountTx(tx, discussionID)
	if err != nil {
		return fmt.Errorf("failed to update discussion likes count: %w", err)
	}

	return tx.Commit()
}

// IsDiscussionLikedByUser checks if a user has liked a discussion
func (r *DiscussionRepository) IsDiscussionLikedByUser(discussionID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM discussion_likes WHERE discussion_id = $1 AND user_id = $2)`

	var exists bool
	err := r.db.QueryRow(query, discussionID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if discussion is liked: %w", err)
	}

	return exists, nil
}

// LikeDiscussionReply likes a reply
func (r *DiscussionRepository) LikeDiscussionReply(replyID uuid.UUID, userID uuid.UUID) error {
	log.Printf("LikeDiscussionReply called: replyID=%s, userID=%s", replyID, userID)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert like (ignore if already exists)
	query := `
		INSERT INTO discussion_reply_likes (reply_id, user_id, is_test_data)
		VALUES ($1, $2, true)
		ON CONFLICT (reply_id, user_id) DO NOTHING
	`

	_, err = tx.Exec(query, replyID, userID)
	if err != nil {
		return fmt.Errorf("failed to like reply: %w", err)
	}

	// Update likes count
	err = r.updateReplyLikesCountTx(tx, replyID)
	if err != nil {
		return fmt.Errorf("failed to update reply likes count: %w", err)
	}

	return tx.Commit()
}

// UnlikeDiscussionReply unlikes a reply
func (r *DiscussionRepository) UnlikeDiscussionReply(replyID uuid.UUID, userID uuid.UUID) error {
	log.Printf("UnlikeDiscussionReply called: replyID=%s, userID=%s", replyID, userID)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete like
	query := `DELETE FROM discussion_reply_likes WHERE reply_id = $1 AND user_id = $2`

	_, err = tx.Exec(query, replyID, userID)
	if err != nil {
		return fmt.Errorf("failed to unlike reply: %w", err)
	}

	// Update likes count
	err = r.updateReplyLikesCountTx(tx, replyID)
	if err != nil {
		return fmt.Errorf("failed to update reply likes count: %w", err)
	}

	return tx.Commit()
}

// IsReplyLikedByUser checks if a user has liked a reply
func (r *DiscussionRepository) IsReplyLikedByUser(replyID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM discussion_reply_likes WHERE reply_id = $1 AND user_id = $2)`

	var exists bool
	err := r.db.QueryRow(query, replyID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if reply is liked: %w", err)
	}

	return exists, nil
}

// Helper methods for updating counts

func (r *DiscussionRepository) UpdateDiscussionLikesCount(discussionID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = r.updateDiscussionLikesCountTx(tx, discussionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *DiscussionRepository) UpdateDiscussionRepliesCount(discussionID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = r.updateDiscussionRepliesCountTx(tx, discussionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *DiscussionRepository) UpdateReplyLikesCount(replyID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = r.updateReplyLikesCountTx(tx, replyID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *DiscussionRepository) updateDiscussionLikesCountTx(tx *sql.Tx, discussionID uuid.UUID) error {
	query := `
		UPDATE discussions
		SET likes_count = (
			SELECT COUNT(*)
			FROM discussion_likes
			WHERE discussion_id = $1
		)
		WHERE id = $1
	`

	_, err := tx.Exec(query, discussionID)
	return err
}

func (r *DiscussionRepository) updateDiscussionRepliesCountTx(tx *sql.Tx, discussionID uuid.UUID) error {
	query := `
		UPDATE discussions
		SET replies_count = (
			SELECT COUNT(*)
			FROM discussion_replies
			WHERE discussion_id = $1
		)
		WHERE id = $1
	`

	_, err := tx.Exec(query, discussionID)
	return err
}

func (r *DiscussionRepository) updateReplyLikesCountTx(tx *sql.Tx, replyID uuid.UUID) error {
	query := `
		UPDATE discussion_replies
		SET likes_count = (
			SELECT COUNT(*)
			FROM discussion_reply_likes
			WHERE reply_id = $1
		)
		WHERE id = $1
	`

	_, err := tx.Exec(query, replyID)
	return err
}

// GetDiscussionStats retrieves discussion statistics
func (r *DiscussionRepository) GetDiscussionStats(userID *uuid.UUID) (*models.DiscussionStatsResponse, error) {
	log.Printf("GetDiscussionStats called for userID: %v", userID)

	var whereClause string
	var args []interface{}

	if userID != nil {
		whereClause = "WHERE d.user_id = $1"
		args = append(args, *userID)
	}

	query := fmt.Sprintf(`
		SELECT
			COUNT(d.id) as total_discussions,
			COALESCE(SUM(d.replies_count), 0) as total_replies,
			COALESCE(SUM(d.likes_count), 0) as total_likes,
			COALESCE(AVG(d.replies_count), 0) as avg_replies_per_post,
			COALESCE(AVG(d.likes_count), 0) as avg_likes_per_post
		FROM discussions d
		%s
	`, whereClause)

	var stats models.DiscussionStatsResponse
	err := r.db.QueryRow(query, args...).Scan(
		&stats.TotalDiscussions,
		&stats.TotalReplies,
		&stats.TotalLikes,
		&stats.AverageRepliesPerPost,
		&stats.AverageLikesPerPost,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get discussion stats: %w", err)
	}

	// Get most active quiz (overall, not user-specific for now)
	mostActiveQuery := `
		SELECT
			d.quiz_id,
			q.title,
			COUNT(d.id) as discussion_count
		FROM discussions d
		JOIN quizzes q ON d.quiz_id = q.id
		GROUP BY d.quiz_id, q.title
		ORDER BY discussion_count DESC
		LIMIT 1
	`

	var quizID uuid.UUID
	var quizTitle string
	var discussionCount int

	err = r.db.QueryRow(mostActiveQuery).Scan(&quizID, &quizTitle, &discussionCount)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get most active quiz: %w", err)
	}

	if err == nil {
		stats.MostActiveQuiz = &struct {
			QuizID          uuid.UUID `json:"quiz_id"`
			QuizTitle       string    `json:"quiz_title"`
			DiscussionCount int       `json:"discussion_count"`
		}{
			QuizID:          quizID,
			QuizTitle:       quizTitle,
			DiscussionCount: discussionCount,
		}
	}

	return &stats, nil
}