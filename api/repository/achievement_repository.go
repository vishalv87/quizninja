package repository

import (
	"database/sql"
	"fmt"
	"time"

	"quizninja-api/database"
	"quizninja-api/models"

	"github.com/google/uuid"
)

type AchievementRepository struct {
	db *sql.DB
}

func NewAchievementRepository() *AchievementRepository {
	return &AchievementRepository{
		db: database.DB,
	}
}

// GetAllAchievements retrieves all active achievements
func (ar *AchievementRepository) GetAllAchievements() ([]models.Achievement, error) {
	query := `
		SELECT id, key, title, description, icon, color, points_reward, category, is_rare, is_active, created_at, updated_at, is_test_data
		FROM achievements
		WHERE is_active = true
		ORDER BY category, title
	`

	rows, err := ar.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query achievements: %w", err)
	}
	defer rows.Close()

	var achievements []models.Achievement
	for rows.Next() {
		var achievement models.Achievement
		var icon sql.NullString

		err := rows.Scan(
			&achievement.ID,
			&achievement.Key,
			&achievement.Title,
			&achievement.Description,
			&icon,
			&achievement.Color,
			&achievement.PointsReward,
			&achievement.Category,
			&achievement.IsRare,
			&achievement.IsActive,
			&achievement.CreatedAt,
			&achievement.UpdatedAt,
			&achievement.IsTestData,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan achievement: %w", err)
		}

		if icon.Valid {
			achievement.Icon = &icon.String
		}

		achievements = append(achievements, achievement)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating achievements: %w", err)
	}

	return achievements, nil
}

// GetAchievementByKey retrieves an achievement by its unique key
func (ar *AchievementRepository) GetAchievementByKey(key string) (*models.Achievement, error) {
	query := `
		SELECT id, key, title, description, icon, color, points_reward, category, is_rare, is_active, created_at, updated_at, is_test_data
		FROM achievements
		WHERE key = $1 AND is_active = true
	`

	var achievement models.Achievement
	var icon sql.NullString

	err := ar.db.QueryRow(query, key).Scan(
		&achievement.ID,
		&achievement.Key,
		&achievement.Title,
		&achievement.Description,
		&icon,
		&achievement.Color,
		&achievement.PointsReward,
		&achievement.Category,
		&achievement.IsRare,
		&achievement.IsActive,
		&achievement.CreatedAt,
		&achievement.UpdatedAt,
		&achievement.IsTestData,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("achievement with key %s not found", key)
		}
		return nil, fmt.Errorf("failed to get achievement by key: %w", err)
	}

	if icon.Valid {
		achievement.Icon = &icon.String
	}

	return &achievement, nil
}

// GetUserAchievements retrieves all achievements unlocked by a user
func (ar *AchievementRepository) GetUserAchievements(userID uuid.UUID) ([]models.UserAchievement, error) {
	query := `
		SELECT
			ua.id, ua.user_id, ua.achievement_id, ua.unlocked_at, ua.points_awarded,
			a.id, a.key, a.title, a.description, a.icon, a.color, a.points_reward, a.category, a.is_rare, a.is_active, a.created_at, a.updated_at, a.is_test_data
		FROM user_achievements ua
		JOIN achievements a ON ua.achievement_id = a.id
		WHERE ua.user_id = $1
		ORDER BY ua.unlocked_at DESC
	`

	rows, err := ar.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user achievements: %w", err)
	}
	defer rows.Close()

	var userAchievements []models.UserAchievement
	for rows.Next() {
		var ua models.UserAchievement
		var achievement models.Achievement
		var icon sql.NullString

		err := rows.Scan(
			&ua.ID,
			&ua.UserID,
			&ua.AchievementID,
			&ua.UnlockedAt,
			&ua.PointsAwarded,
			&achievement.ID,
			&achievement.Key,
			&achievement.Title,
			&achievement.Description,
			&icon,
			&achievement.Color,
			&achievement.PointsReward,
			&achievement.Category,
			&achievement.IsRare,
			&achievement.IsActive,
			&achievement.CreatedAt,
			&achievement.UpdatedAt,
			&achievement.IsTestData,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user achievement: %w", err)
		}

		if icon.Valid {
			achievement.Icon = &icon.String
		}

		ua.Achievement = &achievement
		userAchievements = append(userAchievements, ua)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user achievements: %w", err)
	}

	return userAchievements, nil
}

