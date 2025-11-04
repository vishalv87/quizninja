package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Notification represents a unified notification that can be of any type
type Notification struct {
	ID                uuid.UUID        `json:"id" db:"id"`
	UserID            uuid.UUID        `json:"userId" db:"user_id"`
	Type              string           `json:"type" db:"type"`
	Title             string           `json:"title" db:"title"`
	Message           *string          `json:"message,omitempty" db:"message"`
	Data              NotificationData `json:"data" db:"data"`
	RelatedUserID     *uuid.UUID       `json:"relatedUserId,omitempty" db:"related_user_id"`
	RelatedEntityID   *uuid.UUID       `json:"relatedEntityId,omitempty" db:"related_entity_id"`
	RelatedEntityType *string          `json:"relatedEntityType,omitempty" db:"related_entity_type"`
	IsRead            bool             `json:"isRead" db:"is_read"`
	IsDeleted         bool             `json:"isDeleted" db:"is_deleted"`
	CreatedAt         time.Time        `json:"timestamp" db:"created_at"`
	ReadAt            *time.Time       `json:"readAt,omitempty" db:"read_at"`
	DeletedAt         *time.Time       `json:"deletedAt,omitempty" db:"deleted_at"`
	ExpiresAt         *time.Time       `json:"expiresAt,omitempty" db:"expires_at"`
	RelatedUser       *User            `json:"relatedUser,omitempty"`
}

// NotificationData represents the JSONB data field as a custom type
type NotificationData map[string]interface{}

// Value implements the driver.Valuer interface for database storage
func (nd NotificationData) Value() (driver.Value, error) {
	if len(nd) == 0 {
		return "{}", nil
	}
	return json.Marshal(nd)
}

// Scan implements the sql.Scanner interface for database retrieval
func (nd *NotificationData) Scan(value interface{}) error {
	if value == nil {
		*nd = NotificationData{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, nd)
	case string:
		return json.Unmarshal([]byte(v), nd)
	default:
		return fmt.Errorf("cannot scan %T into NotificationData", value)
	}
}

// NotificationType constants
const (
	NotificationTypeFriendRequest       = "friend_request"
	NotificationTypeFriendAccepted      = "friend_accepted"
	NotificationTypeFriendRejected      = "friend_rejected"
	NotificationTypeChallengeReceived   = "challenge_received"
	NotificationTypeChallengeAccepted   = "challenge_accepted"
	NotificationTypeChallengeDeclined   = "challenge_declined"
	NotificationTypeChallengeCompleted  = "challenge_completed"
	NotificationTypeAchievementUnlocked = "achievement_unlocked"
	NotificationTypeGeneral             = "general"
	NotificationTypeSystemAnnouncement  = "system_announcement"
)

