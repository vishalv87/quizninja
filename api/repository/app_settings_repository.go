package repository

import (
	"database/sql"

	"quizninja-api/database"
	"quizninja-api/models"

	"github.com/google/uuid"
)

type AppSettingsRepository struct {
	db *sql.DB
}

func NewAppSettingsRepository() *AppSettingsRepository {
	return &AppSettingsRepository{
		db: database.DB,
	}
}

func (asr *AppSettingsRepository) GetPublicSettings() (map[string]any, error) {
	// For now, returning mock data since we don't have the full database schema
	// In production, this would query the database for settings where is_public = true
	settings := map[string]any{
		"app_name":            "QuizNinja",
		"app_version":         "1.0.0",
		"max_questions_per_quiz": 20,
		"default_quiz_duration": 300,
		"supported_languages":   []string{"en", "es", "fr", "de"},
		"difficulty_levels": []map[string]any{
			{"id": "easy", "name": "Easy", "description": "Basic questions"},
			{"id": "medium", "name": "Medium", "description": "Intermediate questions"},
			{"id": "hard", "name": "Hard", "description": "Advanced questions"},
		},
		"notification_frequencies": []map[string]any{
			{"id": "daily", "name": "Daily", "description": "Once per day"},
			{"id": "weekly", "name": "Weekly", "description": "Once per week"},
			{"id": "never", "name": "Never", "description": "No notifications"},
		},
		"quiz_categories_enabled": true,
		"leaderboard_enabled":     true,
		"achievements_enabled":    true,
	}

	return settings, nil
}

func (asr *AppSettingsRepository) GetAllSettings() ([]models.AppSettings, error) {
	// Mock implementation - in production this would query the database
	settings := []models.AppSettings{
		{
			ID:             uuid.New(),
			Key:            "app_name",
			Value:          "QuizNinja",
			Description:    "Application name",
			Category:       "general",
			IsPublic:       true,
			LastModifiedBy: uuid.New(),
		},
		{
			ID:             uuid.New(),
			Key:            "max_questions_per_quiz",
			Value:          "20",
			Description:    "Maximum questions per quiz",
			Category:       "quiz",
			IsPublic:       true,
			LastModifiedBy: uuid.New(),
		},
	}

	return settings, nil
}

func (asr *AppSettingsRepository) GetSettingByKey(key string) (*models.AppSettings, error) {
	settings, err := asr.GetAllSettings()
	if err != nil {
		return nil, err
	}

	for _, setting := range settings {
		if setting.Key == key {
			return &setting, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (asr *AppSettingsRepository) UpdateSetting(setting *models.AppSettings) error {
	// Mock implementation - in production this would update the database
	return nil
}