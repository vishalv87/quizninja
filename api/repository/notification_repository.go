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

type NotificationRepository struct {
	db *sql.DB
}

// NewNotificationRepository creates a new notification repository instance
func NewNotificationRepository() NotificationRepositoryInterface {
	return &NotificationRepository{
		db: database.DB,
	}
}

// CreateNotification creates a new notification
func (r *NotificationRepository) CreateNotification(notification *models.CreateNotificationRequest) (*models.Notification, error) {
	log.Printf("CreateNotification called: userID=%s, type=%s", notification.UserID, notification.Type)

	query := `
		INSERT INTO notifications (
			user_id, type, title, message, data, related_user_id,
			related_entity_id, related_entity_type, expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, type, title, message, data, related_user_id,
				  related_entity_id, related_entity_type, is_read, is_deleted, created_at,
				  read_at, deleted_at, expires_at
	`

	var created models.Notification
	err := r.db.QueryRow(
		query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Data,
		notification.RelatedUserID,
		notification.RelatedEntityID,
		notification.RelatedEntityType,
		notification.ExpiresAt,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.Type,
		&created.Title,
		&created.Message,
		&created.Data,
		&created.RelatedUserID,
		&created.RelatedEntityID,
		&created.RelatedEntityType,
		&created.IsRead,
		&created.IsDeleted,
		&created.CreatedAt,
		&created.ReadAt,
		&created.DeletedAt,
		&created.ExpiresAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	log.Printf("Notification created successfully: id=%s", created.ID)
	return &created, nil
}

// GetNotifications retrieves notifications for a user with pagination and filtering
func (r *NotificationRepository) GetNotifications(userID uuid.UUID, filters *models.NotificationFilters) ([]models.Notification, int, error) {
	log.Printf("GetNotifications called: userID=%s, filters=%+v", userID, filters)

	// Build WHERE clause
	var whereConditions []string
	var args []interface{}
	argCount := 1

	whereConditions = append(whereConditions, fmt.Sprintf("n.user_id = $%d", argCount))
	args = append(args, userID)
	argCount++

	// Always filter out soft-deleted notifications
	whereConditions = append(whereConditions, "n.is_deleted = FALSE")

	// Filter by type
	if filters.Type != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("n.type = $%d", argCount))
		args = append(args, filters.Type)
		argCount++
	}

	// Filter by read status
	if filters.IsRead != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("n.is_read = $%d", argCount))
		args = append(args, *filters.IsRead)
		argCount++
	}

	// Filter by date range
	if filters.StartDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("n.created_at >= $%d", argCount))
		args = append(args, *filters.StartDate)
		argCount++
	}

	if filters.EndDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("n.created_at <= $%d", argCount))
		args = append(args, *filters.EndDate)
		argCount++
	}

	whereClause := "WHERE " + strings.Join(whereConditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM notifications n
		%s
	`, whereClause)

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get notification count: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "n.created_at DESC"
	if filters.SortBy != "" {
		orderBy = fmt.Sprintf("n.%s %s", filters.SortBy, strings.ToUpper(filters.SortOrder))
	}

	// Calculate offset
	offset := (filters.Page - 1) * filters.PageSize

	// Get notifications with user details
	query := fmt.Sprintf(`
		SELECT
			n.id, n.user_id, n.type, n.title, n.message, n.data,
			n.related_user_id, n.related_entity_id, n.related_entity_type,
			n.is_read, n.is_deleted, n.created_at, n.read_at, n.deleted_at, n.expires_at,
			u.id, u.name, u.email, u.avatar_url, u.level, u.total_points,
			u.current_streak, u.best_streak, u.total_quizzes_completed,
			u.average_score, u.is_online, u.last_active, u.created_at, u.updated_at
		FROM notifications n
		LEFT JOIN users u ON n.related_user_id = u.id
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argCount, argCount+1)

	args = append(args, filters.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	notifications := make([]models.Notification, 0)
	for rows.Next() {
		var notification models.Notification
		var relatedUser models.User

		// Use nullable types for user fields that can be NULL from LEFT JOIN
		var relatedUserID, relatedUserName, relatedUserEmail, relatedUserAvatarURL, relatedUserLevel sql.NullString
		var relatedUserTotalPoints, relatedUserCurrentStreak, relatedUserBestStreak, relatedUserTotalQuizzes sql.NullInt64
		var relatedUserAverageScore sql.NullFloat64
		var relatedUserIsOnline sql.NullBool
		var relatedUserLastActive, relatedUserCreatedAt, relatedUserUpdatedAt sql.NullTime

		// Scan notification fields and optional related user
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.Data,
			&notification.RelatedUserID,
			&notification.RelatedEntityID,
			&notification.RelatedEntityType,
			&notification.IsRead,
			&notification.IsDeleted,
			&notification.CreatedAt,
			&notification.ReadAt,
			&notification.DeletedAt,
			&notification.ExpiresAt,
			&relatedUserID,
			&relatedUserName,
			&relatedUserEmail,
			&relatedUserAvatarURL,
			&relatedUserLevel,
			&relatedUserTotalPoints,
			&relatedUserCurrentStreak,
			&relatedUserBestStreak,
			&relatedUserTotalQuizzes,
			&relatedUserAverageScore,
			&relatedUserIsOnline,
			&relatedUserLastActive,
			&relatedUserCreatedAt,
			&relatedUserUpdatedAt,
		)

		if err != nil {
			log.Printf("Error scanning notification row: %v", err)
			continue
		}

		// Set related user if it exists and has valid data
		if notification.RelatedUserID != nil && relatedUserID.Valid {
			userUUID, err := uuid.Parse(relatedUserID.String)
			if err == nil {
				relatedUser.ID = userUUID

				// Populate user fields from nullable values
				if relatedUserName.Valid {
					relatedUser.Name = relatedUserName.String
				}
				if relatedUserEmail.Valid {
					relatedUser.Email = relatedUserEmail.String
				}
				if relatedUserAvatarURL.Valid {
					avatarURL := relatedUserAvatarURL.String
					relatedUser.AvatarURL = &avatarURL
				}
				if relatedUserLevel.Valid {
					relatedUser.Level = relatedUserLevel.String
				}
				if relatedUserTotalPoints.Valid {
					relatedUser.TotalPoints = int(relatedUserTotalPoints.Int64)
				}
				if relatedUserCurrentStreak.Valid {
					relatedUser.CurrentStreak = int(relatedUserCurrentStreak.Int64)
				}
				if relatedUserBestStreak.Valid {
					relatedUser.BestStreak = int(relatedUserBestStreak.Int64)
				}
				if relatedUserTotalQuizzes.Valid {
					relatedUser.TotalQuizzesCompleted = int(relatedUserTotalQuizzes.Int64)
				}
				if relatedUserAverageScore.Valid {
					relatedUser.AverageScore = relatedUserAverageScore.Float64
				}
				if relatedUserIsOnline.Valid {
					relatedUser.IsOnline = relatedUserIsOnline.Bool
				}
				if relatedUserLastActive.Valid {
					relatedUser.LastActive = relatedUserLastActive.Time
				}
				if relatedUserCreatedAt.Valid {
					relatedUser.CreatedAt = relatedUserCreatedAt.Time
				}
				if relatedUserUpdatedAt.Valid {
					relatedUser.UpdatedAt = relatedUserUpdatedAt.Time
				}

				notification.RelatedUser = &relatedUser
			}
		}

		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating notification rows: %w", err)
	}

	log.Printf("Retrieved %d notifications for user %s", len(notifications), userID)
	return notifications, total, nil
}

