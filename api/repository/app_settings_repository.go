package repository

import (
	"database/sql"
	"strconv"
	"strings"

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
	query := `
		SELECT key, value, description
		FROM app_settings
		ORDER BY key ASC
	`

	rows, err := asr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]any)

	for rows.Next() {
		var key, value, description string
		err := rows.Scan(&key, &value, &description)
		if err != nil {
			continue
		}

		// Convert string values to appropriate types based on key
		convertedValue := convertSettingValue(key, value)
		settings[key] = convertedValue
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Add computed values for backwards compatibility
	asr.addComputedSettings(settings)

	return settings, nil
}

func (asr *AppSettingsRepository) GetAllSettings() ([]models.AppSettings, error) {
	query := `
		SELECT key, value, description, updated_at
		FROM app_settings
		ORDER BY key ASC
	`

	rows, err := asr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []models.AppSettings
	for rows.Next() {
		var setting models.AppSettings
		err := rows.Scan(
			&setting.Key,
			&setting.Value,
			&setting.Description,
			&setting.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// Set default values for fields not in the database
		setting.ID = uuid.New() // Generate UUID for API compatibility
		setting.Category = "general"
		setting.IsPublic = true
		setting.LastModifiedBy = uuid.New() // Default UUID
		setting.CreatedAt = setting.UpdatedAt

		settings = append(settings, setting)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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
	query := `
		INSERT INTO app_settings (key, value, description, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (key) DO UPDATE SET
			value = EXCLUDED.value,
			description = EXCLUDED.description,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := asr.db.Exec(query, setting.Key, setting.Value, setting.Description)
	return err
}

// convertSettingValue converts string values from database to appropriate types
func convertSettingValue(key, value string) any {
	// Handle boolean values
	if strings.Contains(key, "enabled") || key == "maintenance_mode" || key == "force_update_required" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}

	// Handle integer values
	if strings.Contains(key, "questions") || strings.Contains(key, "points") ||
		strings.Contains(key, "limit") || strings.Contains(key, "threshold") ||
		strings.Contains(key, "level") || strings.Contains(key, "duration") {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}

	// Handle float values
	if strings.Contains(key, "multiplier") {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}

	// Return as string by default
	return value
}

// addComputedSettings adds settings that are computed or derived for backwards compatibility
func (asr *AppSettingsRepository) addComputedSettings(settings map[string]any) {
	// Add default values if not present in database
	if _, exists := settings["supported_languages"]; !exists {
		settings["supported_languages"] = []string{"en", "es", "fr", "de"}
	}

	if _, exists := settings["difficulty_levels"]; !exists {
		settings["difficulty_levels"] = []map[string]any{
			{"id": "easy", "name": "Easy", "description": "Basic questions"},
			{"id": "medium", "name": "Medium", "description": "Intermediate questions"},
			{"id": "hard", "name": "Hard", "description": "Advanced questions"},
		}
	}

	if _, exists := settings["notification_frequencies"]; !exists {
		settings["notification_frequencies"] = []map[string]any{
			{"id": "daily", "name": "Daily", "description": "Once per day"},
			{"id": "weekly", "name": "Weekly", "description": "Once per week"},
			{"id": "never", "name": "Never", "description": "No notifications"},
		}
	}

	// Set default app name if not in database
	if _, exists := settings["app_name"]; !exists {
		settings["app_name"] = "QuizNinja"
	}

	// Enable key features by default
	if _, exists := settings["quiz_categories_enabled"]; !exists {
		settings["quiz_categories_enabled"] = true
	}
	if _, exists := settings["leaderboard_enabled"]; !exists {
		settings["leaderboard_enabled"] = true
	}
	if _, exists := settings["achievements_enabled"]; !exists {
		settings["achievements_enabled"] = true
	}
}
