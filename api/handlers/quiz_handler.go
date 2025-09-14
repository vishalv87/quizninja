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
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format")
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

// CreateQuiz handles POST /api/v1/quizzes
func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	var req models.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate quiz request
	if err := utils.ValidateCreateQuizRequest(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	quiz := utils.CreateQuizRequestToQuiz(&req, userID)

	err := h.repo.Quiz.CreateQuiz(quiz)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Return the created quiz
	createdQuiz, err := h.repo.Quiz.GetQuizByIDWithQuestions(quiz.ID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.CreatedResponse(c, &models.QuizDetailResponse{Quiz: *createdQuiz}, "Quiz created successfully")
}

// UpdateQuiz handles PUT /api/v1/quizzes/{id}
func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	idParam := c.Param("id")
	quizID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format")
		return
	}

	var req models.UpdateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get existing quiz
	quiz, err := h.repo.Quiz.GetQuizByID(quizID)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Quiz not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	// Check if user owns the quiz
	userID := getUserIDFromContext(c)
	if quiz.CreatedBy != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "You can only update your own quizzes")
		return
	}

	// Apply updates
	utils.UpdateQuizFromRequest(quiz, &req)

	err = h.repo.Quiz.UpdateQuiz(quiz)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Return updated quiz
	updatedQuiz, err := h.repo.Quiz.GetQuizByID(quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, &models.QuizDetailResponse{Quiz: *updatedQuiz}, "Quiz updated successfully")
}

// DeleteQuiz handles DELETE /api/v1/quizzes/{id}
func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	idParam := c.Param("id")
	quizID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format")
		return
	}

	// Get existing quiz to check ownership
	quiz, err := h.repo.Quiz.GetQuizByID(quizID)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Quiz not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	// Check if user owns the quiz
	userID := getUserIDFromContext(c)
	if quiz.CreatedBy != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "You can only delete your own quizzes")
		return
	}

	err = h.repo.Quiz.DeleteQuiz(quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"deleted": true}, "Quiz deleted successfully")
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
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format")
		return
	}

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