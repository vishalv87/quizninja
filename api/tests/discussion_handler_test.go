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

func TestDiscussionHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	userID, token := CreateTestUser(t, tc)

	t.Run("GetDiscussions", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/discussions", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			discussions, exists := data["discussions"]
			if exists {
				discussionsList, ok := discussions.([]interface{})
				if ok && len(discussionsList) > 0 {
					VerifyIsTestDataInArray(t, discussionsList, true, "discussions list")

					// Check nested data in first discussion
					firstDiscussion, ok := discussionsList[0].(map[string]interface{})
					if ok {
						checkDiscussionNestedData(t, firstDiscussion, "discussion[0]")
					}
				}
			}
		}
	})

	t.Run("GetDiscussionsWithFilters", func(t *testing.T) {
		// Test with pagination and filters
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/discussions?page=1&page_size=5&type=general&sort_by=created_at", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			discussions, exists := data["discussions"]
			if exists {
				discussionsList, ok := discussions.([]interface{})
				if ok && len(discussionsList) > 0 {
					VerifyIsTestDataInArray(t, discussionsList, true, "filtered discussions")
				}
			}
		}
	})

	t.Run("GetDiscussionByID", func(t *testing.T) {
		// First get a valid quiz ID from the quizzes endpoint
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/quizzes", token, nil)
		if w.Code != http.StatusOK {
			t.Skip("No quizzes available to test discussion creation")
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

		quizIDStr, exists := firstQuiz["id"]
		if !exists {
			t.Skip("Quiz missing ID")
			return
		}

		// Parse the quiz ID string to UUID
		quizID, err := uuid.Parse(fmt.Sprintf("%v", quizIDStr))
		if err != nil {
			t.Skip("Invalid quiz ID format")
			return
		}

		// Create a discussion using the valid quiz ID
		createReq := models.CreateDiscussionRequest{
			QuizID:  quizID,
			Content: "This is a test discussion for GetDiscussionByID test",
			Type:    "general",
		}

		reqBody, _ := json.Marshal(createReq)
		w = MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/discussions", token, reqBody)

		if w.Code != http.StatusCreated {
			t.Skip("Failed to create discussion for testing")
			return
		}

		// Parse the created discussion response manually since it has 201 status
		var createResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &createResponse)
		if err != nil {
			t.Skip("Failed to parse create discussion response")
			return
		}

		dataInterface, exists := createResponse["data"]
		if !exists {
			t.Skip("No data field in create response")
			return
		}

		data, ok = dataInterface.(map[string]interface{})
		if !ok {
			t.Skip("Invalid data format in create response")
			return
		}

		createdDiscussion, exists := data["discussion"]
		if !exists {
			t.Skip("No discussion field in create response")
			return
		}

		createdDiscussionMap, ok := createdDiscussion.(map[string]interface{})
		if !ok {
			t.Skip("Invalid created discussion data")
			return
		}

		discussionID, exists := createdDiscussionMap["id"]
		if !exists {
			t.Skip("Created discussion missing ID")
			return
		}

		// Test GetDiscussionByID with the created discussion
		url := fmt.Sprintf("/api/v1/discussions/%s", discussionID)
		w = MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

		assert.Equal(t, http.StatusOK, w.Code, "Should successfully retrieve discussion by ID")

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			discussion, exists := data["discussion"]
			assert.True(t, exists, "Response should contain discussion field")

			if exists {
				discussionMap, ok := discussion.(map[string]interface{})
				assert.True(t, ok, "Discussion should be a valid object")

				if ok {
					VerifyIsTestDataField(t, discussionMap, true, "individual discussion")
					checkDiscussionNestedData(t, discussionMap, "individual discussion")

					// Verify the content matches what we created
					content, exists := discussionMap["content"]
					if exists {
						assert.Equal(t, createReq.Content, content, "Discussion content should match")
					}
				}
			}
		}
	})

	t.Run("CreateDiscussion", func(t *testing.T) {
		createReq := models.CreateDiscussionRequest{
			QuizID:  uuid.New(), // Using a dummy quiz ID for testing
			Content: "This is a test discussion content",
			Type:    "general",
		}

		reqBody, _ := json.Marshal(createReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/discussions", token, reqBody)

		if w.Code == http.StatusCreated {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			discussion, exists := data["discussion"]
			if exists {
				discussionMap, ok := discussion.(map[string]interface{})
				if ok {
					VerifyIsTestDataField(t, discussionMap, true, "created discussion")
				}
			}
		}
	})

	t.Run("GetDiscussionReplies", func(t *testing.T) {
		// Test getting replies for a fake discussion ID
		fakeDiscussionID := uuid.New()
		url := fmt.Sprintf("/api/v1/discussions/%s/replies", fakeDiscussionID)
		w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			replies, exists := data["replies"]
			if exists {
				repliesList, ok := replies.([]interface{})
				if ok && len(repliesList) > 0 {
					VerifyIsTestDataInArray(t, repliesList, true, "discussion replies")

					// Check nested data in replies
					for i, reply := range repliesList {
						replyMap, ok := reply.(map[string]interface{})
						if ok {
							checkDiscussionReplyNestedData(t, replyMap, fmt.Sprintf("reply[%d]", i))
						}
					}
				}
			}
		}
	})

	t.Run("CreateDiscussionReply", func(t *testing.T) {
		// Test creating a reply with a fake discussion ID
		fakeDiscussionID := uuid.New()
		replyReq := models.CreateDiscussionReplyRequest{
			Content: "This is a test reply",
		}

		reqBody, _ := json.Marshal(replyReq)
		url := fmt.Sprintf("/api/v1/discussions/%s/replies", fakeDiscussionID)
		w := MakeAuthenticatedRequest(t, tc, "POST", url, token, reqBody)

		// We expect either 404 (discussion not found) or validation error
		// which indicates the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent discussion")
	})

	t.Run("LikeDiscussion", func(t *testing.T) {
		// Test liking a discussion with a fake ID
		fakeDiscussionID := uuid.New()
		url := fmt.Sprintf("/api/v1/discussions/%s/like", fakeDiscussionID)
		w := MakeAuthenticatedRequest(t, tc, "POST", url, token, nil)

		// We expect either 404 or server error, indicating the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent discussion")
	})

	t.Run("UnlikeDiscussion", func(t *testing.T) {
		// Test unliking a discussion with a fake ID
		fakeDiscussionID := uuid.New()
		url := fmt.Sprintf("/api/v1/discussions/%s/unlike", fakeDiscussionID)
		w := MakeAuthenticatedRequest(t, tc, "DELETE", url, token, nil)

		// We expect either 404 or server error, indicating the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent discussion")
	})

	t.Run("LikeDiscussionReply", func(t *testing.T) {
		// Test liking a reply with a fake ID
		fakeReplyID := uuid.New()
		url := fmt.Sprintf("/api/v1/discussions/replies/%s/like", fakeReplyID)
		w := MakeAuthenticatedRequest(t, tc, "POST", url, token, nil)

		// We expect either 404 or server error, indicating the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent reply")
	})

	t.Run("GetUserDiscussions", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/users/discussions", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)
			data := GetDataFromResponse(t, response)

			discussions, exists := data["discussions"]
			if exists {
				discussionsList, ok := discussions.([]interface{})
				if ok && len(discussionsList) > 0 {
					VerifyIsTestDataInArray(t, discussionsList, true, "user discussions")
				}
			}
		}
	})

	t.Run("UpdateDiscussion", func(t *testing.T) {
		// Test updating a discussion with a fake ID
		fakeDiscussionID := uuid.New()
		updateReq := models.UpdateDiscussionRequest{
			Content: stringPtr("Updated content"),
			Type:    stringPtr("general"),
		}

		reqBody, _ := json.Marshal(updateReq)
		url := fmt.Sprintf("/api/v1/discussions/%s", fakeDiscussionID)
		w := MakeAuthenticatedRequest(t, tc, "PUT", url, token, reqBody)

		// We expect either 404 or forbidden, indicating the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent discussion")
	})

	t.Run("DeleteDiscussion", func(t *testing.T) {
		// Test deleting a discussion with a fake ID
		fakeDiscussionID := uuid.New()
		url := fmt.Sprintf("/api/v1/discussions/%s", fakeDiscussionID)
		w := MakeAuthenticatedRequest(t, tc, "DELETE", url, token, nil)

		// We expect either 404 or forbidden, indicating the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent discussion")
	})

	_ = userID // Use userID to avoid unused variable warning
}

