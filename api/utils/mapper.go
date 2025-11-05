package utils

import (
	"fmt"
	"math"
	"quizninja-api/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

// QuizToSummary converts Quiz model to QuizSummary DTO
func QuizToSummary(quiz *models.Quiz) models.QuizSummary {
	tags := make([]string, len(quiz.Tags))
	copy(tags, quiz.Tags)

	summary := models.QuizSummary{
		ID:            quiz.ID,
		Title:         quiz.Title,
		Description:   quiz.Description,
		Category:      quiz.Category,
		Difficulty:    quiz.Difficulty,
		TimeLimit:     quiz.TimeLimit,
		QuestionCount: quiz.QuestionCount,
		Points:        quiz.Points,
		IsFeatured:    quiz.IsFeatured,
		Tags:          tags,
		ThumbnailURL:  quiz.ThumbnailURL,
		CreatedAt:     quiz.CreatedAt,
	}

	if quiz.Statistics != nil {
		summary.Statistics = &models.QuizStatisticsSummary{
			TotalAttempts: quiz.Statistics.TotalAttempts,
			AverageScore:  quiz.Statistics.AverageScore,
			AverageTime:   quiz.Statistics.AverageTime,
		}
	}

	return summary
}

// QuizzesToSummaries converts slice of Quiz models to slice of QuizSummary DTOs
func QuizzesToSummaries(quizzes []models.Quiz) []models.QuizSummary {
	summaries := make([]models.QuizSummary, len(quizzes))
	for i, quiz := range quizzes {
		summaries[i] = QuizToSummary(&quiz)
	}
	return summaries
}

// CreateQuizRequestToQuiz converts CreateQuizRequest DTO to Quiz model
func CreateQuizRequestToQuiz(req *models.CreateQuizRequest, userID uuid.UUID) *models.Quiz {
	quiz := &models.Quiz{
		ID:            uuid.New(),
		Title:         req.Title,
		Description:   req.Description,
		Category:      req.Category,
		Difficulty:    req.Difficulty,
		TimeLimit:     req.TimeLimit,
		QuestionCount: len(req.Questions),
		IsFeatured:    req.IsFeatured,
		IsPublic:      req.IsPublic,
		CreatedBy:     userID,
		Tags:          models.StringArray(req.Tags),
		ThumbnailURL:  req.ThumbnailURL,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Convert questions
	quiz.Questions = make([]models.Question, len(req.Questions))
	for i, qReq := range req.Questions {
		quiz.Questions[i] = models.Question{
			ID:            uuid.New(),
			QuizID:        quiz.ID,
			QuestionText:  qReq.QuestionText,
			QuestionType:  qReq.QuestionType,
			Options:       models.StringArray(qReq.Options),
			CorrectAnswer: qReq.CorrectAnswer,
			Explanation:   qReq.Explanation,
			Points:        qReq.Points,
			Order:         qReq.Order,
			ImageURL:      qReq.ImageURL,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	}

	return quiz
}

// UpdateQuizFromRequest applies updates from UpdateQuizRequest to Quiz model
func UpdateQuizFromRequest(quiz *models.Quiz, req *models.UpdateQuizRequest) {
	if req.Title != nil {
		quiz.Title = *req.Title
	}
	if req.Description != nil {
		quiz.Description = *req.Description
	}
	if req.Category != nil {
		quiz.Category = *req.Category
	}
	if req.Difficulty != nil {
		quiz.Difficulty = *req.Difficulty
	}
	if req.TimeLimit != nil {
		quiz.TimeLimit = *req.TimeLimit
	}
	if req.IsFeatured != nil {
		quiz.IsFeatured = *req.IsFeatured
	}
	if req.IsPublic != nil {
		quiz.IsPublic = *req.IsPublic
	}
	if req.Tags != nil {
		quiz.Tags = models.StringArray(req.Tags)
	}
	if req.ThumbnailURL != nil {
		quiz.ThumbnailURL = req.ThumbnailURL
	}
	quiz.UpdatedAt = time.Now()
}

// BuildQuizListResponse creates a paginated response for quiz list
func BuildQuizListResponse(quizzes []models.Quiz, total, page, pageSize int) *models.QuizListResponse {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.QuizListResponse{
		Quizzes:    QuizzesToSummaries(quizzes),
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// ParseTagsFilter parses comma-separated tags string into slice
func ParseTagsFilter(tags string) []string {
	if tags == "" {
		return nil
	}

	tagsList := strings.Split(tags, ",")
	result := make([]string, 0, len(tagsList))

	for _, tag := range tagsList {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// ValidateQuestionRequest validates a question request based on its type
func ValidateQuestionRequest(req *models.CreateQuestionRequest) error {
	switch req.QuestionType {
	case "multiple_choice":
		if len(req.Options) < 2 {
			return fmt.Errorf("multiple choice questions must have at least 2 options")
		}
		if len(req.Options) > 6 {
			return fmt.Errorf("multiple choice questions can have at most 6 options")
		}
		// Check if correct answer is one of the options
		found := false
		for _, option := range req.Options {
			if option == req.CorrectAnswer {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("correct answer must be one of the provided options")
		}
	case "true_false":
		if req.CorrectAnswer != "true" && req.CorrectAnswer != "false" {
			return fmt.Errorf("true/false questions must have 'true' or 'false' as correct answer")
		}
	case "short_answer":
		// Short answer questions don't need options
		if len(req.Options) > 0 {
			return fmt.Errorf("short answer questions should not have options")
		}
	}

	return nil
}

// ValidateCreateQuizRequest validates the entire quiz creation request
func ValidateCreateQuizRequest(req *models.CreateQuizRequest) error {
	// Validate questions
	for i, question := range req.Questions {
		if err := ValidateQuestionRequest(&question); err != nil {
			return fmt.Errorf("question %d: %s", i+1, err.Error())
		}
	}

	// Check for duplicate question orders
	orders := make(map[int]bool)
	for i, question := range req.Questions {
		if orders[question.Order] {
			return fmt.Errorf("question %d: duplicate order value %d", i+1, question.Order)
		}
		orders[question.Order] = true
	}

	return nil
}
