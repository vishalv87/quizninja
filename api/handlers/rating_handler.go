package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"quizninja-api/models"
	"quizninja-api/repository"
	"quizninja-api/utils"
)

type RatingHandler struct {
	repo *repository.Repository
}

func NewRatingHandler() *RatingHandler {
	return &RatingHandler{
		repo: repository.NewRepository(),
	}
}

// CreateRating handles POST /api/v1/quizzes/:id/ratings
func (h *RatingHandler) CreateRating(c *gin.Context) {
	quizIDParam := c.Param("id")
	quizID, err := uuid.Parse(quizIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+quizIDParam)
		return
	}

	// Get user ID from context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		if utils.HandleAuthError(c, err) {
			return
		}
	}

	// Check if quiz exists
	_, err = h.repo.Quiz.GetQuizByID(quizID)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Quiz not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	// Check if user already rated this quiz
	existingRating, err := h.repo.Rating.GetUserRating(userID, quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if existingRating != nil {
		utils.ErrorResponse(c, http.StatusConflict, "You have already rated this quiz. Use PUT to update your rating.")
		return
	}

	// Bind request
	var req models.CreateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Create rating
	rating := &models.QuizRating{
		ID:     uuid.New(),
		QuizID: quizID,
		UserID: userID,
		Rating: req.Rating,
		Review: req.Review,
	}

	if err := h.repo.Rating.CreateRating(rating); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, rating)
}

// GetQuizRatings handles GET /api/v1/quizzes/:id/ratings
func (h *RatingHandler) GetQuizRatings(c *gin.Context) {
	quizIDParam := c.Param("id")
	quizID, err := uuid.Parse(quizIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+quizIDParam)
		return
	}

	// Parse pagination parameters
	pageParam := c.DefaultQuery("page", "1")
	pageSizeParam := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get ratings
	ratings, total, err := h.repo.Rating.GetQuizRatings(quizID, page, pageSize)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Get average rating
	avgRating, totalRatings, err := h.repo.Rating.GetAverageRating(quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	response := models.RatingListResponse{
		Ratings:       ratings,
		AverageRating: avgRating,
		TotalRatings:  totalRatings,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		TotalPages:    totalPages,
	}

	utils.SuccessResponse(c, response)
}

// GetAverageRating handles GET /api/v1/quizzes/:id/ratings/average
func (h *RatingHandler) GetAverageRating(c *gin.Context) {
	quizIDParam := c.Param("id")
	quizID, err := uuid.Parse(quizIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+quizIDParam)
		return
	}

	avgRating, totalRatings, err := h.repo.Rating.GetAverageRating(quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	response := models.AverageRatingResponse{
		QuizID:        quizID,
		AverageRating: avgRating,
		TotalRatings:  totalRatings,
	}

	utils.SuccessResponse(c, response)
}

// GetUserRating handles GET /api/v1/quizzes/:id/ratings/user
func (h *RatingHandler) GetUserRating(c *gin.Context) {
	quizIDParam := c.Param("id")
	quizID, err := uuid.Parse(quizIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+quizIDParam)
		return
	}

	// Get user ID from context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		if utils.HandleAuthError(c, err) {
			return
		}
	}

	rating, err := h.repo.Rating.GetUserRating(userID, quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if rating == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "You have not rated this quiz yet")
		return
	}

	utils.SuccessResponse(c, rating)
}

// UpdateRating handles PUT /api/v1/quizzes/:id/ratings/:ratingId
func (h *RatingHandler) UpdateRating(c *gin.Context) {
	ratingIDParam := c.Param("ratingId")
	ratingID, err := uuid.Parse(ratingIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid rating ID format: "+ratingIDParam)
		return
	}

	// Get user ID from context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		if utils.HandleAuthError(c, err) {
			return
		}
	}

	// Get existing rating
	existingRating, err := h.repo.Rating.GetRatingByID(ratingID)
	if err != nil {
		if err.Error() == "rating not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Rating not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	// Check ownership
	if existingRating.UserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only update your own ratings")
		return
	}

	// Bind request
	var req models.UpdateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Update rating fields
	if req.Rating != nil {
		existingRating.Rating = *req.Rating
	}
	if req.Review != nil {
		existingRating.Review = req.Review
	}

	if err := h.repo.Rating.UpdateRating(existingRating); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, existingRating)
}

// DeleteRating handles DELETE /api/v1/quizzes/:id/ratings/:ratingId
func (h *RatingHandler) DeleteRating(c *gin.Context) {
	ratingIDParam := c.Param("ratingId")
	ratingID, err := uuid.Parse(ratingIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid rating ID format: "+ratingIDParam)
		return
	}

	// Get user ID from context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		if utils.HandleAuthError(c, err) {
			return
		}
	}

	// Get existing rating
	existingRating, err := h.repo.Rating.GetRatingByID(ratingID)
	if err != nil {
		if err.Error() == "rating not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Rating not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	// Check ownership
	if existingRating.UserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "You can only delete your own ratings")
		return
	}

	if err := h.repo.Rating.DeleteRating(ratingID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating deleted successfully"})
}