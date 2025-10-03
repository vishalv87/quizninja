package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"quizninja-api/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChallengesHandlerComprehensive(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	// Create test users
	challengerID, challengerToken := CreateTestUser(t, tc)
	challengedID, challengedToken := CreateTestUser(t, tc)
	thirdUserID, thirdUserToken := CreateTestUser(t, tc)

	// Create a friendship between challenger and challenged user
	setupFriendship(t, tc, challengerID, challengedID, challengerToken, challengedToken)

	// Get a quiz for testing
	quizID := getFirstAvailableQuiz(t, tc, challengerToken)
	if quizID == uuid.Nil {
		t.Skip("No quizzes available for testing")
		return
	}

	t.Run("Challenge Creation Flow", func(t *testing.T) {
		t.Run("Create Valid Challenge", func(t *testing.T) {
			createReq := models.CreateChallengeRequest{
				ChallengedUserID: challengedID,
				QuizID:           quizID,
				Message:          stringPtr("Let's see who's smarter! 😄"),
				IsGroupChallenge: false,
			}

			reqBody, _ := json.Marshal(createReq)
			w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)

			assert.Equal(t, http.StatusCreated, w.Code, "Should create challenge successfully")

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			challenge, exists := response["challenge"]
			require.True(t, exists)

			challengeMap := challenge.(map[string]interface{})
			assert.Equal(t, "pending", challengeMap["status"])
			assert.Equal(t, "Let's see who's smarter! 😄", challengeMap["message"])
		})

		t.Run("Prevent Self Challenge", func(t *testing.T) {
			createReq := models.CreateChallengeRequest{
				ChallengedUserID: challengerID, // Same as challenger
				QuizID:           quizID,
				Message:          stringPtr("Challenge myself"),
				IsGroupChallenge: false,
			}

			reqBody, _ := json.Marshal(createReq)
			w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)

			assert.Equal(t, http.StatusBadRequest, w.Code, "Should prevent self-challenge")
		})

		t.Run("Prevent Challenge Non-Friend", func(t *testing.T) {
			createReq := models.CreateChallengeRequest{
				ChallengedUserID: thirdUserID, // Not a friend
				QuizID:           quizID,
				Message:          stringPtr("Challenge stranger"),
				IsGroupChallenge: false,
			}

			reqBody, _ := json.Marshal(createReq)
			w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)

			assert.Equal(t, http.StatusBadRequest, w.Code, "Should prevent challenging non-friends")
		})

		t.Run("Prevent Duplicate Pending Challenge", func(t *testing.T) {
			// Create first challenge
			createReq := models.CreateChallengeRequest{
				ChallengedUserID: challengedID,
				QuizID:           quizID,
				Message:          stringPtr("First challenge"),
				IsGroupChallenge: false,
			}

			reqBody, _ := json.Marshal(createReq)
			w1 := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)
			assert.Equal(t, http.StatusCreated, w1.Code)

			// Try to create duplicate
			w2 := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)
			assert.Equal(t, http.StatusBadRequest, w2.Code, "Should prevent duplicate pending challenges")
		})

		t.Run("Invalid Quiz ID", func(t *testing.T) {
			createReq := models.CreateChallengeRequest{
				ChallengedUserID: challengedID,
				QuizID:           uuid.New(), // Non-existent quiz
				Message:          stringPtr("Invalid quiz challenge"),
				IsGroupChallenge: false,
			}

			reqBody, _ := json.Marshal(createReq)
			w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)

			assert.Equal(t, http.StatusBadRequest, w.Code, "Should reject invalid quiz ID")
		})
	})

	t.Run("Challenge State Management", func(t *testing.T) {
		// Create a fresh challenge for state testing
		challengeID := createTestChallenge(t, tc, challengerToken, challengedID, quizID)

		t.Run("Accept Challenge", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/challenges/%s/accept", challengeID)
			w := MakeAuthenticatedRequest(t, tc, "PUT", url, challengedToken, nil)

			assert.Equal(t, http.StatusOK, w.Code, "Challenged user should be able to accept")

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			challenge := response["challenge"].(map[string]interface{})
			assert.Equal(t, "accepted", challenge["status"])
		})

		t.Run("Prevent Non-Challenged User Accept", func(t *testing.T) {
			newChallengeID := createTestChallenge(t, tc, challengerToken, challengedID, quizID)
			url := fmt.Sprintf("/api/v1/challenges/%s/accept", newChallengeID)
			w := MakeAuthenticatedRequest(t, tc, "PUT", url, thirdUserToken, nil)

			assert.Equal(t, http.StatusBadRequest, w.Code, "Non-challenged user should not be able to accept")
		})

		t.Run("Decline Challenge", func(t *testing.T) {
			newChallengeID := createTestChallenge(t, tc, challengerToken, challengedID, quizID)
			url := fmt.Sprintf("/api/v1/challenges/%s/decline", newChallengeID)
			w := MakeAuthenticatedRequest(t, tc, "PUT", url, challengedToken, nil)

			assert.Equal(t, http.StatusOK, w.Code, "Challenged user should be able to decline")
		})

		t.Run("Update Challenge Scores", func(t *testing.T) {
			// Create and accept a challenge
			newChallengeID := createTestChallenge(t, tc, challengerToken, challengedID, quizID)
			acceptUrl := fmt.Sprintf("/api/v1/challenges/%s/accept", newChallengeID)
			MakeAuthenticatedRequest(t, tc, "PUT", acceptUrl, challengedToken, nil)

			// Update challenger score
			scoreReq := models.UpdateChallengeScoreRequest{UserScore: 85.5}
			reqBody, _ := json.Marshal(scoreReq)
			scoreUrl := fmt.Sprintf("/api/v1/challenges/%s/score", newChallengeID)
			w1 := MakeAuthenticatedRequest(t, tc, "PUT", scoreUrl, challengerToken, reqBody)
			assert.Equal(t, http.StatusOK, w1.Code, "Should update challenger score")

			// Update challenged user score
			scoreReq2 := models.UpdateChallengeScoreRequest{UserScore: 92.0}
			reqBody2, _ := json.Marshal(scoreReq2)
			w2 := MakeAuthenticatedRequest(t, tc, "PUT", scoreUrl, challengedToken, reqBody2)
			assert.Equal(t, http.StatusOK, w2.Code, "Should update challenged user score")

			// Verify challenge is completed
			var response map[string]interface{}
			err := json.Unmarshal(w2.Body.Bytes(), &response)
			require.NoError(t, err)

			challenge := response["challenge"].(map[string]interface{})
			assert.Equal(t, "completed", challenge["status"], "Challenge should auto-complete when both users scored")
			assert.Equal(t, 85.5, challenge["challenger_score"])
			assert.Equal(t, 92.0, challenge["challenged_score"])
		})

		t.Run("Prevent Non-Participant Score Update", func(t *testing.T) {
			newChallengeID := createTestChallenge(t, tc, challengerToken, challengedID, quizID)
			scoreReq := models.UpdateChallengeScoreRequest{UserScore: 75.0}
			reqBody, _ := json.Marshal(scoreReq)
			scoreUrl := fmt.Sprintf("/api/v1/challenges/%s/score", newChallengeID)
			w := MakeAuthenticatedRequest(t, tc, "PUT", scoreUrl, thirdUserToken, reqBody)

			assert.Equal(t, http.StatusBadRequest, w.Code, "Non-participants should not update scores")
		})
	})

	t.Run("Challenge Retrieval and Filtering", func(t *testing.T) {
		// Create challenges in different states for testing
		pendingChallengeID := createTestChallenge(t, tc, challengerToken, challengedID, quizID)
		activeChallengeID := createTestChallenge(t, tc, challengerToken, challengedID, quizID)
		completedChallengeID := createTestChallenge(t, tc, challengerToken, challengedID, quizID)

		// Accept the active challenge
		acceptUrl := fmt.Sprintf("/api/v1/challenges/%s/accept", activeChallengeID)
		MakeAuthenticatedRequest(t, tc, "PUT", acceptUrl, challengedToken, nil)

		// Complete the completed challenge
		acceptUrl2 := fmt.Sprintf("/api/v1/challenges/%s/accept", completedChallengeID)
		MakeAuthenticatedRequest(t, tc, "PUT", acceptUrl2, challengedToken, nil)

		// Add scores to complete it
		scoreReq := models.UpdateChallengeScoreRequest{UserScore: 80.0}
		reqBody, _ := json.Marshal(scoreReq)
		scoreUrl := fmt.Sprintf("/api/v1/challenges/%s/score", completedChallengeID)
		MakeAuthenticatedRequest(t, tc, "PUT", scoreUrl, challengerToken, reqBody)
		MakeAuthenticatedRequest(t, tc, "PUT", scoreUrl, challengedToken, reqBody)

		t.Run("Get All Challenges", func(t *testing.T) {
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges", challengerToken, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			var response models.ChallengeListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Greater(t, len(response.Challenges), 0, "Should return challenges")
		})

		t.Run("Get Pending Challenges", func(t *testing.T) {
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/pending", challengerToken, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			response := ParseJSONResponse(t, w)
			challenges := response["challenges"].([]interface{})

			// Should have at least the pending challenge we created
			found := false
			for _, challenge := range challenges {
				challengeMap := challenge.(map[string]interface{})
				if challengeMap["id"].(string) == pendingChallengeID.String() {
					found = true
					assert.Equal(t, "pending", challengeMap["status"])
				}
			}
			assert.True(t, found, "Should find our pending challenge")
		})

		t.Run("Get Active Challenges", func(t *testing.T) {
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/active", challengerToken, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			response := ParseJSONResponse(t, w)
			challenges := response["challenges"].([]interface{})

			found := false
			for _, challenge := range challenges {
				challengeMap := challenge.(map[string]interface{})
				if challengeMap["id"].(string) == activeChallengeID.String() {
					found = true
					assert.Equal(t, "accepted", challengeMap["status"])
				}
			}
			assert.True(t, found, "Should find our active challenge")
		})

		t.Run("Get Completed Challenges", func(t *testing.T) {
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/completed", challengerToken, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			response := ParseJSONResponse(t, w)
			challenges := response["challenges"].([]interface{})

			found := false
			for _, challenge := range challenges {
				challengeMap := challenge.(map[string]interface{})
				if challengeMap["id"].(string) == completedChallengeID.String() {
					found = true
					assert.Equal(t, "completed", challengeMap["status"])
				}
			}
			assert.True(t, found, "Should find our completed challenge")
		})

		t.Run("Get Challenge Stats", func(t *testing.T) {
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/stats", challengerToken, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			response := ParseJSONResponse(t, w)

			// Verify stats structure
			assert.Contains(t, response, "total_challenges")
			assert.Contains(t, response, "pending_challenges")
			assert.Contains(t, response, "active_challenges")
			assert.Contains(t, response, "completed_challenges")
			assert.Contains(t, response, "won_challenges")
			assert.Contains(t, response, "lost_challenges")

			// Verify stats are numbers
			assert.IsType(t, float64(0), response["total_challenges"])
			assert.IsType(t, float64(0), response["pending_challenges"])
		})

		t.Run("Get Challenge by ID", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/challenges/%s", pendingChallengeID)
			w := MakeAuthenticatedRequest(t, tc, "GET", url, challengerToken, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			var response models.ChallengeDetailResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, pendingChallengeID, response.Challenge.ID)
		})

		t.Run("Prevent Access to Non-Participant Challenge", func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/challenges/%s", pendingChallengeID)
			w := MakeAuthenticatedRequest(t, tc, "GET", url, thirdUserToken, nil)
			assert.Equal(t, http.StatusForbidden, w.Code, "Non-participants should not access challenge details")
		})
	})

	t.Run("Challenge Filtering and Pagination", func(t *testing.T) {
		t.Run("Filter by Status", func(t *testing.T) {
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges?status=pending", challengerToken, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			var response models.ChallengeListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			for _, challenge := range response.Challenges {
				assert.Equal(t, "pending", string(challenge.Status))
			}
		})

		t.Run("Pagination", func(t *testing.T) {
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges?page=1&page_size=2", challengerToken, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			var response models.ChallengeListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.LessOrEqual(t, len(response.Challenges), 2, "Should respect page size")
			assert.Equal(t, 1, response.Page)
			assert.Equal(t, 2, response.PageSize)
		})

		t.Run("Filter by User Type", func(t *testing.T) {
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges?user_type=challenger", challengerToken, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			var response models.ChallengeListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			for _, challenge := range response.Challenges {
				assert.Equal(t, challengerID, challenge.ChallengerID)
			}
		})
	})

	t.Run("Authentication and Authorization", func(t *testing.T) {
		t.Run("Unauthenticated Request", func(t *testing.T) {
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges", "", nil)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("Invalid Token", func(t *testing.T) {
			w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges", "invalid.token.here", nil)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	})

	// Cleanup test users
	defer func() {
		CleanupTestUser(challengerID)
		CleanupTestUser(challengedID)
		CleanupTestUser(thirdUserID)
	}()
}

// Helper function to create a friendship between two users

// Helper function to get the first available quiz ID

// Helper function to create a test challenge and return its ID
func createTestChallenge(t *testing.T, tc *TestConfig, challengerToken string, challengedID, quizID uuid.UUID) uuid.UUID {
	createReq := models.CreateChallengeRequest{
		ChallengedUserID: challengedID,
		QuizID:           quizID,
		Message:          stringPtr("Test challenge"),
		IsGroupChallenge: false,
	}

	reqBody, _ := json.Marshal(createReq)
	w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", challengerToken, reqBody)
	require.Equal(t, http.StatusCreated, w.Code, "Should create test challenge")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	challenge := response["challenge"].(map[string]interface{})
	challengeIDStr := challenge["id"].(string)
	challengeID, err := uuid.Parse(challengeIDStr)
	require.NoError(t, err)

	return challengeID
}
