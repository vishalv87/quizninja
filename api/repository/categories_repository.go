package repository

import (
	"database/sql"
	"strings"

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
	// Query to get all interests from the database
	// Note: is_test_data field is only used for test isolation, not production filtering
	query := `
		SELECT id, name, description,
		       CONCAT('/icons/', COALESCE(icon_name, 'default'), '.png') as icon_url,
		       created_at, updated_at, is_test_data
		FROM interests
		ORDER BY name ASC
	`

	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interests []models.Interest
	for rows.Next() {
		var interest models.Interest
		err := rows.Scan(
			&interest.ID,
			&interest.Name,
			&interest.Description,
			&interest.IconURL,
			&interest.CreatedAt,
			&interest.UpdatedAt,
			&interest.IsTestData,
		)
		if err != nil {
			continue
		}

		// Set computed fields
		interest.DisplayName = toDisplayName(interest.Name)
		interest.IsActive = true

		interests = append(interests, interest)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Group interests into logical categories
	categories := cr.groupInterestsIntoCategories(interests)

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

// Helper function to convert snake_case names to Display Names
func toDisplayName(name string) string {
	words := strings.Split(strings.ReplaceAll(name, "_", " "), " ")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
}

// groupInterestsIntoCategories groups interests into logical categories
func (cr *CategoriesRepository) groupInterestsIntoCategories(interests []models.Interest) []models.Category {
	// Define category mappings based on interest types
	categoryMappings := map[string][]string{
		"general":       {"general_knowledge", "history", "geography", "literature"},
		"science":       {"science", "biology", "chemistry", "physics", "technology"},
		"sports":        {"sports", "football", "basketball"},
		"entertainment": {"movies_tv", "music", "art"},
	}

	// Create interest lookup map
	interestMap := make(map[string]models.Interest)
	for _, interest := range interests {
		interestMap[interest.ID] = interest
	}

	var categories []models.Category

	// Create categories based on mappings
	for categoryID, interestIDs := range categoryMappings {
		var categoryInterests []models.Interest
		for _, interestID := range interestIDs {
			if interest, exists := interestMap[interestID]; exists {
				categoryInterests = append(categoryInterests, interest)
			}
		}

		if len(categoryInterests) > 0 {
			category := models.Category{
				ID:          categoryID,
				Name:        categoryID,
				DisplayName: toDisplayName(categoryID),
				Description: getCategoryDescription(categoryID),
				IconURL:     "/icons/" + categoryID + ".png",
				IsActive:    true,
				Interests:   categoryInterests,
			}
			categories = append(categories, category)
		}
	}

	// Add any remaining interests to a "miscellaneous" category
	var miscInterests []models.Interest
	for _, interest := range interests {
		found := false
		for _, interestIDs := range categoryMappings {
			for _, id := range interestIDs {
				if id == interest.ID {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			miscInterests = append(miscInterests, interest)
		}
	}

	if len(miscInterests) > 0 {
		category := models.Category{
			ID:          "miscellaneous",
			Name:        "miscellaneous",
			DisplayName: "Miscellaneous",
			Description: "Other topics and interests",
			IconURL:     "/icons/miscellaneous.png",
			IsActive:    true,
			Interests:   miscInterests,
		}
		categories = append(categories, category)
	}

	return categories
}

// getCategoryDescription returns a description for each category
func getCategoryDescription(categoryID string) string {
	descriptions := map[string]string{
		"general":       "Questions covering a wide range of topics",
		"science":       "Scientific topics and discoveries",
		"sports":        "Sports trivia and athletics",
		"entertainment": "Movies, TV shows, music, and art",
		"miscellaneous": "Other topics and interests",
	}

	if desc, exists := descriptions[categoryID]; exists {
		return desc
	}
	return "Various quiz topics"
}
