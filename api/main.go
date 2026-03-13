package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"quizninja-api/config"
	"quizninja-api/database"
	"quizninja-api/middleware"
	"quizninja-api/routes"
	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Load()

	// Initialize structured logger
	utils.InitLogger(cfg)
	utils.Infof("Starting QuizNinja API with log level: %s, format: %s", cfg.LogLevel, cfg.LogFormat)

	// Validate configuration
	if err := cfg.ValidateConfig(); err != nil {
		utils.Fatal("Configuration validation failed:", err)
	}

	// Log authentication strategy
	utils.WithFields(logrus.Fields{
		"auth_strategy": cfg.GetAuthStrategy(),
		"gin_mode":      cfg.GinMode,
	}).Info("Application configuration loaded")

	gin.SetMode(cfg.GinMode)

	database.Connect(cfg)
	defer database.Close()

	// Initialize rate limiters if enabled
	if cfg.RateLimitEnabled {
		middleware.InitRateLimiters(cfg)
		utils.WithFields(logrus.Fields{
			"global":   cfg.RateLimitGlobal,
			"auth":     cfg.RateLimitAuth,
			"write":    cfg.RateLimitWrite,
			"per_user": cfg.RateLimitPerUser,
		}).Info("Rate limiting enabled")
	} else {
		utils.Info("Rate limiting disabled")
	}

	// Initialize request size limits if enabled
	if cfg.RequestSizeLimitEnabled {
		middleware.InitRequestSizeLimits(cfg)
		utils.WithFields(logrus.Fields{
			"default_mb": cfg.RequestSizeDefault / (1024 * 1024),
			"auth_mb":    cfg.RequestSizeAuth / (1024 * 1024),
			"write_mb":   cfg.RequestSizeWrite / (1024 * 1024),
		}).Info("Request size limiting enabled")
	} else {
		utils.Info("Request size limiting disabled")
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

	// Cloud Run provides PORT via environment variable
	// Use PORT env var if available, otherwise use config port
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	utils.WithFields(logrus.Fields{
		"port":     port,
		"gin_mode": cfg.GinMode,
	}).Info("Server starting")

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatal("Server error:", err)
		}
	}()

	utils.Info("Server started successfully. Press Ctrl+C to stop.")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// Accept SIGINT (Ctrl+C) and SIGTERM (Cloud Run shutdown)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Info("Shutting down server...")

	// Give active requests time to finish (max 30 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.Fatal("Server forced shutdown:", err)
	}

	utils.Info("Server stopped gracefully")
}