// NotificationFilters represents filters for notification queries
type NotificationFilters struct {
	Type      string     `form:"type" binding:"omitempty,oneof=friend_request friend_accepted friend_rejected challenge_received challenge_accepted challenge_declined challenge_completed achievement_unlocked general system_announcement"`
	IsRead    *bool      `form:"is_read"`
	StartDate *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"end_date" time_format:"2006-01-02"`
	Page      int        `form:"page,default=1" binding:"min=1"`
	PageSize  int        `form:"page_size,default=20" binding:"min=1,max=100"`
	SortBy    string     `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at read_at"`
	SortOrder string     `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// NotificationResponse represents the API response for notifications
type NotificationResponse struct {
	Notifications []Notification `json:"notifications"`
	UnreadCount   int            `json:"unread_count"`
	Total         int            `json:"total"`
	Page          int            `json:"page"`
	PageSize      int            `json:"page_size"`
	TotalPages    int            `json:"total_pages"`
}

// MarkNotificationReadRequest represents the request to mark a notification as read
type MarkNotificationReadRequest struct {
	// No body needed - notification ID comes from URL path
}

// MarkAllNotificationsReadRequest represents the request to mark all notifications as read
type MarkAllNotificationsReadRequest struct {
	// No body needed
}

// CreateNotificationRequest represents the request to create a new notification (admin/system use)
type CreateNotificationRequest struct {
	UserID            uuid.UUID        `json:"user_id" binding:"required"`
	Type              string           `json:"type" binding:"required"`
	Title             string           `json:"title" binding:"required,min=1,max=255"`
	Message           *string          `json:"message,omitempty"`
	Data              NotificationData `json:"data,omitempty"`
	RelatedUserID     *uuid.UUID       `json:"related_user_id,omitempty"`
	RelatedEntityID   *uuid.UUID       `json:"related_entity_id,omitempty"`
	RelatedEntityType *string          `json:"related_entity_type,omitempty"`
	ExpiresAt         *time.Time       `json:"expires_at,omitempty"`
}

// NotificationStatsResponse represents notification statistics
type NotificationStatsResponse struct {
	TotalNotifications  int                    `json:"total_notifications"`
	UnreadNotifications int                    `json:"unread_notifications"`
	NotificationsByType map[string]int         `json:"notifications_by_type"`
	RecentNotifications []Notification         `json:"recent_notifications"`
	NotificationCounts  NotificationTypeCounts `json:"notification_counts"`
}

// NotificationTypeCounts represents counts by notification type
type NotificationTypeCounts struct {
	FriendRequests      int `json:"friend_requests"`
	FriendResponses     int `json:"friend_responses"`
	Challenges          int `json:"challenges"`
	Achievements        int `json:"achievements"`
	General             int `json:"general"`
	SystemAnnouncements int `json:"system_announcements"`
}

// IsExpired returns true if the notification has expired
func (n *Notification) IsExpired() bool {
	if n.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*n.ExpiresAt)
}

// IsSoftDeleted returns true if the notification has been soft deleted
func (n *Notification) IsSoftDeleted() bool {
	return n.IsDeleted && n.DeletedAt != nil
}

// CanBeShown returns true if the notification should be visible to users
func (n *Notification) CanBeShown() bool {
	return !n.IsDeleted && !n.IsExpired()
}

// GetDisplayIcon returns an appropriate icon for the notification type
func (n *Notification) GetDisplayIcon() string {
	switch n.Type {
	case NotificationTypeFriendRequest, NotificationTypeFriendAccepted, NotificationTypeFriendRejected:
		return "👥"
	case NotificationTypeChallengeReceived, NotificationTypeChallengeAccepted, NotificationTypeChallengeDeclined, NotificationTypeChallengeCompleted:
		return "⚔️"
	case NotificationTypeAchievementUnlocked:
		return "🏆"
	case NotificationTypeSystemAnnouncement:
		return "📢"
	case NotificationTypeGeneral:
	default:
		return "📝"
	}
	return "📝"
}

// GetDisplayCategory returns a human-readable category for the notification
func (n *Notification) GetDisplayCategory() string {
	switch n.Type {
	case NotificationTypeFriendRequest, NotificationTypeFriendAccepted, NotificationTypeFriendRejected:
		return "Friends"
	case NotificationTypeChallengeReceived, NotificationTypeChallengeAccepted, NotificationTypeChallengeDeclined, NotificationTypeChallengeCompleted:
		return "Challenges"
	case NotificationTypeAchievementUnlocked:
		return "Achievements"
	case NotificationTypeSystemAnnouncement:
		return "System"
	case NotificationTypeGeneral:
	default:
		return "General"
	}
	return "General"
}

// GetActionURL returns the appropriate URL for handling the notification action
func (n *Notification) GetActionURL() *string {
	switch n.Type {
	case NotificationTypeFriendRequest:
		if entityID := n.RelatedEntityID; entityID != nil {
			url := fmt.Sprintf("/friends/requests/%s", entityID.String())
			return &url
		}
	case NotificationTypeChallengeReceived, NotificationTypeChallengeAccepted, NotificationTypeChallengeDeclined, NotificationTypeChallengeCompleted:
		if entityID := n.RelatedEntityID; entityID != nil {
			url := fmt.Sprintf("/challenges/%s", entityID.String())
			return &url
		}
	case NotificationTypeAchievementUnlocked:
		url := "/achievements"
		return &url
	}
	return nil
}

// GetDataValue safely retrieves a value from the notification data
func (n *Notification) GetDataValue(key string) interface{} {
	if n.Data == nil {
		return nil
	}
	return n.Data[key]
}

// GetDataString safely retrieves a string value from the notification data
func (n *Notification) GetDataString(key string) string {
	if value := n.GetDataValue(key); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetDataInt safely retrieves an int value from the notification data
func (n *Notification) GetDataInt(key string) int {
	if value := n.GetDataValue(key); value != nil {
		switch v := value.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	return 0
}

// GetDataFloat safely retrieves a float64 value from the notification data
func (n *Notification) GetDataFloat(key string) float64 {
	if value := n.GetDataValue(key); value != nil {
		switch v := value.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		}
	}
	return 0.0
}

// GetDataBool safely retrieves a bool value from the notification data
func (n *Notification) GetDataBool(key string) bool {
	if value := n.GetDataValue(key); value != nil {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return false
}

// IsActionable returns true if the notification requires user action
func (n *Notification) IsActionable() bool {
	switch n.Type {
	case NotificationTypeFriendRequest:
		return true
	case NotificationTypeChallengeReceived:
		return true
	default:
		return false
	}
}

// GetTimeAgo returns a human-readable time difference string
func (n *Notification) GetTimeAgo() string {
	now := time.Now()
	diff := now.Sub(n.CreatedAt)

	if diff < time.Minute {
		return "Just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		} else if days < 7 {
			return fmt.Sprintf("%d days ago", days)
		} else {
			return n.CreatedAt.Format("Jan 2, 2006")
		}
	}
}

// Backward compatibility with existing friend notification models
// These types ensure the existing friend-related code continues to work

// FriendNotificationCompat represents a friend notification for backward compatibility
type FriendNotificationCompat struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	Type            string     `json:"type" db:"type"`
	Title           string     `json:"title" db:"title"`
	Message         *string    `json:"message,omitempty" db:"message"`
	RelatedUserID   *uuid.UUID `json:"related_user_id,omitempty" db:"related_user_id"`
	FriendRequestID *uuid.UUID `json:"friend_request_id,omitempty" db:"friend_request_id"`
	IsRead          bool       `json:"is_read" db:"is_read"`
	IsDeleted       bool       `json:"is_deleted" db:"is_deleted"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	ReadAt          *time.Time `json:"read_at,omitempty" db:"read_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	RelatedUser     *User      `json:"related_user,omitempty"`
}

// ToFriendNotificationCompat converts a unified notification to the old friend notification format
func (n *Notification) ToFriendNotificationCompat() *FriendNotificationCompat {
	var friendRequestID *uuid.UUID
	if n.RelatedEntityType != nil && *n.RelatedEntityType == "friend_request" && n.RelatedEntityID != nil {
		friendRequestID = n.RelatedEntityID
	}

	return &FriendNotificationCompat{
		ID:              n.ID,
		UserID:          n.UserID,
		Type:            n.Type,
		Title:           n.Title,
		Message:         n.Message,
		RelatedUserID:   n.RelatedUserID,
		FriendRequestID: friendRequestID,
		IsRead:          n.IsRead,
		IsDeleted:       n.IsDeleted,
		CreatedAt:       n.CreatedAt,
		ReadAt:          n.ReadAt,
		DeletedAt:       n.DeletedAt,
		RelatedUser:     n.RelatedUser,
	}
}

// FriendNotificationsResponseCompat represents the backward compatible response
type FriendNotificationsResponseCompat struct {
	Notifications []FriendNotificationCompat `json:"notifications"`
	UnreadCount   int                        `json:"unread_count"`
	Total         int                        `json:"total"`
}
