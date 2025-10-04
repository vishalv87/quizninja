package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"quizninja-api/database"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Database verification helper functions

func getSessionFromDB(t *testing.T, sessionID uuid.UUID) map[string]interface{} {
	query := `
		SELECT id, attempt_id, user_id, quiz_id, current_question_index,
		       current_answers, session_state, time_remaining, time_spent_so_far,
		       last_activity_at, paused_at, created_at, updated_at
		FROM quiz_sessions WHERE id = $1`

	row := database.DB.QueryRow(query, sessionID)

	var id, attemptID, userID, quizID uuid.UUID
	var currentQuestionIndex, timeSpentSoFar int
	var currentAnswersJSON []byte
	var sessionState string
	var timeRemaining *int
	var lastActivityAt, createdAt, updatedAt time.Time
	var pausedAt *time.Time

	err := row.Scan(&id, &attemptID, &userID, &quizID, &currentQuestionIndex,
		&currentAnswersJSON, &sessionState, &timeRemaining, &timeSpentSoFar,
		&lastActivityAt, &pausedAt, &createdAt, &updatedAt)
	require.NoError(t, err, "Failed to get session from database")

	var currentAnswers []map[string]interface{}
	if len(currentAnswersJSON) > 0 {
		err = json.Unmarshal(currentAnswersJSON, &currentAnswers)
		require.NoError(t, err, "Failed to unmarshal current answers")
	}

	return map[string]interface{}{
		"id":                     id,
		"attempt_id":             attemptID,
		"user_id":                userID,
		"quiz_id":                quizID,
		"current_question_index": currentQuestionIndex,
		"current_answers":        currentAnswers,
		"session_state":          sessionState,
		"time_remaining":         timeRemaining,
		"time_spent_so_far":      timeSpentSoFar,
		"last_activity_at":       lastActivityAt,
		"paused_at":              pausedAt,
		"created_at":             createdAt,
		"updated_at":             updatedAt,
	}
}

func getAttemptFromDB(t *testing.T, attemptID uuid.UUID) map[string]interface{} {
	query := `
		SELECT id, quiz_id, user_id, status, is_completed, created_at, updated_at
		FROM quiz_attempts WHERE id = $1`

	row := database.DB.QueryRow(query, attemptID)

	var id, quizID, userID uuid.UUID
	var status string
	var isCompleted bool
	var createdAt, updatedAt time.Time

	err := row.Scan(&id, &quizID, &userID, &status, &isCompleted, &createdAt, &updatedAt)
	require.NoError(t, err, "Failed to get attempt from database")

	return map[string]interface{}{
		"id":           id,
		"quiz_id":      quizID,
		"user_id":      userID,
		"status":       status,
		"is_completed": isCompleted,
		"created_at":   createdAt,
		"updated_at":   updatedAt,
	}
}

func verifySessionStateInDB(t *testing.T, sessionID uuid.UUID, expectedState string) {
	session := getSessionFromDB(t, sessionID)
	actualState := session["session_state"].(string)
	assert.Equal(t, expectedState, actualState, "Session state in database should match expected")
}

func verifyAttemptStatusInDB(t *testing.T, attemptID uuid.UUID, expectedStatus string) {
	attempt := getAttemptFromDB(t, attemptID)
	actualStatus := attempt["status"].(string)
	assert.Equal(t, expectedStatus, actualStatus, "Attempt status in database should match expected")
}

// Independent test helper functions

func createActiveSessionForTest(t *testing.T, tc *TestConfig) (string, uuid.UUID, uuid.UUID, uuid.UUID) {
	userID, token := CreateTestUser(t, tc)

	// Get a quiz to start
	w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes", token, nil)
	require.Equal(t, http.StatusOK, w.Code)

	response := ParseJSONResponse(t, w)
	data := GetDataFromResponse(t, response)
	quizzes := data["quizzes"].([]interface{})
	require.Greater(t, len(quizzes), 0, "Should have at least one quiz")

	firstQuiz := quizzes[0].(map[string]interface{})
	quizID := firstQuiz["id"].(string)
	quizUUID, err := uuid.Parse(quizID)
	require.NoError(t, err)

	// Start quiz attempt to create active session
	startURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts", quizID)
	startResponse := MakeAuthenticatedRequest(t, tc, "POST", startURL, token, nil)
	require.True(t, startResponse.Code == http.StatusCreated || startResponse.Code == http.StatusOK)

	startData := ParseJSONResponse(t, startResponse)
	startDataMap := GetDataFromResponse(t, startData)

	attemptID := startDataMap["attempt_id"].(string)

	attemptUUID, err := uuid.Parse(attemptID)
	require.NoError(t, err)

	return token, userID, quizUUID, attemptUUID
}

