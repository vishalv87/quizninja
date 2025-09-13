package repository

import (
	"database/sql"

	"quizninja-api/database"
	"quizninja-api/models"
)

type CategoriesRepository struct {
	db *sql.DB
}

func NewCategoriesRepository() *CategoriesRepository {
	return &CategoriesRepository{
		db: database.DB,
	}
}

func (cr *CategoriesRepository) GetAllCategories() ([]models.Category, error) {
	// For now, returning mock data since we don't have the full database schema
	// In production, this would query the database
	categories := []models.Category{
		{
			ID:          "general",
			Name:        "general",
			DisplayName: "General Knowledge",
			Description: "Questions covering a wide range of topics",
			IconURL:     "/icons/general.png",
			IsActive:    true,
			Interests: []models.Interest{
				{
					ID:          "history",
					Name:        "history",
					DisplayName: "History",
					Description: "Historical events and figures",
					IconURL:     "/icons/history.png",
					IsActive:    true,
				},
				{
					ID:          "geography",
					Name:        "geography",
					DisplayName: "Geography",
					Description: "World geography and locations",
					IconURL:     "/icons/geography.png",
					IsActive:    true,
				},
			},
		},
		{
			ID:          "science",
			Name:        "science",
			DisplayName: "Science",
			Description: "Scientific topics and discoveries",
			IconURL:     "/icons/science.png",
			IsActive:    true,
			Interests: []models.Interest{
				{
					ID:          "biology",
					Name:        "biology",
					DisplayName: "Biology",
					Description: "Life sciences and organisms",
					IconURL:     "/icons/biology.png",
					IsActive:    true,
				},
				{
					ID:          "chemistry",
					Name:        "chemistry",
					DisplayName: "Chemistry",
					Description: "Chemical elements and reactions",
					IconURL:     "/icons/chemistry.png",
					IsActive:    true,
				},
				{
					ID:          "physics",
					Name:        "physics",
					DisplayName: "Physics",
					Description: "Physical laws and phenomena",
					IconURL:     "/icons/physics.png",
					IsActive:    true,
				},
			},
		},
		{
			ID:          "sports",
			Name:        "sports",
			DisplayName: "Sports",
			Description: "Sports trivia and athletics",
			IconURL:     "/icons/sports.png",
			IsActive:    true,
			Interests: []models.Interest{
				{
					ID:          "football",
					Name:        "football",
					DisplayName: "Football",
					Description: "American football trivia",
					IconURL:     "/icons/football.png",
					IsActive:    true,
				},
				{
					ID:          "basketball",
					Name:        "basketball",
					DisplayName: "Basketball",
					Description: "Basketball trivia and stats",
					IconURL:     "/icons/basketball.png",
					IsActive:    true,
				},
			},
		},
	}

	return categories, nil
}

func (cr *CategoriesRepository) GetCategoryByID(id string) (*models.Category, error) {
	categories, err := cr.GetAllCategories()
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		if category.ID == id {
			return &category, nil
		}
	}

	return nil, sql.ErrNoRows
}