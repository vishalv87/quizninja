package handlers

import (
	"log"
	"net/http"
	"strconv"

	"quizninja-api/config"
	"quizninja-api/models"
	"quizninja-api/repository"
	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DiscussionHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewDiscussionHandler creates a new discussion handler
func NewDiscussionHandler(cfg *config.Config) *DiscussionHandler {
	return &DiscussionHandler{
		repo: repository.NewRepository(),
		cfg:  cfg,
	}
}

// GetDiscussions retrieves discussions with filtering
// GET /api/v1/discussions
func (h *DiscussionHandler) GetDiscussions(c *gin.Context) {
	log.Println("GetDiscussions called")

	// Get user ID if authenticated (optional for public discussions)
	var userID *uuid.UUID
	if authUserID, exists := c.Get("user_id"); exists {
		uid := authUserID.(uuid.UUID)
		userID = &uid
	}

	// Parse filters manually to handle UUID strings properly
	var filters models.DiscussionFilters

	// Parse quiz_id if provided
	if quizIDStr := c.Query("quiz_id"); quizIDStr != "" {
		if quizID, err := uuid.Parse(quizIDStr); err == nil {
			filters.QuizID = &quizID
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz_id format")
			return
		}
	}

	// Parse question_id if provided
	if questionIDStr := c.Query("question_id"); questionIDStr != "" {
		if questionID, err := uuid.Parse(questionIDStr); err == nil {
			filters.QuestionID = &questionID
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid question_id format")
			return
		}
	}

	// Parse user_id if provided
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			filters.UserID = &userID
		} else {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user_id format")
			return
		}
	}

	// Parse other string/int parameters
	filters.Type = c.Query("type")
	filters.Search = c.Query("search")
	filters.SortBy = c.DefaultQuery("sort_by", "created_at")
	filters.SortOrder = c.DefaultQuery("sort_order", "desc")

	// Parse pagination parameters
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && page > 0 {
		filters.Page = page
	} else {
		filters.Page = 1
	}

	if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10")); err == nil && pageSize > 0 && pageSize <= 100 {
		filters.PageSize = pageSize
	} else {
		filters.PageSize = 10
	}

	// Validate type if provided
	if filters.Type != "" {
		validTypes := map[string]bool{"general": true, "question": true, "explanation": true, "help": true}
		if !validTypes[filters.Type] {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid type. Must be one of: general, question, explanation, help")
			return
		}
	}

	// Validate sort_by if provided
	if filters.SortBy != "" {
		validSortBy := map[string]bool{"created_at": true, "likes_count": true, "replies_count": true}
		if !validSortBy[filters.SortBy] {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid sort_by. Must be one of: created_at, likes_count, replies_count")
			return
		}
	}

	// Validate sort_order if provided
	if filters.SortOrder != "" {
		validSortOrder := map[string]bool{"asc": true, "desc": true}
		if !validSortOrder[filters.SortOrder] {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid sort_order. Must be one of: asc, desc")
			return
		}
	}


	discussions, total, err := h.repo.Discussion.GetDiscussions(&filters, userID)
	if err != nil {
		log.Printf("Error getting discussions: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve discussions")
		return
	}

	totalPages := (total + filters.PageSize - 1) / filters.PageSize

	response := models.DiscussionListResponse{
		Discussions: discussions,
		Total:       total,
		Page:        filters.Page,
		PageSize:    filters.PageSize,
		TotalPages:  totalPages,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Discussions retrieved successfully",
	})
}

// GetDiscussion retrieves a single discussion by ID
// GET /api/v1/discussions/:id
func (h *DiscussionHandler) GetDiscussion(c *gin.Context) {
	log.Println("GetDiscussion called")

	discussionIDStr := c.Param("id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid discussion ID")
		return
	}

	// Get user ID if authenticated (optional)
	var userID *uuid.UUID
	if authUserID, exists := c.Get("user_id"); exists {
		uid := authUserID.(uuid.UUID)
		userID = &uid
	}

	discussion, err := h.repo.Discussion.GetDiscussionWithDetails(discussionID, userID)
	if err != nil {
		log.Printf("Error getting discussion: %v", err)
		if err.Error() == "discussion not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Discussion not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve discussion")
		return
	}

	response := models.DiscussionDetailResponse{
		Discussion: models.Discussion{
			ID:           discussion.ID,
			QuizID:       discussion.QuizID,
			QuestionID:   discussion.QuestionID,
			UserID:       discussion.UserID,
			Content:      discussion.Content,
			LikesCount:   discussion.LikesCount,
			RepliesCount: discussion.RepliesCount,
			Type:         discussion.Type,
			CreatedAt:    discussion.CreatedAt,
			UpdatedAt:    discussion.UpdatedAt,
			IsLikedByUser: discussion.IsLikedByUser,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Discussion retrieved successfully",
	})
}

