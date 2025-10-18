package handlers

import (
	"net/http"

	"quizninja-api/config"
	"quizninja-api/models"
	"quizninja-api/repository"
	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LeaderboardHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewLeaderboardHandler creates a new leaderboard handler
func NewLeaderboardHandler(cfg *config.Config) *LeaderboardHandler {
	return &LeaderboardHandler{
		repo: repository.NewRepository(),
		cfg:  cfg,
	}
}

// GetLeaderboard retrieves the leaderboard with optional filtering
// GET /api/v1/leaderboard
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	// Parse query parameters
	var filters models.LeaderboardFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// Get user ID from context if available (for friends filtering)
	var userID uuid.UUID
	var isAuthenticated bool
	if userIDInterface, exists := c.Get("user_id"); exists {
		userID = userIDInterface.(uuid.UUID)
		isAuthenticated = true
	}

	var leaderboard []models.LeaderboardEntry
	var total int
	var err error

	// If friends only filter is requested but user is not authenticated, return error
	if filters.FriendsOnly && !isAuthenticated {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication required for friends leaderboard")
		return
	}

	// Get leaderboard data based on filters
	if filters.FriendsOnly && isAuthenticated {
		leaderboard, total, err = h.repo.Leaderboard.GetFriendsLeaderboard(userID, filters.Period, filters.Limit, filters.Offset)
	} else {
		leaderboard, total, err = h.repo.Leaderboard.GetGlobalLeaderboard(filters.Period, filters.Limit, filters.Offset)
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve leaderboard")
		return
	}

	// Get user's rank information if authenticated
	var userRank *models.UserRankInfo
	if isAuthenticated {
		userRank, err = h.repo.Leaderboard.GetUserRank(userID, filters.Period)
		if err != nil {
			// Don't fail the request, just log the error
			userRank = nil
		}
	}

	response := models.LeaderboardResponse{
		Leaderboard: leaderboard,
		UserRank:    userRank,
		Total:       total,
		Period:      filters.Period,
		FriendsOnly: filters.FriendsOnly,
	}

	utils.SuccessResponse(c, response)
}

// GetUserRank retrieves the current user's rank information
// GET /api/v1/leaderboard/rank
func (h *LeaderboardHandler) GetUserRank(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse query parameters
	period := c.DefaultQuery("period", "alltime")

	// Validate period
	validPeriods := map[string]bool{
		"today":   true,
		"week":    true,
		"month":   true,
		"alltime": true,
	}

	if !validPeriods[period] {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid period. Must be one of: today, week, month, alltime")
		return
	}

	userRank, err := h.repo.Leaderboard.GetUserRank(userID.(uuid.UUID), period)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve user rank")
		return
	}

	utils.SuccessResponse(c, map[string]interface{}{
		"user_rank": userRank,
		"period":    period,
	})
}

// UpdateUserScore updates user's score after quiz completion
// This is typically called internally by the quiz completion handler
// POST /api/v1/leaderboard/score
func (h *LeaderboardHandler) UpdateUserScore(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req struct {
		Points int       `json:"points" binding:"required,min=0"`
		QuizID uuid.UUID `json:"quiz_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Update user's score
	err := h.repo.Leaderboard.UpdateUserScore(userID.(uuid.UUID), req.Points, req.QuizID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user score")
		return
	}

	// Recalculate user level
	err = h.repo.Leaderboard.RecalculateUserLevel(userID.(uuid.UUID))
	if err != nil {
		// Don't fail the request for level calculation errors
	}

	utils.SuccessResponse(c, map[string]interface{}{
		"message":        "User score updated successfully",
		"points_awarded": req.Points,
		"quiz_id":        req.QuizID,
	})
}

// GetLeaderboardStats retrieves overall leaderboard statistics
// GET /api/v1/leaderboard/stats
func (h *LeaderboardHandler) GetLeaderboardStats(c *gin.Context) {
	// Parse query parameters
	period := c.DefaultQuery("period", "alltime")

	// Validate period
	validPeriods := map[string]bool{
		"today":   true,
		"week":    true,
		"month":   true,
		"alltime": true,
	}

	if !validPeriods[period] {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid period. Must be one of: today, week, month, alltime")
		return
	}

	// Get basic leaderboard to calculate stats
	leaderboard, total, err := h.repo.Leaderboard.GetGlobalLeaderboard(period, 100, 0)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve leaderboard statistics")
		return
	}

	// Calculate statistics
	var totalPoints int
	var totalQuizzes int
	var totalStreaks int
	averageScore := 0.0

	for _, entry := range leaderboard {
		totalPoints += entry.Points
		totalQuizzes += entry.QuizzesCompleted
		totalStreaks += entry.CurrentStreak
		averageScore += entry.AverageScore
	}

	if len(leaderboard) > 0 {
		averageScore = averageScore / float64(len(leaderboard))
	}

	stats := map[string]interface{}{
		"period":         period,
		"total_users":    total,
		"active_users":   len(leaderboard),
		"total_points":   totalPoints,
		"total_quizzes":  totalQuizzes,
		"average_score":  averageScore,
		"average_streak": float64(totalStreaks) / float64(len(leaderboard)),
		"top_performer":  nil,
		"most_active":    nil,
		"longest_streak": 0,
	}

	// Find top performers
	if len(leaderboard) > 0 {
		topPerformer := leaderboard[0]
		stats["top_performer"] = map[string]interface{}{
			"name":   topPerformer.Name,
			"points": topPerformer.Points,
			"level":  topPerformer.Level,
		}

		// Find most active user (most quizzes completed)
		mostActive := leaderboard[0]
		longestStreak := leaderboard[0].CurrentStreak

		for _, entry := range leaderboard {
			if entry.QuizzesCompleted > mostActive.QuizzesCompleted {
				mostActive = entry
			}
			if entry.CurrentStreak > longestStreak {
				longestStreak = entry.CurrentStreak
			}
		}

		stats["most_active"] = map[string]interface{}{
			"name":              mostActive.Name,
			"quizzes_completed": mostActive.QuizzesCompleted,
		}
		stats["longest_streak"] = longestStreak
	}

	utils.SuccessResponse(c, stats)
}
