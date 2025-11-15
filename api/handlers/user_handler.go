package handlers

import (
	"database/sql"
	"net/http"
	"quizninja-api/models"
	"quizninja-api/repository"
	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userRepo    *repository.UserRepository
	friendsRepo repository.FriendsRepositoryInterface
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userRepo:    repository.NewUserRepository(),
		friendsRepo: repository.NewFriendsRepository(),
	}
}

// GetUserProfile handles GET /users/:userId
// Returns another user's profile with privacy-aware data
func (uh *UserHandler) GetUserProfile(c *gin.Context) {
	// Get the requesting user's ID from context (set by auth middleware)
	currentUserID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	currentUserUUID, ok := currentUserID.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	// Get the target user ID from URL parameter
	userIDStr := c.Param("userId")
	targetUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	utils.WithFields(logrus.Fields{
		"current_user_id": currentUserUUID,
		"target_user_id":  targetUserID,
	}).Info("GetUserProfile: Fetching user profile")

	// Fetch the target user
	targetUser, err := uh.userRepo.GetUserByID(targetUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(c, http.StatusNotFound, "User not found")
			return
		}
		utils.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("GetUserProfile: Failed to fetch user")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch user profile")
		return
	}

	// Fetch target user's preferences for privacy settings
	preferences, err := uh.userRepo.GetUserPreferences(targetUserID)
	if err != nil && err != sql.ErrNoRows {
		utils.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("GetUserProfile: Failed to fetch user preferences, using defaults")
	}

	// Default privacy settings if preferences not found
	if preferences == nil {
		preferences = &models.UserPreferences{
			ProfileVisibility:   true, // Default to public
			AllowFriendRequests: true,
		}
	}

	// Check if users are friends
	areFriends, err := uh.friendsRepo.AreFriends(currentUserUUID, targetUserID)
	if err != nil {
		utils.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("GetUserProfile: Failed to check friendship status")
		areFriends = false
	}

	// Check friendship request status
	friendRequestStatus := "none"
	if areFriends {
		friendRequestStatus = "friends"
	} else {
		// Check for pending requests
		pendingRequests, err := uh.friendsRepo.GetPendingFriendRequests(currentUserUUID)
		if err == nil {
			for _, req := range pendingRequests {
				if req.RequestedID == targetUserID && req.Status == "pending" {
					friendRequestStatus = "pending_sent"
					break
				}
				if req.RequesterID == targetUserID && req.Status == "pending" {
					friendRequestStatus = "pending_received"
					break
				}
			}
		}
	}

	// Check privacy permissions
	// If profile is private (profile_visibility = false) and not friends, deny access
	if !preferences.ProfileVisibility && !areFriends && currentUserUUID != targetUserID {
		utils.ErrorResponse(c, http.StatusForbidden, "This profile is private")
		return
	}

	// Build response
	response := models.UserProfileResponse{
		ID:                  targetUser.ID,
		UserID:              targetUser.ID,
		Name:                targetUser.Name,
		FullName:            targetUser.Name,
		Email:               targetUser.Email,
		AvatarURL:           targetUser.AvatarURL,
		Bio:                 nil, // TODO: Add bio field to users table if needed
		CreatedAt:           targetUser.CreatedAt,
		UpdatedAt:           targetUser.UpdatedAt,
		IsFriend:            areFriends,
		FriendRequestStatus: friendRequestStatus,
		Preferences: &models.UserProfilePreferences{
			ProfileVisibility:  preferences.ProfileVisibility,
			ShowAchievements:   true,  // Default for now, can be added to user_preferences table
			ShowStats:          true,  // Default for now, can be added to user_preferences table
			AllowFriendRequest: preferences.AllowFriendRequests,
		},
	}

	// Include stats if privacy allows or if viewing own profile
	// Stats are visible if profile is public OR users are friends OR viewing own profile
	if preferences.ProfileVisibility || areFriends || currentUserUUID == targetUserID {
		stats, err := uh.userRepo.GetUserStatistics(targetUserID)
		if err == nil && stats != nil {
			response.Stats = &models.UserStats{
				UserID:                targetUserID,
				TotalQuizzesTaken:     stats.TotalAttempts,
				TotalQuizzesCompleted: stats.CompletedQuizzes,
				TotalPoints:           stats.TotalPoints,
				AverageScore:          stats.AverageScore,
				TotalTimeSpentMinutes: stats.AverageCompletionTime / 60, // Convert seconds to minutes
				CurrentStreak:         stats.CurrentStreak,
				LongestStreak:         stats.BestStreak,
				AchievementsUnlocked:  0, // TODO: Add achievement count query
				ChallengesWon:         0, // TODO: Add challenges won count query
				ChallengesLost:        0, // TODO: Add challenges lost count query
				Rank:                  0, // TODO: Add rank calculation
			}
		} else if err != nil {
			utils.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("GetUserProfile: Failed to fetch user statistics")
		}
	}

	utils.WithFields(logrus.Fields{
		"target_user_id": targetUserID,
		"is_friend":      areFriends,
		"includes_stats": response.Stats != nil,
		"request_status": friendRequestStatus,
	}).Info("GetUserProfile: Successfully fetched user profile")

	c.JSON(http.StatusOK, gin.H{
		"message": "User profile fetched successfully",
		"data":    response,
	})
}
