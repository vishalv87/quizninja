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

type FriendsRepository struct {
	db *sql.DB
}

// NewFriendsRepository creates a new friends repository instance
func NewFriendsRepository() FriendsRepositoryInterface {
	return &FriendsRepository{
		db: database.DB,
	}
}

// SendFriendRequest creates a new friend request
func (r *FriendsRepository) SendFriendRequest(requesterID, requestedID uuid.UUID, message *string) (*models.FriendRequest, error) {
	log.Printf("SendFriendRequest called: requesterID=%s, requestedID=%s", requesterID, requestedID)
	query := `
		INSERT INTO friend_requests (requester_id, requested_id, message)
		VALUES ($1, $2, $3)
		RETURNING id, requester_id, requested_id, status, message, created_at, responded_at
	`

	var friendRequest models.FriendRequest
	err := r.db.QueryRow(query, requesterID, requestedID, message).Scan(
		&friendRequest.ID,
		&friendRequest.RequesterID,
		&friendRequest.RequestedID,
		&friendRequest.Status,
		&friendRequest.Message,
		&friendRequest.CreatedAt,
		&friendRequest.RespondedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to send friend request: %w", err)
	}

	return &friendRequest, nil
}

// GetFriendRequest retrieves a friend request by ID
func (r *FriendsRepository) GetFriendRequest(id uuid.UUID) (*models.FriendRequest, error) {
	log.Printf("GetFriendRequest called: id=%s", id)
	query := `
		SELECT fr.id, fr.requester_id, fr.requested_id, fr.status, fr.message, fr.created_at, fr.responded_at,
			   u1.id, u1.name, u1.email, u1.avatar_url, u1.level, u1.total_points, u1.is_online, u1.last_active,
			   u2.id, u2.name, u2.email, u2.avatar_url, u2.level, u2.total_points, u2.is_online, u2.last_active
		FROM friend_requests fr
		JOIN users u1 ON fr.requester_id = u1.id
		JOIN users u2 ON fr.requested_id = u2.id
		WHERE fr.id = $1
	`

	var friendRequest models.FriendRequest
	var requester models.User
	var requested models.User

	err := r.db.QueryRow(query, id).Scan(
		&friendRequest.ID,
		&friendRequest.RequesterID,
		&friendRequest.RequestedID,
		&friendRequest.Status,
		&friendRequest.Message,
		&friendRequest.CreatedAt,
		&friendRequest.RespondedAt,
		&requester.ID,
		&requester.Name,
		&requester.Email,
		&requester.AvatarURL,
		&requester.Level,
		&requester.TotalPoints,
		&requester.IsOnline,
		&requester.LastActive,
		&requested.ID,
		&requested.Name,
		&requested.Email,
		&requested.AvatarURL,
		&requested.Level,
		&requested.TotalPoints,
		&requested.IsOnline,
		&requested.LastActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("friend request not found")
		}
		return nil, fmt.Errorf("failed to get friend request: %w", err)
	}

	friendRequest.Requester = &requester
	friendRequest.Requested = &requested

	return &friendRequest, nil
}

