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

			users := protected.Group("/users")
			{
				users.PUT("/preferences", preferencesHandler.UpdatePreferences)
				users.GET("/preferences", preferencesHandler.GetPreferences)
				users.GET("/quizzes", quizHandler.GetUserQuizzes)
			}

			// Protected quiz endpoints
			protectedQuizzes := protected.Group("/quizzes")
			{
				protectedQuizzes.POST("/:id/attempts", quizHandler.StartQuizAttempt)
				protectedQuizzes.POST("/:id/attempts/:attemptId/submit", quizHandler.SubmitQuizAttempt)
			}

			// Admin endpoints for cache management
			admin := protected.Group("/admin")
			{
				admin.DELETE("/cache/app-settings", appSettingsHandler.ClearCache)
			}
		}
	}
}
