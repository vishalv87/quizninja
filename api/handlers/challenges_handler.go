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

type ChallengesHandler struct {
	repo *repository.Repository
	cfg  *config.Config
}

// NewChallengesHandler creates a new challenges handler
func NewChallengesHandler(cfg *config.Config) *ChallengesHandler {
	return &ChallengesHandler{
		repo: repository.NewRepository(),
		cfg:  cfg,
	}
}

// CreateChallenge creates a new challenge
// POST /api/v1/challenges
func (h *ChallengesHandler) CreateChallenge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	challengerID := userID.(uuid.UUID)

	var req models.CreateChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	//  SECURITY: Sanitize and validate message if provided
	if req.Message != nil {
		sanitizedMessage := utils.SanitizeHTML(*req.Message)
		if err := utils.ValidateMessage(sanitizedMessage); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		req.Message = &sanitizedMessage
	}

	// Check if user is trying to challenge themselves
	if challengerID == req.ChallengeeUserID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cannot challenge yourself")
		return
	}

	// Check if users are friends
	canChallenge, err := h.repo.Challenges.CanUserChallenge(challengerID, req.ChallengeeUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check challenge eligibility")
		return
	}

	if !canChallenge {
		utils.ErrorResponse(c, http.StatusBadRequest, "You can only challenge friends")
		return
	}

	// Check if there's already a pending challenge for this quiz between these users
	hasPending, err := h.repo.Challenges.HasPendingChallenge(challengerID, req.ChallengeeUserID, req.QuizID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to check pending challenges")
		return
	}

	if hasPending {
		utils.ErrorResponse(c, http.StatusBadRequest, "You already have a pending challenge with this user for this quiz")
		return
	}

	// Verify quiz exists
	_, err = h.repo.Quiz.GetQuizByID(req.QuizID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Quiz not found")
		return
	}

	// Create challenge
	challenge := &models.Challenge{
		ChallengerID:     challengerID,
		ChallengeeID:     req.ChallengeeUserID,
		QuizID:           req.QuizID,
		Message:          req.Message,
		ExpiresAt:        req.ExpiresAt,
		IsGroupChallenge: req.IsGroupChallenge,
		ParticipantIDs:   req.ParticipantIDs,
		IsTestData:       true,
	}

	if err := h.repo.Challenges.CreateChallenge(challenge); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create challenge")
		return
	}

	// Get the created challenge with details
	challengeDetails, err := h.repo.Challenges.GetChallengeWithDetails(challenge.ID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve challenge details")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Challenge created successfully",
		"challenge": challengeDetails,
	})
}

// GetChallenges retrieves challenges for the current user
// GET /api/v1/challenges
func (h *ChallengesHandler) GetChallenges(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	// Parse query parameters
	var filters models.ChallengeFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// Set defaults
	if filters.Page == 0 {
		filters.Page = 1
	}
	if filters.PageSize == 0 {
		filters.PageSize = 10
	}

	challenges, total, err := h.repo.Challenges.GetUserChallenges(currentUserID, &filters)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve challenges")
		return
	}

	totalPages := (total + filters.PageSize - 1) / filters.PageSize

	response := models.ChallengeListResponse{
		Challenges: challenges,
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetChallengeByID retrieves a specific challenge
// GET /api/v1/challenges/:id
func (h *ChallengesHandler) GetChallengeByID(c *gin.Context) {
	challengeIDStr := c.Param("id")
	challengeID, err := uuid.Parse(challengeIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid challenge ID")
		return
	}

	challenge, err := h.repo.Challenges.GetChallengeWithDetails(challengeID)
	if err != nil {
		if err.Error() == "challenge not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Challenge not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve challenge")
		return
	}

	//  Use centralized authorization check
	err = utils.RequireAnyOwnership(c,
		[]uuid.UUID{challenge.ChallengerID, challenge.ChallengeeID},
		"challenge",
	)
	if utils.HandleAuthError(c, err) {
		return
	}

	response := models.ChallengeDetailResponse{
		Challenge: *challenge,
	}

	c.JSON(http.StatusOK, response)
}

// AcceptChallenge accepts a challenge
// PUT /api/v1/challenges/:id/accept
func (h *ChallengesHandler) AcceptChallenge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	challengeIDStr := c.Param("id")
	challengeID, err := uuid.Parse(challengeIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid challenge ID")
		return
	}

	// Accept the challenge
	if err := h.repo.Challenges.AcceptChallenge(challengeID, currentUserID); err != nil {
		if err.Error() == "challenge not found or cannot be accepted" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Challenge not found or cannot be accepted")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to accept challenge")
		return
	}

	// Get updated challenge details
	challenge, err := h.repo.Challenges.GetChallengeWithDetails(challengeID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve updated challenge")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Challenge accepted successfully",
		"challenge": challenge,
	})
}

// DeclineChallenge declines a challenge
// PUT /api/v1/challenges/:id/decline
func (h *ChallengesHandler) DeclineChallenge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	challengeIDStr := c.Param("id")
	challengeID, err := uuid.Parse(challengeIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid challenge ID")
		return
	}

	// Decline the challenge
	if err := h.repo.Challenges.DeclineChallenge(challengeID, currentUserID); err != nil {
		if err.Error() == "challenge not found or cannot be declined" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Challenge not found or cannot be declined")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to decline challenge")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Challenge declined successfully",
	})
}

