package handlers

import (
	"net/http"

	"quizninja-api/config"
	internalmodels "quizninja-api/internal/models"
	"quizninja-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// QuizInternalHandler handles internal quiz-related endpoints
type QuizInternalHandler struct {
	repo   *repository.Repository
	config *config.Config
}

// NewQuizInternalHandler creates a new QuizInternalHandler
func NewQuizInternalHandler(cfg *config.Config) *QuizInternalHandler {
	return &QuizInternalHandler{
		repo:   repository.NewRepository(),
		config: cfg,
	}
}

// GetQuestionsInternal returns quiz questions with correct answers (for internal scoring)
// GET /internal/v1/quizzes/:quizId/questions
func (h *QuizInternalHandler) GetQuestionsInternal(c *gin.Context) {
	quizIDParam := c.Param("quizId")
	quizID, err := uuid.Parse(quizIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_QUIZ_ID",
				"message": "Invalid quiz ID format",
			},
		})
		return
	}

	// Get quiz with questions (includes correct answers)
	quiz, err := h.repo.Quiz.GetQuizByIDWithQuestions(quizID)
	if err != nil {
		if err.Error() == "quiz not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "QUIZ_NOT_FOUND",
					"message": "Quiz not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to fetch quiz",
			},
		})
		return
	}

	// Convert to internal format with correct answers included
	questions := make([]internalmodels.QuestionWithAnswer, len(quiz.Questions))
	for i, q := range quiz.Questions {
		questions[i] = internalmodels.QuestionWithAnswer{
			ID:            q.ID,
			QuizID:        q.QuizID,
			QuestionText:  q.QuestionText,
			QuestionType:  q.QuestionType,
			Options:       q.Options,
			CorrectAnswer: q.CorrectAnswer, // Include correct answer for internal use
			Points:        q.Points,
			OrderIndex:    q.Order,
		}
	}

	response := internalmodels.GetQuestionsResponse{
		QuizID:    quizID,
		Questions: questions,
	}

	c.JSON(http.StatusOK, response)
}
