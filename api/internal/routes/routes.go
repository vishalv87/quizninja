package routes

import (
	"quizninja-api/config"
	"quizninja-api/internal/handlers"
	"quizninja-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupInternalRoutes configures the internal API routes
func SetupInternalRoutes(r *gin.Engine, cfg *config.Config) {
	// Create handlers
	attemptHandler := handlers.NewAttemptHandler(cfg)
	quizHandler := handlers.NewQuizInternalHandler(cfg)
	scoringHandler := handlers.NewScoringHandler()
	statisticsHandler := handlers.NewStatisticsHandler(cfg)
	achievementHandler := handlers.NewAchievementInternalHandler(cfg)

	// Internal API routes group
	internal := r.Group("/internal/v1")
	internal.Use(middleware.InternalAuthMiddleware())
	{
		// Attempt endpoints
		attempts := internal.Group("/attempts")
		{
			attempts.POST("/:attemptId/validate", attemptHandler.ValidateAttempt)
			attempts.PUT("/:attemptId", attemptHandler.UpdateAttempt)
		}

		// Quiz endpoints (internal - includes correct answers)
		quizzes := internal.Group("/quizzes")
		{
			quizzes.GET("/:quizId/questions", quizHandler.GetQuestionsInternal)
		}

		// Scoring endpoint
		scoring := internal.Group("/scoring")
		{
			scoring.POST("/calculate", scoringHandler.CalculateScore)
		}

		// User statistics endpoints
		users := internal.Group("/users")
		{
			users.POST("/:userId/statistics", statisticsHandler.UpdateStatistics)
			users.POST("/:userId/achievements/check", achievementHandler.CheckAchievements)
		}
	}
}
