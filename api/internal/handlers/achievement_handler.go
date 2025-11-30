package handlers

import (
	"net/http"

	"quizninja-api/config"
	internalmodels "quizninja-api/internal/models"
	"quizninja-api/repository"
	"quizninja-api/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AchievementInternalHandler handles internal achievement endpoints
type AchievementInternalHandler struct {
	repo               *repository.Repository
	config             *config.Config
	achievementService *services.AchievementService
}

// NewAchievementInternalHandler creates a new AchievementInternalHandler
func NewAchievementInternalHandler(cfg *config.Config) *AchievementInternalHandler {
	repo := repository.NewRepository()
	return &AchievementInternalHandler{
		repo:               repo,
		config:             cfg,
		achievementService: services.NewAchievementService(repo),
	}
}

// CheckAchievements checks and unlocks achievements for a user
// POST /internal/v1/users/:userId/achievements/check
func (h *AchievementInternalHandler) CheckAchievements(c *gin.Context) {
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

	var req internalmodels.CheckAchievementsRequest
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

	// Map trigger string to AchievementTrigger type
	var trigger services.AchievementTrigger
	switch req.Trigger {
	case "quiz_completed":
		trigger = services.TriggerQuizCompleted
	case "streak_updated":
		trigger = services.TriggerStreakUpdated
	case "friend_added":
		trigger = services.TriggerFriendAdded
	case "perfect_score":
		trigger = services.TriggerPerfectScore
	case "level_up":
		trigger = services.TriggerLevelUp
	case "leaderboard_rank":
		trigger = services.TriggerLeaderboardRank
	default:
		trigger = services.TriggerQuizCompleted // Default to quiz_completed
	}

	// Check achievements using the service
	result, err := h.achievementService.CheckAchievementsForUser(userID, trigger)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "CHECK_FAILED",
				"message": "Failed to check achievements",
				"details": err.Error(),
			},
		})
		return
	}

	// Convert result to internal response format
	notifications := make([]internalmodels.AchievementNotification, len(result.Notifications))
	for i, n := range result.Notifications {
		icon := ""
		if n.Icon != nil {
			icon = *n.Icon
		}
		notifications[i] = internalmodels.AchievementNotification{
			AchievementID: n.AchievementID,
			Title:         n.Title,
			Description:   n.Description,
			Icon:          icon,
			Color:         n.Color,
			PointsAwarded: n.PointsAwarded,
			IsRare:        n.IsRare,
		}
	}

	response := internalmodels.CheckAchievementsResponse{
		NewAchievements: notifications,
		TotalChecked:    result.TotalChecked,
		TotalUnlocked:   result.TotalUnlocked,
	}

	c.JSON(http.StatusOK, response)
}