// GetNotificationByID retrieves a single notification by ID
func (r *NotificationRepository) GetNotificationByID(notificationID uuid.UUID, userID uuid.UUID) (*models.Notification, error) {
	log.Printf("GetNotificationByID called: notificationID=%s, userID=%s", notificationID, userID)

	query := `
		SELECT
			n.id, n.user_id, n.type, n.title, n.message, n.data,
			n.related_user_id, n.related_entity_id, n.related_entity_type,
			n.is_read, n.is_deleted, n.created_at, n.read_at, n.deleted_at, n.expires_at,
			u.id, u.name, u.email, u.avatar_url, u.level, u.total_points,
			u.current_streak, u.best_streak, u.total_quizzes_completed,
			u.average_score, u.is_online, u.last_active, u.created_at, u.updated_at
		FROM notifications n
		LEFT JOIN users u ON n.related_user_id = u.id
		WHERE n.id = $1 AND n.user_id = $2 AND n.is_deleted = FALSE
	`

	var notification models.Notification
	var relatedUser models.User

	// Use nullable types for user fields that can be NULL from LEFT JOIN
	var relatedUserID, relatedUserName, relatedUserEmail, relatedUserAvatarURL, relatedUserLevel sql.NullString
	var relatedUserTotalPoints, relatedUserCurrentStreak, relatedUserBestStreak, relatedUserTotalQuizzes sql.NullInt64
	var relatedUserAverageScore sql.NullFloat64
	var relatedUserIsOnline sql.NullBool
	var relatedUserLastActive, relatedUserCreatedAt, relatedUserUpdatedAt sql.NullTime

	err := r.db.QueryRow(query, notificationID, userID).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&notification.Data,
		&notification.RelatedUserID,
		&notification.RelatedEntityID,
		&notification.RelatedEntityType,
		&notification.IsRead,
		&notification.IsDeleted,
		&notification.CreatedAt,
		&notification.ReadAt,
		&notification.DeletedAt,
		&notification.ExpiresAt,
		&relatedUserID,
		&relatedUserName,
		&relatedUserEmail,
		&relatedUserAvatarURL,
		&relatedUserLevel,
		&relatedUserTotalPoints,
		&relatedUserCurrentStreak,
		&relatedUserBestStreak,
		&relatedUserTotalQuizzes,
		&relatedUserAverageScore,
		&relatedUserIsOnline,
		&relatedUserLastActive,
		&relatedUserCreatedAt,
		&relatedUserUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	// Set related user if it exists and has valid data
	if notification.RelatedUserID != nil && relatedUserID.Valid {
		userUUID, err := uuid.Parse(relatedUserID.String)
		if err == nil {
			relatedUser.ID = userUUID

			// Populate user fields from nullable values
			if relatedUserName.Valid {
				relatedUser.Name = relatedUserName.String
			}
			if relatedUserEmail.Valid {
				relatedUser.Email = relatedUserEmail.String
			}
			if relatedUserAvatarURL.Valid {
				avatarURL := relatedUserAvatarURL.String
				relatedUser.AvatarURL = &avatarURL
			}
			if relatedUserLevel.Valid {
				relatedUser.Level = relatedUserLevel.String
			}
			if relatedUserTotalPoints.Valid {
				relatedUser.TotalPoints = int(relatedUserTotalPoints.Int64)
			}
			if relatedUserCurrentStreak.Valid {
				relatedUser.CurrentStreak = int(relatedUserCurrentStreak.Int64)
			}
			if relatedUserBestStreak.Valid {
				relatedUser.BestStreak = int(relatedUserBestStreak.Int64)
			}
			if relatedUserTotalQuizzes.Valid {
				relatedUser.TotalQuizzesCompleted = int(relatedUserTotalQuizzes.Int64)
			}
			if relatedUserAverageScore.Valid {
				relatedUser.AverageScore = relatedUserAverageScore.Float64
			}
			if relatedUserIsOnline.Valid {
				relatedUser.IsOnline = relatedUserIsOnline.Bool
			}
			if relatedUserLastActive.Valid {
				relatedUser.LastActive = relatedUserLastActive.Time
			}
			if relatedUserCreatedAt.Valid {
				relatedUser.CreatedAt = relatedUserCreatedAt.Time
			}
			if relatedUserUpdatedAt.Valid {
				relatedUser.UpdatedAt = relatedUserUpdatedAt.Time
			}

			notification.RelatedUser = &relatedUser
		}
	}

	log.Printf("Retrieved notification: id=%s, relatedUser=%t", notification.ID, notification.RelatedUser != nil)
	return &notification, nil
}

