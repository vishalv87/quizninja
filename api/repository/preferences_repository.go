package repository

import (
	"database/sql"
	"log"

	"quizninja-api/database"
	"quizninja-api/models"
)

type PreferencesRepository struct {
	db *sql.DB
}

func NewPreferencesRepository() *PreferencesRepository {
	return &PreferencesRepository{
		db: database.DB,
	}
}

// GetCategories retrieves all categories
func (pr *PreferencesRepository) GetCategories() ([]models.Category, error) {
	query := `
		SELECT id, name, name as display_name, description,
		       CONCAT('/icons/', icon_name, '.png') as icon_url,
		       true as is_active, created_at, created_at as updated_at, is_test_data
		FROM categories
		ORDER BY name ASC
	`

	rows, err := pr.db.Query(query)
	if err != nil {
		log.Printf("GetCategories: Failed to query categories: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.DisplayName,
			&category.Description,
			&category.IconURL,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.IsTestData,
		)
		if err != nil {
			log.Printf("GetCategories: Failed to scan category: %v", err)
			continue
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		log.Printf("GetCategories: Row iteration error: %v", err)
		return nil, err
	}

	log.Printf("GetCategories: Successfully retrieved %d categories", len(categories))
	return categories, nil
}

// GetDifficultyLevels retrieves all difficulty levels
func (pr *PreferencesRepository) GetDifficultyLevels() ([]models.DifficultyLevel, error) {
	query := `
		SELECT id, name, name as display_name, description, 1 as order_level, true as is_active, is_test_data
		FROM difficulty_levels
		ORDER BY name ASC
	`

	rows, err := pr.db.Query(query)
	if err != nil {
		log.Printf("GetDifficultyLevels: Failed to query difficulty levels: %v", err)
		return nil, err
	}
	defer rows.Close()

	var levels []models.DifficultyLevel
	for rows.Next() {
		var level models.DifficultyLevel
		err := rows.Scan(
			&level.ID,
			&level.Name,
			&level.DisplayName,
			&level.Description,
			&level.Order,
			&level.IsActive,
			&level.IsTestData,
		)
		if err != nil {
			log.Printf("GetDifficultyLevels: Failed to scan difficulty level: %v", err)
			continue
		}
		levels = append(levels, level)
	}

	if err = rows.Err(); err != nil {
		log.Printf("GetDifficultyLevels: Row iteration error: %v", err)
		return nil, err
	}

	log.Printf("GetDifficultyLevels: Successfully retrieved %d difficulty levels", len(levels))
	return levels, nil
}

// GetNotificationFrequencies retrieves all notification frequencies
func (pr *PreferencesRepository) GetNotificationFrequencies() ([]models.NotificationFrequency, error) {
	query := `
		SELECT id, name, name as display_name, description, true as is_active, is_test_data
		FROM notification_frequencies
		ORDER BY name ASC
	`

	rows, err := pr.db.Query(query)
	if err != nil {
		log.Printf("GetNotificationFrequencies: Failed to query notification frequencies: %v", err)
		return nil, err
	}
	defer rows.Close()

	var frequencies []models.NotificationFrequency
	for rows.Next() {
		var frequency models.NotificationFrequency
		err := rows.Scan(
			&frequency.ID,
			&frequency.Name,
			&frequency.DisplayName,
			&frequency.Description,
			&frequency.IsActive,
			&frequency.IsTestData,
		)
		if err != nil {
			log.Printf("GetNotificationFrequencies: Failed to scan notification frequency: %v", err)
			continue
		}
		frequencies = append(frequencies, frequency)
	}

	if err = rows.Err(); err != nil {
		log.Printf("GetNotificationFrequencies: Row iteration error: %v", err)
		return nil, err
	}

	log.Printf("GetNotificationFrequencies: Successfully retrieved %d notification frequencies", len(frequencies))
	return frequencies, nil
}