// UnlockAchievement awards an achievement to a user
func (ar *AchievementRepository) UnlockAchievement(userID uuid.UUID, achievementKey string) (*models.UserAchievement, error) {
	// First, get the achievement details
	achievement, err := ar.GetAchievementByKey(achievementKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get achievement: %w", err)
	}

	// Check if user already has this achievement
	if hasAchievement, err := ar.HasUserAchievement(userID, achievement.ID); err != nil {
		return nil, fmt.Errorf("failed to check existing achievement: %w", err)
	} else if hasAchievement {
		return nil, fmt.Errorf("user already has this achievement")
	}

	// Start transaction
	tx, err := ar.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert user achievement
	insertQuery := `
		INSERT INTO user_achievements (user_id, achievement_id, unlocked_at, points_awarded, is_test_data)
		VALUES ($1, $2, $3, $4, true)
		RETURNING id
	`

	var userAchievementID uuid.UUID
	unlockedAt := time.Now()

	err = tx.QueryRow(insertQuery, userID, achievement.ID, unlockedAt, achievement.PointsReward).Scan(&userAchievementID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user achievement: %w", err)
	}

	// Update user's total points
	updatePointsQuery := `
		UPDATE users
		SET total_points = total_points + $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err = tx.Exec(updatePointsQuery, achievement.PointsReward, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user points: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the created user achievement
	userAchievement := &models.UserAchievement{
		ID:             userAchievementID,
		UserID:         userID,
		AchievementID:  achievement.ID,
		UnlockedAt:     unlockedAt,
		PointsAwarded:  achievement.PointsReward,
		Achievement:    achievement,
	}

	return userAchievement, nil
}

// HasUserAchievement checks if a user has already unlocked a specific achievement
func (ar *AchievementRepository) HasUserAchievement(userID, achievementID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_achievements
			WHERE user_id = $1 AND achievement_id = $2
		)
	`

	var exists bool
	err := ar.db.QueryRow(query, userID, achievementID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user achievement: %w", err)
	}

	return exists, nil
}

// HasUserAchievementByKey checks if a user has already unlocked a specific achievement by key
func (ar *AchievementRepository) HasUserAchievementByKey(userID uuid.UUID, achievementKey string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_achievements ua
			JOIN achievements a ON ua.achievement_id = a.id
			WHERE ua.user_id = $1 AND a.key = $2
		)
	`

	var exists bool
	err := ar.db.QueryRow(query, userID, achievementKey).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user achievement by key: %w", err)
	}

	return exists, nil
}

// GetAchievementProgress calculates progress towards achievements for a user
func (ar *AchievementRepository) GetAchievementProgress(userID uuid.UUID) ([]models.AchievementProgress, error) {
	// Get all achievements
	achievements, err := ar.GetAllAchievements()
	if err != nil {
		return nil, fmt.Errorf("failed to get achievements: %w", err)
	}

	// Get user's current stats
	userStats, err := ar.getUserStats(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	// Get user's unlocked achievements
	userAchievements, err := ar.GetUserAchievements(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user achievements: %w", err)
	}

	// Create a map of unlocked achievements for quick lookup
	unlockedMap := make(map[uuid.UUID]bool)
	for _, ua := range userAchievements {
		unlockedMap[ua.AchievementID] = true
	}

	var progress []models.AchievementProgress
	for _, achievement := range achievements {
		currentValue, targetValue := ar.getProgressValues(achievement.Key, userStats)

		progressPercent := float64(0)
		if targetValue > 0 {
			progressPercent = float64(currentValue) / float64(targetValue) * 100
			if progressPercent > 100 {
				progressPercent = 100
			}
		}

		isUnlocked := unlockedMap[achievement.ID]

		achievementProgress := models.AchievementProgress{
			AchievementID: achievement.ID,
			Title:         achievement.Title,
			Description:   achievement.Description,
			Icon:          achievement.Icon,
			Color:         achievement.Color,
			Category:      achievement.Category,
			IsRare:        achievement.IsRare,
			CurrentValue:  currentValue,
			TargetValue:   targetValue,
			Progress:      progressPercent,
			IsUnlocked:    isUnlocked,
		}

		progress = append(progress, achievementProgress)
	}

	return progress, nil
}

// UserStats represents user statistics for achievement calculation
type UserStats struct {
	QuizzesCompleted int
	TotalPoints      int
	CurrentStreak    int
	BestStreak       int
	AverageScore     float64
	FriendsCount     int
}

// getUserStats retrieves user statistics for achievement progress calculation
func (ar *AchievementRepository) getUserStats(userID uuid.UUID) (*UserStats, error) {
	query := `
		SELECT
			total_quizzes_completed,
			total_points,
			current_streak,
			best_streak,
			average_score
		FROM users
		WHERE id = $1
	`

	var stats UserStats
	err := ar.db.QueryRow(query, userID).Scan(
		&stats.QuizzesCompleted,
		&stats.TotalPoints,
		&stats.CurrentStreak,
		&stats.BestStreak,
		&stats.AverageScore,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	// Get friends count
	friendsQuery := `
		SELECT COUNT(*) FROM friendships
		WHERE user1_id = $1 OR user2_id = $1
	`
	err = ar.db.QueryRow(friendsQuery, userID).Scan(&stats.FriendsCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get friends count: %w", err)
	}

	return &stats, nil
}

// getProgressValues returns current and target values for achievement progress
func (ar *AchievementRepository) getProgressValues(achievementKey string, stats *UserStats) (int, int) {
	switch achievementKey {
	case "first_win":
		return stats.QuizzesCompleted, 1
	case "week_warrior":
		return stats.CurrentStreak, 7
	case "quiz_master":
		return stats.QuizzesCompleted, 100
	case "streak_legend":
		return stats.BestStreak, 30
	case "social_butterfly":
		return stats.FriendsCount, 5
	case "perfect_score":
		// For perfect score, we'd need to check if user has achieved 100% on any quiz
		// This is simplified - in reality you'd query quiz attempts
		if stats.AverageScore >= 100 {
			return 1, 1
		}
		return 0, 1
	default:
		return 0, 1
	}
}