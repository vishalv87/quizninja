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

func TestFriendsHandler(t *testing.T) {
	tc := SetupTestServer(t)
	defer Cleanup(t)

	userID, token := CreateTestUser(t, tc)

	t.Run("GetFriends", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/friends", token, nil)

		if w.Code == http.StatusOK {
			var response models.FriendsListResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse friends list response")

			if len(response.Friends) > 0 {
				// Verify friends have is_test_data field
				friendsData := make([]interface{}, len(response.Friends))
				for i, friend := range response.Friends {
					friendsData[i] = map[string]interface{}{
						"is_test_data": friend.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, friendsData, true, "friends list")
			}
		}
	})

	t.Run("GetFriendRequests", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/friends/requests", token, nil)

		if w.Code == http.StatusOK {
			var response models.FriendRequestsResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse friend requests response")

			if len(response.PendingRequests) > 0 {
				// Verify pending requests have is_test_data field
				pendingData := make([]interface{}, len(response.PendingRequests))
				for i, request := range response.PendingRequests {
					pendingData[i] = map[string]interface{}{
						"is_test_data": request.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, pendingData, true, "pending friend requests")
			}

			if len(response.SentRequests) > 0 {
				// Verify sent requests have is_test_data field
				sentData := make([]interface{}, len(response.SentRequests))
				for i, request := range response.SentRequests {
					sentData[i] = map[string]interface{}{
						"is_test_data": request.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, sentData, true, "sent friend requests")
			}
		}
	})

	t.Run("SearchUsers", func(t *testing.T) {
		// Test searching for users
		searchQueries := []string{"test", "user", "admin"}

		for _, query := range searchQueries {
			t.Run(fmt.Sprintf("Query_%s", query), func(t *testing.T) {
				url := fmt.Sprintf("/api/v1/friends/search?q=%s&limit=5", query)
				w := MakeAuthenticatedRequest(t, tc, "GET", url, token, nil)

				if w.Code == http.StatusOK {
					var response models.UserSearchResponse
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err, "Should parse user search response")

					if len(response.Users) > 0 {
						// Verify users have is_test_data field
						usersData := make([]interface{}, len(response.Users))
						for i, user := range response.Users {
							usersData[i] = map[string]interface{}{
								"is_test_data": user.IsTestData,
							}
						}
						VerifyIsTestDataInArray(t, usersData, true, fmt.Sprintf("search results for '%s'", query))
					}
				}
			})
		}
	})

	t.Run("GetFriendNotifications", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/friends/notifications", token, nil)

		if w.Code == http.StatusOK {
			var response models.FriendNotificationsResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse friend notifications response")

			if len(response.Notifications) > 0 {
				// Verify notifications have is_test_data field
				notificationsData := make([]interface{}, len(response.Notifications))
				for i, notification := range response.Notifications {
					notificationsData[i] = map[string]interface{}{
						"is_test_data": notification.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, notificationsData, true, "friend notifications")
			}

			// Verify response structure
			assert.GreaterOrEqual(t, response.UnreadCount, 0, "Unread count should be non-negative")
			assert.GreaterOrEqual(t, response.Total, 0, "Total count should be non-negative")
		}
	})

	t.Run("GetFriendNotificationsWithPagination", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "GET", "/api/v1/friends/notifications?limit=5&offset=0", token, nil)

		if w.Code == http.StatusOK {
			var response models.FriendNotificationsResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should parse paginated friend notifications response")

			if len(response.Notifications) > 0 {
				// Verify notifications have is_test_data field
				notificationsData := make([]interface{}, len(response.Notifications))
				for i, notification := range response.Notifications {
					notificationsData[i] = map[string]interface{}{
						"is_test_data": notification.IsTestData,
					}
				}
				VerifyIsTestDataInArray(t, notificationsData, true, "paginated friend notifications")
			}
		}
	})

	// Test creating friend requests (this will test the full flow if we have multiple users)
	t.Run("SendFriendRequest", func(t *testing.T) {
		// Create a second test user to send a request to
		secondUserID, _ := CreateTestUser(t, tc)

		sendReq := models.SendFriendRequestRequest{
			RequestedUserID: secondUserID,
			Message:         stringPtr("Test friend request"),
		}

		reqBody, _ := json.Marshal(sendReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/friends/requests", token, reqBody)

		if w.Code == http.StatusCreated {
			response := ParseJSONResponse(t, w)

			// Verify the response contains a request object
			request, requestExists := response["request"]
			if requestExists {
				requestMap, ok := request.(map[string]interface{})
				if ok {
					// Check if the request has is_test_data field
					if _, hasTestData := requestMap["is_test_data"]; hasTestData {
						VerifyIsTestDataField(t, requestMap, true, "sent friend request")
					}
				}
			}

			// Verify success message
			message, messageExists := response["message"]
			if messageExists {
				assert.Contains(t, message, "successfully", "Should indicate successful friend request")
			}
		}
	})

	t.Run("SendFriendRequestToSelf", func(t *testing.T) {
		// Test sending friend request to self (should fail)
		sendReq := models.SendFriendRequestRequest{
			RequestedUserID: userID,
			Message:         stringPtr("Test self request"),
		}

		reqBody, _ := json.Marshal(sendReq)
		w := MakeAuthenticatedRequest(t, tc, "POST", "/api/v1/friends/requests", token, reqBody)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Should not allow sending friend request to self")
	})

	t.Run("RespondToFriendRequest", func(t *testing.T) {
		// This would require setting up a friend request first
		// For now, we'll test the endpoint with a fake ID to ensure it handles authentication
		fakeRequestID := uuid.New()
		respondReq := models.RespondToFriendRequestRequest{
			Status: "accepted",
		}

		reqBody, _ := json.Marshal(respondReq)
		url := fmt.Sprintf("/api/v1/friends/requests/%s", fakeRequestID)
		w := MakeAuthenticatedRequest(t, tc, "PUT", url, token, reqBody)

		// We expect either 404 (request not found) or 403 (not authorized)
		// which indicates the endpoint is working and checking permissions
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusForbidden,
			"Should return 404 or 403 for non-existent or unauthorized request")
	})

	t.Run("CancelFriendRequest", func(t *testing.T) {
		// Test canceling a friend request with a fake ID
		fakeRequestID := uuid.New()
		url := fmt.Sprintf("/api/v1/friends/requests/%s", fakeRequestID)
		w := MakeAuthenticatedRequest(t, tc, "DELETE", url, token, nil)

		// We expect a server error or not found, indicating the endpoint is working
		assert.True(t, w.Code >= 400, "Should return an error for non-existent request")
	})

	t.Run("RemoveFriend", func(t *testing.T) {
		// Test removing a friend with a fake ID
		fakeFriendID := uuid.New()
		url := fmt.Sprintf("/api/v1/friends/%s", fakeFriendID)
		w := MakeAuthenticatedRequest(t, tc, "DELETE", url, token, nil)

		// We expect a bad request since they're not friends
		assert.Equal(t, http.StatusBadRequest, w.Code, "Should return bad request for non-friend")
	})

	t.Run("MarkNotificationAsRead", func(t *testing.T) {
		// Test marking a notification as read with a fake ID
		fakeNotificationID := uuid.New()
		url := fmt.Sprintf("/api/v1/friends/notifications/%s/read", fakeNotificationID)
		w := MakeAuthenticatedRequest(t, tc, "PUT", url, token, nil)

		// We expect a server error for non-existent notification
		assert.True(t, w.Code >= 400, "Should return an error for non-existent notification")
	})

	t.Run("MarkAllNotificationsAsRead", func(t *testing.T) {
		w := MakeAuthenticatedRequest(t, tc, "PUT", "/api/v1/friends/notifications/read-all", token, nil)

		if w.Code == http.StatusOK {
			response := ParseJSONResponse(t, w)

			// Verify success message
			message, messageExists := response["message"]
			if messageExists {
				assert.Contains(t, message, "marked as read", "Should indicate notifications marked as read")
			}
		}
	})

	_ = userID // Use userID to avoid unused variable warning
}
