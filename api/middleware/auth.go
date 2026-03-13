package middleware

import (
	"net/http"
	"strings"

	"quizninja-api/config"
	"quizninja-api/repository"
	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		var user *utils.SupabaseUser
		var err error

		// Check if we should use mock auth for testing
		if cfg.IsMockAuthEnabled() {
			// Use mock token validation
			user, err = utils.ValidateMockJWT(token, utils.DefaultMockJWTConfig)
			if err != nil {
				utils.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Warn("Mock authentication failed")
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Authentication failed - invalid or expired mock token",
				})
				c.Abort()
				return
			}
		} else {
			// Use real Supabase validation
			var supabaseErr *utils.SupabaseAuthError
			user, supabaseErr = utils.ValidateSupabaseTokenHTTP(token, cfg.SupabaseURL, cfg.SupabaseAnonKey)
			if supabaseErr != nil {
				if !utils.IsSupabaseErrorRetryable(supabaseErr) {
					utils.WithFields(logrus.Fields{
						"error":     supabaseErr.Message,
						"retryable": false,
					}).Warn("Supabase authentication failed")
				} else {
					utils.WithFields(logrus.Fields{
						"error":     supabaseErr.Message,
						"retryable": true,
					}).Warn("Supabase authentication failed")
				}

				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Authentication failed - invalid or expired Supabase token",
				})
				c.Abort()
				return
			}
		}

		// SUCCESS: Auth worked - now lookup database user ID
		supabaseUserID, parseErr := utils.ConvertSupabaseIDToUUID(user.ID)
		if parseErr != nil {
			utils.WithFields(logrus.Fields{
				"supabase_user_id": user.ID,
				"error":            parseErr.Error(),
			}).Error("Failed to parse Supabase user ID as UUID")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Supabase user ID format",
			})
			c.Abort()
			return
		}

		// Look up the database user by Supabase ID
		userRepo := repository.NewUserRepository()
		dbUser, err := userRepo.GetUserBySupabaseID(supabaseUserID.String())
		if err != nil {
			utils.WithFields(logrus.Fields{
				"supabase_user_id": supabaseUserID,
				"error":            err.Error(),
			}).Error("Failed to find database user by Supabase ID")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found in database",
			})
			c.Abort()
			return
		}

		// Set the database user ID (not Supabase ID) in context
		c.Set("user_id", dbUser.ID)
		if cfg.IsMockAuthEnabled() {
			c.Set("auth_method", "mock")
		} else {
			c.Set("auth_method", "supabase")
		}
		c.Set("supabase_user", user)
		c.Next()
	}
}
