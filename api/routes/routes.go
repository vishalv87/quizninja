package routes

import (
	"net/http"

	"quizninja-api/config"
	"quizninja-api/handlers"
	"quizninja-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
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
	userHandler := handlers.NewUserHandler()
	friendsHandler := handlers.NewFriendsHandler(cfg)
	challengesHandler := handlers.NewChallengesHandler(cfg)
	leaderboardHandler := handlers.NewLeaderboardHandler(cfg)
	achievementHandler := handlers.NewAchievementHandler(cfg)
	favoritesHandler := handlers.NewFavoritesHandler(cfg)
	discussionHandler := handlers.NewDiscussionHandler(cfg)
	notificationHandler := handlers.NewNotificationHandler(cfg)
	ratingHandler := handlers.NewRatingHandler()

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
			quizzes.GET("/featured", quizHandler.GetFeaturedQuizzes)
			quizzes.GET("/category/:category", quizHandler.GetQuizzesByCategory)
			quizzes.GET("/categories", categoriesHandler.GetCategoryGroups)
		}

		// Categories endpoint for simple flat list
		categories := api.Group("/categories")
		{
			categories.GET("", categoriesHandler.GetCategories)
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
		// Apply strict rate limiting to auth endpoints if enabled
		if cfg.RateLimitEnabled {
			auth.Use(middleware.AuthRateLimit())
		}
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected endpoints
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg))
		// Apply per-user rate limiting if enabled
		if cfg.RateLimitEnabled {
			protected.Use(middleware.PerUserRateLimit())
		}
		{
			protected.POST("/auth/logout", authHandler.Logout)
			protected.GET("/profile", authHandler.GetProfile)
			protected.PUT("/profile", authHandler.UpdateProfile)

			users := protected.Group("/users")
			{
				users.GET("/:userId", userHandler.GetUserProfile) // Get another user's profile
				users.PUT("/preferences", preferencesHandler.UpdatePreferences)
				users.GET("/preferences", preferencesHandler.GetPreferences)
				users.POST("/onboarding/complete", preferencesHandler.CompleteOnboarding)
				users.GET("/onboarding/status", preferencesHandler.GetOnboardingStatus)
				users.GET("/quizzes", quizHandler.GetUserQuizzes)
				users.GET("/quizzes/:quizId/attempt", quizHandler.GetUserQuizAttempt)
				users.GET("/quizzes/:quizId/completed-attempt", quizHandler.GetUserLatestCompletedAttempt)
				users.GET("/stats", authHandler.GetUserStats)
				users.GET("/attempts", quizHandler.GetUserAttempts)
				users.GET("/attempts/:attemptId", quizHandler.GetAttemptDetails)
			}

			// Protected quiz endpoints
			protectedQuizzes := protected.Group("/quizzes")
			{
				protectedQuizzes.GET("/:id", quizHandler.GetQuizByID)
				protectedQuizzes.GET("/:id/questions", quizHandler.GetQuizQuestions)
				protectedQuizzes.POST("/:id/attempts", quizHandler.StartQuizAttempt)
				protectedQuizzes.POST("/:id/attempts/:attemptId/submit", quizHandler.SubmitQuizAttempt)
				protectedQuizzes.PUT("/:id/attempts/:attemptId", quizHandler.UpdateQuizAttempt)
				protectedQuizzes.DELETE("/:id/attempts/:attemptId/abandon", quizHandler.AbandonQuizAttempt)

				// Rating endpoints
				protectedQuizzes.POST("/:id/ratings", ratingHandler.CreateRating)
				protectedQuizzes.GET("/:id/ratings", ratingHandler.GetQuizRatings)
				protectedQuizzes.GET("/:id/ratings/average", ratingHandler.GetAverageRating)
				protectedQuizzes.GET("/:id/ratings/user", ratingHandler.GetUserRating)
				protectedQuizzes.PUT("/:id/ratings/:ratingId", ratingHandler.UpdateRating)
				protectedQuizzes.DELETE("/:id/ratings/:ratingId", ratingHandler.DeleteRating)
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
				// Friend notification endpoints (backward compatibility)
				friends.GET("/notifications", notificationHandler.GetFriendNotifications)
				friends.PUT("/notifications/:id/read", notificationHandler.MarkFriendNotificationAsRead)
				friends.PUT("/notifications/read-all", notificationHandler.MarkAllFriendNotificationsAsRead)
			}

			// Unified notification endpoints
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.GetNotifications)
				notifications.GET("/stats", notificationHandler.GetNotificationStats)
				notifications.GET("/:id", notificationHandler.GetNotificationByID)
				notifications.PUT("/:id/read", notificationHandler.MarkNotificationAsRead)
				notifications.PUT("/:id/unread", notificationHandler.MarkNotificationAsUnread)
				notifications.PUT("/read-all", notificationHandler.MarkAllNotificationsAsRead)
				notifications.DELETE("/:id", notificationHandler.DeleteNotification)
				notifications.POST("", notificationHandler.CreateNotification)                  // Admin/system endpoint
				notifications.POST("/cleanup", notificationHandler.CleanupExpiredNotifications) // Admin endpoint
			}

			// Challenge endpoints
			challenges := protected.Group("/challenges")
			// Apply stricter rate limiting for write operations if enabled
			if cfg.RateLimitEnabled {
				challenges.Use(middleware.StrictRateLimit())
			}
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
				challenges.PUT("/:id/cancel", challengesHandler.CancelChallenge)
				challenges.PUT("/:id/score", challengesHandler.UpdateChallengeScore)
				challenges.POST("/:id/link-attempt", challengesHandler.LinkAttemptToChallenge)
				challenges.PUT("/:id/complete", challengesHandler.CompleteChallengeAttempt)
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

			// Favorites endpoints
			favorites := protected.Group("/favorites")
			{
				favorites.POST("", favoritesHandler.AddFavorite)
				favorites.DELETE("/:quizId", favoritesHandler.RemoveFavorite)
				favorites.GET("", favoritesHandler.GetFavorites)
				favorites.GET("/check/:quizId", favoritesHandler.CheckFavorite)
			}

			// Discussion endpoints
			discussions := protected.Group("/discussions")
			{
				discussions.GET("", discussionHandler.GetDiscussions)
				discussions.POST("", discussionHandler.CreateDiscussion)
				discussions.GET("/stats", discussionHandler.GetDiscussionStats)
				discussions.GET("/:id", discussionHandler.GetDiscussion)
				discussions.PUT("/:id", discussionHandler.UpdateDiscussion)
				discussions.DELETE("/:id", discussionHandler.DeleteDiscussion)
				discussions.PUT("/:id/like", discussionHandler.LikeDiscussion)
				discussions.GET("/:id/replies", discussionHandler.GetDiscussionReplies)
				discussions.POST("/:id/replies", discussionHandler.CreateDiscussionReply)
				discussions.PUT("/replies/:replyId", discussionHandler.UpdateDiscussionReply)
				discussions.DELETE("/replies/:replyId", discussionHandler.DeleteDiscussionReply)
				discussions.PUT("/replies/:replyId/like", discussionHandler.LikeDiscussionReply)
			}

			// Admin endpoints for cache management
			admin := protected.Group("/admin")
			{
				admin.DELETE("/cache/app-settings", appSettingsHandler.ClearCache)
			}
		}
	}
}
