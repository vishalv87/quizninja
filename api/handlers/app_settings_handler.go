package handlers

import (
	"net/http"
	"time"

	"quizninja-api/config"
	"quizninja-api/models"
	"quizninja-api/repository"

	"github.com/gin-gonic/gin"
)

type AppSettingsHandler struct {
	appSettingsRepo *repository.AppSettingsRepository
	config          *config.Config
	cache           map[string]any
	lastCacheUpdate time.Time
	cacheTTL        time.Duration
}

func NewAppSettingsHandler(config *config.Config) *AppSettingsHandler {
	return &AppSettingsHandler{
		appSettingsRepo: repository.NewAppSettingsRepository(),
		config:          config,
		cache:           make(map[string]any),
		cacheTTL:        5 * time.Minute, // Cache for 5 minutes
	}
}

func (ash *AppSettingsHandler) GetAppSettings(c *gin.Context) {
	// Check cache first
	if time.Since(ash.lastCacheUpdate) < ash.cacheTTL && len(ash.cache) > 0 {
		c.JSON(http.StatusOK, models.AppSettingsResponse{
			Settings:    ash.cache,
			Version:     "1.0.0",
			LastUpdated: ash.lastCacheUpdate,
		})
		return
	}

	// Fetch from repository
	settings, err := ash.appSettingsRepo.GetPublicSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch app settings",
		})
		return
	}

	// Update cache
	ash.cache = settings
	ash.lastCacheUpdate = time.Now()

	c.JSON(http.StatusOK, models.AppSettingsResponse{
		Settings:    settings,
		Version:     "1.0.0",
		LastUpdated: ash.lastCacheUpdate,
	})
}

func (ash *AppSettingsHandler) ClearCache(c *gin.Context) {
	ash.cache = make(map[string]any)
	ash.lastCacheUpdate = time.Time{}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cache cleared successfully",
	})
}