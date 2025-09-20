package handlers

import (
	"net/http"
	"strconv"

	"quizninja-api/config"
	"quizninja-api/models"
	"quizninja-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AchievementHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewAchievementHandler(cfg *config.Config) *AchievementHandler {
	return &AchievementHandler{
		repo: repository.NewRepository(),
		cfg:  cfg,
	}
}

// GetAllAchievements retrieves all available achievements
func (ah *AchievementHandler) GetAllAchievements(c *gin.Context) {
	achievements, err := ah.repo.Achievement.GetAllAchievements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve achievements",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"achievements": achievements,
		"total":        len(achievements),
	})
}

// GetUserAchievements retrieves all achievements unlocked by the current user
func (ah *AchievementHandler) GetUserAchievements(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	achievements, err := ah.repo.Achievement.GetUserAchievements(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user achievements",
		})
		return
	}

	response := models.AchievementListResponse{
		Achievements: achievements,
		Total:        len(achievements),
	}

	c.JSON(http.StatusOK, response)
}

// GetAchievementProgress retrieves achievement progress for the current user
func (ah *AchievementHandler) GetAchievementProgress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	progress, err := ah.repo.Achievement.GetAchievementProgress(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve achievement progress",
		})
		return
	}

	response := models.AchievementProgressResponse{
		Progress: progress,
		Total:    len(progress),
	}

	c.JSON(http.StatusOK, response)
}

// GetUserAchievementsByUserID retrieves achievements for a specific user (for friend profiles, etc.)
func (ah *AchievementHandler) GetUserAchievementsByUserID(c *gin.Context) {
	userIDParam := c.Param("userId")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Optional: Check if current user is authorized to view this user's achievements
	// For now, we'll allow viewing anyone's achievements (public profile data)

	achievements, err := ah.repo.Achievement.GetUserAchievements(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user achievements",
		})
		return
	}

	response := models.AchievementListResponse{
		Achievements: achievements,
		Total:        len(achievements),
	}

	c.JSON(http.StatusOK, response)
}

// UnlockAchievement manually unlocks an achievement for a user (admin/testing endpoint)
func (ah *AchievementHandler) UnlockAchievement(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	achievementKey := c.Param("key")
	if achievementKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Achievement key is required",
		})
		return
	}

	// Check if user already has this achievement
	hasAchievement, err := ah.repo.Achievement.HasUserAchievementByKey(uid, achievementKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check existing achievements",
		})
		return
	}

	if hasAchievement {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already has this achievement",
		})
		return
	}

	userAchievement, err := ah.repo.Achievement.UnlockAchievement(uid, achievementKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to unlock achievement",
		})
		return
	}

	// Create notification data
	notification := models.AchievementNotification{
		AchievementID: userAchievement.Achievement.ID,
		Title:         userAchievement.Achievement.Title,
		Description:   userAchievement.Achievement.Description,
		Icon:          userAchievement.Achievement.Icon,
		Color:         userAchievement.Achievement.Color,
		PointsAwarded: userAchievement.PointsAwarded,
		IsRare:        userAchievement.Achievement.IsRare,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Achievement unlocked successfully",
		"achievement":  userAchievement,
		"notification": notification,
	})
}

// GetAchievementsByCategory retrieves achievements filtered by category
func (ah *AchievementHandler) GetAchievementsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Category is required",
		})
		return
	}

	achievements, err := ah.repo.Achievement.GetAllAchievements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve achievements",
		})
		return
	}

	// Filter by category
	var filteredAchievements []models.Achievement
	for _, achievement := range achievements {
		if achievement.Category == category {
			filteredAchievements = append(filteredAchievements, achievement)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"achievements": filteredAchievements,
		"total":        len(filteredAchievements),
		"category":     category,
	})
}

