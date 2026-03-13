package models

import (
	"time"

	"github.com/google/uuid"
)

// FriendRequest represents a friend request
type FriendRequest struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	RequesterID uuid.UUID  `json:"requester_id" db:"requester_id"`
	RequestedID uuid.UUID  `json:"requested_id" db:"requested_id"`
	Status      string     `json:"status" db:"status"`
	Message     *string    `json:"message,omitempty" db:"message"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	RespondedAt *time.Time `json:"responded_at,omitempty" db:"responded_at"`
	Requester   *User      `json:"requester,omitempty"`
	Requested   *User      `json:"requested,omitempty"`
}

// Friendship represents an accepted friendship between two users
type Friendship struct {
	ID        uuid.UUID `json:"id" db:"id"`
	User1ID   uuid.UUID `json:"user1_id" db:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id" db:"user2_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	User1     *User     `json:"user1,omitempty"`
	User2     *User     `json:"user2,omitempty"`
}

// FriendNotification represents a notification related to friends
type FriendNotification struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	Type            string     `json:"type" db:"type"`
	Title           string     `json:"title" db:"title"`
	Message         *string    `json:"message,omitempty" db:"message"`
	RelatedUserID   *uuid.UUID `json:"related_user_id,omitempty" db:"related_user_id"`
	FriendRequestID *uuid.UUID `json:"friend_request_id,omitempty" db:"friend_request_id"`
	IsRead          bool       `json:"is_read" db:"is_read"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	ReadAt          *time.Time `json:"read_at,omitempty" db:"read_at"`
	RelatedUser     *User      `json:"related_user,omitempty"`
}

// Friend represents a user's friend with additional metadata
type Friend struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	AvatarURL             *string   `json:"avatar_url,omitempty"`
	Level                 string    `json:"level"`
	TotalPoints           int       `json:"total_points"`
	CurrentStreak         int       `json:"current_streak"`
	BestStreak            int       `json:"best_streak"`
	TotalQuizzesCompleted int       `json:"total_quizzes_completed"`
	AverageScore          float64   `json:"average_score"`
	IsOnline              bool      `json:"is_online"`
	LastActive            time.Time `json:"last_active"`
	FriendsSince          time.Time `json:"friends_since"`
}

// Friends DTOs

// SendFriendRequestRequest represents the request to send a friend request
type SendFriendRequestRequest struct {
	RequestedUserID uuid.UUID `json:"requested_user_id" binding:"required"`
	Message         *string   `json:"message,omitempty"`
}

// RespondToFriendRequestRequest represents the request to respond to a friend request
type RespondToFriendRequestRequest struct {
	Status string `json:"status" binding:"required,oneof=accepted rejected"`
}

// FriendRequestsResponse represents the response for friend requests
type FriendRequestsResponse struct {
	PendingRequests []FriendRequest `json:"pending_requests"`
	SentRequests    []FriendRequest `json:"sent_requests"`
	Total           int             `json:"total"`
}

// FriendsListResponse represents the response for friends list
type FriendsListResponse struct {
	Friends []Friend `json:"friends"`
	Total   int      `json:"total"`
}

// UserSearchResponse represents the response for user search
type UserSearchResponse struct {
	Users []UserSearchResult `json:"users"`
	Total int                `json:"total"`
}

// UserSearchResult represents a user in search results
type UserSearchResult struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	AvatarURL             *string   `json:"avatar_url,omitempty"`
	Level                 string    `json:"level"`
	TotalPoints           int       `json:"total_points"`
	TotalQuizzesCompleted int       `json:"total_quizzes_completed"`
	AverageScore          float64   `json:"average_score"`
	IsOnline              bool      `json:"is_online"`
	IsFriend              bool      `json:"is_friend"`
	HasPendingRequest     bool      `json:"has_pending_request"`
	RequestSentByMe       bool      `json:"request_sent_by_me"`
}

// FriendNotificationsResponse represents the response for friend notifications
type FriendNotificationsResponse struct {
	Notifications []FriendNotification `json:"notifications"`
	UnreadCount   int                  `json:"unread_count"`
	Total         int                  `json:"total"`
}