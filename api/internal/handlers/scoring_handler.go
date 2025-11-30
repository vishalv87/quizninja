package handlers

import (
	"net/http"
	"strings"

	internalmodels "quizninja-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ScoringHandler handles score calculation endpoints
type ScoringHandler struct{}

// NewScoringHandler creates a new ScoringHandler
func NewScoringHandler() *ScoringHandler {
	return &ScoringHandler{}
}

// CalculateScore calculates the score for submitted answers
// POST /internal/v1/scoring/calculate
func (h *ScoringHandler) CalculateScore(c *gin.Context) {
	var req internalmodels.CalculateScoreRequest
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

	// Create a map for easy question lookup
	questionMap := make(map[uuid.UUID]internalmodels.QuestionForScoring)
	for _, q := range req.Questions {
		questionMap[q.ID] = q
	}

	totalQuestions := len(req.Questions)
	correctAnswers := 0
	validatedAnswers := make([]internalmodels.ValidatedAnswer, 0, len(req.Answers))

	// Validate and score each answer
	for _, answer := range req.Answers {
		question, exists := questionMap[answer.QuestionID]
		if !exists {
			// Skip answers for questions not in the quiz
			continue
		}

		isCorrect := false
		pointsEarned := 0
		selectedOptionIndex := -1
		textAnswer := ""

		// Handle different answer formats
		if answer.SelectedOptionIndex != nil {
			selectedOptionIndex = *answer.SelectedOptionIndex
			// Validate the index is within bounds
			if selectedOptionIndex >= 0 && selectedOptionIndex < len(question.Options) {
				selectedAnswerText := question.Options[selectedOptionIndex]
				textAnswer = selectedAnswerText
				if selectedAnswerText == question.CorrectAnswer {
					isCorrect = true
					pointsEarned = 1
					correctAnswers++
				}
			}
		} else if answer.SelectedOption != nil {
			// Legacy support: if SelectedOption is provided as string
			textAnswer = *answer.SelectedOption
			if textAnswer == question.CorrectAnswer {
				isCorrect = true
				pointsEarned = 1
				correctAnswers++
			}
			// Try to find the index of the selected option
			for i, opt := range question.Options {
				if opt == textAnswer {
					selectedOptionIndex = i
					break
				}
			}
		} else if answer.TextAnswer != "" {
			// For text-based answers (non-multiple choice)
			textAnswer = answer.TextAnswer
			if strings.TrimSpace(strings.ToLower(textAnswer)) == strings.TrimSpace(strings.ToLower(question.CorrectAnswer)) {
				isCorrect = true
				pointsEarned = 1
				correctAnswers++
			}
		}

		validatedAnswer := internalmodels.ValidatedAnswer{
			QuestionID:     answer.QuestionID,
			SelectedOption: selectedOptionIndex,
			TextAnswer:     textAnswer,
			IsCorrect:      isCorrect,
			PointsEarned:   pointsEarned,
		}
		validatedAnswers = append(validatedAnswers, validatedAnswer)
	}

	// Calculate final score
	score := float64(correctAnswers)
	percentageScore := float64(0)
	if totalQuestions > 0 {
		percentageScore = float64(correctAnswers) / float64(totalQuestions) * 100
	}
	passed := percentageScore >= 60 // Passing threshold is 60%

	response := internalmodels.CalculateScoreResponse{
		TotalQuestions:   totalQuestions,
		CorrectAnswers:   correctAnswers,
		Score:            score,
		PercentageScore:  percentageScore,
		Passed:           passed,
		ValidatedAnswers: validatedAnswers,
	}

	c.JSON(http.StatusOK, response)
}
