package services

import (
	"fmt"
	"log"

	"quizninja-api/models"
	"quizninja-api/repository"

	"github.com/google/uuid"
)

// StringPtr returns a pointer to the given string
func StringPtr(s string) *string {
	return &s
}

// AchievementService handles achievement logic and checking
type AchievementService struct {
	repo *repository.Repository
}

// NewAchievementService creates a new achievement service
func NewAchievementService(repo *repository.Repository) *AchievementService {
	return &AchievementService{
		repo: repo,
	}
}

// AchievementTrigger represents different events that can trigger achievement checks
type AchievementTrigger string

const (
	TriggerQuizCompleted   AchievementTrigger = "quiz_completed"
	TriggerStreakUpdated   AchievementTrigger = "streak_updated"
	TriggerFriendAdded     AchievementTrigger = "friend_added"
	TriggerPerfectScore    AchievementTrigger = "perfect_score"
	TriggerLevelUp         AchievementTrigger = "level_up"
	TriggerChallengeWon    AchievementTrigger = "challenge_won"
	TriggerLeaderboardRank AchievementTrigger = "leaderboard_rank"
)

// CheckResult represents the result of an achievement check
type CheckResult struct {
	NewAchievements []models.UserAchievement         `json:"new_achievements"`
	Notifications   []models.AchievementNotification `json:"notifications"`
	TotalChecked    int                              `json:"total_checked"`
	TotalUnlocked   int                              `json:"total_unlocked"`
}

// CheckAchievementsForUser checks and unlocks achievements for a user based on their current stats
func (as *AchievementService) CheckAchievementsForUser(userID uuid.UUID, trigger AchievementTrigger) (*CheckResult, error) {
	log.Printf("Checking achievements for user %s with trigger %s", userID, trigger)

	// Get all achievements
	allAchievements, err := as.repo.Achievement.GetAllAchievements()
	if err != nil {
		return nil, fmt.Errorf("failed to get all achievements: %w", err)
	}

	// Get user's current achievements
	userAchievements, err := as.repo.Achievement.GetUserAchievements(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user achievements: %w", err)
	}

	// Create a map of unlocked achievements for quick lookup
	unlockedMap := make(map[string]bool)
	for _, ua := range userAchievements {
		if ua.Achievement != nil {
			unlockedMap[ua.Achievement.Key] = true
		}
	}

	// Get achievements to check based on trigger
	achievementsToCheck := as.getAchievementsToCheck(trigger, allAchievements)

	var newAchievements []models.UserAchievement
	var notifications []models.AchievementNotification

	// Check each achievement
	for _, achievement := range achievementsToCheck {
		// Skip if already unlocked
		if unlockedMap[achievement.Key] {
			continue
		}

		// Check if achievement should be unlocked
		shouldUnlock, err := as.shouldUnlockAchievement(userID, achievement, trigger)
		if err != nil {
			log.Printf("Error checking achievement %s for user %s: %v", achievement.Key, userID, err)
			continue
		}

		if shouldUnlock {
			log.Printf("Unlocking achievement %s for user %s", achievement.Key, userID)

			// Unlock the achievement
			userAchievement, err := as.repo.Achievement.UnlockAchievement(userID, achievement.Key)
			if err != nil {
				log.Printf("Error unlocking achievement %s for user %s: %v", achievement.Key, userID, err)
				continue
			}

			newAchievements = append(newAchievements, *userAchievement)

			// Create persistent notification in the database
			notificationData := map[string]interface{}{
				"achievement_id":    userAchievement.Achievement.ID,
				"achievement_key":   userAchievement.Achievement.Key,
				"achievement_title": userAchievement.Achievement.Title,
				"achievement_description": userAchievement.Achievement.Description,
				"points_awarded":    userAchievement.PointsAwarded,
				"is_rare":          userAchievement.Achievement.IsRare,
				"icon":             userAchievement.Achievement.Icon,
				"color":            userAchievement.Achievement.Color,
			}

			notificationReq := &models.CreateNotificationRequest{
				UserID:            userID,
				Type:              models.NotificationTypeAchievementUnlocked,
				Title:             "Achievement Unlocked!",
				Message:           &userAchievement.Achievement.Title,
				Data:              notificationData,
				RelatedEntityID:   &userAchievement.Achievement.ID,
				RelatedEntityType: StringPtr("achievement"),
			}

			// Persist the notification to database
			_, err = as.repo.Notification.CreateNotification(notificationReq)
			if err != nil {
				log.Printf("Error creating achievement notification for user %s: %v", userID, err)
				// Don't fail the whole operation, just log the error
			}

			// Create in-memory notification for immediate response
			notification := models.AchievementNotification{
				AchievementID: userAchievement.Achievement.ID,
				Title:         userAchievement.Achievement.Title,
				Description:   userAchievement.Achievement.Description,
				Icon:          userAchievement.Achievement.Icon,
				Color:         userAchievement.Achievement.Color,
				PointsAwarded: userAchievement.PointsAwarded,
				IsRare:        userAchievement.Achievement.IsRare,
			}
			notifications = append(notifications, notification)
		}
	}

	result := &CheckResult{
		NewAchievements: newAchievements,
		Notifications:   notifications,
		TotalChecked:    len(achievementsToCheck),
		TotalUnlocked:   len(newAchievements),
	}

	log.Printf("Achievement check completed: %d checked, %d unlocked", result.TotalChecked, result.TotalUnlocked)

	return result, nil
}