// GetFriendRequestBetweenUsers retrieves a friend request between two users
func (r *FriendsRepository) GetFriendRequestBetweenUsers(requesterID, requestedID uuid.UUID) (*models.FriendRequest, error) {
	log.Printf("GetFriendRequestBetweenUsers called: requesterID=%s, requestedID=%s", requesterID, requestedID)
	query := `
		SELECT id, requester_id, requested_id, status, message, created_at, responded_at
		FROM friend_requests
		WHERE requester_id = $1 AND requested_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var friendRequest models.FriendRequest
	err := r.db.QueryRow(query, requesterID, requestedID).Scan(
		&friendRequest.ID,
		&friendRequest.RequesterID,
		&friendRequest.RequestedID,
		&friendRequest.Status,
		&friendRequest.Message,
		&friendRequest.CreatedAt,
		&friendRequest.RespondedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No request found, return nil without error
		}
		return nil, fmt.Errorf("failed to get friend request between users: %w", err)
	}

	return &friendRequest, nil
}

// RespondToFriendRequest updates the status of a friend request
func (r *FriendsRepository) RespondToFriendRequest(requestID uuid.UUID, status string) error {
	log.Printf("RespondToFriendRequest called: requestID=%s, status=%s", requestID, status)
	query := `
		UPDATE friend_requests
		SET status = $1
		WHERE id = $2 AND status = 'pending'
	`

	result, err := r.db.Exec(query, status, requestID)
	if err != nil {
		return fmt.Errorf("failed to respond to friend request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("friend request not found or already responded to")
	}

	return nil
}

// CancelFriendRequest cancels a pending friend request
func (r *FriendsRepository) CancelFriendRequest(requestID uuid.UUID, requesterID uuid.UUID) error {
	log.Printf("CancelFriendRequest called: requestID=%s, requesterID=%s", requestID, requesterID)
	query := `
		UPDATE friend_requests
		SET status = 'cancelled'
		WHERE id = $1 AND requester_id = $2 AND status = 'pending'
	`

	result, err := r.db.Exec(query, requestID, requesterID)
	if err != nil {
		return fmt.Errorf("failed to cancel friend request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("friend request not found or cannot be cancelled")
	}

	return nil
}

// GetPendingFriendRequests retrieves pending friend requests for a user
func (r *FriendsRepository) GetPendingFriendRequests(userID uuid.UUID) ([]models.FriendRequest, error) {
	log.Printf("GetPendingFriendRequests called: userID=%s", userID)
	query := `
		SELECT fr.id, fr.requester_id, fr.requested_id, fr.status, fr.message, fr.created_at, fr.responded_at,
			   u.id, u.name, u.email, u.avatar_url, u.level, u.total_points, u.is_online, u.last_active
		FROM friend_requests fr
		JOIN users u ON fr.requester_id = u.id
		WHERE fr.requested_id = $1 AND fr.status = 'pending'
		ORDER BY fr.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending friend requests: %w", err)
	}
	defer rows.Close()

	var requests []models.FriendRequest
	for rows.Next() {
		var friendRequest models.FriendRequest
		var requester models.User

		err := rows.Scan(
			&friendRequest.ID,
			&friendRequest.RequesterID,
			&friendRequest.RequestedID,
			&friendRequest.Status,
			&friendRequest.Message,
			&friendRequest.CreatedAt,
			&friendRequest.RespondedAt,
			&requester.ID,
			&requester.Name,
			&requester.Email,
			&requester.AvatarURL,
			&requester.Level,
			&requester.TotalPoints,
			&requester.IsOnline,
			&requester.LastActive,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan pending friend request: %w", err)
		}

		friendRequest.Requester = &requester
		requests = append(requests, friendRequest)
	}

	return requests, nil
}

// GetSentFriendRequests retrieves sent friend requests for a user
func (r *FriendsRepository) GetSentFriendRequests(userID uuid.UUID) ([]models.FriendRequest, error) {
	log.Printf("GetSentFriendRequests called: userID=%s", userID)
	query := `
		SELECT fr.id, fr.requester_id, fr.requested_id, fr.status, fr.message, fr.created_at, fr.responded_at,
			   u.id, u.name, u.email, u.avatar_url, u.level, u.total_points, u.is_online, u.last_active
		FROM friend_requests fr
		JOIN users u ON fr.requested_id = u.id
		WHERE fr.requester_id = $1 AND fr.status = 'pending'
		ORDER BY fr.created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent friend requests: %w", err)
	}
	defer rows.Close()

	var requests []models.FriendRequest
	for rows.Next() {
		var friendRequest models.FriendRequest
		var requested models.User

		err := rows.Scan(
			&friendRequest.ID,
			&friendRequest.RequesterID,
			&friendRequest.RequestedID,
			&friendRequest.Status,
			&friendRequest.Message,
			&friendRequest.CreatedAt,
			&friendRequest.RespondedAt,
			&requested.ID,
			&requested.Name,
			&requested.Email,
			&requested.AvatarURL,
			&requested.Level,
			&requested.TotalPoints,
			&requested.IsOnline,
			&requested.LastActive,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan sent friend request: %w", err)
		}

		friendRequest.Requested = &requested
		requests = append(requests, friendRequest)
	}

	return requests, nil
}

// GetFriends retrieves the friends list for a user
func (r *FriendsRepository) GetFriends(userID uuid.UUID) ([]models.Friend, error) {
	log.Printf("GetFriends called: userID=%s", userID)
	query := `
		SELECT u.id, u.name, u.email, u.avatar_url, u.level, u.total_points, u.current_streak,
			   u.best_streak, u.total_quizzes_completed, u.average_score, u.is_online, u.last_active,
			   f.created_at as friends_since
		FROM friendships f
		JOIN users u ON (CASE WHEN f.user1_id = $1 THEN f.user2_id ELSE f.user1_id END) = u.id
		WHERE f.user1_id = $1 OR f.user2_id = $1
		ORDER BY u.name
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friends: %w", err)
	}
	defer rows.Close()

	var friends []models.Friend
	for rows.Next() {
		var friend models.Friend

		err := rows.Scan(
			&friend.ID,
			&friend.Name,
			&friend.Email,
			&friend.AvatarURL,
			&friend.Level,
			&friend.TotalPoints,
			&friend.CurrentStreak,
			&friend.BestStreak,
			&friend.TotalQuizzesCompleted,
			&friend.AverageScore,
			&friend.IsOnline,
			&friend.LastActive,
			&friend.FriendsSince,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan friend: %w", err)
		}

		friends = append(friends, friend)
	}

	return friends, nil
}

