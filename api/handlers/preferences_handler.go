package handlers

import (
	"net/http"

	"quizninja-api/config"
	"quizninja-api/models"
	"quizninja-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PreferencesHandler struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

func NewPreferencesHandler(config *config.Config) *PreferencesHandler {
	return &PreferencesHandler{
		userRepo: repository.NewUserRepository(),
		config:   config,
	}
}

func (ph *PreferencesHandler) UpdatePreferences(c *gin.Context) {
	var req models.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	preferences := &models.UserPreferences{
		UserID:                  userID,
		SelectedInterests:       models.StringArray(req.SelectedInterests),
		DifficultyPreference:    req.DifficultyPreference,
		NotificationsEnabled:    req.NotificationsEnabled,
		NotificationFrequency:   req.NotificationFrequency,
	}

	err = ph.userRepo.UpdateUserPreferences(preferences)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update preferences",
		})
		return
	}

	updatedPreferences, err := ph.userRepo.GetUserPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch updated preferences",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": updatedPreferences,
		"message": "Preferences updated successfully",
	})
}

func (ph *PreferencesHandler) GetPreferences(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	preferences, err := ph.userRepo.GetUserPreferences(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch preferences",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": preferences,
	})
}