// MarkNotificationAsRead marks a specific notification as read
func (r *NotificationRepository) MarkNotificationAsRead(notificationID uuid.UUID, userID uuid.UUID) error {
	log.Printf("MarkNotificationAsRead called: notificationID=%s, userID=%s", notificationID, userID)

	query := `
		UPDATE notifications
		SET is_read = true, read_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2 AND is_read = false AND is_deleted = false
	`

	result, err := r.db.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found or already read")
	}

	log.Printf("Notification marked as read: id=%s", notificationID)
	return nil
}

// MarkNotificationAsUnread marks a specific notification as unread
func (r *NotificationRepository) MarkNotificationAsUnread(notificationID uuid.UUID, userID uuid.UUID) error {
	log.Printf("MarkNotificationAsUnread called: notificationID=%s, userID=%s", notificationID, userID)

	query := `
		UPDATE notifications
		SET is_read = false, read_at = NULL
		WHERE id = $1 AND user_id = $2 AND is_read = true AND is_deleted = false
	`

	result, err := r.db.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as unread: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found, not owned by user, or already unread")
	}

	log.Printf("Notification marked as unread: id=%s", notificationID)
	return nil
}

// MarkAllNotificationsAsRead marks all notifications as read for a user
func (r *NotificationRepository) MarkAllNotificationsAsRead(userID uuid.UUID) error {
	log.Printf("MarkAllNotificationsAsRead called: userID=%s", userID)

	query := `
		UPDATE notifications
		SET is_read = true, read_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND is_read = false AND is_deleted = false
	`

	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	log.Printf("Marked %d notifications as read for user %s", rowsAffected, userID)
	return nil
}

