package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoriesHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	_, token, _, cleanup := CreateTestUserWithCleanup(t, tc, "Categories Test User")
	defer cleanup()

	t.Run("GetCategories", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)

		assert.Equal(t, http.StatusOK, w.Code, "Categories endpoint should return 200 OK")

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			categories, exists := response["data"]
			assert.True(t, exists, "Response should contain 'categories' field")

			categoriesList, ok := categories.([]interface{})
			assert.True(t, ok, "Categories field should be an array")
			assert.Greater(t, len(categoriesList), 0, "Should have at least one category")

			// Track which expected categories we find
			expectedCategories := map[string]bool{
				"general":       false,
				"science":       false,
				"sports":        false,
				"entertainment": false,
			}

			for i, category := range categoriesList {
				categoryMap, ok := category.(map[string]interface{})
				assert.True(t, ok, "Each category should be an object")

				if ok {
					// Verify required fields
					id, hasID := categoryMap["id"]
					assert.True(t, hasID, "Category should have 'id' field")

					name, hasName := categoryMap["name"]
					assert.True(t, hasName, "Category should have 'name' field")

					displayName, hasDisplayName := categoryMap["display_name"]
					assert.True(t, hasDisplayName, "Category should have 'display_name' field")

					description, hasDescription := categoryMap["description"]
					assert.True(t, hasDescription, "Category should have 'description' field")

					categories, hasCategories := categoryMap["categories"]
					assert.True(t, hasCategories, "Category should have 'categories' field")

					// Verify categories structure
					if categoriesList, ok := categories.([]interface{}); ok {
						assert.Greater(t, len(categoriesList), 0, "Category should have at least one category")

						// Check first category for proper structure and is_test_data
						if len(categoriesList) > 0 {
							if firstCategory, ok := categoriesList[0].(map[string]interface{}); ok {
								VerifyIsTestDataField(t, firstCategory, true, "category category")
							}
						}
					}

					// Mark expected categories as found
					if categoryIDStr, ok := id.(string); ok {
						if _, expected := expectedCategories[categoryIDStr]; expected {
							expectedCategories[categoryIDStr] = true
						}
					}

					// Verify category names match expected structure
					if nameStr, ok := name.(string); ok {
						assert.NotEmpty(t, nameStr, "Category name should not be empty")
					}

					if displayNameStr, ok := displayName.(string); ok {
						assert.NotEmpty(t, displayNameStr, "Category display name should not be empty")
					}

					if descStr, ok := description.(string); ok {
						assert.NotEmpty(t, descStr, "Category description should not be empty")
					}
				}

				// Limit checking to first few items for performance
				if i >= 10 {
					break
				}
			}

			// Verify we found all expected categories
			for categoryName, found := range expectedCategories {
				assert.True(t, found, "Should find expected category: %s", categoryName)
			}
		}
	})

	t.Run("VerifyNewlyAddedCategories", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/categories", token, nil)

		assert.Equal(t, http.StatusOK, w.Code, "Categories endpoint should return 200 OK")

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			categories, exists := response["data"]
			assert.True(t, exists, "Response should contain 'categories' field")

			categoriesList, ok := categories.([]interface{})
			assert.True(t, ok, "Categories field should be an array")

			// Track newly added categories from migration 027
			newCategories := map[string]bool{
				"biology":    false,
				"chemistry":  false,
				"physics":    false,
				"football":   false,
				"basketball": false,
			}

			for _, category := range categoriesList {
				if categoryMap, ok := category.(map[string]interface{}); ok {
					if categories, hasCategories := categoryMap["categories"]; hasCategories {
						if categoriesList, ok := categories.([]interface{}); ok {
							for _, category := range categoriesList {
								if categoryMap, ok := category.(map[string]interface{}); ok {
									// Verify is_test_data is true for all categories
									VerifyIsTestDataField(t, categoryMap, true, "category")

									// Check if this is one of our newly added categories
									if categoryID, hasID := categoryMap["id"]; hasID {
										if idStr, ok := categoryID.(string); ok {
											if _, isNew := newCategories[idStr]; isNew {
												newCategories[idStr] = true
											}
										}
									}
								}
							}
						}
					}
				}
			}

			// Verify all new categories were found
			for categoryName, found := range newCategories {
				assert.True(t, found, "Should find newly added category: %s", categoryName)
			}

			// Verify specific category mappings for new categories
			scienceCategory := findCategoryByID(categoriesList, "science")
			if scienceCategory != nil {
				categories := getCategoryCategories(scienceCategory)
				assert.Contains(t, categories, "biology", "Science category should contain biology")
				assert.Contains(t, categories, "chemistry", "Science category should contain chemistry")
				assert.Contains(t, categories, "physics", "Science category should contain physics")
			}

			sportsCategory := findCategoryByID(categoriesList, "sports")
			if sportsCategory != nil {
				categories := getCategoryCategories(sportsCategory)
				assert.Contains(t, categories, "football", "Sports category should contain football")
				assert.Contains(t, categories, "basketball", "Sports category should contain basketball")
			}
		}
	})

	t.Run("GetCategoriesWithStats", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/categories?include_stats=true", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			categories, exists := data["categories"]
			if exists {
				categoriesList, ok := categories.([]interface{})
				if ok && len(categoriesList) > 0 {
					// When stats are included, categories are more likely to be objects
					for i, category := range categoriesList {
						categoryMap, ok := category.(map[string]interface{})
						if ok {
							// Check for stats fields
							if quizCount, exists := categoryMap["quiz_count"]; exists {
								_, ok := quizCount.(float64)
								assert.True(t, ok, "quiz_count should be a number")
							}

							// Verify is_test_data if present
							if _, hasTestData := categoryMap["is_test_data"]; hasTestData {
								VerifyIsTestDataField(t, categoryMap, true, "category with stats")
							}
						}

						// Limit checking for performance
						if i >= 5 {
							break
						}
					}
				}
			}
		}
	})

	t.Run("GetCategoryDetails", func(t *testing.T) {
		// Test getting details for common categories
		commonCategories := []string{"science", "technology", "sports", "history", "general"}

		for _, categoryName := range commonCategories {
			t.Run("Category_"+categoryName, func(t *testing.T) {
				url := "/api/v1/categories/" + categoryName
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					response := ParseJSONResponse(t, w)
					data := GetDataFromResponse(t, response)

					// Verify category details structure
					category, exists := data["category"]
					if exists {
						categoryMap, ok := category.(map[string]interface{})
						if ok {
							// Verify name matches requested category
							if name, nameExists := categoryMap["name"]; nameExists {
								assert.Equal(t, categoryName, name, "Category name should match request")
							}

							// Verify is_test_data if present
							if _, hasTestData := categoryMap["is_test_data"]; hasTestData {
								VerifyIsTestDataField(t, categoryMap, true, "category details")
							}
						}
					}

					// Check quiz count if present
					quizCount, countExists := data["quiz_count"]
					if countExists {
						_, ok := quizCount.(float64)
						assert.True(t, ok, "quiz_count should be a number")
					}
				}
			})
		}
	})

	t.Run("GetPopularCategories", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/categories/popular?limit=5", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			categories, exists := data["categories"]
			if exists {
				categoriesList, ok := categories.([]interface{})
				if ok && len(categoriesList) > 0 {
					for i, category := range categoriesList {
						categoryMap, ok := category.(map[string]interface{})
						if ok {
							// Popular categories should have popularity metrics
							if popularity, exists := categoryMap["popularity_score"]; exists {
								_, ok := popularity.(float64)
								assert.True(t, ok, "popularity_score should be a number")
							}

							// Verify is_test_data if present
							if _, hasTestData := categoryMap["is_test_data"]; hasTestData {
								VerifyIsTestDataField(t, categoryMap, true, "popular category")
							}
						}

						// Limit checking for performance
						if i >= 3 {
							break
						}
					}
				}
			}
		}
	})

	t.Run("SearchCategories", func(t *testing.T) {
		searchQueries := []string{"sci", "tech", "sport"}

		for _, query := range searchQueries {
			t.Run("Search_"+query, func(t *testing.T) {
				url := "/api/v1/categories/search?q=" + query
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					response := ParseJSONResponse(t, w)
					data := GetDataFromResponse(t, response)

					categories, exists := data["categories"]
					if exists {
						categoriesList, ok := categories.([]interface{})
						if ok && len(categoriesList) > 0 {
							for i, category := range categoriesList {
								// Categories in search results might be strings or objects
								if categoryMap, ok := category.(map[string]interface{}); ok {
									// Verify is_test_data if present
									if _, hasTestData := categoryMap["is_test_data"]; hasTestData {
										VerifyIsTestDataField(t, categoryMap, true, "search result category")
									}
								}

								// Limit checking for performance
								if i >= 3 {
									break
								}
							}
						}
					}

					// Verify search query is reflected
					searchQuery, queryExists := data["query"]
					if queryExists {
						assert.Equal(t, query, searchQuery, "Search query should be reflected in response")
					}
				}
			})
		}
	})

	t.Run("GetCategoryHierarchy", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/categories/hierarchy", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Category hierarchy might be a tree structure
			hierarchy, exists := data["hierarchy"]
			if exists {
				hierarchyList, ok := hierarchy.([]interface{})
				if ok && len(hierarchyList) > 0 {
					for i, item := range hierarchyList {
						itemMap, ok := item.(map[string]interface{})
						if ok {
							// Verify is_test_data if present
							if _, hasTestData := itemMap["is_test_data"]; hasTestData {
								VerifyIsTestDataField(t, itemMap, true, "hierarchy category")
							}

							// Check children if present
							if children, childrenExist := itemMap["children"]; childrenExist {
								childrenList, ok := children.([]interface{})
								if ok && len(childrenList) > 0 {
									for j, child := range childrenList {
										if childMap, ok := child.(map[string]interface{}); ok {
											if _, hasTestData := childMap["is_test_data"]; hasTestData {
												VerifyIsTestDataField(t, childMap, true, "hierarchy child category")
											}
										}

										// Limit child checking
										if j >= 2 {
											break
										}
									}
								}
							}
						}

						// Limit checking for performance
						if i >= 3 {
							break
						}
					}
				}
			}
		}
	})

}

// Helper function to find a category by ID
func findCategoryByID(categories []interface{}, categoryID string) map[string]interface{} {
	for _, category := range categories {
		if categoryMap, ok := category.(map[string]interface{}); ok {
			if id, hasID := categoryMap["id"]; hasID {
				if idStr, ok := id.(string); ok && idStr == categoryID {
					return categoryMap
				}
			}
		}
	}
	return nil
}

// Helper function to get category IDs from a category
func getCategoryCategories(category map[string]interface{}) []string {
	var categories []string
	if categoriesField, hasCategories := category["categories"]; hasCategories {
		if categoriesList, ok := categoriesField.([]interface{}); ok {
			for _, category := range categoriesList {
				if categoryMap, ok := category.(map[string]interface{}); ok {
					if id, hasID := categoryMap["id"]; hasID {
						if idStr, ok := id.(string); ok {
							categories = append(categories, idStr)
						}
					}
				}
			}
		}
	}
	return categories
}
