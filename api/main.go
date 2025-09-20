package main

import (
	"log"
	"net/http"

	"quizninja-api/config"
	"quizninja-api/database"
	"quizninja-api/handlers"
	"quizninja-api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)

	database.Connect(cfg)
	defer database.Close()

	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	setupRoutes(r, cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(r *gin.Engine, cfg *config.Config) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "QuizNinja API is running",
		})
	})

	authHandler := handlers.NewAuthHandler(cfg)
	categoriesHandler := handlers.NewCategoriesHandler(cfg)
	preferencesHandler := handlers.NewPreferencesHandler(cfg)
	appSettingsHandler := handlers.NewAppSettingsHandler(cfg)
	quizHandler := handlers.NewQuizHandler(cfg)
	friendsHandler := handlers.NewFriendsHandler(cfg)
	challengesHandler := handlers.NewChallengesHandler(cfg)
	leaderboardHandler := handlers.NewLeaderboardHandler(cfg)
	achievementHandler := handlers.NewAchievementHandler(cfg)

	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Public endpoints
		quizzes := api.Group("/quizzes")
		{
			quizzes.GET("", quizHandler.GetQuizzes)
			quizzes.GET("/:id", quizHandler.GetQuizByID)
			quizzes.GET("/featured", quizHandler.GetFeaturedQuizzes)
			quizzes.GET("/category/:category", quizHandler.GetQuizzesByCategory)
			quizzes.GET("/categories", categoriesHandler.GetCategories)
		}

		config := api.Group("/config")
		{
			config.GET("/app-settings", appSettingsHandler.GetAppSettings)
		}

		// Preferences endpoints (public)
		preferences := api.Group("/preferences")
		{
			preferences.GET("/categories", preferencesHandler.GetCategories)
			preferences.GET("/difficulty-levels", preferencesHandler.GetDifficultyLevels)
			preferences.GET("/notification-frequencies", preferencesHandler.GetNotificationFrequencies)
		}

		// Note: Leaderboard endpoints moved to protected section for authentication

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// Protected endpoints
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			protected.GET("/profile", authHandler.GetProfile)
			protected.PUT("/profile", authHandler.UpdateProfile)

			users := protected.Group("/users")
			{
				users.PUT("/preferences", preferencesHandler.UpdatePreferences)
				users.GET("/preferences", preferencesHandler.GetPreferences)
				users.GET("/quizzes", quizHandler.GetUserQuizzes)
				users.GET("/stats", authHandler.GetUserStats)
				users.GET("/attempts", quizHandler.GetUserAttempts)
				users.GET("/attempts/:attemptId", quizHandler.GetAttemptDetails)
			}

			// Protected quiz endpoints
			protectedQuizzes := protected.Group("/quizzes")
			{
				protectedQuizzes.POST("/:id/attempts", quizHandler.StartQuizAttempt)
				protectedQuizzes.POST("/:id/attempts/:attemptId/submit", quizHandler.SubmitQuizAttempt)
			}

			// Friends endpoints
			friends := protected.Group("/friends")
			{
				friends.POST("/requests", friendsHandler.SendFriendRequest)
				friends.GET("/requests", friendsHandler.GetFriendRequests)
				friends.PUT("/requests/:id", friendsHandler.RespondToFriendRequest)
				friends.DELETE("/requests/:id", friendsHandler.CancelFriendRequest)
				friends.GET("", friendsHandler.GetFriends)
				friends.DELETE("/:id", friendsHandler.RemoveFriend)
				friends.GET("/search", friendsHandler.SearchUsers)
				friends.GET("/notifications", friendsHandler.GetFriendNotifications)
				friends.PUT("/notifications/:id/read", friendsHandler.MarkNotificationAsRead)
				friends.PUT("/notifications/read-all", friendsHandler.MarkAllNotificationsAsRead)
			}

			// Challenge endpoints
			challenges := protected.Group("/challenges")
			{
				challenges.POST("", challengesHandler.CreateChallenge)
				challenges.GET("", challengesHandler.GetChallenges)
				challenges.GET("/stats", challengesHandler.GetChallengeStats)
				challenges.GET("/pending", challengesHandler.GetPendingChallenges)
				challenges.GET("/active", challengesHandler.GetActiveChallenges)
				challenges.GET("/completed", challengesHandler.GetCompletedChallenges)
				challenges.GET("/:id", challengesHandler.GetChallengeByID)
				challenges.PUT("/:id/accept", challengesHandler.AcceptChallenge)
				challenges.PUT("/:id/decline", challengesHandler.DeclineChallenge)
				challenges.PUT("/:id/score", challengesHandler.UpdateChallengeScore)
				challenges.POST("/expire", challengesHandler.ExpireChallenges) // Admin endpoint
			}

			// All leaderboard endpoints (require authentication)
			leaderboard := protected.Group("/leaderboard")
			{
				leaderboard.GET("", leaderboardHandler.GetLeaderboard)
				leaderboard.GET("/stats", leaderboardHandler.GetLeaderboardStats)
				leaderboard.GET("/rank", leaderboardHandler.GetUserRank)
				leaderboard.POST("/score", leaderboardHandler.UpdateUserScore)
				leaderboard.GET("/achievements", achievementHandler.GetLeaderboardWithAchievements)
			}

			// Achievement endpoints
			achievements := protected.Group("/achievements")
			{
				achievements.GET("", achievementHandler.GetAllAchievements)
				achievements.GET("/progress", achievementHandler.GetAchievementProgress)
				achievements.GET("/stats", achievementHandler.GetAchievementStats)
				achievements.POST("/check", achievementHandler.CheckAchievements)
				achievements.GET("/category/:category", achievementHandler.GetAchievementsByCategory)
				achievements.POST("/unlock/:key", achievementHandler.UnlockAchievement) // Admin/testing endpoint
			}

			// User achievement endpoints
			users.GET("/achievements", achievementHandler.GetUserAchievements)
			users.GET("/:userId/achievements", achievementHandler.GetUserAchievementsByUserID)

			// Admin endpoints for cache management
			admin := protected.Group("/admin")
			{
				admin.DELETE("/cache/app-settings", appSettingsHandler.ClearCache)
			}
		}
	}
}
