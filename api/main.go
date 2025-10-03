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
	if cfg.IsSupabaseAuthEnabled() {
		log.Printf("Supabase Auth enabled - no JWT fallback (strict mode)")
	} else {
		log.Printf("JWT-only authentication enabled")
	}

	gin.SetMode(cfg.GinMode)

	database.Connect(cfg)
	defer database.Close()

	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	routes.SetupRoutes(r, cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