func createPausedSessionForTest(t *testing.T, tc *TestConfig) (string, uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID) {
	token, userID, quizID, attemptID := createActiveSessionForTest(t, tc)

	// Pause the session
	pauseURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/pause", quizID, attemptID)
	pausePayload := map[string]interface{}{
		"current_question_index": 1,
		"current_answers": []map[string]interface{}{
			{
				"question_id":     uuid.New().String(),
				"selected_option": 0,
				"text_answer":     "test answer",
				"is_correct":      false,
				"points_earned":   0,
			},
		},
		"time_spent_so_far": 60,
		"time_remaining":    300,
	}

	pauseBody, _ := json.Marshal(pausePayload)
	pauseResponse := MakeAuthenticatedRequest(t, tc, "POST", pauseURL, token, pauseBody)
	require.Equal(t, http.StatusOK, pauseResponse.Code)

	// Get session ID from database
	var sessionID uuid.UUID
	query := `SELECT id FROM quiz_sessions WHERE attempt_id = $1`
	err := database.DB.QueryRow(query, attemptID).Scan(&sessionID)
	require.NoError(t, err)

	return token, userID, quizID, attemptID, sessionID
}

// Test cases

func TestQuizSessionPauseFlow(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	// Setup: Create active session
	token, _, quizID, attemptID := createActiveSessionForTest(t, tc)

	// Get session ID before pause
	var sessionID uuid.UUID
	query := `SELECT id FROM quiz_sessions WHERE attempt_id = $1`
	err := database.DB.QueryRow(query, attemptID).Scan(&sessionID)
	require.NoError(t, err)

	// Verify initial state in DB
	verifySessionStateInDB(t, sessionID, "active")
	verifyAttemptStatusInDB(t, attemptID, "started")

	// Action: Pause the session
	pauseURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/pause", quizID, attemptID)
	pausePayload := map[string]interface{}{
		"current_question_index": 2,
		"current_answers": []map[string]interface{}{
			{
				"question_id":     uuid.New().String(),
				"selected_option": 1,
				"text_answer":     "test answer 1",
				"is_correct":      true,
				"points_earned":   1,
			},
			{
				"question_id":     uuid.New().String(),
				"selected_option": 0,
				"text_answer":     "test answer 2",
				"is_correct":      false,
				"points_earned":   0,
			},
		},
		"time_spent_so_far": 120,
		"time_remaining":    240,
	}

	pauseBody, _ := json.Marshal(pausePayload)
	w := MakeAuthenticatedRequest(t, tc, "POST", pauseURL, token, pauseBody)

	// Assert API Response
	assert.Equal(t, http.StatusOK, w.Code)

	response := ParseJSONResponse(t, w)
	data := GetDataFromResponse(t, response)

	assert.Equal(t, sessionID.String(), data["session_id"])
	assert.Equal(t, "paused", data["action"])
	assert.Equal(t, "paused", data["session_state"])
	assert.Contains(t, data["message"].(string), "paused")

	// Assert DB State
	verifySessionStateInDB(t, sessionID, "paused")
	verifyAttemptStatusInDB(t, attemptID, "paused")

	// Verify specific progress data in DB
	sessionData := getSessionFromDB(t, sessionID)
	assert.Equal(t, 2, sessionData["current_question_index"])
	assert.Equal(t, 120, sessionData["time_spent_so_far"])
	assert.Equal(t, 240, *sessionData["time_remaining"].(*int))
	assert.NotNil(t, sessionData["paused_at"])

	currentAnswers := sessionData["current_answers"].([]map[string]interface{})
	assert.Len(t, currentAnswers, 2)
	assert.Equal(t, float64(1), currentAnswers[0]["selected_option"]) // JSON numbers become float64
	assert.Equal(t, "test answer 1", currentAnswers[0]["text_answer"])
}

