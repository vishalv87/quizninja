package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"quizninja-api/models"

	"github.com/stretchr/testify/assert"
)

func TestDigestHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	userID, token := CreateTestUser(t, tc)

	t.Run("GetTodaysDigest", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/digest/today", token, nil)

		if w.Code == http.StatusOK {
			var response models.DigestResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse digest response")

			// Verify digest has is_test_data field
			digestMap := map[string]interface{}{
				"is_test_data": response.Digest.IsTestData,
			}
			VerifyIsTestDataField(t, digestMap, true, "today's digest")

			// Verify articles if they exist
			if len(response.Digest.Articles) > 0 {
				articlesData := make([]interface{}, len(response.Digest.Articles))
				for i, article := range response.Digest.Articles {
					articlesData[i] = map[string]interface{}{
						"is_test_data": article.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, articlesData, true, "today's digest articles")
			}
		}
	})

	t.Run("GetDigestList", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/digest", token, nil)

		if w.Code == http.StatusOK {
			var response models.DigestListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse digest list response")

			if len(response.Digests) > 0 {
				// Verify digests have is_test_data field
				digestsData := make([]interface{}, len(response.Digests))
				for i, digest := range response.Digests {
					digestsData[i] = map[string]interface{}{
						"is_test_data": digest.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, digestsData, true, "digest list")

				// Test pagination
				wPage2 := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/digest?page=2&page_size=5", token, nil)
				if wPage2.Code == http.StatusOK {
					var page2Response models.DigestListResponse
					err := json.Unmarshal(wPage2.Body.Bytes(), &page2Response)
					assert.NoError(t, err, "Should parse page 2 digest list response")

					if len(page2Response.Digests) > 0 {
						digestsData2 := make([]interface{}, len(page2Response.Digests))
						for i, digest := range page2Response.Digests {
							digestsData2[i] = map[string]interface{}{
								"is_test_data": digest.IsTestData,
							}
						}
						VerifyIsTestDataInArray(t, digestsData2, true, "digest list page 2")
					}
				}
			}
		}
	})

	t.Run("GetDigestByDate", func(t *testing.T) {
		// Try common date formats
		dates := []string{"2024-01-01", "2024-02-15", "2024-03-30"}

		for _, date := range dates {
			t.Run(fmt.Sprintf("Date_%s", date), func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/digest/%s", date)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					var response models.DigestResponse
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err, "Should parse digest response")

					// Verify digest has is_test_data field
					digestMap := map[string]interface{}{
						"is_test_data": response.Digest.IsTestData,
					}
					VerifyIsTestDataField(t, digestMap, true, fmt.Sprintf("digest for date %s", date))

					// Verify articles if they exist
					if len(response.Digest.Articles) > 0 {
						articlesData := make([]interface{}, len(response.Digest.Articles))
						for i, article := range response.Digest.Articles {
							articlesData[i] = map[string]interface{}{
								"is_test_data": article.IsTestData,
							}
						}
						VerifyIsTestDataInArray(t, articlesData, true, fmt.Sprintf("articles for digest %s", date))
					}
				}
			})
		}
	})

	t.Run("GetDigestCategories", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/digest/categories", token, nil)

		if w.Code == http.StatusOK {
			var response models.CategoriesResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse categories response")

			assert.True(t, response.Success, "Categories response should be successful")
			// Categories are just strings, no is_test_data field to verify
		}
	})

	t.Run("GetDigestStats", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/digest/stats", token, nil)

		if w.Code == http.StatusOK {
			var response models.DigestStatsResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse digest stats response")

			assert.True(t, response.Success, "Digest stats response should be successful")
			// Stats are aggregated data, no direct is_test_data field to verify
		}
	})

	t.Run("GetTrendingArticles", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/digest/trending", token, nil)

		if w.Code == http.StatusOK {
			var response models.TrendingArticlesResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse trending articles response")

			if len(response.Articles) > 0 {
				// Verify articles have is_test_data field
				articlesData := make([]interface{}, len(response.Articles))
				for i, article := range response.Articles {
					articlesData[i] = map[string]interface{}{
						"is_test_data": article.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, articlesData, true, "trending articles")

				// Test pagination
				wPage2 := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/digest/trending?page=2&page_size=5", token, nil)
				if wPage2.Code == http.StatusOK {
					var page2Response models.TrendingArticlesResponse
					err := json.Unmarshal(wPage2.Body.Bytes(), &page2Response)
					assert.NoError(t, err, "Should parse page 2 trending articles response")

					if len(page2Response.Articles) > 0 {
						articlesData2 := make([]interface{}, len(page2Response.Articles))
						for i, article := range page2Response.Articles {
							articlesData2[i] = map[string]interface{}{
								"is_test_data": article.IsTestData,
							}
						}
						VerifyIsTestDataInArray(t, articlesData2, true, "trending articles page 2")
					}
				}
			}
		}
	})

	t.Run("GetTrendingArticlesByCategory", func(t *testing.T) {
		categories := []string{"technology", "science", "business", "sports", "entertainment"}

		for _, category := range categories {
			t.Run(fmt.Sprintf("Category_%s", category), func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/digest/trending/%s", category)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					var response models.TrendingArticlesResponse
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err, "Should parse trending articles by category response")

					if len(response.Articles) > 0 {
						// Verify articles have is_test_data field
						articlesData := make([]interface{}, len(response.Articles))
						for i, article := range response.Articles {
							articlesData[i] = map[string]interface{}{
								"is_test_data": article.IsTestData,
							}
						}
						VerifyIsTestDataInArray(t, articlesData, true, fmt.Sprintf("trending articles for category %s", category))
					}
				}
			})
		}
	})

	t.Run("GetArticleByID", func(t *testing.T) {
		// First get a list of digests to find articles
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/digest", token, nil)
		if w.Code != http.StatusOK {
			t.Skip("No digests available to test article retrieval")
			return
		}

		var response models.DigestListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil || len(response.Digests) == 0 {
			t.Skip("No valid digests to test")
			return
		}

		// Try to find a digest with articles
		for _, digest := range response.Digests {
			if len(digest.Articles) > 0 {
				articleID := digest.Articles[0].ID
				url := fmt.Sprintf("/api/v1/digest/articles/%s", articleID)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					var articleResponse models.ArticleResponse
					err := json.Unmarshal(w.Body.Bytes(), &articleResponse)
					assert.NoError(t, err, "Should parse article response")

					// Verify article has is_test_data field
					articleMap := map[string]interface{}{
						"is_test_data": articleResponse.Article.IsTestData,
					}
					VerifyIsTestDataField(t, articleMap, true, "individual article")
					break // Found one article, that's enough for the test
				}
			}
		}
	})

	t.Run("GetOrCreateTodaysDigest", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/digest/today/ensure", token, nil)

		if w.Code == http.StatusOK {
			var response models.DigestResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse get or create digest response")

			// Verify digest has is_test_data field
			digestMap := map[string]interface{}{
				"is_test_data": response.Digest.IsTestData,
			}
			VerifyIsTestDataField(t, digestMap, true, "get or create today's digest")
		}
	})

	// Test CRUD operations (Create, Update, Delete) - these might be admin-only endpoints
	t.Run("CreateDigest", func(t *testing.T) {
		testDate, _ := time.Parse("2006-01-02", "2024-12-01")
		createReq := models.DigestRequest{
			Date:    testDate,
			Title:   "Test Digest",
			Summary: stringPtr("Test digest summary"),
			IsDummy: true,
		}

		reqBody, _ := json.Marshal(createReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/digest", token, reqBody)

		if w.Code == http.StatusCreated {
			var response models.DigestResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse create digest response")

			// Verify created digest has is_test_data field
			digestMap := map[string]interface{}{
				"is_test_data": response.Digest.IsTestData,
			}
			VerifyIsTestDataField(t, digestMap, true, "created digest")
		}
	})

	_ = userID // Use userID to avoid unused variable warning
}
