package handlers

import (
	"fmt"
	"net/http"

	"quizninja-api/config"
	internalmodels "quizninja-api/internal/models"
	"quizninja-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StatisticsHandler handles user statistics endpoints
type StatisticsHandler struct {
	repo   *repository.Repository
	config *config.Config
}

// NewStatisticsHandler creates a new StatisticsHandler
func NewStatisticsHandler(cfg *config.Config) *StatisticsHandler {
	return &StatisticsHandler{
		repo:   repository.NewRepository(),
		config: cfg,
	}
}

// UpdateStatistics updates user statistics after quiz completion
// POST /internal/v1/users/:userId/statistics
func (h *StatisticsHandler) UpdateStatistics(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "Invalid user ID format",
			},
		})
		return
	}

	var req internalmodels.UpdateStatisticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "Invalid request body",
				"details": err.Error(),
			},
		})
		return
	}

	// Update user statistics using the existing repository method
	err = h.repo.User.UpdateUserStatistics(userID, req.PercentageScore)
	if err != nil {
		fmt.Printf("Failed to update user statistics for user %s: %v\n", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "UPDATE_FAILED",
				"message": "Failed to update user statistics",
			},
		})
		return
	}

	// Fetch updated user stats
	user, err := h.repo.User.GetUserByID(userID)
	if err != nil {
		// Stats were updated, but we couldn't fetch the updated values
		// Return success with zero values
		c.JSON(http.StatusOK, internalmodels.UpdateStatisticsResponse{
			Success: true,
		})
		return
	}

	response := internalmodels.UpdateStatisticsResponse{
		Success:               true,
		TotalQuizzesCompleted: user.TotalQuizzesCompleted,
		AverageScore:          user.AverageScore,
		CurrentStreak:         user.CurrentStreak,
		BestStreak:            user.BestStreak,
	}

	c.JSON(http.StatusOK, response)
}