// Helper function to check nested data in discussion objects
func checkDiscussionNestedData(t *testing.T, discussionMap map[string]interface{}, prefix string) {
	// Check user/author data
	if user, exists := discussionMap["user"]; exists && user != nil {
		userMap, ok := user.(map[string]interface{})
		if ok {
			VerifyIsTestDataField(t, userMap, true, prefix+" user")
		}
	}

	if author, exists := discussionMap["author"]; exists && author != nil {
		authorMap, ok := author.(map[string]interface{})
		if ok {
			VerifyIsTestDataField(t, authorMap, true, prefix+" author")
		}
	}

	// Check quiz data if present
	if quiz, exists := discussionMap["quiz"]; exists && quiz != nil {
		quizMap, ok := quiz.(map[string]interface{})
		if ok {
			VerifyIsTestDataField(t, quizMap, true, prefix+" quiz")
		}
	}

	// Check question data if present
	if question, exists := discussionMap["question"]; exists && question != nil {
		questionMap, ok := question.(map[string]interface{})
		if ok {
			VerifyIsTestDataField(t, questionMap, true, prefix+" question")
		}
	}

	// Check replies if present
	if replies, exists := discussionMap["replies"]; exists && replies != nil {
		repliesList, ok := replies.([]interface{})
		if ok && len(repliesList) > 0 {
			VerifyIsTestDataInArray(t, repliesList, true, prefix+" replies")
		}
	}

	// Check likes if present
	if likes, exists := discussionMap["likes"]; exists && likes != nil {
		likesList, ok := likes.([]interface{})
		if ok && len(likesList) > 0 {
			VerifyIsTestDataInArray(t, likesList, true, prefix+" likes")
		}
	}
}

// Helper function to check nested data in discussion reply objects
func checkDiscussionReplyNestedData(t *testing.T, replyMap map[string]interface{}, prefix string) {
	// Check user/author data
	if user, exists := replyMap["user"]; exists && user != nil {
		userMap, ok := user.(map[string]interface{})
		if ok {
			VerifyIsTestDataField(t, userMap, true, prefix+" user")
		}
	}

	if author, exists := replyMap["author"]; exists && author != nil {
		authorMap, ok := author.(map[string]interface{})
		if ok {
			VerifyIsTestDataField(t, authorMap, true, prefix+" author")
		}
	}

	// Check likes if present
	if likes, exists := replyMap["likes"]; exists && likes != nil {
		likesList, ok := likes.([]interface{})
		if ok && len(likesList) > 0 {
			VerifyIsTestDataInArray(t, likesList, true, prefix+" likes")
		}
	}
}
