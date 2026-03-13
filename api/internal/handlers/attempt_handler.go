package handlers

import (
	"net/http"
	"time"

	"quizninja-api/config"
	internalmodels "quizninja-api/internal/models"
	"quizninja-api/models"
	"quizninja-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AttemptHandler handles internal attempt-related endpoints
type AttemptHandler struct {
	repo   *repository.Repository
	config *config.Config
}

// NewAttemptHandler creates a new AttemptHandler
func NewAttemptHandler(cfg *config.Config) *AttemptHandler {
	return &AttemptHandler{
		repo:   repository.NewRepository(),
		config: cfg,
	}
}

// ValidateAttempt validates an attempt for submission
// POST /internal/v1/attempts/:attemptId/validate
func (h *AttemptHandler) ValidateAttempt(c *gin.Context) {
	attemptIDParam := c.Param("attemptId")
	attemptID, err := uuid.Parse(attemptIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ATTEMPT_ID",
				"message": "Invalid attempt ID format",
			},
		})
		return
	}

	var req internalmodels.ValidateAttemptRequest
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

	// Fetch the attempt
	attempt, err := h.repo.Quiz.GetQuizAttempt(attemptID)
	if err != nil {
		if err.Error() == "quiz attempt not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"valid":       false,
				"error_code":  "ATTEMPT_NOT_FOUND",
				"error_message": "Quiz attempt not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch attempt",
			},
		})
		return
	}

	// Validate ownership
	if attempt.UserID != req.UserID {
		c.JSON(http.StatusOK, gin.H{
			"valid":       false,
			"error_code":  "UNAUTHORIZED",
			"error_message": "User does not own this attempt",
		})
		return
	}

	// Validate quiz ID matches
	if attempt.QuizID != req.QuizID {
		c.JSON(http.StatusOK, gin.H{
			"valid":       false,
			"error_code":  "QUIZ_MISMATCH",
			"error_message": "Attempt does not belong to specified quiz",
		})
		return
	}

	// Check if already completed
	if attempt.IsCompleted {
		c.JSON(http.StatusOK, gin.H{
			"valid":       false,
			"error_code":  "ALREADY_COMPLETED",
			"error_message": "Quiz attempt already completed",
		})
		return
	}

	// Build response
	response := internalmodels.ValidateAttemptResponse{
		Valid: true,
		Attempt: &internalmodels.AttemptInfo{
			ID:          attempt.ID,
			QuizID:      attempt.QuizID,
			UserID:      attempt.UserID,
			Status:      attempt.Status,
			IsCompleted: attempt.IsCompleted,
			StartedAt:   attempt.StartedAt,
			CompletedAt: attempt.CompletedAt,
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateAttempt updates an attempt with score and completion data
// PUT /internal/v1/attempts/:attemptId
func (h *AttemptHandler) UpdateAttempt(c *gin.Context) {
	attemptIDParam := c.Param("attemptId")
	attemptID, err := uuid.Parse(attemptIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ATTEMPT_ID",
				"message": "Invalid attempt ID format",
			},
		})
		return
	}

	var req internalmodels.UpdateAttemptRequest
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

	// Fetch the attempt
	attempt, err := h.repo.Quiz.GetQuizAttempt(attemptID)
	if err != nil {
		if err.Error() == "quiz attempt not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "ATTEMPT_NOT_FOUND",
					"message": "Quiz attempt not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch attempt",
			},
		})
		return
	}

	// Validate ownership
	if attempt.UserID != req.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "User does not own this attempt",
			},
		})
		return
	}

	// Convert internal validated answers to model format
	attemptAnswers := make([]models.AttemptAnswer, len(req.Answers))
	for i, va := range req.Answers {
		attemptAnswers[i] = models.AttemptAnswer{
			QuestionID:     va.QuestionID,
			SelectedOption: va.SelectedOption,
			TextAnswer:     va.TextAnswer,
			IsCorrect:      va.IsCorrect,
			PointsEarned:   va.PointsEarned,
		}
	}

	// Update the attempt
	timeNow := time.Now()
	attempt.Answers = attemptAnswers
	attempt.Score = req.Score
	attempt.TotalPoints = req.TotalPoints
	attempt.TimeSpent = req.TimeSpent
	attempt.PercentageScore = req.PercentageScore
	attempt.Passed = req.Passed
	attempt.Status = req.Status
	attempt.UpdatedAt = timeNow

	if req.Status == "completed" {
		attempt.IsCompleted = true
		attempt.CompletedAt = &timeNow
	}

	err = h.repo.Quiz.UpdateQuizAttempt(attempt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "UPDATE_FAILED",
				"message": "Failed to update attempt",
			},
		})
		return
	}

	response := internalmodels.UpdateAttemptResponse{
		Success:     true,
		AttemptID:   attempt.ID,
		CompletedAt: attempt.CompletedAt,
	}

	c.JSON(http.StatusOK, response)
}