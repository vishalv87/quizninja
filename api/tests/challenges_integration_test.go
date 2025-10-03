package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"quizninja-api/database"
	"quizninja-api/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChallengesIntegrationEndToEnd(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	// Create test users with comprehensive cleanup
	challengerID, challengerToken, challengerSupabaseID, challengerCleanup := CreateTestUserWithCleanup(t, tc, "Integration Test Challenger")
	defer challengerCleanup()

	challengedID, challengedToken, challengedSupabaseID, challengedCleanup := CreateTestUserWithCleanup(t, tc, "Integration Test Challenged")
	defer challengedCleanup()

	// Create friendship
	setupFriendship(t, tc, challengerID, challengedID, challengerToken, challengedToken)

	// Get a quiz
	quizID := getFirstAvailableQuiz(t, tc, challengerToken)
	if quizID == uuid.Nil {
		t.Skip("No quizzes available for integration test")
		return
	}

	// Additional cleanup for any test data created during this specific test
	defer func() {
		// Clean up any challenges created during this test
		if database.DB != nil {
			database.DB.Exec("DELETE FROM challenges WHERE (challenger_id = $1 OR challenged_id = $1 OR challenger_id = $2 OR challenged_id = $2) AND is_test_data = true", challengerID, challengedID)
		}
	}()

	// Suppress unused variable warnings
	_ = challengerSupabaseID
	_ = challengedSupabaseID

	var challengeID uuid.UUID

	t.Run("Complete Challenge Workflow", func(t *testing.T) {
		t.Run("Step 1: Create Challenge", func(t *testing.T) {
			createReq := models.CreateChallengeRequest{
				ChallengedUserID: challengedID,
				QuizID:           quizID,
				Message:          stringPtr("End-to-end test challenge! 🚀"),
				IsGroupChallenge: false,
			}

			reqBody, _ := json.Marshal(createReq)
			w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)

			require.Equal(t, http.StatusCreated, w.Code, "Should create challenge")

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			challenge := response["challenge"].(map[string]interface{})
			challengeIDStr := challenge["id"].(string)
			challengeID, err = uuid.Parse(challengeIDStr)
			require.NoError(t, err)

			assert.Equal(t, "pending", challenge["status"])
			assert.Equal(t, "End-to-end test challenge! 🚀", challenge["message"])
		})

		t.Run("Step 2: Verify Challenge in Pending List", func(t *testing.T) {
			// Check challenger's sent challenges
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges?user_type=challenger&status=pending", challengerToken, nil)
			require.Equal(t, http.StatusOK, w.Code)

			var response models.ChallengeListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			found := false
			for _, challenge := range response.Challenges {
				if challenge.ID == challengeID {
					found = true
					assert.Equal(t, "pending", string(challenge.Status))
					assert.Equal(t, challengerID, challenge.ChallengerID)
					assert.Equal(t, challengedID, challenge.ChallengedID)
				}
			}
			assert.True(t, found, "Should find created challenge in challenger's pending list")

			// Check challenged user's received challenges
			w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/pending", challengedToken, nil)
			require.Equal(t, http.StatusOK, w.Code)

			response2 := ParseJSONResponse(t, w)
			challenges := response2["challenges"].([]interface{})

			found = false
			for _, challengeInterface := range challenges {
				challenge := challengeInterface.(map[string]interface{})
				if challenge["id"].(string) == challengeID.String() {
					found = true
					assert.Equal(t, "pending", challenge["status"])
				}
			}
			assert.True(t, found, "Should find challenge in challenged user's pending list")
		})

		t.Run("Step 3: Accept Challenge", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/challenges/%s/accept", challengeID)
			w := MakeAuthenticatedRequest(t, tc, "PUT", url, challengedToken, nil)

			require.Equal(t, http.StatusOK, w.Code, "Should accept challenge")

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			challenge := response["challenge"].(map[string]interface{})
			assert.Equal(t, "accepted", challenge["status"])
		})

		t.Run("Step 4: Verify Challenge in Active Lists", func(t *testing.T) {
			// Check challenger's active challenges
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/active", challengerToken, nil)
			require.Equal(t, http.StatusOK, w.Code)

			response := ParseJSONResponse(t, w)
			challenges := response["challenges"].([]interface{})

			found := false
			for _, challengeInterface := range challenges {
				challenge := challengeInterface.(map[string]interface{})
				if challenge["id"].(string) == challengeID.String() {
					found = true
					assert.Equal(t, "accepted", challenge["status"])
				}
			}
			assert.True(t, found, "Should find accepted challenge in challenger's active list")

			// Check challenged user's active challenges
			w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/active", challengedToken, nil)
			require.Equal(t, http.StatusOK, w.Code)

			response = ParseJSONResponse(t, w)
			challenges = response["challenges"].([]interface{})

			found = false
			for _, challengeInterface := range challenges {
				challenge := challengeInterface.(map[string]interface{})
				if challenge["id"].(string) == challengeID.String() {
					found = true
					assert.Equal(t, "accepted", challenge["status"])
				}
			}
			assert.True(t, found, "Should find accepted challenge in challenged user's active list")
		})

		t.Run("Step 5: Submit Challenger Score", func(t *testing.T) {
			scoreReq := models.UpdateChallengeScoreRequest{
				UserScore: 87.5,
			}

			reqBody, _ := json.Marshal(scoreReq)
			url := fmt.Sprintf("/api/v1/challenges/%s/score", challengeID)
			w := MakeAuthenticatedRequest(t, tc, "PUT", url, challengerToken, reqBody)

			require.Equal(t, http.StatusOK, w.Code, "Should update challenger score")

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			challenge := response["challenge"].(map[string]interface{})
			assert.Equal(t, "accepted", challenge["status"], "Should remain accepted until both scores submitted")
			assert.Equal(t, 87.5, challenge["challenger_score"])
			assert.Nil(t, challenge["challenged_score"], "Challenged score should still be nil")
		})

		t.Run("Step 6: Submit Challenged User Score (Auto-Complete)", func(t *testing.T) {
			scoreReq := models.UpdateChallengeScoreRequest{
				UserScore: 92.0,
			}

			reqBody, _ := json.Marshal(scoreReq)
			url := fmt.Sprintf("/api/v1/challenges/%s/score", challengeID)
			w := MakeAuthenticatedRequest(t, tc, "PUT", url, challengedToken, reqBody)

			require.Equal(t, http.StatusOK, w.Code, "Should update challenged user score")

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			challenge := response["challenge"].(map[string]interface{})
			assert.Equal(t, "completed", challenge["status"], "Should auto-complete when both scores submitted")
			assert.Equal(t, 87.5, challenge["challenger_score"])
			assert.Equal(t, 92.0, challenge["challenged_score"])
		})

		t.Run("Step 7: Verify Challenge in Completed Lists", func(t *testing.T) {
			// Check challenger's completed challenges
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/completed", challengerToken, nil)
			require.Equal(t, http.StatusOK, w.Code)

			response := ParseJSONResponse(t, w)
			challenges := response["challenges"].([]interface{})

			found := false
			for _, challengeInterface := range challenges {
				challenge := challengeInterface.(map[string]interface{})
				if challenge["id"].(string) == challengeID.String() {
					found = true
					assert.Equal(t, "completed", challenge["status"])
					assert.Equal(t, 87.5, challenge["challenger_score"])
					assert.Equal(t, 92.0, challenge["challenged_score"])
				}
			}
			assert.True(t, found, "Should find completed challenge in challenger's completed list")

			// Check challenged user's completed challenges
			w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/completed", challengedToken, nil)
			require.Equal(t, http.StatusOK, w.Code)

			response = ParseJSONResponse(t, w)
			challenges = response["challenges"].([]interface{})

			found = false
			for _, challengeInterface := range challenges {
				challenge := challengeInterface.(map[string]interface{})
				if challenge["id"].(string) == challengeID.String() {
					found = true
					assert.Equal(t, "completed", challenge["status"])
				}
			}
			assert.True(t, found, "Should find completed challenge in challenged user's completed list")
		})

		t.Run("Step 8: Verify Updated Statistics", func(t *testing.T) {
			// Check challenger's stats
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/stats", challengerToken, nil)
			require.Equal(t, http.StatusOK, w.Code)

			response := ParseJSONResponse(t, w)
			assert.Greater(t, response["total_challenges"], float64(0), "Should have at least 1 challenge")
			assert.Greater(t, response["completed_challenges"], float64(0), "Should have at least 1 completed challenge")

			// Challenger scored 87.5 vs 92.0, so they lost
			assert.Greater(t, response["lost_challenges"], float64(0), "Challenger should have at least 1 loss")

			// Check challenged user's stats
			w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/stats", challengedToken, nil)
			require.Equal(t, http.StatusOK, w.Code)

			response = ParseJSONResponse(t, w)
			assert.Greater(t, response["total_challenges"], float64(0), "Should have at least 1 challenge")
			assert.Greater(t, response["completed_challenges"], float64(0), "Should have at least 1 completed challenge")

			// Challenged user scored 92.0 vs 87.5, so they won
			assert.Greater(t, response["won_challenges"], float64(0), "Challenged user should have at least 1 win")
		})

		t.Run("Step 9: Verify Challenge Details", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/challenges/%s", challengeID)
			w := MakeAuthenticatedRequest(t, tc, "GET", url, challengerToken, nil)

			require.Equal(t, http.StatusOK, w.Code)

			var response models.ChallengeDetailResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			challenge := response.Challenge
			assert.Equal(t, challengeID, challenge.ID)
			assert.Equal(t, "completed", string(challenge.Status))
			assert.Equal(t, "End-to-end test challenge! 🚀", *challenge.Message)
			assert.Equal(t, 87.5, *challenge.ChallengerScore)
			assert.Equal(t, 92.0, *challenge.ChallengedScore)
			assert.NotEmpty(t, challenge.ChallengerName)
			assert.NotEmpty(t, challenge.ChallengedName)
			assert.NotEmpty(t, challenge.QuizTitle)
		})
	})

	t.Run("Alternative Flows", func(t *testing.T) {
		t.Run("Challenge Decline Flow", func(t *testing.T) {
			// Create another challenge
			createReq := models.CreateChallengeRequest{
				ChallengedUserID: challengedID,
				QuizID:           quizID,
				Message:          stringPtr("This one will be declined"),
				IsGroupChallenge: false,
			}

			reqBody, _ := json.Marshal(createReq)
			w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)
			require.Equal(t, http.StatusCreated, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			challenge := response["challenge"].(map[string]interface{})
			declineChallengeID, err := uuid.Parse(challenge["id"].(string))
			require.NoError(t, err)

			// Decline the challenge
			url := fmt.Sprintf("/api/v1/challenges/%s/decline", declineChallengeID)
			w = MakeAuthenticatedRequest(t, tc, "PUT", url, challengedToken, nil)
			require.Equal(t, http.StatusOK, w.Code, "Should decline challenge")

			// Verify it's no longer in pending lists
			w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/pending", challengedToken, nil)
			require.Equal(t, http.StatusOK, w.Code)

			response2 := ParseJSONResponse(t, w)
			challenges := response2["challenges"].([]interface{})

			found := false
			for _, challengeInterface := range challenges {
				challengeMap := challengeInterface.(map[string]interface{})
				if challengeMap["id"].(string) == declineChallengeID.String() {
					found = true
				}
			}
			assert.False(t, found, "Declined challenge should not appear in pending list")
		})

		t.Run("Prevent Duplicate Pending Challenges", func(t *testing.T) {
			// Try to create the same challenge again
			createReq := models.CreateChallengeRequest{
				ChallengedUserID: challengedID,
				QuizID:           quizID,
				Message:          stringPtr("Duplicate attempt"),
				IsGroupChallenge: false,
			}

			reqBody, _ := json.Marshal(createReq)
			w1 := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)
			require.Equal(t, http.StatusCreated, w1.Code, "First challenge should succeed")

			// Second attempt should fail
			w2 := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)
			assert.Equal(t, http.StatusBadRequest, w2.Code, "Duplicate challenge should be prevented")
		})

		t.Run("Error Handling", func(t *testing.T) {
			// Try to accept non-existent challenge
			fakeID := uuid.New()
			url := fmt.Sprintf("/api/v1/challenges/%s/accept", fakeID)
			w := MakeAuthenticatedRequest(t, tc, "PUT", url, challengedToken, nil)
			assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent challenge")

			// Try to update score for non-existent challenge
			scoreReq := models.UpdateChallengeScoreRequest{UserScore: 80.0}
			reqBody, _ := json.Marshal(scoreReq)
			url = fmt.Sprintf("/api/v1/challenges/%s/score", fakeID)
			w = MakeAuthenticatedRequest(t, tc, "PUT", url, challengerToken, reqBody)
			assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent challenge score update")

			// Try to get non-existent challenge details
			url = fmt.Sprintf("/api/v1/challenges/%s", fakeID)
			w = MakeAuthenticatedRequest(t, tc, "GET", url, challengerToken, nil)
			assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent challenge details")
		})
	})

	t.Run("Performance and Load Test", func(t *testing.T) {
		t.Run("Multiple Concurrent Challenges", func(t *testing.T) {
			// Create multiple challenges in sequence to test system load
			for i := 0; i < 5; i++ {
				createReq := models.CreateChallengeRequest{
					ChallengedUserID: challengedID,
					QuizID:           quizID,
					Message:          stringPtr(fmt.Sprintf("Load test challenge #%d", i+1)),
					IsGroupChallenge: false,
				}

				reqBody, _ := json.Marshal(createReq)
				w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)

				if i == 0 {
					assert.Equal(t, http.StatusCreated, w.Code, "First challenge should succeed")
				} else {
					// Subsequent challenges should fail due to pending challenge prevention
					assert.Equal(t, http.StatusBadRequest, w.Code, "Subsequent challenges should be prevented")
				}
			}
		})

		t.Run("Pagination Performance", func(t *testing.T) {
			// Test pagination with different page sizes
			pageSizes := []int{1, 5, 10, 20}

			for _, pageSize := range pageSizes {
				url := fmt.Sprintf("/api/v1/challenges?page=1&page_size=%d", pageSize)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, challengerToken, nil)
				assert.Equal(t, http.StatusOK, w.Code, "Pagination should work for page size %d", pageSize)

				var response models.ChallengeListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.LessOrEqual(t, len(response.Challenges), pageSize, "Should respect page size limit")
			}
		})
	})
}