// CreateDiscussion creates a new discussion
// POST /api/v1/discussions
func (h *DiscussionHandler) CreateDiscussion(c *gin.Context) {
	log.Println("CreateDiscussion called")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.CreateDiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Set default type if not provided
	if req.Type == "" {
		req.Type = "general"
	}

	discussion := &models.Discussion{
		QuizID:     req.QuizID,
		QuestionID: req.QuestionID,
		UserID:     userID.(uuid.UUID),
		Content:    req.Content,
		Type:       req.Type,
	}

	err := h.repo.Discussion.CreateDiscussion(discussion)
	if err != nil {
		log.Printf("Error creating discussion: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create discussion")
		return
	}

	response := models.DiscussionDetailResponse{
		Discussion: *discussion,
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    response,
		"message": "Discussion created successfully",
	})
}

// UpdateDiscussion updates an existing discussion
// PUT /api/v1/discussions/:id
func (h *DiscussionHandler) UpdateDiscussion(c *gin.Context) {
	log.Println("UpdateDiscussion called")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	discussionIDStr := c.Param("id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid discussion ID")
		return
	}

	var req models.UpdateDiscussionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Get existing discussion
	discussion, err := h.repo.Discussion.GetDiscussionByID(discussionID)
	if err != nil {
		log.Printf("Error getting discussion: %v", err)
		utils.ErrorResponse(c, http.StatusNotFound, "Discussion not found")
		return
	}

	// Check ownership
	if discussion.UserID != userID.(uuid.UUID) {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only edit your own discussions")
		return
	}

	// Update fields if provided
	if req.Content != nil {
		discussion.Content = *req.Content
	}
	if req.Type != nil {
		discussion.Type = *req.Type
	}

	err = h.repo.Discussion.UpdateDiscussion(discussion)
	if err != nil {
		log.Printf("Error updating discussion: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update discussion")
		return
	}

	response := models.DiscussionDetailResponse{
		Discussion: *discussion,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Discussion updated successfully",
	})
}

// DeleteDiscussion deletes a discussion
// DELETE /api/v1/discussions/:id
func (h *DiscussionHandler) DeleteDiscussion(c *gin.Context) {
	log.Println("DeleteDiscussion called")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	discussionIDStr := c.Param("id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid discussion ID")
		return
	}

	err = h.repo.Discussion.DeleteDiscussion(discussionID, userID.(uuid.UUID))
	if err != nil {
		log.Printf("Error deleting discussion: %v", err)
		if err.Error() == "discussion not found or unauthorized" {
			utils.ErrorResponse(c, http.StatusNotFound, "Discussion not found or unauthorized")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete discussion")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Discussion deleted successfully",
	})
}

// GetDiscussionReplies retrieves replies for a discussion
// GET /api/v1/discussions/:id/replies
func (h *DiscussionHandler) GetDiscussionReplies(c *gin.Context) {
	log.Println("GetDiscussionReplies called")

	discussionIDStr := c.Param("id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid discussion ID")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get user ID if authenticated (optional)
	var userID *uuid.UUID
	if authUserID, exists := c.Get("user_id"); exists {
		uid := authUserID.(uuid.UUID)
		userID = &uid
	}

	replies, total, err := h.repo.Discussion.GetDiscussionReplies(discussionID, userID, pageSize, offset)
	if err != nil {
		log.Printf("Error getting discussion replies: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve replies")
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	response := models.DiscussionRepliesResponse{
		Replies:    replies,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Replies retrieved successfully",
	})
}

// CreateDiscussionReply creates a new reply to a discussion
// POST /api/v1/discussions/:id/replies
func (h *DiscussionHandler) CreateDiscussionReply(c *gin.Context) {
	log.Println("CreateDiscussionReply called")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	discussionIDStr := c.Param("id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid discussion ID")
		return
	}

	var req models.CreateDiscussionReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Verify discussion exists
	_, err = h.repo.Discussion.GetDiscussionByID(discussionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Discussion not found")
		return
	}

	reply := &models.DiscussionReply{
		DiscussionID: discussionID,
		UserID:       userID.(uuid.UUID),
		Content:      req.Content,
	}

	err = h.repo.Discussion.CreateDiscussionReply(reply)
	if err != nil {
		log.Printf("Error creating discussion reply: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create reply")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    reply,
		"message": "Reply created successfully",
	})
}

// UpdateDiscussionReply updates an existing reply
// PUT /api/v1/discussions/replies/:replyId
func (h *DiscussionHandler) UpdateDiscussionReply(c *gin.Context) {
	log.Println("UpdateDiscussionReply called")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	replyIDStr := c.Param("replyId")
	replyID, err := uuid.Parse(replyIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid reply ID")
		return
	}

	var req models.UpdateDiscussionReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Get existing reply
	reply, err := h.repo.Discussion.GetReplyByID(replyID)
	if err != nil {
		log.Printf("Error getting reply: %v", err)
		utils.ErrorResponse(c, http.StatusNotFound, "Reply not found")
		return
	}

	// Check ownership
	if reply.UserID != userID.(uuid.UUID) {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only edit your own replies")
		return
	}

	// Update content if provided
	if req.Content != nil {
		reply.Content = *req.Content
	}

	err = h.repo.Discussion.UpdateDiscussionReply(reply)
	if err != nil {
		log.Printf("Error updating reply: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update reply")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    reply,
		"message": "Reply updated successfully",
	})
}

