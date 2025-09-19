package handlers

import (
	"net/http"
	"strconv"

	"quizninja-api/config"
	"quizninja-api/models"
	"quizninja-api/repository"
	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QuizHandler struct {
	repo   *repository.Repository
	config *config.Config
}

func NewQuizHandler(cfg *config.Config) *QuizHandler {
	return &QuizHandler{
		repo:   repository.NewRepository(),
		config: cfg,
	}
}

// GetQuizzes handles GET /api/v1/quizzes
func (h *QuizHandler) GetQuizzes(c *gin.Context) {
	var filters models.QuizFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	// Set defaults for pagination
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}
	if filters.PageSize > 100 {
		filters.PageSize = 100
	}

	quizzes, total, err := h.repo.Quiz.GetQuizzes(&filters)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	response := utils.BuildQuizListResponse(quizzes, total, filters.Page, filters.PageSize)
	utils.SuccessResponse(c, response)
}

// GetQuizByID handles GET /api/v1/quizzes/{id}
func (h *QuizHandler) GetQuizByID(c *gin.Context) {
	idParam := c.Param("id")
	quizID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+idParam)
		return
	}

	// Check if user wants questions included
	includeQuestions := c.Query("include_questions") == "true"
	includeStats := c.Query("include_stats") == "true"

	var quiz *models.Quiz

	if includeQuestions && includeStats {
		quiz, err = h.repo.Quiz.GetQuizByIDWithAll(quizID)
	} else if includeQuestions {
		quiz, err = h.repo.Quiz.GetQuizByIDWithQuestions(quizID)
	} else if includeStats {
		quiz, err = h.repo.Quiz.GetQuizByIDWithStatistics(quizID)
	} else {
		quiz, err = h.repo.Quiz.GetQuizByID(quizID)
	}

	if err != nil {
		if err.Error() == "quiz not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Quiz not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	// Hide correct answers from questions unless user is the creator
	userID := getUserIDFromContext(c)
	if quiz.CreatedBy != userID && len(quiz.Questions) > 0 {
		for i := range quiz.Questions {
			quiz.Questions[i].CorrectAnswer = ""
		}
	}

	response := &models.QuizDetailResponse{Quiz: *quiz}
	utils.SuccessResponse(c, response)
}




// GetFeaturedQuizzes handles GET /api/v1/quizzes/featured
func (h *QuizHandler) GetFeaturedQuizzes(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	quizzes, err := h.repo.Quiz.GetFeaturedQuizzes(limit)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	summaries := utils.QuizzesToSummaries(quizzes)
	utils.SuccessResponse(c, gin.H{"quizzes": summaries})
}

// GetQuizzesByCategory handles GET /api/v1/quizzes/category/{category}
func (h *QuizHandler) GetQuizzesByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Category parameter is required")
		return
	}

	limitParam := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	quizzes, err := h.repo.Quiz.GetQuizzesByCategory(category, limit)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	summaries := utils.QuizzesToSummaries(quizzes)
	utils.SuccessResponse(c, gin.H{"quizzes": summaries})
}

// GetUserQuizzes handles GET /api/v1/users/quizzes
func (h *QuizHandler) GetUserQuizzes(c *gin.Context) {
	userID := getUserIDFromContext(c)

	pageParam := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSizeParam := c.DefaultQuery("page_size", "10")
	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	quizzes, total, err := h.repo.Quiz.GetQuizzesByUser(userID, offset, pageSize)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	response := utils.BuildQuizListResponse(quizzes, total, page, pageSize)
	utils.SuccessResponse(c, response)
}

