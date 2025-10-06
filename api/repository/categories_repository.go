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

// GetAllCategories returns a simple flat list of all categories
func (cr *CategoriesRepository) GetAllCategories() ([]models.Category, error) {
	query := `
		SELECT id, name, description,
		       CONCAT('/icons/', COALESCE(icon_name, 'default'), '.png') as icon_url,
		       created_at, updated_at, is_test_data
		FROM categories
		ORDER BY name ASC
	`

	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.IconURL,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.IsTestData,
		)
		if err != nil {
			continue
		}

		// Set computed fields
		category.DisplayName = toDisplayName(category.Name)
		category.IsActive = true

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (cr *CategoriesRepository) GetAllCategoryGroups() ([]models.CategoryGroup, error) {
	// Query to get all categories from the database
	// Note: is_test_data field is only used for test isolation, not production filtering
	query := `
		SELECT id, name, description,
		       CONCAT('/icons/', COALESCE(icon_name, 'default'), '.png') as icon_url,
		       created_at, updated_at, is_test_data
		FROM categories
		ORDER BY name ASC
	`

	rows, err := cr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.IconURL,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.IsTestData,
		)
		if err != nil {
			continue
		}

		// Set computed fields
		category.DisplayName = toDisplayName(category.Name)
		category.IsActive = true

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Group categories into logical category groups
	categoryGroups := cr.groupCategoriesIntoCategoryGroups(categories)

	return categoryGroups, nil
}

func (cr *CategoriesRepository) GetCategoryGroupByID(id string) (*models.CategoryGroup, error) {
	categoryGroups, err := cr.GetAllCategoryGroups()
	if err != nil {
		return nil, err
	}

	for _, categoryGroup := range categoryGroups {
		if categoryGroup.ID == id {
			return &categoryGroup, nil
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

// groupCategoriesIntoCategoryGroups groups categories into logical category groups
func (cr *CategoriesRepository) groupCategoriesIntoCategoryGroups(categories []models.Category) []models.CategoryGroup {
	// Define category group mappings based on category types
	categoryGroupMappings := map[string][]string{
		"general":       {"general_knowledge", "history", "geography", "literature"},
		"science":       {"science", "biology", "chemistry", "physics", "technology"},
		"sports":        {"sports", "football", "basketball"},
		"entertainment": {"movies_tv", "music", "art"},
	}

	// Create category lookup map
	categoryMap := make(map[string]models.Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	var categoryGroups []models.CategoryGroup

	// Create category groups based on mappings
	for categoryGroupID, categoryIDs := range categoryGroupMappings {
		var groupCategories []models.Category
		for _, categoryID := range categoryIDs {
			if category, exists := categoryMap[categoryID]; exists {
				groupCategories = append(groupCategories, category)
			}
		}

		if len(groupCategories) > 0 {
			categoryGroup := models.CategoryGroup{
				ID:          categoryGroupID,
				Name:        categoryGroupID,
				DisplayName: toDisplayName(categoryGroupID),
				Description: getCategoryGroupDescription(categoryGroupID),
				IconURL:     "/icons/" + categoryGroupID + ".png",
				IsActive:    true,
				Categories:  groupCategories,
			}
			categoryGroups = append(categoryGroups, categoryGroup)
		}
	}

	// Add any remaining categories to a "miscellaneous" category group
	var miscCategories []models.Category
	for _, category := range categories {
		found := false
		for _, categoryIDs := range categoryGroupMappings {
			for _, id := range categoryIDs {
				if id == category.ID {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			miscCategories = append(miscCategories, category)
		}
	}

	if len(miscCategories) > 0 {
		categoryGroup := models.CategoryGroup{
			ID:          "miscellaneous",
			Name:        "miscellaneous",
			DisplayName: "Miscellaneous",
			Description: "Other topics and categories",
			IconURL:     "/icons/miscellaneous.png",
			IsActive:    true,
			Categories:  miscCategories,
		}
		categoryGroups = append(categoryGroups, categoryGroup)
	}

	return categoryGroups
}

// getCategoryGroupDescription returns a description for each category group
func getCategoryGroupDescription(categoryGroupID string) string {
	descriptions := map[string]string{
		"general":       "Questions covering a wide range of topics",
		"science":       "Scientific topics and discoveries",
		"sports":        "Sports trivia and athletics",
		"entertainment": "Movies, TV shows, music, and art",
		"miscellaneous": "Other topics and categories",
	}

	if desc, exists := descriptions[categoryGroupID]; exists {
		return desc
	}
	return "Various quiz topics"
}