// GetUnreadNotificationCount gets the count of unread notifications for a user
func (r *NotificationRepository) GetUnreadNotificationCount(userID uuid.UUID) (int, error) {
	log.Printf("GetUnreadNotificationCount called: userID=%s", userID)

	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = false AND is_deleted = false
	`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread notification count: %w", err)
	}

	log.Printf("Unread notification count for user %s: %d", userID, count)
	return count, nil
}

// DeleteNotification soft deletes a notification (for user's own notifications only)
func (r *NotificationRepository) DeleteNotification(notificationID uuid.UUID, userID uuid.UUID) error {
	log.Printf("DeleteNotification (soft delete) called: notificationID=%s, userID=%s", notificationID, userID)

	query := `
		UPDATE notifications
		SET is_deleted = true, deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2 AND is_deleted = false
	`

	result, err := r.db.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to soft delete notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found or not owned by user")
	}

	log.Printf("Notification soft deleted: id=%s", notificationID)
	return nil
}

// HardDeleteNotification permanently deletes a notification (admin use)
func (r *NotificationRepository) HardDeleteNotification(notificationID uuid.UUID, userID uuid.UUID) error {
	log.Printf("HardDeleteNotification called: notificationID=%s, userID=%s", notificationID, userID)

	query := `
		DELETE FROM notifications
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to hard delete notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found or not owned by user")
	}

	log.Printf("Notification hard deleted: id=%s", notificationID)
	return nil
}

// RestoreNotification restores a soft deleted notification
func (r *NotificationRepository) RestoreNotification(notificationID uuid.UUID, userID uuid.UUID) error {
	log.Printf("RestoreNotification called: notificationID=%s, userID=%s", notificationID, userID)

	query := `
		UPDATE notifications
		SET is_deleted = false, deleted_at = NULL
		WHERE id = $1 AND user_id = $2 AND is_deleted = true
	`

	result, err := r.db.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to restore notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification not found, not owned by user, or not deleted")
	}

	log.Printf("Notification restored: id=%s", notificationID)
	return nil
}

// GetNotificationStats gets statistics about notifications for a user
func (r *NotificationRepository) GetNotificationStats(userID uuid.UUID) (*models.NotificationStatsResponse, error) {
	log.Printf("GetNotificationStats called: userID=%s", userID)

	// Combined query: get counts by type AND total/unread in a single scan
	combinedQuery := `
		SELECT type, COUNT(*) as type_count,
			COUNT(CASE WHEN is_read = false THEN 1 END) as type_unread
		FROM notifications
		WHERE user_id = $1 AND is_deleted = false
		GROUP BY type
	`

	rows, err := r.db.Query(combinedQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification counts: %w", err)
	}
	defer rows.Close()

	var totalCount, unreadCount int
	notificationsByType := make(map[string]int)
	for rows.Next() {
		var notificationType string
		var typeCount, typeUnread int
		if err := rows.Scan(&notificationType, &typeCount, &typeUnread); err != nil {
			continue
		}
		notificationsByType[notificationType] = typeCount
		totalCount += typeCount
		unreadCount += typeUnread
	}

	// Get recent notifications (last 5, excluding soft-deleted)
	recentQuery := `
		SELECT
			id, user_id, type, title, message, data,
			related_user_id, related_entity_id, related_entity_type,
			is_read, is_deleted, created_at, read_at, deleted_at, expires_at
		FROM notifications
		WHERE user_id = $1 AND is_deleted = false
		ORDER BY created_at DESC
		LIMIT 5
	`

	recentRows, err := r.db.Query(recentQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent notifications: %w", err)
	}
	defer recentRows.Close()

	var recentNotifications []models.Notification
	for recentRows.Next() {
		var notification models.Notification
		err := recentRows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.Data,
			&notification.RelatedUserID,
			&notification.RelatedEntityID,
			&notification.RelatedEntityType,
			&notification.IsRead,
			&notification.IsDeleted,
			&notification.CreatedAt,
			&notification.ReadAt,
			&notification.DeletedAt,
			&notification.ExpiresAt,
		)
		if err != nil {
			continue
		}
		recentNotifications = append(recentNotifications, notification)
	}

	// Build type counts
	typeCounts := models.NotificationTypeCounts{
		FriendRequests:      notificationsByType[models.NotificationTypeFriendRequest],
		FriendResponses:     notificationsByType[models.NotificationTypeFriendAccepted] + notificationsByType[models.NotificationTypeFriendRejected],
		Achievements:        notificationsByType[models.NotificationTypeAchievementUnlocked],
		General:             notificationsByType[models.NotificationTypeGeneral],
		SystemAnnouncements: notificationsByType[models.NotificationTypeSystemAnnouncement],
	}

	stats := &models.NotificationStatsResponse{
		TotalNotifications:  totalCount,
		UnreadNotifications: unreadCount,
		NotificationsByType: notificationsByType,
		RecentNotifications: recentNotifications,
		NotificationCounts:  typeCounts,
	}

	log.Printf("Retrieved notification stats for user %s: total=%d, unread=%d", userID, totalCount, unreadCount)
	return stats, nil
}

