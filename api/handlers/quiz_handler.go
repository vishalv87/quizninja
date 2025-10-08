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

	// Check for retake request in the request body
	var retakeRequest models.CreateRetakeRequest
	isRetakeRequest := false
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&retakeRequest); err == nil && retakeRequest.IsRetake {
			isRetakeRequest = true
		}
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

	// Check if this is a test user to properly set IsTestData fields
	isTestUser, err := h.repo.User.IsTestUser(userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

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

	// Handle retake logic
	var originalAttemptID *uuid.UUID
	var retakeCount int
	var performanceComparison map[string]interface{}

	if isRetakeRequest {
		// Parse and validate original attempt ID
		originalID, err := uuid.Parse(retakeRequest.OriginalAttemptID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid original attempt ID format")
			return
		}
		originalAttemptID = &originalID

		// Verify the original attempt exists and belongs to the user
		originalAttempt, err := h.repo.Quiz.GetQuizAttempt(originalID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Original attempt not found")
			return
		}

		if originalAttempt.UserID != userID {
			utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to retake this attempt")
			return
		}

		if originalAttempt.QuizID != quizID {
			utils.ErrorResponse(c, http.StatusBadRequest, "Original attempt does not belong to this quiz")
			return
		}

		// Validate retake limit
		err = h.repo.Quiz.ValidateRetakeLimit(userID, quizID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
			return
		}

		// Get previous attempts for performance comparison
		previousAttempts, err := h.repo.Quiz.GetQuizAttemptsForComparison(userID, quizID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		// Set retake count (increment from the maximum found)
		retakeCount = originalAttempt.RetakeCount + 1
		for _, attempt := range previousAttempts {
			if attempt.RetakeCount >= retakeCount {
				retakeCount = attempt.RetakeCount + 1
			}
		}

		// Calculate performance comparison for completed attempts
		if len(previousAttempts) > 0 {
			performanceComparison = h.repo.Quiz.CalculatePerformanceComparison(
				&models.QuizAttempt{PercentageScore: 0, TimeSpent: 0}, // Placeholder for new attempt
				previousAttempts,
			)
		}
	}

	// Create quiz attempt
	attempt := &models.QuizAttempt{
		ID:                    uuid.New(),
		QuizID:                quizID,
		UserID:                userID,
		Answers:               []models.AttemptAnswer{},
		Score:                 0,
		TotalPoints:           0,
		TimeSpent:             0,
		PercentageScore:       0,
		Passed:                false,
		Status:                "started",
		IsCompleted:           false,
		StartedAt:             time.Now(),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		RetakeCount:           retakeCount,
		OriginalAttemptID:     originalAttemptID,
		PerformanceComparison: performanceComparison,
		IsTestData:            isTestUser,
	}

	err = h.repo.Quiz.CreateQuizAttempt(attempt)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Create a quiz session for the attempt
	timeLimit := quiz.TimeLimit * 60 // Convert minutes to seconds
	session := &models.QuizSession{
		ID:                   uuid.New(),
		AttemptID:            attempt.ID,
		UserID:               userID,
		QuizID:               quizID,
		CurrentQuestionIndex: 0,
		CurrentAnswers:       []models.AttemptAnswer{},
		SessionState:         "active",
		TimeRemaining:        &timeLimit,
		TimeSpentSoFar:       0,
		LastActivityAt:       time.Now(),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		IsTestData:           isTestUser,
	}

	err = h.repo.QuizSession.CreateSession(session)
	if err != nil {
		// If session creation fails, we should clean up the attempt
		_ = h.repo.Quiz.DeleteActiveQuizAttempt(userID, quizID)
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
		"attempt_id":             attempt.ID,
		"session_id":             session.ID,
		"quiz":                   quizWithQuestions,
		"started_at":             attempt.StartedAt,
		"time_remaining":         session.TimeRemaining,
		"is_retake":              attempt.IsRetake(),
		"retake_count":           attempt.RetakeCount,
		"original_attempt_id":    attempt.OriginalAttemptID,
		"performance_comparison": attempt.PerformanceComparison,
	}

	var message string
	if attempt.IsRetake() {
		message = fmt.Sprintf("Quiz retake #%d started successfully", attempt.RetakeCount)
	} else {
		message = "Quiz attempt started successfully"
	}

	utils.CreatedResponse(c, response, message)
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

	// Complete the quiz session
	err = h.repo.QuizSession.CompleteSession(attemptID)
	if err != nil {
		// Log error but don't fail the request since the attempt is already completed
		fmt.Printf("Failed to complete quiz session for attempt %s: %v\n", attemptID, err)
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

	// Return the results
	response := gin.H{
		"attempt_id":                attempt.ID,
		"score":                     int(finalScore), // Convert back to int for response
		"total_questions":           totalQuestions,
		"correct_answers":           correctAnswers,
		"time_spent":                submitRequest.TimeSpent,
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

	// Calculate score and percentage
	basePoints := totalQuestions // 1 point per question as default
	scorePercentage := float64(correctAnswers) / float64(totalQuestions) * 100
	finalScore := scorePercentage * float64(basePoints) / 100
	passed := scorePercentage >= 60.0 // 60% passing threshold

	// Update the attempt
	attempt.Answers = validatedAnswers
	attempt.Score = finalScore
	attempt.TotalPoints = basePoints
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

// Quiz Session Management Handlers

// PauseQuizSession handles POST /api/v1/quizzes/{id}/attempts/{attemptId}/pause
func (h *QuizHandler) PauseQuizSession(c *gin.Context) {
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
	fmt.Printf("[DEBUG] PauseQuizSession: quizID=%s, attemptID=%s, userID=%s\n", quizID, attemptID, userID)

	var pauseRequest models.PauseSessionRequest
	if err := c.ShouldBindJSON(&pauseRequest); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}
	fmt.Printf("[DEBUG] PauseQuizSession: Received request - QuestionIndex=%d, AnswersCount=%d, TimeSpent=%d, TimeRemaining=%v\n",
		pauseRequest.CurrentQuestionIndex, len(pauseRequest.CurrentAnswers),
		pauseRequest.TimeSpentSoFar, pauseRequest.TimeRemaining)

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
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to pause this attempt")
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

	// Pause the session
	err = h.repo.QuizSession.PauseSession(attemptID, &pauseRequest)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Get the updated session details
	session, err := h.repo.QuizSession.GetSessionByAttemptID(attemptID)
	if err != nil {
		fmt.Printf("[ERROR] PauseQuizSession: Failed to get session - %v\n", err)
		utils.HandleError(c, err)
		return
	}
	fmt.Printf("[DEBUG] PauseQuizSession: Retrieved session - ID=%s, State=%s, TimeRemaining=%v\n",
		session.ID, session.SessionState, session.TimeRemaining)

	response := models.SessionActionResponse{
		SessionID:     session.ID,
		Action:        "paused",
		SessionState:  session.SessionState,
		Message:       "Quiz session paused successfully",
		TimeRemaining: session.TimeRemaining,
		Progress:      session.GetProgress(),
	}

	utils.SuccessResponse(c, response)
}

// ResumeQuizSession handles POST /api/v1/quizzes/{id}/attempts/{attemptId}/resume
func (h *QuizHandler) ResumeQuizSession(c *gin.Context) {
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
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to resume this attempt")
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

	// Check if session can be resumed
	canResume, err := h.repo.QuizSession.CanResumeSession(attemptID, userID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if !canResume {
		utils.ErrorResponse(c, http.StatusConflict, "Session cannot be resumed (not paused, already completed, or abandoned)")
		return
	}

	// Resume the session
	err = h.repo.QuizSession.ResumeSession(attemptID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Get the updated session details with quiz information
	sessionDetails, err := h.repo.QuizSession.GetSessionByAttemptID(attemptID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Get quiz with questions for resuming
	quiz, err := h.repo.Quiz.GetQuizByIDWithQuestions(quizID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Remove correct answers from questions
	for i := range quiz.Questions {
		quiz.Questions[i].CorrectAnswer = ""
	}

	response := gin.H{
		"session_id":             sessionDetails.ID,
		"action":                 "resumed",
		"session_state":          sessionDetails.SessionState,
		"message":                "Quiz session resumed successfully",
		"quiz":                   quiz,
		"current_question_index": sessionDetails.CurrentQuestionIndex,
		"current_answers":        sessionDetails.CurrentAnswers,
		"time_remaining":         sessionDetails.TimeRemaining,
		"time_spent_so_far":      sessionDetails.TimeSpentSoFar,
		"progress":               sessionDetails.GetProgress(),
	}

	utils.SuccessResponse(c, response)
}

// GetUserActiveSessions handles GET /api/v1/users/active-sessions
func (h *QuizHandler) GetUserActiveSessions(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var filters models.SessionFilters
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
	if filters.PageSize > 50 {
		filters.PageSize = 50
	}

	// Set defaults for sorting
	if filters.SortBy == "" {
		filters.SortBy = "last_activity_at"
	}
	if filters.SortOrder == "" {
		filters.SortOrder = "desc"
	}

	sessions, total, err := h.repo.QuizSession.GetUserActiveSessions(userID, &filters)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Count active and paused sessions
	activeCount := 0
	pausedCount := 0
	for _, session := range sessions {
		if session.SessionState == "active" {
			activeCount++
		} else if session.SessionState == "paused" {
			pausedCount++
		}
	}

	response := models.ActiveSessionsResponse{
		Sessions:    sessions,
		Total:       total,
		ActiveCount: activeCount,
		PausedCount: pausedCount,
	}

	utils.SuccessResponse(c, response)
}

// SaveQuizProgress handles PUT /api/v1/quizzes/{id}/attempts/{attemptId}/save-progress
func (h *QuizHandler) SaveQuizProgress(c *gin.Context) {
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

	var updateRequest models.UpdateQuizSessionRequest
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
		utils.ErrorResponse(c, http.StatusForbidden, "Not authorized to save progress for this attempt")
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

	// Save the progress
	err = h.repo.QuizSession.SaveSessionProgress(attemptID, &updateRequest)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Get the session for the response
	session, err := h.repo.QuizSession.GetSessionByAttemptID(attemptID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	response := gin.H{
		"session_id": session.ID,
		"message":    "Progress saved successfully",
		"timestamp":  time.Now(),
	}

	utils.SuccessResponse(c, response)
}

// AbandonQuizSession handles DELETE /api/v1/quizzes/{id}/attempts/{attemptId}/abandon
func (h *QuizHandler) AbandonQuizSession(c *gin.Context) {
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
		utils.ErrorResponse(c, http.StatusConflict, "Quiz attempt already completed")
		return
	}

	// Abandon the session
	err = h.repo.QuizSession.AbandonSession(attemptID)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	response := models.SessionActionResponse{
		SessionID:    uuid.Nil, // Session is abandoned
		Action:       "abandoned",
		SessionState: "abandoned",
		Message:      "Quiz session abandoned successfully",
		Progress:     0.0,
	}

	utils.SuccessResponse(c, response)
}