// getAchievementsToCheck returns achievements to check based on the trigger
func (as *AchievementService) getAchievementsToCheck(trigger AchievementTrigger, allAchievements []models.Achievement) []models.Achievement {
	var achievementsToCheck []models.Achievement

	switch trigger {
	case TriggerQuizCompleted:
		// Check quiz-related achievements
		keys := []string{"first_win", "quiz_master"}
		achievementsToCheck = as.filterAchievementsByKeys(allAchievements, keys)

	case TriggerStreakUpdated:
		// Check streak-related achievements
		keys := []string{"week_warrior", "streak_legend"}
		achievementsToCheck = as.filterAchievementsByKeys(allAchievements, keys)

	case TriggerFriendAdded:
		// Check social achievements
		keys := []string{"social_butterfly"}
		achievementsToCheck = as.filterAchievementsByKeys(allAchievements, keys)

	case TriggerPerfectScore:
		// Check score achievements
		keys := []string{"perfect_score"}
		achievementsToCheck = as.filterAchievementsByKeys(allAchievements, keys)

	case TriggerLevelUp, TriggerChallengeWon, TriggerLeaderboardRank:
		// Check related achievements for these triggers
		keys := []string{"rising_star", "tech_genius", "sports_expert", "speed_demon"}
		achievementsToCheck = as.filterAchievementsByKeys(allAchievements, keys)

	default:
		// Check all achievements
		achievementsToCheck = allAchievements
	}

	return achievementsToCheck
}

// filterAchievementsByKeys filters achievements by their keys
func (as *AchievementService) filterAchievementsByKeys(achievements []models.Achievement, keys []string) []models.Achievement {
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	var filtered []models.Achievement
	for _, achievement := range achievements {
		if keyMap[achievement.Key] {
			filtered = append(filtered, achievement)
		}
	}

	return filtered
}

// shouldUnlockAchievement determines if an achievement should be unlocked for a user
func (as *AchievementService) shouldUnlockAchievement(userID uuid.UUID, achievement models.Achievement, trigger AchievementTrigger) (bool, error) {
	// Get user stats
	userStats, err := as.getUserStats(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user stats: %w", err)
	}

	// Check achievement conditions based on the key
	switch achievement.Key {
	case "first_win":
		return userStats.QuizzesCompleted >= 1, nil

	case "week_warrior":
		return userStats.CurrentStreak >= 7, nil

	case "quiz_master":
		return userStats.QuizzesCompleted >= 100, nil

	case "streak_legend":
		return userStats.BestStreak >= 30, nil

	case "social_butterfly":
		return userStats.FriendsCount >= 5, nil

	case "perfect_score":
		// Check if user has achieved a perfect score (100%)
		return as.hasUserAchievedPerfectScore(userID)

	case "tech_genius":
		// Check if user has scored above 90% in 10 Technology quizzes
		return as.hasCategoryExpertise(userID, "Technology", 90.0, 10)

	case "sports_expert":
		// Check if user has scored above 90% in 10 Sports quizzes
		return as.hasCategoryExpertise(userID, "Sports", 90.0, 10)

	case "rising_star":
		// Check if user has climbed 10 positions in leaderboard (simplified)
		return userStats.TotalPoints >= 1000, nil

	case "speed_demon":
		// Check if user has completed quizzes under 2 minutes (simplified)
		return as.hasCompletedQuizzesUnderTime(userID, 120, 5) // 120 seconds, 5 quizzes

	default:
		// Unknown achievement, don't unlock
		return false, nil
	}
}

