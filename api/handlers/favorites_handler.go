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

type FavoritesHandler struct {
	repo   *repository.Repository
	config *config.Config
}

func NewFavoritesHandler(cfg *config.Config) *FavoritesHandler {
	return &FavoritesHandler{
		repo:   repository.NewRepository(),
		config: cfg,
	}
}

// AddFavorite handles POST /api/v1/favorites
func (h *FavoritesHandler) AddFavorite(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var request models.AddFavoriteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Verify quiz exists
	quiz, err := h.repo.Quiz.GetQuizByID(request.QuizID)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Quiz not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	// Check if quiz is public
	if !quiz.IsPublic {
		utils.ErrorResponse(c, http.StatusForbidden, "Quiz is not public")
		return
	}

	// Add to favorites
	err = h.repo.Quiz.AddFavorite(userID, request.QuizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Quiz added to favorites successfully",
		"quiz_id": request.QuizID,
	})
}

// RemoveFavorite handles DELETE /api/v1/favorites/:quizId
func (h *FavoritesHandler) RemoveFavorite(c *gin.Context) {
	userID := getUserIDFromContext(c)

	quizIdParam := c.Param("quizId")
	quizID, err := uuid.Parse(quizIdParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+quizIdParam)
		return
	}

	err = h.repo.Quiz.RemoveFavorite(userID, quizID)
	if err != nil {
		if err.Error() == "favorite not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Quiz not in favorites")
			return
		}
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message": "Quiz removed from favorites successfully",
		"quiz_id": quizID,
	})
}

// GetFavorites handles GET /api/v1/favorites
func (h *FavoritesHandler) GetFavorites(c *gin.Context) {
	userID := getUserIDFromContext(c)

	// Parse pagination parameters
	pageParam := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSizeParam := c.DefaultQuery("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	favorites, total, err := h.repo.Quiz.GetUserFavorites(userID, page, pageSize)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Calculate total pages
	totalPages := (total + pageSize - 1) / pageSize

	response := &models.FavoritesListResponse{
		Favorites:  favorites,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	utils.SuccessResponse(c, response)
}

// CheckFavorite handles GET /api/v1/favorites/check/:quizId
func (h *FavoritesHandler) CheckFavorite(c *gin.Context) {
	userID := getUserIDFromContext(c)

	quizIdParam := c.Param("quizId")
	quizID, err := uuid.Parse(quizIdParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+quizIdParam)
		return
	}

	isFavorite, err := h.repo.Quiz.IsFavorite(userID, quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"quiz_id":     quizID,
		"is_favorite": isFavorite,
	})
}
