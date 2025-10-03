package middleware

import (
	"log"
	"net/http"
	"strings"

	"quizninja-api/config"
	"quizninja-api/repository"
	"quizninja-api/utils"

	"github.com/gin-gonic/gin"
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

		// PRIMARY: Try Supabase Auth first (if enabled)
		if cfg.UseSupabaseAuth {
			user, err := utils.ValidateSupabaseTokenHTTP(token, cfg.SupabaseURL, cfg.SupabaseAnonKey)
			if err == nil {
				// SUCCESS: Supabase auth worked - now lookup database user ID
				supabaseUserID, parseErr := utils.ConvertSupabaseIDToUUID(user.ID)
				if parseErr != nil {
					log.Printf("Failed to parse Supabase user ID as UUID: %v", parseErr)
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
					log.Printf("Failed to find database user by Supabase ID %s: %v", supabaseUserID, err)
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "User not found in database",
					})
					c.Abort()
					return
				}

				// Set the database user ID (not Supabase ID) in context
				c.Set("user_id", dbUser.ID)
				c.Set("auth_method", "supabase")
				c.Set("supabase_user", user)
				c.Next()
				return
			} else {
				// When Supabase auth is enabled, don't fallback to JWT - return error
				if !utils.IsSupabaseErrorRetryable(err) {
					log.Printf("Supabase auth failed (non-retryable): %v", err)
				} else {
					log.Printf("Supabase auth failed (retryable): %v", err)
				}

				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Authentication failed - invalid or expired Supabase token",
				})
				c.Abort()
				return
			}
		}

		// No JWT fallback - Supabase auth is required
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Supabase authentication is required",
		})
		c.Abort()
		return
	}
}
