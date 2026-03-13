package handlers

import (
	"net/http"
	"strconv"

	"quizninja-api/config"
	"quizninja-api/models"
	"quizninja-api/repository"
	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(cfg *config.Config) *NotificationHandler {
	return &NotificationHandler{
		repo: repository.NewRepository(),
		cfg:  cfg,
	}
}

// GetNotifications retrieves notifications for the current user
// GET /api/v1/notifications
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	// Parse query parameters
	var filters models.NotificationFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// Set defaults
	if filters.Page == 0 {
		filters.Page = 1
	}
	if filters.PageSize == 0 {
		filters.PageSize = 20
	}

	// Get notifications
	notifications, total, err := h.repo.Notification.GetNotifications(currentUserID, &filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve notifications")
		return
	}

	// Get unread count
	unreadCount, err := h.repo.Notification.GetUnreadNotificationCount(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get unread count")
		return
	}

	// Calculate pagination
	totalPages := (total + filters.PageSize - 1) / filters.PageSize

	response := models.NotificationResponse{
		Notifications: notifications,
		UnreadCount:   unreadCount,
		Total:         total,
		Page:          filters.Page,
		PageSize:      filters.PageSize,
		TotalPages:    totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetNotificationByID retrieves a specific notification
// GET /api/v1/notifications/:id
func (h *NotificationHandler) GetNotificationByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	// Get notification
	notification, err := h.repo.Notification.GetNotificationByID(notificationID, currentUserID)
	if err != nil {
		if err.Error() == "notification not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Notification not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve notification")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification": notification,
	})
}

// MarkNotificationAsRead marks a specific notification as read
// PUT /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkNotificationAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	err = h.repo.Notification.MarkNotificationAsRead(notificationID, currentUserID)
	if err != nil {
		if err.Error() == "notification not found or already read" {
			utils.ErrorResponse(c, http.StatusNotFound, "Notification not found or already read")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification marked as read",
	})
}

// MarkNotificationAsUnread marks a specific notification as unread
// PUT /api/v1/notifications/:id/unread
func (h *NotificationHandler) MarkNotificationAsUnread(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	err = h.repo.Notification.MarkNotificationAsUnread(notificationID, currentUserID)
	if err != nil {
		if err.Error() == "notification not found, not owned by user, or already unread" {
			utils.ErrorResponse(c, http.StatusNotFound, "Notification not found or already unread")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark notification as unread")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification marked as unread",
	})
}

// MarkAllNotificationsAsRead marks all notifications as read for the user
// PUT /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllNotificationsAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	err := h.repo.Notification.MarkAllNotificationsAsRead(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark all notifications as read")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All notifications marked as read",
	})
}

// DeleteNotification deletes a notification
// DELETE /api/v1/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	err = h.repo.Notification.DeleteNotification(notificationID, currentUserID)
	if err != nil {
		if err.Error() == "notification not found or not owned by user" {
			utils.ErrorResponse(c, http.StatusNotFound, "Notification not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete notification")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification deleted successfully",
	})
}

// GetNotificationStats retrieves notification statistics for the user
// GET /api/v1/notifications/stats
func (h *NotificationHandler) GetNotificationStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	stats, err := h.repo.Notification.GetNotificationStats(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve notification statistics")
		return
	}

	c.JSON(http.StatusOK, stats)
}

// CreateNotification creates a new notification (admin/system endpoint)
// POST /api/v1/notifications
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req models.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Create notification
	notification, err := h.repo.Notification.CreateNotification(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create notification")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Notification created successfully",
		"notification": notification,
	})
}

// GetFriendNotifications retrieves friend notifications for backward compatibility
// GET /api/v1/friends/notifications
func (h *NotificationHandler) GetFriendNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	notifications, total, err := h.repo.Notification.GetFriendNotifications(currentUserID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get notifications")
		return
	}

	unreadCount, err := h.repo.Notification.GetFriendUnreadNotificationCount(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get unread count")
		return
	}

	response := models.FriendNotificationsResponseCompat{
		Notifications: notifications,
		UnreadCount:   unreadCount,
		Total:         total,
	}

	c.JSON(http.StatusOK, response)
}

// MarkFriendNotificationAsRead marks a friend notification as read for backward compatibility
// PUT /api/v1/friends/notifications/:id/read
func (h *NotificationHandler) MarkFriendNotificationAsRead(c *gin.Context) {
	// This just delegates to the unified notification handler
	h.MarkNotificationAsRead(c)
}

// MarkAllFriendNotificationsAsRead marks all friend notifications as read for backward compatibility
// PUT /api/v1/friends/notifications/read-all
func (h *NotificationHandler) MarkAllFriendNotificationsAsRead(c *gin.Context) {
	// This just delegates to the unified notification handler
	h.MarkAllNotificationsAsRead(c)
}

// CleanupExpiredNotifications cleans up expired notifications (admin endpoint)
// POST /api/v1/notifications/cleanup
func (h *NotificationHandler) CleanupExpiredNotifications(c *gin.Context) {
	err := h.repo.Notification.CleanupExpiredNotifications()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cleanup expired notifications")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Expired notifications cleaned up successfully",
	})
}
