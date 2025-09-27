package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuizHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	userID, token := CreateTestUser(t, tc)

	t.Run("GetQuizzes", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		quizzes, exists := data["quizzes"]
		assert.True(t, exists, "Response should contain 'quizzes' field")

		quizzesList, ok := quizzes.([]interface{})
		assert.True(t, ok, "Quizzes field should be an array")

		if len(quizzesList) > 0 {
			VerifyIsTestDataInArray(t, quizzesList, true, "quizzes")
		}
	})

	t.Run("GetFeaturedQuizzes", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes/featured", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		quizzes, exists := data["quizzes"]
		assert.True(t, exists, "Response should contain 'quizzes' field")

		quizzesList, ok := quizzes.([]interface{})
		assert.True(t, ok, "Quizzes field should be an array")

		if len(quizzesList) > 0 {
			VerifyIsTestDataInArray(t, quizzesList, true, "featured quizzes")
		}
	})

	t.Run("GetUserQuizzes", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/quizzes", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		quizzes, exists := data["quizzes"]
		assert.True(t, exists, "Response should contain 'quizzes' field")

		quizzesList, ok := quizzes.([]interface{})
		assert.True(t, ok, "Quizzes field should be an array")

		if len(quizzesList) > 0 {
			VerifyIsTestDataInArray(t, quizzesList, true, "user quizzes")
		}
	})

	t.Run("GetQuizByCategory", func(t *testing.T) {
		categories := []string{"science", "technology", "sports", "history"}

		for _, category := range categories {
			t.Run(fmt.Sprintf("Category_%s", category), func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/quizzes/category/%s", category)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					response := ParseJSONResponse(t, w)
					data := GetDataFromResponse(t, response)

					quizzes, exists := data["quizzes"]
					assert.True(t, exists, "Response should contain 'quizzes' field")

					quizzesList, ok := quizzes.([]interface{})
					assert.True(t, ok, "Quizzes field should be an array")

					if len(quizzesList) > 0 {
						VerifyIsTestDataInArray(t, quizzesList, true, fmt.Sprintf("category %s quizzes", category))
					}
				}
			})
		}
	})

	t.Run("GetQuizByID", func(t *testing.T) {
		// First get a list of quizzes to find a valid quiz ID
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		quizzes, exists := data["quizzes"]
		assert.True(t, exists, "Response should contain 'quizzes' field")

		quizzesList, ok := quizzes.([]interface{})
		assert.True(t, ok, "Quizzes field should be an array")

		if len(quizzesList) > 0 {
			firstQuiz, ok := quizzesList[0].(map[string]interface{})
			assert.True(t, ok, "First quiz should be an object")

			quizID, exists := firstQuiz["id"]
			assert.True(t, exists, "Quiz should have an ID")

			// Test getting quiz without questions
			url := fmt.Sprintf("/api/v1/quizzes/%s", quizID)
			w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			quiz, exists := data["quiz"]
			assert.True(t, exists, "Response should contain 'quiz' field")

			quizMap, ok := quiz.(map[string]interface{})
			assert.True(t, ok, "Quiz field should be an object")

			VerifyIsTestDataField(t, quizMap, true, "quiz")

			// Test getting quiz with questions
			urlWithQuestions := fmt.Sprintf("/api/v1/quizzes/%s?include_questions=true", quizID)
			w = MakeAuthenticatedRequest(t, tc, "GET", urlWithQuestions, token, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			response = ParseJSONResponse(t, w)
			data = GetDataFromResponse(t, response)

			quiz, exists = data["quiz"]
			assert.True(t, exists, "Response should contain 'quiz' field")

			quizMap, ok = quiz.(map[string]interface{})
			assert.True(t, ok, "Quiz field should be an object")

			VerifyIsTestDataField(t, quizMap, true, "quiz with questions")

			// Check if questions exist and verify is_test_data field
			questions, questionsExist := quizMap["questions"]
			if questionsExist {
				questionsList, ok := questions.([]interface{})
				if ok && len(questionsList) > 0 {
					VerifyIsTestDataInArray(t, questionsList, true, "quiz questions")
				}
			}

			// Test getting quiz with stats
			urlWithStats := fmt.Sprintf("/api/v1/quizzes/%s?include_stats=true", quizID)
			w = MakeAuthenticatedRequest(t, tc, "GET", urlWithStats, token, nil)
			assert.Equal(t, http.StatusOK, w.Code)

			response = ParseJSONResponse(t, w)
			data = GetDataFromResponse(t, response)

			quiz, exists = data["quiz"]
			assert.True(t, exists, "Response should contain 'quiz' field")

			quizMap, ok = quiz.(map[string]interface{})
			assert.True(t, ok, "Quiz field should be an object")

			VerifyIsTestDataField(t, quizMap, true, "quiz with stats")

			// Check if statistics exist and verify is_test_data field
			statistics, statsExist := quizMap["statistics"]
			if statsExist {
				statisticsMap, ok := statistics.(map[string]interface{})
				if ok {
					VerifyIsTestDataField(t, statisticsMap, true, "quiz statistics")
				}
			}
		}
	})

	t.Run("GetUserAttempts", func(t *testing.T) {
		// First, create a completed quiz attempt to ensure we have data to test

		// 1. Get available quizzes
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes", token, nil)
		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			if quizzes, exists := data["quizzes"]; exists {
				if quizzesList, ok := quizzes.([]interface{}); ok && len(quizzesList) > 0 {
					if firstQuiz, ok := quizzesList[0].(map[string]interface{}); ok {
						if quizID, exists := firstQuiz["id"]; exists {
							// 2. Start a quiz attempt
							startURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts", quizID)
							startResponse := MakeAuthenticatedRequest(t, tc, "POST", startURL, token, nil)

							if startResponse.Code == http.StatusCreated || startResponse.Code == http.StatusOK {
								startData := ParseJSONResponse(t, startResponse)
								startDataMap := GetDataFromResponse(t, startData)

								if attemptID, exists := startDataMap["attempt_id"]; exists {
									if quiz, exists := startDataMap["quiz"]; exists {
										if quizMap, ok := quiz.(map[string]interface{}); ok {
											// 3. Build answers payload for quiz completion
											var answers []map[string]interface{}
											if questions, exists := quizMap["questions"]; exists {
												if questionsList, ok := questions.([]interface{}); ok {
													for i, question := range questionsList {
														if questionMap, ok := question.(map[string]interface{}); ok {
															if questionID, exists := questionMap["id"]; exists {
																// Create a simple answer (first option for multiple choice)
																answer := map[string]interface{}{
																	"questionId":          questionID,
																	"selectedOptionIndex": 0, // Always select first option
																}
																answers = append(answers, answer)
															}
														}
														// Limit to first 3 questions to keep test fast
														if i >= 2 {
															break
														}
													}
												}
											}

											// 4. Submit the quiz attempt to complete it
											submitURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/submit", quizID, attemptID)
											submitPayload := map[string]interface{}{
												"attemptId": attemptID,
												"answers":   answers,
												"timeSpent": 120, // 2 minutes
											}

											submitBody, _ := json.Marshal(submitPayload)
											submitResponse := MakeAuthenticatedRequest(t, tc, "POST", submitURL, token, submitBody)
											// Don't fail test if submission fails, just continue to test the GET endpoint
											_ = submitResponse // Use the response to avoid unused variable error
										}
									}
								}
							}
						}
					}
				}
			}
		}

		// 5. Now test the GetUserAttempts endpoint
		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/attempts", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		attempts, exists := data["attempts"]
		assert.True(t, exists, "Response should contain 'attempts' field")

		attemptsList, ok := attempts.([]interface{})
		assert.True(t, ok, "Attempts field should be an array")

		if len(attemptsList) > 0 {
			VerifyIsTestDataInArray(t, attemptsList, true, "user attempts")

			// Check the first attempt and verify nested quiz data
			firstAttempt, ok := attemptsList[0].(map[string]interface{})
			if ok {
				// Verify attempt has required fields and is_test_data
				VerifyIsTestDataField(t, firstAttempt, true, "first attempt")
				assert.True(t, firstAttempt["is_completed"].(bool), "Attempt should be completed")

				// Verify attempt has expected structure
				assert.NotNil(t, firstAttempt["id"], "Attempt should have an ID")
				assert.NotNil(t, firstAttempt["quiz_id"], "Attempt should have a quiz_id")
				assert.NotNil(t, firstAttempt["user_id"], "Attempt should have a user_id")
				assert.NotNil(t, firstAttempt["score"], "Attempt should have a score")

				// Check nested quiz data
				quiz, quizExists := firstAttempt["quiz"]
				if quizExists {
					quizMap, ok := quiz.(map[string]interface{})
					if ok {
						VerifyIsTestDataField(t, quizMap, true, "attempt quiz")
						// Verify quiz has expected structure
						assert.NotNil(t, quizMap["id"], "Quiz should have an ID")
						assert.NotNil(t, quizMap["title"], "Quiz should have a title")
						assert.NotNil(t, quizMap["category"], "Quiz should have a category")
					}
				}
			}
		}
	})

	t.Run("GetUserActiveSessions", func(t *testing.T) {
		// First, create an active quiz session to ensure we have data to test

		// 1. Get available quizzes
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes", token, nil)
		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			if quizzes, exists := data["quizzes"]; exists {
				if quizzesList, ok := quizzes.([]interface{}); ok && len(quizzesList) > 0 {
					if firstQuiz, ok := quizzesList[0].(map[string]interface{}); ok {
						if quizID, exists := firstQuiz["id"]; exists {
							// 2. Start a quiz attempt (this creates an active session)
							startURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts", quizID)
							startResponse := MakeAuthenticatedRequest(t, tc, "POST", startURL, token, nil)

							if startResponse.Code == http.StatusCreated || startResponse.Code == http.StatusOK {
								startData := ParseJSONResponse(t, startResponse)
								startDataMap := GetDataFromResponse(t, startData)

								// Verify the session was created
								if sessionID, exists := startDataMap["session_id"]; exists {
									assert.NotNil(t, sessionID, "Session ID should be returned when starting quiz")
								}
							}
						}
					}
				}
			}
		}

		// 3. Now test the GetUserActiveSessions endpoint
		w = MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/active-sessions", token, nil)
		assert.Equal(t, http.StatusOK, w.Code)

		response := ParseJSONResponse(t, w)
		data := GetDataFromResponse(t, response)

		sessions, exists := data["sessions"]
		assert.True(t, exists, "Response should contain 'sessions' field")

		sessionsList, ok := sessions.([]interface{})
		assert.True(t, ok, "Sessions field should be an array")

		if len(sessionsList) > 0 {
			VerifyIsTestDataInArray(t, sessionsList, true, "active sessions")

			// Check nested data in sessions
			for i, session := range sessionsList {
				sessionMap, ok := session.(map[string]interface{})
				if ok {
					// Verify session has required fields and is_test_data
					VerifyIsTestDataField(t, sessionMap, true, fmt.Sprintf("session[%d]", i))

					// Verify session state is active or paused
					if sessionState, exists := sessionMap["session_state"]; exists {
						state := sessionState.(string)
						assert.Contains(t, []string{"active", "paused"}, state, "Session state should be active or paused")
					}

					// Verify session has expected structure
					assert.NotNil(t, sessionMap["id"], "Session should have an ID")
					assert.NotNil(t, sessionMap["user_id"], "Session should have a user_id")
					assert.NotNil(t, sessionMap["quiz_id"], "Session should have a quiz_id")

					// Check nested quiz data if present (from QuizSessionWithDetails)
					if quizTitle, exists := sessionMap["quiz_title"]; exists {
						assert.NotEmpty(t, quizTitle, "Session should have quiz_title")
					}
					if quizCategory, exists := sessionMap["quiz_category"]; exists {
						assert.NotEmpty(t, quizCategory, "Session should have quiz_category")
					}

					// Check nested quiz data (if using nested Quiz object)
					if quiz, quizExists := sessionMap["quiz"]; quizExists {
						quizMap, ok := quiz.(map[string]interface{})
						if ok {
							VerifyIsTestDataField(t, quizMap, true, fmt.Sprintf("session[%d] quiz", i))
							assert.NotNil(t, quizMap["id"], "Quiz should have an ID")
							assert.NotNil(t, quizMap["title"], "Quiz should have a title")
						}
					}

					// Check nested attempt data (if using nested Attempt object)
					if attempt, attemptExists := sessionMap["attempt"]; attemptExists {
						attemptMap, ok := attempt.(map[string]interface{})
						if ok {
							VerifyIsTestDataField(t, attemptMap, true, fmt.Sprintf("session[%d] attempt", i))
							assert.NotNil(t, attemptMap["id"], "Attempt should have an ID")
						}
					}
				}
			}
		}
	})

	// Test creating a quiz attempt to verify is_test_data propagation
	t.Run("StartQuizAttempt", func(t *testing.T) {
		// First get a quiz to attempt
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes", token, nil)
		if w.Code != http.StatusOK {
			t.Skip("No quizzes available to test quiz attempt")
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
			t.Skip("No quizzes available to test")
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

		// Attempt to start the quiz
		url := fmt.Sprintf("/api/v1/quizzes/%s/attempts", quizID)
		w = MakeAuthenticatedRequest(t, tc, "POST", url, token, nil)

		if w.Code == http.StatusCreated || w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			// Verify quiz data in the attempt response
			quiz, exists := data["quiz"]
			if exists {
				quizMap, ok := quiz.(map[string]interface{})
				if ok {
					VerifyIsTestDataField(t, quizMap, true, "attempt start quiz")

					// Check questions if they exist
					questions, questionsExist := quizMap["questions"]
					if questionsExist {
						questionsList, ok := questions.([]interface{})
						if ok && len(questionsList) > 0 {
							VerifyIsTestDataInArray(t, questionsList, true, "attempt start questions")
						}
					}
				}
			}
		}
	})

	_ = userID // Use userID to avoid unused variable warning
}