// DeleteDiscussionReply deletes a reply
// DELETE /api/v1/discussions/replies/:replyId
func (h *DiscussionHandler) DeleteDiscussionReply(c *gin.Context) {
	log.Println("DeleteDiscussionReply called")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	replyIDStr := c.Param("replyId")
	replyID, err := uuid.Parse(replyIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid reply ID")
		return
	}

	err = h.repo.Discussion.DeleteDiscussionReply(replyID, userID.(uuid.UUID))
	if err != nil {
		log.Printf("Error deleting reply: %v", err)
		if err.Error() == "reply not found or unauthorized" {
			utils.ErrorResponse(c, http.StatusNotFound, "Reply not found or unauthorized")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete reply")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reply deleted successfully",
	})
}

// LikeDiscussion likes or unlikes a discussion
// PUT /api/v1/discussions/:id/like
func (h *DiscussionHandler) LikeDiscussion(c *gin.Context) {
	log.Println("LikeDiscussion called")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	discussionIDStr := c.Param("id")
	discussionID, err := uuid.Parse(discussionIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid discussion ID")
		return
	}

	// Verify discussion exists
	discussion, err := h.repo.Discussion.GetDiscussionByID(discussionID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Discussion not found")
		return
	}

	// Check if already liked
	isLiked, err := h.repo.Discussion.IsDiscussionLikedByUser(discussionID, userID.(uuid.UUID))
	if err != nil {
		log.Printf("Error checking if discussion is liked: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check like status")
		return
	}

	var message string
	if isLiked {
		// Unlike
		err = h.repo.Discussion.UnlikeDiscussion(discussionID, userID.(uuid.UUID))
		if err != nil {
			log.Printf("Error unliking discussion: %v", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to unlike discussion")
			return
		}
		message = "Discussion unliked successfully"
		isLiked = false
	} else {
		// Like
		err = h.repo.Discussion.LikeDiscussion(discussionID, userID.(uuid.UUID))
		if err != nil {
			log.Printf("Error liking discussion: %v", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to like discussion")
			return
		}
		message = "Discussion liked successfully"
		isLiked = true
	}

	// Get updated discussion to return new likes count
	updatedDiscussion, err := h.repo.Discussion.GetDiscussionByID(discussionID)
	if err != nil {
		log.Printf("Error getting updated discussion: %v", err)
		// Still return success but with old likes count
		updatedDiscussion = discussion
	}

	response := models.LikeResponse{
		IsLiked:    isLiked,
		LikesCount: updatedDiscussion.LikesCount,
		Message:    message,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": message,
	})
}

// LikeDiscussionReply likes or unlikes a discussion reply
// PUT /api/v1/discussions/replies/:replyId/like
func (h *DiscussionHandler) LikeDiscussionReply(c *gin.Context) {
	log.Println("LikeDiscussionReply called")

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	replyIDStr := c.Param("replyId")
	replyID, err := uuid.Parse(replyIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid reply ID")
		return
	}

	// Verify reply exists
	reply, err := h.repo.Discussion.GetReplyByID(replyID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Reply not found")
		return
	}

	// Check if already liked
	isLiked, err := h.repo.Discussion.IsReplyLikedByUser(replyID, userID.(uuid.UUID))
	if err != nil {
		log.Printf("Error checking if reply is liked: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check like status")
		return
	}

	var message string
	if isLiked {
		// Unlike
		err = h.repo.Discussion.UnlikeDiscussionReply(replyID, userID.(uuid.UUID))
		if err != nil {
			log.Printf("Error unliking reply: %v", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to unlike reply")
			return
		}
		message = "Reply unliked successfully"
		isLiked = false
	} else {
		// Like
		err = h.repo.Discussion.LikeDiscussionReply(replyID, userID.(uuid.UUID))
		if err != nil {
			log.Printf("Error liking reply: %v", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to like reply")
			return
		}
		message = "Reply liked successfully"
		isLiked = true
	}

	// Get updated reply to return new likes count
	updatedReply, err := h.repo.Discussion.GetReplyByID(replyID)
	if err != nil {
		log.Printf("Error getting updated reply: %v", err)
		// Still return success but with old likes count
		updatedReply = reply
	}

	response := models.LikeResponse{
		IsLiked:    isLiked,
		LikesCount: updatedReply.LikesCount,
		Message:    message,
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": message,
	})
}

// GetDiscussionStats retrieves discussion statistics
// GET /api/v1/discussions/stats
func (h *DiscussionHandler) GetDiscussionStats(c *gin.Context) {
	log.Println("GetDiscussionStats called")

	// Get user ID if authenticated (optional)
	var userID *uuid.UUID
	if authUserID, exists := c.Get("user_id"); exists {
		uid := authUserID.(uuid.UUID)
		userID = &uid
	}

	stats, err := h.repo.Discussion.GetDiscussionStats(userID)
	if err != nil {
		log.Printf("Error getting discussion stats: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve discussion statistics")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    stats,
		"message": "Discussion statistics retrieved successfully",
	})
}