// CancelChallenge allows the challenger to cancel a pending challenge they created
// PUT /api/v1/challenges/:id/cancel
func (h *ChallengesHandler) CancelChallenge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	challengeIDStr := c.Param("id")
	challengeID, err := uuid.Parse(challengeIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid challenge ID")
		return
	}

	// Cancel the challenge
	if err := h.repo.Challenges.CancelChallenge(challengeID, currentUserID); err != nil {
		if err.Error() == "challenge not found, already accepted/declined, or you are not the challenger" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Challenge not found, already accepted/declined, or you are not the challenger")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cancel challenge")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Challenge cancelled successfully",
	})
}

// UpdateChallengeScore updates the user's score for a challenge
// PUT /api/v1/challenges/:id/score
func (h *ChallengesHandler) UpdateChallengeScore(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	challengeIDStr := c.Param("id")
	challengeID, err := uuid.Parse(challengeIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid challenge ID")
		return
	}

	var req models.UpdateChallengeScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Update the score
	if err := h.repo.Challenges.UpdateChallengeScore(challengeID, currentUserID, req.UserScore); err != nil {
		if err.Error() == "user is not part of this challenge" {
			utils.ErrorResponse(c, http.StatusBadRequest, "You are not part of this challenge")
			return
		}
		if err.Error() == "challenge not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Challenge not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update challenge score")
		return
	}

	// Get updated challenge details
	challenge, err := h.repo.Challenges.GetChallengeWithDetails(challengeID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve updated challenge")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Score updated successfully",
		"challenge": challenge,
	})
}

// GetChallengeStats retrieves challenge statistics for the current user
// GET /api/v1/challenges/stats
func (h *ChallengesHandler) GetChallengeStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	stats, err := h.repo.Challenges.GetChallengeStats(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve challenge statistics")
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ExpireChallenges manually triggers challenge expiration (admin endpoint)
// POST /api/v1/challenges/expire
func (h *ChallengesHandler) ExpireChallenges(c *gin.Context) {
	if err := h.repo.Challenges.ExpireChallenges(); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to expire challenges")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Challenges expired successfully",
	})
}

// GetPendingChallenges retrieves pending challenges for the current user
// GET /api/v1/challenges/pending
func (h *ChallengesHandler) GetPendingChallenges(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	challenges, err := h.repo.Challenges.GetPendingChallenges(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve pending challenges")
		return
	}

	// Ensure challenges is never nil to avoid null in JSON response
	if challenges == nil {
		challenges = []models.ChallengeWithDetails{}
	}

	c.JSON(http.StatusOK, gin.H{
		"challenges": challenges,
		"total":      len(challenges),
	})
}

// GetActiveChallenges retrieves active challenges for the current user
// GET /api/v1/challenges/active
func (h *ChallengesHandler) GetActiveChallenges(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	challenges, err := h.repo.Challenges.GetActiveChallenges(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve active challenges")
		return
	}

	// Ensure challenges is never nil to avoid null in JSON response
	if challenges == nil {
		challenges = []models.ChallengeWithDetails{}
	}

	c.JSON(http.StatusOK, gin.H{
		"challenges": challenges,
		"total":      len(challenges),
	})
}

// GetCompletedChallenges retrieves completed challenges for the current user
// GET /api/v1/challenges/completed
func (h *ChallengesHandler) GetCompletedChallenges(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	challenges, err := h.repo.Challenges.GetCompletedChallenges(currentUserID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve completed challenges")
		return
	}

	// Ensure challenges is never nil to avoid null in JSON response
	if challenges == nil {
		challenges = []models.ChallengeWithDetails{}
	}

	c.JSON(http.StatusOK, gin.H{
		"challenges": challenges,
		"total":      len(challenges),
	})
}

// LinkAttemptToChallenge links a quiz attempt to a challenge
// POST /api/v1/challenges/:id/link-attempt
func (h *ChallengesHandler) LinkAttemptToChallenge(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	challengeIDStr := c.Param("id")
	challengeID, err := uuid.Parse(challengeIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid challenge ID")
		return
	}

	var req struct {
		AttemptID string `json:"attempt_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	attemptID, err := uuid.Parse(req.AttemptID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid attempt ID")
		return
	}

	// Link the attempt to the challenge
	if err := h.repo.Challenges.LinkAttemptToChallenge(challengeID, attemptID, currentUserID); err != nil {
		if err.Error() == "user is not part of this challenge" {
			utils.ErrorResponse(c, http.StatusBadRequest, "You are not part of this challenge")
			return
		}
		if err.Error() == "challenge not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Challenge not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to link attempt to challenge")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attempt linked to challenge successfully",
	})
}

// CompleteChallengeAttempt marks a user's challenge attempt as complete
// PUT /api/v1/challenges/:id/complete
func (h *ChallengesHandler) CompleteChallengeAttempt(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserID := userID.(uuid.UUID)

	challengeIDStr := c.Param("id")
	challengeID, err := uuid.Parse(challengeIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid challenge ID")
		return
	}

	var req struct {
		Score float64 `json:"score" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Complete the challenge attempt
	if err := h.repo.Challenges.CompleteChallengeAttempt(challengeID, currentUserID, req.Score); err != nil {
		if err.Error() == "user is not part of this challenge" {
			utils.ErrorResponse(c, http.StatusBadRequest, "You are not part of this challenge")
			return
		}
		if err.Error() == "challenge not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Challenge not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to complete challenge attempt")
		return
	}

	// Get updated challenge details
	challenge, err := h.repo.Challenges.GetChallengeWithDetails(challengeID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve updated challenge")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Challenge attempt completed successfully",
		"challenge": challenge,
	})
}
