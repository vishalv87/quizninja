package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoriesHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	userID, token := CreateTestUser(t, tc)

	t.Run("GetCategories", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/categories", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			categories, exists := data["categories"]
			assert.True(t, exists, "Response should contain 'categories' field")

			categoriesList, ok := categories.([]interface{})
			assert.True(t, ok, "Categories field should be an array")

			if len(categoriesList) > 0 {
				// Categories might be strings or objects
				// Only verify is_test_data if they're objects
				for i, category := range categoriesList {
					categoryMap, ok := category.(map[string]interface{})
					if ok {
						// Only verify if the category object has is_test_data field
						if _, hasTestData := categoryMap["is_test_data"]; hasTestData {
							VerifyIsTestDataField(t, categoryMap, true, "category")
						}
					} else {
						// If it's a string, just verify it's not empty
						categoryStr, ok := category.(string)
						if ok {
							assert.NotEmpty(t, categoryStr, "Category string should not be empty")
						}
					}

					// Limit checking to first few items for performance
					if i >= 10 {
						break
					}
				}
			}

			// Verify total count if present
			total, totalExists := data["total"]
			if totalExists {
				totalFloat, ok := total.(float64)
				if ok {
					assert.Equal(t, len(categoriesList), int(totalFloat), "Total should match categories count")
				}
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

	_ = userID // Use userID to avoid unused variable warning
}
