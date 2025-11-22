package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	//  SECURITY CHECK: Verify access to private quizzes
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		if utils.HandleAuthError(c, err) {
			return
		}
	}

	if !quiz.IsPublic && quiz.CreatedBy != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied to private quiz")
		return
	}

	// Hide correct answers from questions unless user is the creator
	if quiz.CreatedBy != userID && len(quiz.Questions) > 0 {
		for i := range quiz.Questions {
			quiz.Questions[i].CorrectAnswer = ""
		}
	}

	response := &models.QuizDetailResponse{Quiz: *quiz}
	utils.SuccessResponse(c, response)
}

// GetQuizQuestions handles GET /api/v1/quizzes/{id}/questions
func (h *QuizHandler) GetQuizQuestions(c *gin.Context) {
	idParam := c.Param("id")
	quizID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+idParam)
		return
	}

	// Get quiz with questions
	quiz, err := h.repo.Quiz.GetQuizByIDWithQuestions(quizID)
	if err != nil {
		if err.Error() == "quiz not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Quiz not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	// SECURITY CHECK: Verify access to private quizzes
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		if utils.HandleAuthError(c, err) {
			return
		}
	}

	if !quiz.IsPublic && quiz.CreatedBy != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied to private quiz")
		return
	}

	// Hide correct answers from questions unless user is the creator
	if quiz.CreatedBy != userID && len(quiz.Questions) > 0 {
		for i := range quiz.Questions {
			quiz.Questions[i].CorrectAnswer = ""
		}
	}

	// Return only the questions array
	utils.SuccessResponse(c, quiz.Questions)
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

	quizzes, total, err := h.repo.Quiz.GetCompletedQuizzesByUser(userID, offset, pageSize)
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

	// Check if user has already completed this quiz (enforce one-attempt-per-quiz policy)
	hasCompleted, err := h.repo.Quiz.HasUserCompletedQuiz(userID, quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if hasCompleted {
		utils.ErrorResponse(c, http.StatusConflict, "This quiz has already been completed. Each quiz can only be attempted once.")
		return
	}

	// Auto-abandon any existing active attempt for this quiz
	existingAttempt, err := h.repo.Quiz.GetActiveQuizAttempt(userID, quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if existingAttempt != nil {
		// Abandon the existing attempt - no resume functionality
		err = h.repo.Quiz.AbandonQuizAttempt(existingAttempt.ID)
		if err != nil {
			fmt.Printf("Warning: Failed to abandon existing attempt %s: %v\n", existingAttempt.ID, err)
			// Continue anyway - we'll create a new attempt
		}
	}

	// Create quiz attempt
	attempt := &models.QuizAttempt{
		ID:              uuid.New(),
		QuizID:          quizID,
		UserID:          userID,
		Answers:         []models.AttemptAnswer{},
		Score:           0,
		TotalPoints:     0,
		TimeSpent:       0,
		PercentageScore: 0,
		Passed:          false,
		Status:          "started",
		IsCompleted:     false,
		StartedAt:       time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
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

	// Calculate time limit in seconds
	timeLimitSeconds := quiz.TimeLimit * 60

	response := gin.H{
		"id":             attempt.ID,
		"quiz":           quizWithQuestions,
		"started_at":     attempt.StartedAt,
		"time_limit":     timeLimitSeconds,
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
			QuestionId          string  `json:"questionId"`
			SelectedOption      *string `json:"selectedOption"`
			SelectedOptionIndex *int    `json:"selectedOptionIndex"`
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

	// Calculate score and build validated answers
	correctAnswers := 0
	totalQuestions := len(quiz.Questions)
	validatedAnswers := make([]models.AttemptAnswer, 0, len(submitRequest.Answers))

	// Create a map for easy question lookup
	questionMap := make(map[uuid.UUID]models.Question)
	for _, q := range quiz.Questions {
		questionMap[q.ID] = q
	}

	// Check each answer and build validated answers array
	for _, answer := range submitRequest.Answers {
		questionID, err := uuid.Parse(answer.QuestionId)
		if err != nil {
			continue // Skip invalid question IDs
		}

		question, exists := questionMap[questionID]
		if !exists {
			continue // Skip questions not found in quiz
		}

		// Determine the selected answer and check if it's correct
		var selectedOptionIndex int
		var selectedAnswerText string
		isCorrect := false
		pointsEarned := 0

		if answer.SelectedOptionIndex != nil {
			selectedOptionIndex = *answer.SelectedOptionIndex
			// Validate the index is within bounds
			if selectedOptionIndex >= 0 && selectedOptionIndex < len(question.Options) {
				selectedAnswerText = question.Options[selectedOptionIndex]
				if selectedAnswerText == question.CorrectAnswer {
					isCorrect = true
					pointsEarned = 1 // Default 1 point per correct answer
					correctAnswers++
				}
			}
		} else if answer.SelectedOption != nil {
			// Legacy support: if SelectedOption is provided as string
			selectedAnswerText = *answer.SelectedOption
			if selectedAnswerText == question.CorrectAnswer {
				isCorrect = true
				pointsEarned = 1
				correctAnswers++
			}
		}

		// Create validated answer with all fields
		validatedAnswer := models.AttemptAnswer{
			QuestionID:     questionID,
			SelectedOption: selectedOptionIndex,
			TextAnswer:     selectedAnswerText,
			IsCorrect:      isCorrect,
			PointsEarned:   pointsEarned,
		}
		validatedAnswers = append(validatedAnswers, validatedAnswer)
	}

	// Log warning if answer count doesn't match expected questions
	if len(validatedAnswers) != totalQuestions {
		fmt.Printf("[WARNING] Answer count mismatch for quiz %s: submitted=%d, expected=%d\n",
			quizID, len(validatedAnswers), totalQuestions)
	}

	// Calculate final score: 1 point per correct answer
	finalScore := float64(correctAnswers)
	scorePercentage := float64(correctAnswers) / float64(totalQuestions) * 100

	// Update the attempt with all required fields including answers
	timeNow := time.Now()
	attempt.Answers = validatedAnswers
	attempt.Score = finalScore
	attempt.TotalPoints = totalQuestions
	attempt.TimeSpent = submitRequest.TimeSpent
	attempt.PercentageScore = scorePercentage
	attempt.Passed = scorePercentage >= 60 // Passing threshold is 60%
	attempt.Status = "completed"
	attempt.IsCompleted = true
	attempt.CompletedAt = &timeNow
	attempt.UpdatedAt = timeNow

	err = h.repo.Quiz.UpdateQuizAttempt(attempt)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Update user statistics (total_quizzes_completed, average_score, etc.)
	err = h.repo.User.UpdateUserStatistics(userID, scorePercentage)
	if err != nil {
		// Log error but don't fail the request since the quiz is already completed
		fmt.Printf("Failed to update user statistics for user %s: %v\n", userID, err)
	}

	// Check for achievements after quiz completion
	var achievementNotifications []models.AchievementNotification
	// Note: Achievement checking is commented out for now to avoid dependencies
	// In a future implementation, we would check achievements here:
	// achievementService := services.NewAchievementService(h.repo)
	// result, err := achievementService.CheckAchievementsForUser(userID, services.TriggerQuizCompleted)
	// if err == nil && len(result.NewAchievements) > 0 {
	//     achievementNotifications = result.Notifications
	// }

	// Return the results including validated answers
	response := gin.H{
		"attempt_id":                attempt.ID,
		"score":                     int(finalScore), // Convert back to int for response
		"total_questions":           totalQuestions,
		"correct_answers":           correctAnswers,
		"time_spent":                submitRequest.TimeSpent,
		"answers":                   validatedAnswers,         // Include validated answers with isCorrect status
		"achievement_notifications": achievementNotifications, // Include achievement notifications
	}

	utils.SuccessResponse(c, response)
}

// GetUserAttempts handles GET /api/v1/users/attempts
func (h *QuizHandler) GetUserAttempts(c *gin.Context) {
	userID := getUserIDFromContext(c)
	fmt.Printf("[QuizHandler] GetUserAttempts called for user: %s\n", userID)

	var filters models.AttemptFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		fmt.Printf("[QuizHandler] Failed to bind query parameters: %v\n", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	fmt.Printf("[QuizHandler] Parsed filters: page=%d, pageSize=%d, category=%s, difficulty=%s, sortBy=%s, sortOrder=%s\n",
		filters.Page, filters.PageSize, filters.Category, filters.Difficulty, filters.SortBy, filters.SortOrder)

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

	// Set defaults for sorting
	if filters.SortBy == "" {
		filters.SortBy = "completed_at"
	}
	if filters.SortOrder == "" {
		filters.SortOrder = "desc"
	}

	attempts, total, err := h.repo.Quiz.GetUserAttempts(userID, &filters)
	if err != nil {
		fmt.Printf("[QuizHandler] Error getting user attempts: %v\n", err)
		utils.HandleError(c, err)
		return
	}

	fmt.Printf("[QuizHandler] Successfully retrieved %d attempts out of %d total for user %s\n", len(attempts), total, userID)

	// Calculate total pages
	totalPages := (total + filters.PageSize - 1) / filters.PageSize

	response := &models.AttemptHistoryResponse{
		Attempts:   attempts,
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	fmt.Printf("[QuizHandler] Returning attempt history response: total=%d, page=%d/%d, attempts_count=%d\n",
		response.Total, response.Page, response.TotalPages, len(response.Attempts))

	utils.SuccessResponse(c, response)
}

// GetAttemptDetails handles GET /api/v1/users/attempts/{attemptId}
func (h *QuizHandler) GetAttemptDetails(c *gin.Context) {
	attemptIdParam := c.Param("attemptId")
	fmt.Printf("[QuizHandler] GetAttemptDetails called for attempt ID: %s\n", attemptIdParam)

	attemptID, err := uuid.Parse(attemptIdParam)
	if err != nil {
		fmt.Printf("[QuizHandler] Invalid attempt ID format: %s, error: %v\n", attemptIdParam, err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid attempt ID format: "+attemptIdParam)
		return
	}

	userID := getUserIDFromContext(c)
	fmt.Printf("[QuizHandler] Getting attempt details for user: %s, attempt: %s\n", userID, attemptID)

	attempt, err := h.repo.Quiz.GetAttemptWithDetails(attemptID)
	if err != nil {
		fmt.Printf("[QuizHandler] Error getting attempt details: %v\n", err)
		if err.Error() == "attempt not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Attempt not found")
			return
		}
		utils.HandleError(c, err)
		return
	}

	// Verify the attempt belongs to the user
	if attempt.UserID != userID {
		fmt.Printf("[QuizHandler] Unauthorized access attempt: user %s tried to access attempt %s belonging to user %s\n",
			userID, attemptID, attempt.UserID)
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to view this attempt")
		return
	}

	fmt.Printf("[QuizHandler] Successfully retrieved attempt details for user %s: quiz='%s', score=%.2f, completed=%t\n",
		userID, attempt.Quiz.Title, attempt.Score, attempt.IsCompleted)

	response := &models.AttemptDetailResponse{
		Attempt: *attempt,
	}

	utils.SuccessResponse(c, response)
}

// GetUserQuizAttempt handles GET /api/v1/users/quizzes/{quizId}/attempt
func (h *QuizHandler) GetUserQuizAttempt(c *gin.Context) {
	quizIDParam := c.Param("quizId")
	quizID, err := uuid.Parse(quizIDParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid quiz ID format: "+quizIDParam)
		return
	}

	userID := getUserIDFromContext(c)

	attempt, err := h.repo.Quiz.GetUserQuizAttemptWithDetails(userID, quizID)
	if err != nil {
		if err.Error() == "no attempt found for this quiz" {
			utils.ErrorResponse(c, http.StatusNotFound, "No attempt found for this quiz")
			return
		}
		utils.HandleError(c, err)
		return
	}

	// Verify the attempt belongs to the user (double check)
	if attempt.UserID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to view this attempt")
		return
	}

	utils.SuccessResponse(c, gin.H{"attempt": attempt})
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

// UpdateQuizAttempt handles PUT /api/v1/quizzes/{id}/attempts/{attemptId}
func (h *QuizHandler) UpdateQuizAttempt(c *gin.Context) {
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
	var updateRequest models.UpdateAttemptRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
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
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to update this attempt")
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

	// Get quiz with questions to validate answers and calculate score
	quiz, err := h.repo.Quiz.GetQuizByIDWithQuestions(quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Validate answers and calculate score
	validatedAnswers, correctAnswers, totalQuestions, err := h.validateAndScoreAnswers(updateRequest.Answers, quiz.Questions)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid answers", err.Error())
		return
	}

	// Calculate score: 1 point per correct answer
	finalScore := float64(correctAnswers)
	scorePercentage := float64(correctAnswers) / float64(totalQuestions) * 100
	passed := scorePercentage >= 60.0 // 60% passing threshold

	// Update the attempt
	attempt.Answers = validatedAnswers
	attempt.Score = finalScore
	attempt.TotalPoints = totalQuestions
	attempt.TimeSpent = updateRequest.TimeSpent
	attempt.PercentageScore = scorePercentage
	attempt.Passed = passed
	attempt.Status = updateRequest.Status

	if updateRequest.Status == "completed" {
		attempt.IsCompleted = true
		now := time.Now()
		attempt.CompletedAt = &now
	}

	attempt.UpdatedAt = time.Now()

	err = h.repo.Quiz.UpdateQuizAttempt(attempt)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Update quiz statistics if attempt is completed
	if updateRequest.Status == "completed" {
		err = h.updateQuizStatistics(quizID, finalScore, updateRequest.TimeSpent)
		if err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to update quiz statistics: %v\n", err)
		}
	}

	// Return the updated attempt
	response := gin.H{
		"attempt_id":       attempt.ID,
		"answers":          validatedAnswers,
		"score":            int(finalScore),
		"percentage_score": scorePercentage,
		"passed":           passed,
		"total_questions":  totalQuestions,
		"correct_answers":  correctAnswers,
		"time_spent":       updateRequest.TimeSpent,
		"status":           updateRequest.Status,
		"completed_at":     attempt.CompletedAt,
	}

	utils.SuccessResponse(c, response)
}

// AbandonQuizAttempt handles DELETE /api/v1/quizzes/{id}/attempts/{attemptId}/abandon
func (h *QuizHandler) AbandonQuizAttempt(c *gin.Context) {
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
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to abandon this attempt")
		return
	}

	if attempt.QuizID != quizID {
		utils.ErrorResponse(c, http.StatusBadRequest, "Attempt does not belong to specified quiz")
		return
	}

	if attempt.IsCompleted {
		utils.ErrorResponse(c, http.StatusConflict, "Quiz attempt already completed, cannot abandon")
		return
	}

	// Abandon the attempt
	err = h.repo.Quiz.AbandonQuizAttempt(attemptID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"message":    "Quiz attempt abandoned successfully",
		"attempt_id": attemptID,
	})
}

// validateAndScoreAnswers validates user answers against quiz questions and calculates score
func (h *QuizHandler) validateAndScoreAnswers(answers []models.AttemptAnswer, questions []models.Question) ([]models.AttemptAnswer, int, int, error) {
	// Create a map for easy question lookup
	questionMap := make(map[uuid.UUID]models.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	totalQuestions := len(questions)
	correctAnswers := 0
	validatedAnswers := make([]models.AttemptAnswer, 0, len(answers))

	// Validate and score each answer
	for _, answer := range answers {
		question, exists := questionMap[answer.QuestionID]
		if !exists {
			return nil, 0, 0, fmt.Errorf("invalid question ID: %s", answer.QuestionID)
		}

		// Score the answer
		isCorrect := false
		pointsEarned := 0

		// Check if answer is correct based on question type
		if question.QuestionType == "multiple_choice" || question.QuestionType == "multipleChoice" {
			// For multiple choice, check selected option against correct answer
			if answer.SelectedOption >= 0 && answer.SelectedOption < len(question.Options) {
				selectedAnswer := question.Options[answer.SelectedOption]
				if selectedAnswer == question.CorrectAnswer {
					isCorrect = true
					pointsEarned = 1 // Default 1 point per correct answer
				}
			}
		} else {
			// For other question types, check text answer
			if strings.TrimSpace(strings.ToLower(answer.TextAnswer)) == strings.TrimSpace(strings.ToLower(question.CorrectAnswer)) {
				isCorrect = true
				pointsEarned = 1
			}
		}

		if isCorrect {
			correctAnswers++
		}

		// Create validated answer
		validatedAnswer := models.AttemptAnswer{
			QuestionID:     answer.QuestionID,
			SelectedOption: answer.SelectedOption,
			TextAnswer:     answer.TextAnswer,
			IsCorrect:      isCorrect,
			PointsEarned:   pointsEarned,
		}

		validatedAnswers = append(validatedAnswers, validatedAnswer)
	}

	return validatedAnswers, correctAnswers, totalQuestions, nil
}

// updateQuizStatistics updates quiz statistics when an attempt is completed
func (h *QuizHandler) updateQuizStatistics(quizID uuid.UUID, score float64, timeSpent int) error {
	// Use the repository method to create or update quiz statistics
	err := h.repo.Quiz.CreateOrUpdateQuizStatistics(quizID, score, timeSpent)
	if err != nil {
		return fmt.Errorf("failed to update quiz statistics: %w", err)
	}

	fmt.Printf("Successfully updated quiz statistics for quiz %s: score=%.2f, time=%d\n",
		quizID, score, timeSpent)

	return nil
}

