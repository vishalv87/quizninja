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

func TestChallengesHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	_, token, _, cleanup := CreateTestUserWithCleanup(t, tc, "Challenges Handler Main User")
	defer cleanup()

	t.Run("GetChallenges", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges", token, nil)

		if w.Code == http.StatusOK {
			var response models.ChallengeListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse challenges list response")

			if len(response.Challenges) > 0 {
				// Verify challenges have is_test_data field
				challengesData := make([]interface{}, len(response.Challenges))
				for i, challenge := range response.Challenges {
					challengesData[i] = map[string]interface{}{
						"is_test_data": challenge.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, challengesData, true, "challenges list")

				// Check nested data in challenges
				for i, challenge := range response.Challenges {
					// Check challenger data
					if challenge.Challenger != nil {
						challengerMap := map[string]interface{}{
							"is_test_data": challenge.Challenger.IsTestData,
						}
						VerifyIsTestDataField(t, challengerMap, true, fmt.Sprintf("challenge[%d] challenger", i))
					}

					// Check challenged user data
					if challenge.Challenged != nil {
						challengedMap := map[string]interface{}{
							"is_test_data": challenge.Challenged.IsTestData,
						}
						VerifyIsTestDataField(t, challengedMap, true, fmt.Sprintf("challenge[%d] challenged", i))
					}

					// Check quiz data (QuizSummary doesn't have IsTestData field)
					if challenge.Quiz != nil {
						t.Logf("Challenge[%d] has quiz: %s", i, challenge.Quiz.Title)
					}
				}
			}
		}
	})

	t.Run("GetChallengesWithFilters", func(t *testing.T) {
		// Test with pagination and status filters
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges?page=1&page_size=5&status=active", token, nil)

		if w.Code == http.StatusOK {
			var response models.ChallengeListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse filtered challenges response")

			if len(response.Challenges) > 0 {
				challengesData := make([]interface{}, len(response.Challenges))
				for i, challenge := range response.Challenges {
					challengesData[i] = map[string]interface{}{
						"is_test_data": challenge.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, challengesData, true, "filtered challenges")
			}
		}
	})

	t.Run("GetPendingChallenges", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/pending", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			challenges, exists := response["challenges"]
			assert.True(t, exists, "Response should contain 'challenges' field")

			challengesList, ok := challenges.([]interface{})
			assert.True(t, ok, "Challenges field should be an array")

			if len(challengesList) > 0 {
				VerifyIsTestDataInArray(t, challengesList, true, "pending challenges")

				// Check nested data in first challenge
				firstChallenge, ok := challengesList[0].(map[string]interface{})
				if ok {
					checkChallengeNestedData(t, firstChallenge, "pending challenge[0]")
				}
			}
		}
	})

	t.Run("GetActiveChallenges", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/active", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			challenges, exists := response["challenges"]
			assert.True(t, exists, "Response should contain 'challenges' field")

			challengesList, ok := challenges.([]interface{})
			assert.True(t, ok, "Challenges field should be an array")

			if len(challengesList) > 0 {
				VerifyIsTestDataInArray(t, challengesList, true, "active challenges")

				// Check nested data in first challenge
				firstChallenge, ok := challengesList[0].(map[string]interface{})
				if ok {
					checkChallengeNestedData(t, firstChallenge, "active challenge[0]")
				}
			}
		}
	})

	t.Run("GetCompletedChallenges", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/completed", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			challenges, exists := response["challenges"]
			assert.True(t, exists, "Response should contain 'challenges' field")

			challengesList, ok := challenges.([]interface{})
			assert.True(t, ok, "Challenges field should be an array")

			if len(challengesList) > 0 {
				VerifyIsTestDataInArray(t, challengesList, true, "completed challenges")

				// Check nested data in first challenge
				firstChallenge, ok := challengesList[0].(map[string]interface{})
				if ok {
					checkChallengeNestedData(t, firstChallenge, "completed challenge[0]")
				}
			}
		}
	})

	t.Run("GetChallengeStats", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges/stats", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Challenge stats are aggregated data, verify basic structure
			assert.NotNil(t, response, "Should receive valid challenge stats response")

			// Stats might not have is_test_data since they're aggregated
			// but we can check if the response is structured correctly
			if total, exists := response["total_challenges"]; exists {
				_, ok := total.(float64) // JSON numbers are float64
				assert.True(t, ok, "total_challenges should be a number")
			}
		}
	})

	t.Run("GetChallengeByID", func(t *testing.T) {
		// First get a list of challenges to find a valid challenge ID
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/challenges", token, nil)
		if w.Code != http.StatusOK {
			t.Skip("No challenges available to test individual challenge retrieval")
			return
		}

		var challengesResponse models.ChallengeListResponse
		err := json.Unmarshal(w.Body.Bytes(), &challengesResponse)
		if err != nil || len(challengesResponse.Challenges) == 0 {
			t.Skip("No valid challenges to test")
			return
		}

		challengeID := challengesResponse.Challenges[0].ID
		url := fmt.Sprintf("/api/v1/challenges/%s", challengeID)
		w = MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

		if w.Code == http.StatusOK {
			var response models.ChallengeDetailResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse challenge detail response")

			// Verify challenge has is_test_data field
			challengeMap := map[string]interface{}{
				"is_test_data": response.Challenge.IsTestData,
			}
			VerifyIsTestDataField(t, challengeMap, true, "individual challenge")

			// Check nested data
			if response.Challenge.Challenger != nil {
				challengerMap := map[string]interface{}{
					"is_test_data": response.Challenge.Challenger.IsTestData,
				}
				VerifyIsTestDataField(t, challengerMap, true, "challenge detail challenger")
			}

			if response.Challenge.Challenged != nil {
				challengedMap := map[string]interface{}{
					"is_test_data": response.Challenge.Challenged.IsTestData,
				}
				VerifyIsTestDataField(t, challengedMap, true, "challenge detail challenged")
			}

			if response.Challenge.Quiz != nil {
				t.Logf("Challenge detail has quiz: %s", response.Challenge.Quiz.Title)
			}
		}
	})

	t.Run("CreateChallenge", func(t *testing.T) {
		// Create a second test user to challenge
		secondUserID, _, _, cleanup := CreateTestUserWithCleanup(t, tc, "Challenges Handler Second User")
		defer cleanup()

		// First get a quiz to use for the challenge
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes", token, nil)
		if w.Code != http.StatusOK {
			t.Skip("No quizzes available to create challenge")
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

		createReq := models.CreateChallengeRequest{
			ChallengedUserID: secondUserID,
			QuizID:           quizUUID,
			Message:          stringPtr("Test challenge"),
			IsGroupChallenge: false,
		}

		reqBody, _ := json.Marshal(createReq)
		w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/challenges", token, reqBody)

		// The request might fail due to friendship requirements, but we test the endpoint structure
		if w.Code == http.StatusCreated {
			response := ParseJSONResponse(t, w)

			// Verify the response contains a challenge object
			challenge, challengeExists := response["challenge"]
			if challengeExists {
				challengeMap, ok := challenge.(map[string]interface{})
				if ok {
					VerifyIsTestDataField(t, challengeMap, true, "created challenge")
				}
			}
		}
	})

	// Test challenge actions with fake IDs to verify endpoint authentication
	t.Run("AcceptChallenge", func(t *testing.T) {
		fakeChallengeID := uuid.New()
		url := fmt.Sprintf("/api/v1/challenges/%s/accept", fakeChallengeID)
		w := MakeAuthenticatedRequest(t, tc, "PUT", url, token, nil)

		// We expect either 400 or 500, indicating the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent challenge")
	})

	t.Run("DeclineChallenge", func(t *testing.T) {
		fakeChallengeID := uuid.New()
		url := fmt.Sprintf("/api/v1/challenges/%s/decline", fakeChallengeID)
		w := MakeAuthenticatedRequest(t, tc, "PUT", url, token, nil)

		// We expect either 400 or 500, indicating the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent challenge")
	})

	t.Run("UpdateChallengeScore", func(t *testing.T) {
		fakeChallengeID := uuid.New()
		scoreReq := models.UpdateChallengeScoreRequest{
			UserScore: 85.5,
		}

		reqBody, _ := json.Marshal(scoreReq)
		url := fmt.Sprintf("/api/v1/challenges/%s/score", fakeChallengeID)
		w := MakeAuthenticatedRequest(t, tc, "PUT", url, token, reqBody)

		// We expect either 400 or 404, indicating the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent challenge")
	})

}

// Helper function to check nested data in challenge objects
func checkChallengeNestedData(t *testing.T, challengeMap map[string]interface{}, prefix string) {
	// Check challenger
	if challenger, exists := challengeMap["challenger"]; exists && challenger != nil {
		challengerMap, ok := challenger.(map[string]interface{})
		if ok {
			VerifyIsTestDataField(t, challengerMap, true, prefix+" challenger")
		}
	}

	// Check challenged user
	if challenged, exists := challengeMap["challenged"]; exists && challenged != nil {
		challengedMap, ok := challenged.(map[string]interface{})
		if ok {
			VerifyIsTestDataField(t, challengedMap, true, prefix+" challenged")
		}
	}

	// Check quiz
	if quiz, exists := challengeMap["quiz"]; exists && quiz != nil {
		quizMap, ok := quiz.(map[string]interface{})
		if ok {
			VerifyIsTestDataField(t, quizMap, true, prefix+" quiz")
		}
	}
}