func TestQuizSessionResumeFlow(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	// Setup: Create paused session
	token, _, quizID, attemptID, sessionID := createPausedSessionForTest(t, tc)

	// Verify initial state in DB
	verifySessionStateInDB(t, sessionID, "paused")
	verifyAttemptStatusInDB(t, attemptID, "paused")

	// Action: Resume the session
	resumeURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/resume", quizID, attemptID)
	w := MakeAuthenticatedRequest(t, tc, "POST", resumeURL, token, nil)

	// Assert API Response
	assert.Equal(t, http.StatusOK, w.Code)

	response := ParseJSONResponse(t, w)
	data := GetDataFromResponse(t, response)

	assert.Equal(t, sessionID.String(), data["session_id"])
	assert.Equal(t, "resumed", data["action"])
	assert.Equal(t, "active", data["session_state"])
	assert.Contains(t, data["message"].(string), "resumed")

	// Verify quiz data is included for resumption
	quiz, exists := data["quiz"]
	assert.True(t, exists, "Resume response should include quiz data")
	quizMap := quiz.(map[string]interface{})
	assert.NotNil(t, quizMap["id"])
	assert.NotNil(t, quizMap["title"])

	// Verify session progress data is included
	assert.NotNil(t, data["current_question_index"])
	assert.NotNil(t, data["current_answers"])
	assert.NotNil(t, data["time_remaining"])

	// Assert DB State
	verifySessionStateInDB(t, sessionID, "active")
	verifyAttemptStatusInDB(t, attemptID, "started")

	// Verify paused_at timestamp was cleared
	sessionData := getSessionFromDB(t, sessionID)
	assert.Nil(t, sessionData["paused_at"])
}

func TestQuizSessionSaveProgressFlow(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	// Setup: Create active session
	token, _, quizID, attemptID := createActiveSessionForTest(t, tc)

	// Get session ID
	var sessionID uuid.UUID
	query := `SELECT id FROM quiz_sessions WHERE attempt_id = $1`
	err := database.DB.QueryRow(query, attemptID).Scan(&sessionID)
	require.NoError(t, err)

	// Verify initial state
	verifySessionStateInDB(t, sessionID, "active")

	// Get initial last_activity_at timestamp
	initialSessionData := getSessionFromDB(t, sessionID)
	initialLastActivity := initialSessionData["last_activity_at"].(time.Time)

	// Wait a moment to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Action: Save progress
	saveProgressURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/save-progress", quizID, attemptID)
	savePayload := map[string]interface{}{
		"current_question_index": 3,
		"current_answers": []map[string]interface{}{
			{
				"question_id":     uuid.New().String(),
				"selected_option": 2,
				"text_answer":     "saved answer 1",
				"is_correct":      true,
				"points_earned":   1,
			},
			{
				"question_id":     uuid.New().String(),
				"selected_option": 1,
				"text_answer":     "saved answer 2",
				"is_correct":      false,
				"points_earned":   0,
			},
			{
				"question_id":     uuid.New().String(),
				"selected_option": 0,
				"text_answer":     "saved answer 3",
				"is_correct":      true,
				"points_earned":   1,
			},
		},
		"time_spent_so_far": 180,
		"time_remaining":    420,
	}

	saveBody, _ := json.Marshal(savePayload)
	w := MakeAuthenticatedRequest(t, tc, "PUT", saveProgressURL, token, saveBody)

	// Assert API Response
	assert.Equal(t, http.StatusOK, w.Code)

	response := ParseJSONResponse(t, w)
	data := GetDataFromResponse(t, response)

	assert.Equal(t, sessionID.String(), data["session_id"])
	assert.Contains(t, data["message"].(string), "saved")
	assert.NotNil(t, data["timestamp"])

	// Assert DB State - session should still be active
	verifySessionStateInDB(t, sessionID, "active")
	verifyAttemptStatusInDB(t, attemptID, "started")

	// Verify progress data was saved correctly
	sessionData := getSessionFromDB(t, sessionID)
	assert.Equal(t, 3, sessionData["current_question_index"])
	assert.Equal(t, 180, sessionData["time_spent_so_far"])
	assert.Equal(t, 420, *sessionData["time_remaining"].(*int))

	// Verify last_activity_at was updated
	updatedLastActivity := sessionData["last_activity_at"].(time.Time)
	assert.True(t, updatedLastActivity.After(initialLastActivity), "last_activity_at should be updated")

	// Verify current_answers were saved correctly
	currentAnswers := sessionData["current_answers"].([]map[string]interface{})
	assert.Len(t, currentAnswers, 3)
	assert.Equal(t, float64(2), currentAnswers[0]["selected_option"])
	assert.Equal(t, "saved answer 1", currentAnswers[0]["text_answer"])
	assert.Equal(t, true, currentAnswers[0]["is_correct"])
	assert.Equal(t, float64(1), currentAnswers[0]["points_earned"])
}

