package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"quizninja-api/config"
	"quizninja-api/models"
	"quizninja-api/repository"
	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

func NewAuthHandler(config *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo: repository.NewUserRepository(),
		config:   config,
	}
}

func (ah *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	existingUser, err := ah.userRepo.GetUserByEmail(req.Email)
	if err != sql.ErrNoRows {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check existing user",
			})
			return
		}
		if existingUser != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
			return
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		Age:          req.Age,
		IsTestData:   true,
	}

	if err := ah.userRepo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Create user preferences if provided
	if req.Preferences != nil {
		now := time.Now()
		preferences := &models.UserPreferences{
			UserID:                  user.ID,
			SelectedInterests:       models.StringArray(req.Preferences.SelectedInterests),
			DifficultyPreference:    req.Preferences.DifficultyPreference,
			NotificationsEnabled:    req.Preferences.NotificationsEnabled,
			NotificationFrequency:   req.Preferences.NotificationFrequency,
			OnboardingCompletedAt:   &now, // Mark onboarding as completed
			IsTestData:              true,
		}

		if err := ah.userRepo.CreateUserPreferences(preferences); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user preferences",
			})
			return
		}

		user.Preferences = preferences
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, ah.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, ah.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	if err := ah.userRepo.SaveRefreshToken(refreshTokenModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save refresh token",
		})
		return
	}

	response := models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusCreated, response)
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := ah.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Update user's online status and last active timestamp
	if err := ah.userRepo.UpdateUserOnlineStatus(user.ID, true); err != nil {
		// Log but don't fail the login
	}
	if err := ah.userRepo.UpdateUserLastActive(user.ID); err != nil {
		// Log but don't fail the login
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, ah.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, ah.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	if err := ah.userRepo.SaveRefreshToken(refreshTokenModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save refresh token",
		})
		return
	}

	response := models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, response)
}

func (ah *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	claims, err := utils.ValidateRefreshToken(req.RefreshToken, ah.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	storedToken, err := ah.userRepo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	if storedToken.UserID != claims.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	user, err := ah.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(user.ID, ah.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, ah.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	if err := ah.userRepo.DeleteRefreshToken(req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete old refresh token",
		})
		return
	}

	newRefreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}

	if err := ah.userRepo.SaveRefreshToken(newRefreshTokenModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (ah *AuthHandler) Logout(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get refresh token to find user ID for status update
	storedToken, err := ah.userRepo.GetRefreshToken(req.RefreshToken)
	if err == nil {
		// Update user's online status to false
		ah.userRepo.UpdateUserOnlineStatus(storedToken.UserID, false)
	}

	if err := ah.userRepo.DeleteRefreshToken(req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func (ah *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	user, err := ah.userRepo.GetUserWithPreferences(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user profile",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ah *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get current user
	currentUser, err := ah.userRepo.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get current user",
		})
		return
	}

	// Check if email is being updated and if it's already taken
	if req.Email != nil && *req.Email != currentUser.Email {
		existingUser, err := ah.userRepo.GetUserByEmail(*req.Email)
		if err != sql.ErrNoRows {
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to check email availability",
				})
				return
			}
			if existingUser != nil {
				c.JSON(http.StatusConflict, gin.H{
					"error": "Email is already in use",
				})
				return
			}
		}
	}

	// Update only provided fields
	if req.Name != nil {
		currentUser.Name = *req.Name
	}
	if req.Email != nil {
		currentUser.Email = *req.Email
	}
	if req.Age != nil {
		currentUser.Age = req.Age
	}
	if req.AvatarURL != nil {
		currentUser.AvatarURL = req.AvatarURL
	}

	// Update the user in database
	if err := ah.userRepo.UpdateUser(currentUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update profile",
		})
		return
	}

	// Get updated user with preferences
	updatedUser, err := ah.userRepo.GetUserWithPreferences(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve updated profile",
		})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (ah *AuthHandler) GetUserStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	stats, err := ah.userRepo.GetUserStatistics(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user statistics",
		})
		return
	}

	response := models.UserStatisticsResponse{
		Statistics: *stats,
		Message:    "User statistics retrieved successfully",
	}

	c.JSON(http.StatusOK, response)
}