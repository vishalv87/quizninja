package handlers

import (
	"net/http"
	"time"

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

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Set default notification types if not provided
	notificationTypes := req.NotificationTypes
	if notificationTypes == nil {
		notificationTypes = map[string]interface{}{
			"challenges":           true,
			"achievements":         true,
			"quiz_reminders":       true,
			"friend_activity":      true,
			"leaderboard_updates":  false,
			"system_announcements": true,
		}
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
		NotificationTypes:       notificationTypes,
		IsTestData:              true, // Set for test environment
	}

	err := ph.userRepo.UpdateUserPreferences(preferences)
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
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
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

// CompleteOnboarding handles completing the onboarding process
func (ph *PreferencesHandler) CompleteOnboarding(c *gin.Context) {
	var req models.OnboardingCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	now := time.Now()
	preferences := &models.UserPreferences{
		UserID:                  userID,
		SelectedInterests:       models.StringArray(req.SelectedInterests),
		DifficultyPreference:    req.DifficultyPreference,
		NotificationsEnabled:    req.NotificationsEnabled,
		NotificationFrequency:   req.NotificationFrequency,
		OnboardingCompletedAt:   &now,
		ProfileVisibility:       true,  // Default values for new onboarding
		ShowOnlineStatus:        true,
		AllowFriendRequests:     true,
		ShareActivityStatus:     true,
	}

	// Check if preferences already exist
	existingPrefs, err := ph.userRepo.GetUserPreferences(userID)
	if err == nil && existingPrefs != nil {
		// Update existing preferences
		preferences.ID = existingPrefs.ID
		preferences.CreatedAt = existingPrefs.CreatedAt
		err = ph.userRepo.UpdateUserPreferences(preferences)
	} else {
		// Create new preferences
		err = ph.userRepo.CreateUserPreferences(preferences)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save onboarding preferences",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": preferences,
		"message": "Onboarding completed successfully",
	})
}

// GetOnboardingStatus checks if user has completed onboarding
func (ph *PreferencesHandler) GetOnboardingStatus(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	preferences, err := ph.userRepo.GetUserPreferences(userID)
	if err != nil {
		// No preferences found - onboarding not completed
		response := models.OnboardingStatusResponse{
			IsCompleted: false,
			CompletedAt: nil,
			Preferences: nil,
		}
		c.JSON(http.StatusOK, gin.H{
			"data": response,
		})
		return
	}

	isCompleted := preferences.OnboardingCompletedAt != nil
	response := models.OnboardingStatusResponse{
		IsCompleted: isCompleted,
		CompletedAt: preferences.OnboardingCompletedAt,
		Preferences: preferences,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// Helper function to get boolean value from pointer with default
func getBoolValue(ptr *bool, defaultValue bool) bool {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}