func TestQuizSessionAbandonFlow(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	// Setup: Create active session
	token, _, quizID, attemptID := createActiveSessionForTest(t, tc)

	// Get session ID
	var sessionID uuid.UUID
	query := `SELECT id FROM quiz_sessions WHERE attempt_id = $1`
	err := database.DB.QueryRow(query, attemptID).Scan(&sessionID)
	require.NoError(t, err)

	// Verify initial state
	verifySessionStateInDB(t, sessionID, "active")
	verifyAttemptStatusInDB(t, attemptID, "started")

	// Action: Abandon the session
	abandonURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/abandon", quizID, attemptID)
	w := MakeAuthenticatedRequest(t, tc, "DELETE", abandonURL, token, nil)

	// Assert API Response
	assert.Equal(t, http.StatusOK, w.Code)

	response := ParseJSONResponse(t, w)
	data := GetDataFromResponse(t, response)

	assert.Equal(t, "abandoned", data["action"])
	assert.Equal(t, "abandoned", data["session_state"])
	assert.Contains(t, data["message"].(string), "abandoned")
	assert.Equal(t, float64(0), data["progress"]) // JSON numbers become float64

	// Assert DB State
	verifySessionStateInDB(t, sessionID, "abandoned")
	verifyAttemptStatusInDB(t, attemptID, "abandoned")

	// Verify session cannot be resumed after abandonment
	resumeURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/resume", quizID, attemptID)
	resumeResponse := MakeAuthenticatedRequest(t, tc, "POST", resumeURL, token, nil)
	assert.Equal(t, http.StatusConflict, resumeResponse.Code, "Should not be able to resume abandoned session")
}