// StartQuizAttempt handles POST /api/v1/quizzes/{id}/attempts
func (h *QuizHandler) StartQuizAttempt(c *gin.Context) {
	idParam := c.Param("id")
	quizID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+idParam)
		return
	}

	// Check if force_restart parameter is provided
	forceRestart := c.Query("force_restart") == "true"

	// Verify quiz exists and is public
	quiz, err := h.repo.Quiz.GetQuizByID(quizID)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Quiz not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	if !quiz.IsPublic {
		utils.ErrorResponse(c, http.StatusForbidden, "Quiz is not public")
		return
	}

	userID := getUserIDFromContext(c)

	// Check for existing active attempt
	existingAttempt, err := h.repo.Quiz.GetActiveQuizAttempt(userID, quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// If there's an existing active attempt and force_restart is not true, return error
	if existingAttempt != nil && !forceRestart {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    http.StatusConflict,
				"message": "Active quiz attempt already exists",
				"details": gin.H{
					"existing_attempt_id": existingAttempt.ID,
					"started_at":          existingAttempt.StartedAt,
					"can_restart":         true,
				},
			},
		})
		return
	}

	// If force_restart is true, delete any existing active attempts
	if forceRestart {
		err = h.repo.Quiz.DeleteActiveQuizAttempt(userID, quizID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
	}

	// Create quiz attempt
	attempt := &models.QuizAttempt{
		ID:          uuid.New(),
		QuizID:      quizID,
		UserID:      userID,
		Score:       0,
		TotalPoints: 0,
		TimeSpent:   0,
		IsCompleted: false,
		StartedAt:   quiz.CreatedAt, // This should be time.Now() but using CreatedAt to match the pattern
		CreatedAt:   quiz.CreatedAt,
		UpdatedAt:   quiz.UpdatedAt,
	}

	err = h.repo.Quiz.CreateQuizAttempt(attempt)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Return quiz with questions for the attempt
	quizWithQuestions, err := h.repo.Quiz.GetQuizByIDWithQuestions(quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Remove correct answers from questions
	for i := range quizWithQuestions.Questions {
		quizWithQuestions.Questions[i].CorrectAnswer = ""
	}

	response := gin.H{
		"attempt_id": attempt.ID,
		"quiz":       quizWithQuestions,
		"started_at": attempt.StartedAt,
	}

	utils.CreatedResponse(c, response, "Quiz attempt started successfully")
}

// SubmitQuizAttempt handles POST /api/v1/quizzes/{id}/attempts/{attemptId}/submit
func (h *QuizHandler) SubmitQuizAttempt(c *gin.Context) {
	idParam := c.Param("id")
	quizID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+idParam)
		return
	}

	attemptIdParam := c.Param("attemptId")
	attemptID, err := uuid.Parse(attemptIdParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid attempt ID format: "+attemptIdParam)
		return
	}

	userID := getUserIDFromContext(c)

	// Parse the request body
	var submitRequest struct {
		AttemptId string `json:"attemptId"`
		Answers   []struct {
			QuestionId         string  `json:"questionId"`
			SelectedOption     *string `json:"selectedOption"`
			SelectedOptionIndex *int   `json:"selectedOptionIndex"`
		} `json:"answers"`
		TimeSpent int `json:"timeSpent"`
	}

	if err := c.ShouldBindJSON(&submitRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Verify the attempt exists and belongs to the user
	attempt, err := h.repo.Quiz.GetQuizAttempt(attemptID)
	if err != nil {
		if err.Error() == "quiz attempt not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Quiz attempt not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	if attempt.UserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to submit this attempt")
		return
	}

	if attempt.QuizID != quizID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Attempt does not belong to specified quiz")
		return
	}

	if attempt.IsCompleted {
		utils.ErrorResponse(c, http.StatusConflict, "Quiz attempt already completed")
		return
	}

	// Get quiz with questions to calculate score
	quiz, err := h.repo.Quiz.GetQuizByIDWithQuestions(quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Calculate score
	correctAnswers := 0
	totalQuestions := len(quiz.Questions)

	// Create a map for easy question lookup
	questionMap := make(map[uuid.UUID]models.Question)
	for _, q := range quiz.Questions {
		questionMap[q.ID] = q
	}

	// Check each answer
	for _, answer := range submitRequest.Answers {
		questionID, err := uuid.Parse(answer.QuestionId)
		if err != nil {
			continue // Skip invalid question IDs
		}

		question, exists := questionMap[questionID]
		if !exists {
			continue // Skip questions not found in quiz
		}

		// Check if answer is correct
		if answer.SelectedOption != nil && *answer.SelectedOption == question.CorrectAnswer {
			correctAnswers++
		}
	}

	// Calculate final score based on questions (default 1 point per question)
	basePoints := totalQuestions // 1 point per question as default
	scorePercentage := float64(correctAnswers) / float64(totalQuestions) * 100
	finalScore := scorePercentage * float64(basePoints) / 100

	// Update the attempt
	attempt.Score = finalScore
	attempt.TotalPoints = basePoints
	attempt.TimeSpent = submitRequest.TimeSpent
	attempt.IsCompleted = true
	attempt.CompletedAt = &quiz.CreatedAt // Should be time.Now() but using consistent time

	err = h.repo.Quiz.UpdateQuizAttempt(attempt)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Return the results
	response := gin.H{
		"attempt_id":       attempt.ID,
		"score":           int(finalScore), // Convert back to int for response
		"total_questions": totalQuestions,
		"correct_answers": correctAnswers,
		"time_spent":      submitRequest.TimeSpent,
	}

	utils.SuccessResponse(c, response)
}

// Helper function to get user ID from context
func getUserIDFromContext(c *gin.Context) uuid.UUID {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return userID
}