// GetAchievementStats returns achievement statistics for a user
func (ah *AchievementHandler) GetAchievementStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Get user achievements
	userAchievements, err := ah.repo.Achievement.GetUserAchievements(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user achievements",
		})
		return
	}

	// Get all achievements
	allAchievements, err := ah.repo.Achievement.GetAllAchievements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve all achievements",
		})
		return
	}

	// Calculate statistics
	totalAchievements := len(allAchievements)
	unlockedAchievements := len(userAchievements)
	completionPercentage := float64(0)
	if totalAchievements > 0 {
		completionPercentage = float64(unlockedAchievements) / float64(totalAchievements) * 100
	}

	// Calculate total points from achievements
	totalPointsFromAchievements := 0
	rareAchievements := 0
	categoryStats := make(map[string]int)

	for _, ua := range userAchievements {
		totalPointsFromAchievements += ua.PointsAwarded
		if ua.Achievement != nil {
			if ua.Achievement.IsRare {
				rareAchievements++
			}
			categoryStats[ua.Achievement.Category]++
		}
	}

	stats := gin.H{
		"total_achievements":              totalAchievements,
		"unlocked_achievements":           unlockedAchievements,
		"locked_achievements":             totalAchievements - unlockedAchievements,
		"completion_percentage":           completionPercentage,
		"total_points_from_achievements":  totalPointsFromAchievements,
		"rare_achievements":               rareAchievements,
		"achievements_by_category":        categoryStats,
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// CheckAchievements checks and unlocks achievements for a user based on their current stats
// This is typically called after quiz completion or other achievement-triggering events
func (ah *AchievementHandler) CheckAchievements(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Get query parameter for trigger type (optional)
	triggerType := c.Query("trigger")
	if triggerType == "" {
		triggerType = "manual"
	}

	newAchievements, err := ah.checkAndUnlockAchievements(uid, triggerType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check achievements",
		})
		return
	}

	notifications := make([]models.AchievementNotification, len(newAchievements))
	for i, ua := range newAchievements {
		notifications[i] = models.AchievementNotification{
			AchievementID: ua.Achievement.ID,
			Title:         ua.Achievement.Title,
			Description:   ua.Achievement.Description,
			Icon:          ua.Achievement.Icon,
			Color:         ua.Achievement.Color,
			PointsAwarded: ua.PointsAwarded,
			IsRare:        ua.Achievement.IsRare,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "Achievement check completed",
		"new_achievements":    newAchievements,
		"notifications":       notifications,
		"count":               len(newAchievements),
		"trigger_type":        triggerType,
	})
}

// checkAndUnlockAchievements is a helper method that checks user stats and unlocks applicable achievements
func (ah *AchievementHandler) checkAndUnlockAchievements(userID uuid.UUID, triggerType string) ([]models.UserAchievement, error) {
	var newAchievements []models.UserAchievement

	// Define achievement keys to check based on trigger type
	achievementKeys := []string{}

	switch triggerType {
	case "quiz_completed":
		achievementKeys = []string{"first_win", "quiz_master", "perfect_score"}
	case "streak_updated":
		achievementKeys = []string{"week_warrior", "streak_legend"}
	case "friend_added":
		achievementKeys = []string{"social_butterfly"}
	case "manual":
		// Check all achievements
		achievementKeys = []string{
			"first_win", "week_warrior", "quiz_master", "streak_legend",
			"social_butterfly", "perfect_score",
		}
	default:
		achievementKeys = []string{"first_win", "week_warrior", "quiz_master"}
	}

	// Get achievement progress to determine which achievements should be unlocked
	progress, err := ah.repo.Achievement.GetAchievementProgress(userID)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	progressMap := make(map[string]models.AchievementProgress)
	for _, p := range progress {
		// Find the achievement key for this achievement ID
		achievement, err := ah.repo.Achievement.GetAllAchievements()
		if err != nil {
			continue
		}
		for _, a := range achievement {
			if a.ID == p.AchievementID {
				progressMap[a.Key] = p
				break
			}
		}
	}

	// Check each achievement key
	for _, key := range achievementKeys {
		// Skip if user already has this achievement
		hasAchievement, err := ah.repo.Achievement.HasUserAchievementByKey(userID, key)
		if err != nil || hasAchievement {
			continue
		}

		// Check if achievement should be unlocked based on progress
		if prog, exists := progressMap[key]; exists {
			if prog.Progress >= 100 && !prog.IsUnlocked {
				// Unlock the achievement
				userAchievement, err := ah.repo.Achievement.UnlockAchievement(userID, key)
				if err != nil {
					// Log error but continue with other achievements
					continue
				}
				newAchievements = append(newAchievements, *userAchievement)
			}
		}
	}

	return newAchievements, nil
}

// GetLeaderboardWithAchievements retrieves leaderboard entries with achievement counts
func (ah *AchievementHandler) GetLeaderboardWithAchievements(c *gin.Context) {
	// Parse query parameters
	period := c.DefaultQuery("period", "alltime")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Get leaderboard data
	leaderboard, total, err := ah.repo.Leaderboard.GetGlobalLeaderboard(period, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve leaderboard",
		})
		return
	}

	// Enhance leaderboard entries with achievement counts
	for i := range leaderboard {
		userAchievements, err := ah.repo.Achievement.GetUserAchievements(leaderboard[i].UserID)
		if err != nil {
			// Continue with empty achievements if there's an error
			userAchievements = []models.UserAchievement{}
		}

		// Count achievements by category
		categoryCount := make(map[string]int)
		rareCount := 0

		for _, ua := range userAchievements {
			if ua.Achievement != nil {
				categoryCount[ua.Achievement.Category]++
				if ua.Achievement.IsRare {
					rareCount++
				}
			}
		}

		// Add achievement data to leaderboard entry (this would require extending the model)
		// For now, we'll include it in a separate field in the response
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
		"total":       total,
		"period":      period,
		"limit":       limit,
		"offset":      offset,
	})
}