func TestQuizSessionErrorScenarios(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	userID, token := CreateTestUser(t, tc)

	// Test pause on non-existent attempt
	t.Run("PauseNonExistentAttempt", func(t *testing.T) {
		fakeQuizID := uuid.New().String()
		fakeAttemptID := uuid.New().String()
		pauseURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/pause", fakeQuizID, fakeAttemptID)

		pausePayload := map[string]interface{}{
			"current_question_index": 0,
			"current_answers":        []map[string]interface{}{},
			"time_spent_so_far":      30,
		}

		pauseBody, _ := json.Marshal(pausePayload)
		w := MakeAuthenticatedRequest(t, tc, "POST", pauseURL, token, pauseBody)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// Test resume on already active session
	t.Run("ResumeActiveSession", func(t *testing.T) {
		// Create active session with the same user
		token2, _, quizID, attemptID := createActiveSessionForTest(t, tc)

		// Try to resume already active session
		resumeURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/resume", quizID, attemptID)
		w := MakeAuthenticatedRequest(t, tc, "POST", resumeURL, token2, nil)
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	// Test unauthorized access to other user's session
	t.Run("UnauthorizedSessionAccess", func(t *testing.T) {
		// Create session with first user
		_, _, quizID, attemptID := createActiveSessionForTest(t, tc)

		// Create second user
		_, secondToken := CreateTestUser(t, tc)

		// Try to pause other user's session
		pauseURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/pause", quizID, attemptID)
		pausePayload := map[string]interface{}{
			"current_question_index": 0,
			"current_answers":        []map[string]interface{}{},
			"time_spent_so_far":      30,
		}

		pauseBody, _ := json.Marshal(pausePayload)
		w := MakeAuthenticatedRequest(t, tc, "POST", pauseURL, secondToken, pauseBody)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	// Test abandon on completed attempt
	t.Run("AbandonCompletedAttempt", func(t *testing.T) {
		// Create active session with the same user
		token3, _, quizID, attemptID := createActiveSessionForTest(t, tc)

		// Mark attempt as completed in database
		updateQuery := `UPDATE quiz_attempts SET is_completed = true, status = 'completed' WHERE id = $1`
		_, err := database.DB.Exec(updateQuery, attemptID)
		require.NoError(t, err)

		// Try to abandon completed attempt
		abandonURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/abandon", quizID, attemptID)
		w := MakeAuthenticatedRequest(t, tc, "DELETE", abandonURL, token3, nil)
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	// Test save progress on non-existent session
	t.Run("SaveProgressNonExistentSession", func(t *testing.T) {
		fakeQuizID := uuid.New().String()
		fakeAttemptID := uuid.New().String()
		saveURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/save-progress", fakeQuizID, fakeAttemptID)

		savePayload := map[string]interface{}{
			"current_question_index": 1,
			"current_answers":        []map[string]interface{}{},
			"time_spent_so_far":      60,
		}

		saveBody, _ := json.Marshal(savePayload)
		w := MakeAuthenticatedRequest(t, tc, "PUT", saveURL, token, saveBody)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	_ = userID // Use userID to avoid unused variable warning
}

func TestQuizSessionPauseToResumeCompleteFlow(t *testing.T) {
	tc := SetupTestServer(t)
	defer CleanupWithSupabase(t, tc)

	// Setup: Create active session
	token, _, quizID, attemptID := createActiveSessionForTest(t, tc)

	// Get session ID
	var sessionID uuid.UUID
	query := `SELECT id FROM quiz_sessions WHERE attempt_id = $1`
	err := database.DB.QueryRow(query, attemptID).Scan(&sessionID)
	require.NoError(t, err)

	// Step 1: Pause the session
	pauseURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/pause", quizID, attemptID)
	pausePayload := map[string]interface{}{
		"current_question_index": 1,
		"current_answers": []map[string]interface{}{
			{
				"question_id":     uuid.New().String(),
				"selected_option": 1,
				"text_answer":     "flow test answer",
				"is_correct":      true,
				"points_earned":   1,
			},
		},
		"time_spent_so_far": 90,
		"time_remaining":    270,
	}

	pauseBody, _ := json.Marshal(pausePayload)
	pauseResponse := MakeAuthenticatedRequest(t, tc, "POST", pauseURL, token, pauseBody)
	assert.Equal(t, http.StatusOK, pauseResponse.Code)

	// Verify paused state in DB
	verifySessionStateInDB(t, sessionID, "paused")
	verifyAttemptStatusInDB(t, attemptID, "paused")

	// Step 2: Resume the session
	resumeURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/resume", quizID, attemptID)
	resumeResponse := MakeAuthenticatedRequest(t, tc, "POST", resumeURL, token, nil)
	assert.Equal(t, http.StatusOK, resumeResponse.Code)

	// Verify resumed state in DB
	verifySessionStateInDB(t, sessionID, "active")
	verifyAttemptStatusInDB(t, attemptID, "started")

	// Verify resume response includes saved progress
	resumeData := ParseJSONResponse(t, resumeResponse)
	resumeDataMap := GetDataFromResponse(t, resumeData)
	assert.Equal(t, float64(1), resumeDataMap["current_question_index"])
	assert.Equal(t, float64(90), resumeDataMap["time_spent_so_far"])
	assert.Equal(t, float64(270), resumeDataMap["time_remaining"])

	currentAnswers := resumeDataMap["current_answers"].([]interface{})
	assert.Len(t, currentAnswers, 1)
	firstAnswer := currentAnswers[0].(map[string]interface{})
	assert.Equal(t, "flow test answer", firstAnswer["text_answer"])

	// Step 3: Save additional progress
	saveURL := fmt.Sprintf("/api/v1/quizzes/%s/attempts/%s/save-progress", quizID, attemptID)
	savePayload := map[string]interface{}{
		"current_question_index": 2,
		"current_answers": []map[string]interface{}{
			{
				"question_id":     uuid.New().String(),
				"selected_option": 1,
				"text_answer":     "flow test answer",
				"is_correct":      true,
				"points_earned":   1,
			},
			{
				"question_id":     uuid.New().String(),
				"selected_option": 0,
				"text_answer":     "second flow answer",
				"is_correct":      false,
				"points_earned":   0,
			},
		},
		"time_spent_so_far": 150,
		"time_remaining":    210,
	}

	saveBody, _ := json.Marshal(savePayload)
	saveResponse := MakeAuthenticatedRequest(t, tc, "PUT", saveURL, token, saveBody)
	assert.Equal(t, http.StatusOK, saveResponse.Code)

	// Verify final state in DB
	finalSessionData := getSessionFromDB(t, sessionID)
	assert.Equal(t, "active", finalSessionData["session_state"])
	assert.Equal(t, 2, finalSessionData["current_question_index"])
	assert.Equal(t, 150, finalSessionData["time_spent_so_far"])
	assert.Equal(t, 210, *finalSessionData["time_remaining"].(*int))

	finalAnswers := finalSessionData["current_answers"].([]map[string]interface{})
	assert.Len(t, finalAnswers, 2)
	assert.Equal(t, "second flow answer", finalAnswers[1]["text_answer"])
}
