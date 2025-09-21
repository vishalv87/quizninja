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
	userRepo        *repository.UserRepository
	preferencesRepo *repository.PreferencesRepository
	config          *config.Config
}

func NewPreferencesHandler(config *config.Config) *PreferencesHandler {
	return &PreferencesHandler{
		userRepo:        repository.NewUserRepository(),
		preferencesRepo: repository.NewPreferencesRepository(),
		config:          config,
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
		ProfileVisibility:       getBoolValue(req.ProfileVisibility, true),
		ShowOnlineStatus:        getBoolValue(req.ShowOnlineStatus, true),
		AllowFriendRequests:     getBoolValue(req.AllowFriendRequests, true),
		ShareActivityStatus:     getBoolValue(req.ShareActivityStatus, true),
		NotificationTypes:       req.NotificationTypes,
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

// GetCategories retrieves all available interests/categories
func (ph *PreferencesHandler) GetCategories(c *gin.Context) {
	interests, err := ph.preferencesRepo.GetInterests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": interests,
		"meta": gin.H{
			"total": len(interests),
		},
	})
}

// GetDifficultyLevels retrieves all available difficulty levels
func (ph *PreferencesHandler) GetDifficultyLevels(c *gin.Context) {
	levels, err := ph.preferencesRepo.GetDifficultyLevels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch difficulty levels",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": levels,
		"meta": gin.H{
			"total": len(levels),
		},
	})
}

// GetNotificationFrequencies retrieves all available notification frequencies
func (ph *PreferencesHandler) GetNotificationFrequencies(c *gin.Context) {
	frequencies, err := ph.preferencesRepo.GetNotificationFrequencies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch notification frequencies",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": frequencies,
		"meta": gin.H{
			"total": len(frequencies),
		},
	})
}

// Helper function to get boolean value from pointer with default
func getBoolValue(ptr *bool, defaultValue bool) bool {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}