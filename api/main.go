package main

import (
	"log"
	"net/http"

	"quizninja-api/config"
	"quizninja-api/database"
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

	setupRoutes(r)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "QuizNinja API is running",
		})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}
}