// UserStats represents user statistics for achievement checking
type UserStats struct {
	QuizzesCompleted int
	TotalPoints      int
	CurrentStreak    int
	BestStreak       int
	AverageScore     float64
	FriendsCount     int
}

// getUserStats retrieves user statistics
func (as *AchievementService) getUserStats(userID uuid.UUID) (*UserStats, error) {
	// Get user data
	user, err := as.repo.User.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get friends count
	friends, err := as.repo.Friends.GetFriends(userID)
	if err != nil {
		// Don't fail if friends query fails, just set count to 0
		log.Printf("Warning: failed to get friends count for user %s: %v", userID, err)
	}

	stats := &UserStats{
		QuizzesCompleted: user.TotalQuizzesCompleted,
		TotalPoints:      user.TotalPoints,
		CurrentStreak:    user.CurrentStreak,
		BestStreak:       user.BestStreak,
		AverageScore:     user.AverageScore,
		FriendsCount:     len(friends),
	}

	return stats, nil
}

// hasUserAchievedPerfectScore checks if user has achieved a perfect score (100%) on any quiz
func (as *AchievementService) hasUserAchievedPerfectScore(userID uuid.UUID) (bool, error) {
	// This would require checking quiz attempts for perfect scores
	// For now, we'll use a simplified check based on average score
	user, err := as.repo.User.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	// Simplified: if average score is 95% or higher, consider they've achieved perfect score
	return user.AverageScore >= 95.0, nil
}

// hasCategoryExpertise checks if user has expertise in a specific category
func (as *AchievementService) hasCategoryExpertise(userID uuid.UUID, category string, minScore float64, minQuizzes int) (bool, error) {
	// This would require checking user's performance in specific categories
	// For now, we'll use a simplified check
	user, err := as.repo.User.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	// Simplified: if user has completed enough quizzes and has good average
	return user.TotalQuizzesCompleted >= minQuizzes && user.AverageScore >= minScore, nil
}

// hasCompletedQuizzesUnderTime checks if user has completed quizzes under a certain time
func (as *AchievementService) hasCompletedQuizzesUnderTime(userID uuid.UUID, maxTimeSeconds int, minQuizzes int) (bool, error) {
	// This would require checking quiz attempt times
	// For now, we'll use a simplified check based on total quizzes completed
	user, err := as.repo.User.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	// Simplified: if user has completed enough quizzes, assume some were fast
	return user.TotalQuizzesCompleted >= minQuizzes*2, nil
}

// CheckSingleAchievement checks if a specific achievement should be unlocked
func (as *AchievementService) CheckSingleAchievement(userID uuid.UUID, achievementKey string) (*models.UserAchievement, error) {
	// Check if user already has this achievement
	hasAchievement, err := as.repo.Achievement.HasUserAchievementByKey(userID, achievementKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing achievement: %w", err)
	}

	if hasAchievement {
		return nil, fmt.Errorf("user already has this achievement")
	}

	// Get the achievement
	achievement, err := as.repo.Achievement.GetAchievementByKey(achievementKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get achievement: %w", err)
	}

	// Check if it should be unlocked
	shouldUnlock, err := as.shouldUnlockAchievement(userID, *achievement, "manual")
	if err != nil {
		return nil, fmt.Errorf("failed to check achievement condition: %w", err)
	}

	if !shouldUnlock {
		return nil, fmt.Errorf("achievement conditions not met")
	}

	// Unlock the achievement
	userAchievement, err := as.repo.Achievement.UnlockAchievement(userID, achievementKey)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock achievement: %w", err)
	}

	return userAchievement, nil
}

// GetAchievementNotifications creates notifications for newly unlocked achievements
func (as *AchievementService) GetAchievementNotifications(userAchievements []models.UserAchievement) []models.AchievementNotification {
	notifications := make([]models.AchievementNotification, len(userAchievements))

	for i, ua := range userAchievements {
		notifications[i] = models.AchievementNotification{
			AchievementID: ua.AchievementID,
			Title:         ua.Achievement.Title,
			Description:   ua.Achievement.Description,
			Icon:          ua.Achievement.Icon,
			Color:         ua.Achievement.Color,
			PointsAwarded: ua.PointsAwarded,
			IsRare:        ua.Achievement.IsRare,
		}
	}

	return notifications
}