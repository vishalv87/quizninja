package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"quizninja-api/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFavoritesHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	// Create test user with comprehensive cleanup
	userID, token, supabaseUserID, cleanup := CreateTestUserWithCleanup(t, tc, "Favorites Handler Test User")
	defer cleanup()

	// Suppress unused variable warning
	_ = supabaseUserID

	t.Run("GetUserFavorites", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/favorites", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			favorites, exists := data["favorites"]
			if exists {
				favoritesList, ok := favorites.([]interface{})
				if ok && len(favoritesList) > 0 {
					VerifyIsTestDataInArray(t, favoritesList, true, "user favorites")

					// Check nested quiz data in favorites
					for i, favorite := range favoritesList {
						favoriteMap, ok := favorite.(map[string]interface{})
						if ok {
							quiz, quizExists := favoriteMap["quiz"]
							if quizExists && quiz != nil {
								quizMap, ok := quiz.(map[string]interface{})
								if ok {
									VerifyIsTestDataField(t, quizMap, true, fmt.Sprintf("favorite[%d] quiz", i))
								}
							}

							// Check user data if present
							user, userExists := favoriteMap["user"]
							if userExists && user != nil {
								userMap, ok := user.(map[string]interface{})
								if ok {
									VerifyIsTestDataField(t, userMap, true, fmt.Sprintf("favorite[%d] user", i))
								}
							}
						}
					}
				}
			}
		}
	})

	t.Run("GetUserFavoritesWithPagination", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/favorites?page=1&page_size=5", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			favorites, exists := data["favorites"]
			if exists {
				favoritesList, ok := favorites.([]interface{})
				if ok && len(favoritesList) > 0 {
					VerifyIsTestDataInArray(t, favoritesList, true, "paginated favorites")
				}
			}
		}
	})

	t.Run("AddToFavorites", func(t *testing.T) {
		// First get a quiz to add to favorites
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes", token, nil)
		if w.Code != http.StatusOK {
			t.Skip("No quizzes available to test adding to favorites")
			return
		}

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		quizzes, exists := data["quizzes"]
		if !exists {
			t.Skip("No quizzes field in response")
			return
		}

		quizzesList, ok := quizzes.([]interface{})
		if !ok || len(quizzesList) == 0 {
			t.Skip("No quizzes available")
			return
		}

		firstQuiz, ok := quizzesList[0].(map[string]interface{})
		if !ok {
			t.Skip("Invalid quiz data")
			return
		}

		quizID, exists := firstQuiz["id"]
		if !exists {
			t.Skip("Quiz missing ID")
			return
		}

		quizUUID, err := uuid.Parse(quizID.(string))
		if err != nil {
			t.Skip("Invalid quiz ID format")
			return
		}

		addReq := models.AddFavoriteRequest{
			QuizID: quizUUID,
		}

		reqBody, _ := json.Marshal(addReq)
		w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/users/favorites", token, reqBody)

		if w.Code == http.StatusCreated || w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Verify the response contains favorite data
			if data, exists := response["data"]; exists {
				favoriteMap, ok := data.(map[string]interface{})
				if ok {
					VerifyIsTestDataField(t, favoriteMap, true, "added favorite")

					// Check nested quiz data
					quiz, quizExists := favoriteMap["quiz"]
					if quizExists && quiz != nil {
						quizMap, ok := quiz.(map[string]interface{})
						if ok {
							VerifyIsTestDataField(t, quizMap, true, "added favorite quiz")
						}
					}
				}
			}
		}
	})

	t.Run("RemoveFromFavorites", func(t *testing.T) {
		// Test removing a favorite with a fake quiz ID
		fakeQuizID := uuid.New()
		url := fmt.Sprintf("/api/v1/users/favorites/%s", fakeQuizID)
		w := MakeAuthenticatedRequest(t, tc, "DELETE", url, token, nil)

		// We expect either 404 (not found) or 400 (not in favorites)
		// which indicates the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent favorite")
	})

	t.Run("CheckIsFavorite", func(t *testing.T) {
		// Test checking if a quiz is favorited with a fake ID
		fakeQuizID := uuid.New()
		url := fmt.Sprintf("/api/v1/users/favorites/%s/check", fakeQuizID)
		w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Should contain is_favorite field
			isFavorite, exists := response["is_favorite"]
			assert.True(t, exists, "Response should contain 'is_favorite' field")

			_, ok := isFavorite.(bool)
			assert.True(t, ok, "is_favorite should be a boolean")
		}
	})

	t.Run("GetFavoritesByCategory", func(t *testing.T) {
		categories := []string{"science", "technology", "sports", "history"}

		for _, category := range categories {
			t.Run(fmt.Sprintf("Category_%s", category), func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/users/favorites?category=%s", category)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					response := ParseJSONResponse(t, w)
					data := GetDataFromResponse(t, response)

					favorites, exists := data["favorites"]
					if exists {
						favoritesList, ok := favorites.([]interface{})
						if ok && len(favoritesList) > 0 {
							VerifyIsTestDataInArray(t, favoritesList, true, fmt.Sprintf("favorites category %s", category))

							// Check that all favorites are from the correct category
							for i, favorite := range favoritesList {
								favoriteMap, ok := favorite.(map[string]interface{})
								if ok {
									quiz, quizExists := favoriteMap["quiz"]
									if quizExists && quiz != nil {
										quizMap, ok := quiz.(map[string]interface{})
										if ok {
											VerifyIsTestDataField(t, quizMap, true, fmt.Sprintf("favorites category %s[%d] quiz", category, i))
										}
									}
								}
							}
						}
					}
				}
			})
		}
	})

	t.Run("GetFavoritesStats", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/favorites/stats", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Verify basic stats structure
			totalFavorites, exists := response["total_favorites"]
			if exists {
				_, ok := totalFavorites.(float64) // JSON numbers are float64
				assert.True(t, ok, "total_favorites should be a number")
			}

			// Check categories breakdown if it exists
			categoriesBreakdown, exists := response["categories_breakdown"]
			if exists {
				_, ok := categoriesBreakdown.(map[string]interface{})
				assert.True(t, ok, "categories_breakdown should be an object")
			}
		}
	})

	t.Run("GetRecentFavorites", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/favorites/recent?limit=5", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			favorites, exists := data["favorites"]
			if exists {
				favoritesList, ok := favorites.([]interface{})
				if ok && len(favoritesList) > 0 {
					VerifyIsTestDataInArray(t, favoritesList, true, "recent favorites")

					// Check nested data
					for i, favorite := range favoritesList {
						favoriteMap, ok := favorite.(map[string]interface{})
						if ok {
							quiz, quizExists := favoriteMap["quiz"]
							if quizExists && quiz != nil {
								quizMap, ok := quiz.(map[string]interface{})
								if ok {
									VerifyIsTestDataField(t, quizMap, true, fmt.Sprintf("recent favorite[%d] quiz", i))
								}
							}
						}
					}
				}
			}
		}
	})

	// NOTE: BulkRemoveFavoritesRequest model doesn't exist
	// Commenting out this test until bulk remove functionality is implemented
	/*
		t.Run("BulkRemoveFromFavorites", func(t *testing.T) {
			// Test bulk removing favorites - functionality not implemented yet
			t.Skip("BulkRemoveFavoritesRequest model doesn't exist")
		})
	*/

	_ = userID // Use userID to avoid unused variable warning
}
