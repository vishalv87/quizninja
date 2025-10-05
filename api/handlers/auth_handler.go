package handlers

import (
	"database/sql"
	"log"
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
	// Check for idempotency key
	idempotencyKey := c.GetHeader("X-Idempotency-Key")
	if idempotencyKey != "" {
		store := utils.GetIdempotencyStore()
		if cached, exists := store.Get(idempotencyKey); exists {
			log.Printf("Returning cached response for idempotency key: %s", idempotencyKey)
			c.JSON(cached.StatusCode, cached.ResponseBody)
			return
		}
	}

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if user already exists by Supabase ID
	_, err := ah.userRepo.GetUserBySupabaseID(req.SupabaseUserID)
	if err == nil {
		// User already exists with this Supabase ID
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this Supabase ID already exists",
		})
		return
	}

	// Check by email as secondary check
	existingUserByEmail, err := ah.userRepo.GetUserByEmail(req.Email)
	if err != sql.ErrNoRows {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check existing user",
			})
			return
		}
		if existingUserByEmail != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User with this email already exists",
			})
			return
		}
	}

	// Convert Supabase ID to UUID
	supabaseUserID, parseErr := utils.ConvertSupabaseIDToUUID(req.SupabaseUserID)
	if parseErr != nil {
		log.Printf("Failed to parse Supabase user ID as UUID: %v", parseErr)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Supabase user ID format",
		})
		return
	}

	// Create local user profile with Supabase data
	user := &models.User{
		ID:             supabaseUserID,
		Email:          req.Email,
		Name:           req.Name,
		Age:            req.Age,
		AvatarURL:      req.AvatarURL,
		IsTestData:     true,
		AuthMethod:     "supabase",
		SupabaseID:     &req.SupabaseUserID,
		LastAuthMethod: "supabase",
	}

	if err := ah.userRepo.CreateUser(user); err != nil {
		log.Printf("Failed to create local user profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user profile",
		})
		return
	}

	// Create user preferences if provided
	if req.Preferences != nil {
		ah.createUserPreferences(user, req.Preferences)
	}

	// Update auth method tracking
	ah.userRepo.UpdateUserAuthMethod(user.ID, "supabase")

	response := models.AuthResponse{
		User:    *user,
		Message: "User profile created successfully",
	}

	// Cache the response for idempotency
	if idempotencyKey != "" {
		store := utils.GetIdempotencyStore()
		store.Set(idempotencyKey, http.StatusCreated, response)
		log.Printf("Cached registration response for idempotency key: %s", idempotencyKey)
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

	// Try to find or create user based on Supabase ID
	user, err := ah.userRepo.GetUserBySupabaseID(req.SupabaseUserID)
	if err == sql.ErrNoRows {
		// User doesn't exist locally - create from Supabase data
		supabaseUserID, parseErr := utils.ConvertSupabaseIDToUUID(req.SupabaseUserID)
		if parseErr != nil {
			log.Printf("Failed to parse Supabase user ID as UUID: %v", parseErr)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Supabase user ID format",
			})
			return
		}

		user = &models.User{
			ID:             supabaseUserID,
			Email:          req.Email,
			Name:           req.Name,
			AvatarURL:      req.AvatarURL,
			IsTestData:     true,
			AuthMethod:     "supabase",
			SupabaseID:     &req.SupabaseUserID,
			LastAuthMethod: "supabase",
		}

		if createErr := ah.userRepo.CreateUser(user); createErr != nil {
			log.Printf("Failed to create local user profile for Supabase user: %v", createErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user profile",
			})
			return
		}
	} else if err != nil {
		log.Printf("Failed to get user by Supabase ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	} else {
		// User exists - update profile with latest Supabase data
		user.Email = req.Email
		user.Name = req.Name
		user.AvatarURL = req.AvatarURL
		user.LastAuthMethod = "supabase"

		if updateErr := ah.userRepo.UpdateUser(user); updateErr != nil {
			log.Printf("Failed to update user profile: %v", updateErr)
			// Don't fail the login for update errors, just log
		}
	}

	// Update user activity and auth method
	ah.userRepo.UpdateUserOnlineStatus(user.ID, true)
	ah.userRepo.UpdateUserLastActive(user.ID)
	ah.userRepo.UpdateUserAuthMethod(user.ID, "supabase")

	response := models.AuthResponse{
		User:    *user,
		Message: "User login successful",
	}

	c.JSON(http.StatusOK, response)
}

func (ah *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from auth middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Update user's online status to false
	ah.userRepo.UpdateUserOnlineStatus(userID.(uuid.UUID), false)

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

	c.JSON(http.StatusOK, gin.H{
		"data":    *stats,
		"message": "User statistics retrieved successfully",
	})
}

// Helper method to create user preferences
func (ah *AuthHandler) createUserPreferences(user *models.User, preferencesReq *models.UserPreferencesRequest) {
	now := time.Now()
	preferences := &models.UserPreferences{
		UserID:                user.ID,
		SelectedInterests:     models.StringArray(preferencesReq.SelectedInterests),
		DifficultyPreference:  preferencesReq.DifficultyPreference,
		NotificationsEnabled:  preferencesReq.NotificationsEnabled,
		NotificationFrequency: preferencesReq.NotificationFrequency,
		OnboardingCompletedAt: &now, // Mark onboarding as completed
		IsTestData:            true,
	}

	if err := ah.userRepo.CreateUserPreferences(preferences); err != nil {
		log.Printf("Failed to create user preferences: %v", err)
		return
	}

	user.Preferences = preferences
}