// GetFriendship retrieves a friendship between two users
func (r *FriendsRepository) GetFriendship(user1ID, user2ID uuid.UUID) (*models.Friendship, error) {
	log.Printf("GetFriendship called: user1ID=%s, user2ID=%s", user1ID, user2ID)
	// Ensure consistent ordering
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
		SELECT id, user1_id, user2_id, created_at
		FROM friendships
		WHERE user1_id = $1 AND user2_id = $2
	`

	var friendship models.Friendship
	err := r.db.QueryRow(query, user1ID, user2ID).Scan(
		&friendship.ID,
		&friendship.User1ID,
		&friendship.User2ID,
		&friendship.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No friendship found
		}
		return nil, fmt.Errorf("failed to get friendship: %w", err)
	}

	return &friendship, nil
}

// RemoveFriend removes a friendship between two users
func (r *FriendsRepository) RemoveFriend(userID, friendID uuid.UUID) error {
	log.Printf("RemoveFriend called: userID=%s, friendID=%s", userID, friendID)
	// Ensure consistent ordering
	user1ID, user2ID := userID, friendID
	if user1ID.String() > user2ID.String() {
		user1ID, user2ID = user2ID, user1ID
	}

	query := `
		DELETE FROM friendships
		WHERE user1_id = $1 AND user2_id = $2
	`

	result, err := r.db.Exec(query, user1ID, user2ID)
	if err != nil {
		return fmt.Errorf("failed to remove friend: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("friendship not found")
	}

	return nil
}

// AreFriends checks if two users are friends
func (r *FriendsRepository) AreFriends(user1ID, user2ID uuid.UUID) (bool, error) {
	log.Printf("AreFriends called: user1ID=%s, user2ID=%s", user1ID, user2ID)
	friendship, err := r.GetFriendship(user1ID, user2ID)
	if err != nil {
		return false, err
	}
	return friendship != nil, nil
}

// SearchUsers searches for users by name or email
func (r *FriendsRepository) SearchUsers(searchQuery string, currentUserID uuid.UUID, limit, offset int) ([]models.UserSearchResult, int, error) {
	log.Printf("SearchUsers called: searchQuery=%s, currentUserID=%s, limit=%d, offset=%d", searchQuery, currentUserID, limit, offset)
	searchQuery = strings.ToLower(strings.TrimSpace(searchQuery))
	if searchQuery == "" {
		return []models.UserSearchResult{}, 0, nil
	}

	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM users u
		WHERE u.id != $1
		AND (LOWER(u.name) LIKE $2 OR LOWER(u.email) LIKE $2)
	`

	var total int
	err := r.db.QueryRow(countQuery, currentUserID, "%"+searchQuery+"%").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Main query
	query := `
		SELECT u.id, u.name, u.email, u.avatar_url, u.level, u.total_points,
			   u.total_quizzes_completed, u.average_score, u.is_online,
			   CASE WHEN f.id IS NOT NULL THEN true ELSE false END as is_friend,
			   CASE WHEN fr1.id IS NOT NULL OR fr2.id IS NOT NULL THEN true ELSE false END as has_pending_request,
			   CASE WHEN fr1.id IS NOT NULL THEN true ELSE false END as request_sent_by_me
		FROM users u
		LEFT JOIN friendships f ON ((f.user1_id = $1 AND f.user2_id = u.id) OR (f.user2_id = $1 AND f.user1_id = u.id))
		LEFT JOIN friend_requests fr1 ON (fr1.requester_id = $1 AND fr1.requested_id = u.id AND fr1.status = 'pending')
		LEFT JOIN friend_requests fr2 ON (fr2.requester_id = u.id AND fr2.requested_id = $1 AND fr2.status = 'pending')
		WHERE u.id != $1
		AND (LOWER(u.name) LIKE $2 OR LOWER(u.email) LIKE $2)
		ORDER BY u.name
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, currentUserID, "%"+searchQuery+"%", limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []models.UserSearchResult
	for rows.Next() {
		var user models.UserSearchResult

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.AvatarURL,
			&user.Level,
			&user.TotalPoints,
			&user.TotalQuizzesCompleted,
			&user.AverageScore,
			&user.IsOnline,
			&user.IsFriend,
			&user.HasPendingRequest,
			&user.RequestSentByMe,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user search result: %w", err)
		}

		users = append(users, user)
	}

	return users, total, nil
}

// GetFriendNotifications retrieves friend notifications for a user
func (r *FriendsRepository) GetFriendNotifications(userID uuid.UUID, limit, offset int) ([]models.FriendNotification, int, error) {
	log.Printf("GetFriendNotifications called: userID=%s, limit=%d, offset=%d", userID, limit, offset)
	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM friend_notifications
		WHERE user_id = $1
	`

	var total int
	err := r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	// Main query
	query := `
		SELECT fn.id, fn.user_id, fn.type, fn.title, fn.message, fn.related_user_id,
			   fn.friend_request_id, fn.is_read, fn.created_at, fn.read_at,
			   u.id, u.name, u.email, u.avatar_url, u.level, u.total_points, u.is_online, u.last_active
		FROM friend_notifications fn
		LEFT JOIN users u ON fn.related_user_id = u.id
		WHERE fn.user_id = $1
		ORDER BY fn.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get friend notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.FriendNotification
	for rows.Next() {
		var notification models.FriendNotification
		var relatedUser models.User
		var hasRelatedUser bool

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.RelatedUserID,
			&notification.FriendRequestID,
			&notification.IsRead,
			&notification.CreatedAt,
			&notification.ReadAt,
			&relatedUser.ID,
			&relatedUser.Name,
			&relatedUser.Email,
			&relatedUser.AvatarURL,
			&relatedUser.Level,
			&relatedUser.TotalPoints,
			&relatedUser.IsOnline,
			&relatedUser.LastActive,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan friend notification: %w", err)
		}

		// Check if related user data was populated
		if relatedUser.ID != uuid.Nil {
			notification.RelatedUser = &relatedUser
			hasRelatedUser = true
		}

		if !hasRelatedUser {
			notification.RelatedUser = nil
		}

		notifications = append(notifications, notification)
	}

	return notifications, total, nil
}

// MarkNotificationAsRead marks a specific notification as read
func (r *FriendsRepository) MarkNotificationAsRead(notificationID uuid.UUID, userID uuid.UUID) error {
	log.Printf("MarkNotificationAsRead called: notificationID=%s, userID=%s", notificationID, userID)
	query := `
		UPDATE friend_notifications
		SET is_read = true, read_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2 AND is_read = false
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

	return nil
}

// MarkAllNotificationsAsRead marks all notifications as read for a user
func (r *FriendsRepository) MarkAllNotificationsAsRead(userID uuid.UUID) error {
	log.Printf("MarkAllNotificationsAsRead called: userID=%s", userID)
	query := `
		UPDATE friend_notifications
		SET is_read = true, read_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND is_read = false
	`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

// GetUnreadNotificationCount gets the count of unread notifications for a user
func (r *FriendsRepository) GetUnreadNotificationCount(userID uuid.UUID) (int, error) {
	log.Printf("GetUnreadNotificationCount called: userID=%s", userID)
	query := `
		SELECT COUNT(*)
		FROM friend_notifications
		WHERE user_id = $1 AND is_read = false
	`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread notification count: %w", err)
	}

	return count, nil
}