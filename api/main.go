package main

import (
	"log"

	"quizninja-api/config"
	"quizninja-api/database"
	"quizninja-api/middleware"
	"quizninja-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Validate configuration
	if err := cfg.ValidateConfig(); err != nil {
		log.Fatal("Configuration validation failed:", err)
	}

	// Log authentication strategy
	log.Printf("Authentication strategy: %s", cfg.GetAuthStrategy())

	gin.SetMode(cfg.GinMode)

	database.Connect(cfg)
	defer database.Close()

	// Initialize rate limiters if enabled
	if cfg.RateLimitEnabled {
		middleware.InitRateLimiters(cfg)
		log.Println("Rate limiting enabled")
	} else {
		log.Println("Rate limiting disabled")
	}

	// Initialize request size limits if enabled
	if cfg.RequestSizeLimitEnabled {
		middleware.InitRequestSizeLimits(cfg)
		log.Printf("Request size limiting enabled (default: %dMB, auth: %dMB, write: %dMB)",
			cfg.RequestSizeDefault/(1024*1024),
			cfg.RequestSizeAuth/(1024*1024),
			cfg.RequestSizeWrite/(1024*1024))
	} else {
		log.Println("Request size limiting disabled")
	}

	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	// Apply request size limit if enabled
	if cfg.RequestSizeLimitEnabled {
		r.Use(middleware.DefaultRequestSizeLimit())
	}

	// Apply global rate limiting if enabled
	if cfg.RateLimitEnabled {
		r.Use(middleware.GlobalRateLimit())
	}

	routes.SetupRoutes(r, cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