// CleanupExpiredNotifications removes expired notifications
func (r *NotificationRepository) CleanupExpiredNotifications() error {
	log.Println("CleanupExpiredNotifications called")

	query := `
		DELETE FROM notifications
		WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP
	`

	result, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired notifications: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	log.Printf("Cleaned up %d expired notifications", rowsAffected)
	return nil
}

// GetFriendNotifications retrieves friend notifications in the old format for backward compatibility
func (r *NotificationRepository) GetFriendNotifications(userID uuid.UUID, limit, offset int) ([]models.FriendNotificationCompat, int, error) {
	log.Printf("GetFriendNotifications (compat) called: userID=%s, limit=%d, offset=%d", userID, limit, offset)

	// Get total count of friend notifications
	countQuery := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND type IN ($2, $3, $4)
	`

	var total int
	err := r.db.QueryRow(countQuery, userID, models.NotificationTypeFriendRequest, models.NotificationTypeFriendAccepted, models.NotificationTypeFriendRejected).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get friend notification count: %w", err)
	}

	// Get friend notifications with user details
	query := `
		SELECT
			n.id, n.user_id, n.type, n.title, n.message,
			n.related_user_id, n.related_entity_id, n.is_read, n.created_at, n.read_at,
			u.id, u.name, u.email, u.avatar_url, u.level, u.total_points,
			u.current_streak, u.best_streak, u.total_quizzes_completed,
			u.average_score, u.is_online, u.last_active, u.created_at, u.updated_at
		FROM notifications n
		LEFT JOIN users u ON n.related_user_id = u.id
		WHERE n.user_id = $1 AND n.type IN ($2, $3, $4)
		ORDER BY n.created_at DESC
		LIMIT $5 OFFSET $6
	`

	rows, err := r.db.Query(query, userID, models.NotificationTypeFriendRequest, models.NotificationTypeFriendAccepted, models.NotificationTypeFriendRejected, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get friend notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.FriendNotificationCompat
	for rows.Next() {
		var notification models.FriendNotificationCompat
		var relatedUser models.User
		var relatedUserID sql.NullString

		// Scan notification fields and optional related user
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.RelatedUserID,
			&notification.FriendRequestID, // This maps to related_entity_id
			&notification.IsRead,
			&notification.CreatedAt,
			&notification.ReadAt,
			&relatedUserID,
			&relatedUser.Name,
			&relatedUser.Email,
			&relatedUser.AvatarURL,
			&relatedUser.Level,
			&relatedUser.TotalPoints,
			&relatedUser.CurrentStreak,
			&relatedUser.BestStreak,
			&relatedUser.TotalQuizzesCompleted,
			&relatedUser.AverageScore,
			&relatedUser.IsOnline,
			&relatedUser.LastActive,
			&relatedUser.CreatedAt,
			&relatedUser.UpdatedAt,
		)

		if err != nil {
			log.Printf("Error scanning friend notification row: %v", err)
			continue
		}

		// Set related user if it exists
		if notification.RelatedUserID != nil && relatedUserID.Valid {
			userUUID, err := uuid.Parse(relatedUserID.String)
			if err == nil {
				relatedUser.ID = userUUID
				notification.RelatedUser = &relatedUser
			}
		}

		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating friend notification rows: %w", err)
	}

	log.Printf("Retrieved %d friend notifications for user %s", len(notifications), userID)
	return notifications, total, nil
}

// GetFriendUnreadNotificationCount gets the count of unread friend notifications for a user
func (r *NotificationRepository) GetFriendUnreadNotificationCount(userID uuid.UUID) (int, error) {
	log.Printf("GetFriendUnreadNotificationCount called: userID=%s", userID)

	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = false AND type IN ($2, $3, $4)
	`

	var count int
	err := r.db.QueryRow(query, userID, models.NotificationTypeFriendRequest, models.NotificationTypeFriendAccepted, models.NotificationTypeFriendRejected).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread friend notification count: %w", err)
	}

	log.Printf("Unread friend notification count for user %s: %d", userID, count)
	return count, nil
}
