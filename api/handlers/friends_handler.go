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

type FriendsHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewFriendsHandler creates a new friends handler
func NewFriendsHandler(cfg *config.Config) *FriendsHandler {
	return &FriendsHandler{
		repo: repository.NewRepository(),
		cfg:  cfg,
	}
}

// SendFriendRequest sends a friend request to another user
// POST /api/v1/friends/requests
func (h *FriendsHandler) SendFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	requesterID := userID.(uuid.UUID)

	var req models.SendFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Check if user is trying to send friend request to themselves
	if requesterID == req.RequestedUserID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cannot send friend request to yourself")
		return
	}

	// Check if they are already friends
	areFriends, err := h.repo.Friends.AreFriends(requesterID, req.RequestedUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check friendship status")
		return
	}

	if areFriends {
		utils.ErrorResponse(c, http.StatusBadRequest, "You are already friends with this user")
		return
	}

	// Check if there's already a pending request
	existingRequest, err := h.repo.Friends.GetFriendRequestBetweenUsers(requesterID, req.RequestedUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check existing requests")
		return
	}

	if existingRequest != nil && existingRequest.Status == "pending" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Friend request already sent")
		return
	}

	// Check if there's a reverse pending request
	reverseRequest, err := h.repo.Friends.GetFriendRequestBetweenUsers(req.RequestedUserID, requesterID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check reverse requests")
		return
	}

	if reverseRequest != nil && reverseRequest.Status == "pending" {
		utils.ErrorResponse(c, http.StatusBadRequest, "This user has already sent you a friend request")
		return
	}

	// Send the friend request
	friendRequest, err := h.repo.Friends.SendFriendRequest(requesterID, req.RequestedUserID, req.Message)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to send friend request")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Friend request sent successfully",
		"request": friendRequest,
	})
}

// GetFriendRequests gets pending and sent friend requests for the user
// GET /api/v1/friends/requests
func (h *FriendsHandler) GetFriendRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	// Get pending requests (requests sent to current user)
	pendingRequests, err := h.repo.Friends.GetPendingFriendRequests(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get pending requests")
		return
	}

	// Get sent requests (requests sent by current user)
	sentRequests, err := h.repo.Friends.GetSentFriendRequests(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get sent requests")
		return
	}

	response := models.FriendRequestsResponse{
		PendingRequests: pendingRequests,
		SentRequests:    sentRequests,
		Total:           len(pendingRequests) + len(sentRequests),
	}

	c.JSON(http.StatusOK, response)
}

// RespondToFriendRequest accepts or rejects a friend request
// PUT /api/v1/friends/requests/:id
func (h *FriendsHandler) RespondToFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request ID")
		return
	}

	var req models.RespondToFriendRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Get the friend request to verify it belongs to the current user
	friendRequest, err := h.repo.Friends.GetFriendRequest(requestID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Friend request not found")
		return
	}

	// Verify the current user is the requested user (receiver of the request)
	if friendRequest.RequestedID != currentUserID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only respond to requests sent to you")
		return
	}

	// Respond to the request
	err = h.repo.Friends.RespondToFriendRequest(requestID, req.Status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to respond to friend request")
		return
	}

	var message string
	if req.Status == "accepted" {
		message = "Friend request accepted"
	} else {
		message = "Friend request rejected"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"status":  req.Status,
	})
}

// CancelFriendRequest cancels a sent friend request
// DELETE /api/v1/friends/requests/:id
func (h *FriendsHandler) CancelFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	requestIDStr := c.Param("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request ID")
		return
	}

	// Cancel the request
	err = h.repo.Friends.CancelFriendRequest(requestID, currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cancel friend request")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Friend request cancelled",
	})
}

// GetFriends gets the user's friends list
// GET /api/v1/friends
func (h *FriendsHandler) GetFriends(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	friends, err := h.repo.Friends.GetFriends(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get friends")
		return
	}

	response := models.FriendsListResponse{
		Friends: friends,
		Total:   len(friends),
	}

	c.JSON(http.StatusOK, response)
}

// RemoveFriend removes a friend from the user's friends list
// DELETE /api/v1/friends/:id
func (h *FriendsHandler) RemoveFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	friendIDStr := c.Param("id")
	friendID, err := uuid.Parse(friendIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid friend ID")
		return
	}

	// Check if they are actually friends
	areFriends, err := h.repo.Friends.AreFriends(currentUserID, friendID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check friendship status")
		return
	}

	if !areFriends {
		utils.ErrorResponse(c, http.StatusBadRequest, "You are not friends with this user")
		return
	}

	// Remove the friendship
	err = h.repo.Friends.RemoveFriend(currentUserID, friendID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to remove friend")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Friend removed successfully",
	})
}

// SearchUsers searches for users by name or email
// GET /api/v1/friends/search?q=query&limit=10&offset=0
func (h *FriendsHandler) SearchUsers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Search query is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	users, total, err := h.repo.Friends.SearchUsers(query, currentUserID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search users")
		return
	}

	response := models.UserSearchResponse{
		Users: users,
		Total: total,
	}

	c.JSON(http.StatusOK, response)
}

// GetFriendNotifications gets friend-related notifications for the user
// GET /api/v1/friends/notifications?limit=20&offset=0
func (h *FriendsHandler) GetFriendNotifications(c *gin.Context) {
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

	notifications, total, err := h.repo.Friends.GetFriendNotifications(currentUserID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get notifications")
		return
	}

	unreadCount, err := h.repo.Friends.GetUnreadNotificationCount(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get unread count")
		return
	}

	response := models.FriendNotificationsResponse{
		Notifications: notifications,
		UnreadCount:   unreadCount,
		Total:         total,
	}

	c.JSON(http.StatusOK, response)
}

// MarkNotificationAsRead marks a specific notification as read
// PUT /api/v1/friends/notifications/:id/read
func (h *FriendsHandler) MarkNotificationAsRead(c *gin.Context) {
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

	err = h.repo.Friends.MarkNotificationAsRead(notificationID, currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification marked as read",
	})
}

// MarkAllNotificationsAsRead marks all notifications as read for the user
// PUT /api/v1/friends/notifications/read-all
func (h *FriendsHandler) MarkAllNotificationsAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	err := h.repo.Friends.MarkAllNotificationsAsRead(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to mark all notifications as read")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All notifications marked as read",
	})
}