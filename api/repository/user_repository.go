package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"quizninja-api/database"
	"quizninja-api/models"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, name, age, is_test_data)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at, level, total_points, current_streak,
		          best_streak, total_quizzes_completed, average_score, is_online, last_active, is_test_data
	`
	err := ur.db.QueryRow(query, user.Email, user.PasswordHash, user.Name, user.Age, user.IsTestData).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Level, &user.TotalPoints,
		&user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
		&user.AverageScore, &user.IsOnline, &user.LastActive, &user.IsTestData,
	)
	return err
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, age, level, total_points, current_streak,
		       best_streak, total_quizzes_completed, average_score, is_online,
		       last_active, avatar_url, created_at, updated_at, is_test_data
		FROM users
		WHERE email = $1
	`
	err := ur.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Age, &user.Level,
		&user.TotalPoints, &user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
		&user.AverageScore, &user.IsOnline, &user.LastActive, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt, &user.IsTestData,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, age, level, total_points, current_streak,
		       best_streak, total_quizzes_completed, average_score, is_online,
		       last_active, avatar_url, created_at, updated_at, is_test_data
		FROM users
		WHERE id = $1
	`
	err := ur.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Age, &user.Level,
		&user.TotalPoints, &user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
		&user.AverageScore, &user.IsOnline, &user.LastActive, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt, &user.IsTestData,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) SaveRefreshToken(refreshToken *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	now := time.Now()
	err := ur.db.QueryRow(query, refreshToken.UserID, refreshToken.Token, refreshToken.ExpiresAt, now).Scan(&refreshToken.ID)
	if err != nil {
		return err
	}
	refreshToken.CreatedAt = now
	return nil
}

func (ur *UserRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	refreshToken := &models.RefreshToken{}
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1 AND expires_at > NOW()
	`
	err := ur.db.QueryRow(query, token).Scan(
		&refreshToken.ID, &refreshToken.UserID, &refreshToken.Token, &refreshToken.ExpiresAt, &refreshToken.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (ur *UserRepository) DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := ur.db.Exec(query, token)
	return err
}

func (ur *UserRepository) DeleteUserRefreshTokens(userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := ur.db.Exec(query, userID)
	return err
}

// UpdateUser updates an existing user
func (ur *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users SET
			email = $2, name = $3, age = $4, level = $5, total_points = $6,
			current_streak = $7, best_streak = $8, total_quizzes_completed = $9,
			average_score = $10, is_online = $11, last_active = $12, avatar_url = $13
		WHERE id = $1
	`
	_, err := ur.db.Exec(query, user.ID, user.Email, user.Name, user.Age, user.Level,
		user.TotalPoints, user.CurrentStreak, user.BestStreak, user.TotalQuizzesCompleted,
		user.AverageScore, user.IsOnline, user.LastActive, user.AvatarURL)
	return err
}

// DeleteUser deletes a user by ID
func (ur *UserRepository) DeleteUser(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := ur.db.Exec(query, id)
	return err
}

// CreateUserPreferences creates user preferences
func (ur *UserRepository) CreateUserPreferences(preferences *models.UserPreferences) error {
	query := `
		INSERT INTO user_preferences (user_id, selected_interests, difficulty_preference,
		                              notifications_enabled, notification_frequency, is_test_data)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err := ur.db.QueryRow(query, preferences.UserID, preferences.SelectedInterests,
		preferences.DifficultyPreference, preferences.NotificationsEnabled,
		preferences.NotificationFrequency, preferences.IsTestData).Scan(&preferences.ID, &preferences.CreatedAt)
	return err
}

// GetUserPreferences retrieves user preferences by user ID
func (ur *UserRepository) GetUserPreferences(userID uuid.UUID) (*models.UserPreferences, error) {
	preferences := &models.UserPreferences{}
	query := `
		SELECT id, user_id, selected_interests, difficulty_preference,
		       notifications_enabled, notification_frequency, created_at, is_test_data
		FROM user_preferences
		WHERE user_id = $1
	`
	err := ur.db.QueryRow(query, userID).Scan(
		&preferences.ID, &preferences.UserID, &preferences.SelectedInterests,
		&preferences.DifficultyPreference, &preferences.NotificationsEnabled,
		&preferences.NotificationFrequency, &preferences.CreatedAt, &preferences.IsTestData,
	)
	if err != nil {
		return nil, err
	}
	return preferences, nil
}

// UpdateUserPreferences updates or inserts user preferences (UPSERT)
func (ur *UserRepository) UpdateUserPreferences(preferences *models.UserPreferences) error {
	query := `
		INSERT INTO user_preferences (
			user_id, selected_interests, difficulty_preference,
			notifications_enabled, notification_frequency,
			profile_visibility, show_online_status,
			allow_friend_requests, share_activity_status,
			notification_types, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
		ON CONFLICT (user_id) DO UPDATE SET
			selected_interests = EXCLUDED.selected_interests,
			difficulty_preference = EXCLUDED.difficulty_preference,
			notifications_enabled = EXCLUDED.notifications_enabled,
			notification_frequency = EXCLUDED.notification_frequency,
			profile_visibility = EXCLUDED.profile_visibility,
			show_online_status = EXCLUDED.show_online_status,
			allow_friend_requests = EXCLUDED.allow_friend_requests,
			share_activity_status = EXCLUDED.share_activity_status,
			notification_types = EXCLUDED.notification_types,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := ur.db.Exec(query,
		preferences.UserID,
		preferences.SelectedInterests,
		preferences.DifficultyPreference,
		preferences.NotificationsEnabled,
		preferences.NotificationFrequency,
		preferences.ProfileVisibility,
		preferences.ShowOnlineStatus,
		preferences.AllowFriendRequests,
		preferences.ShareActivityStatus,
		preferences.NotificationTypes)
	return err
}

// DeleteUserPreferences deletes user preferences by user ID
func (ur *UserRepository) DeleteUserPreferences(userID uuid.UUID) error {
	query := `DELETE FROM user_preferences WHERE user_id = $1`
	_, err := ur.db.Exec(query, userID)
	return err
}

// GetUserWithPreferences retrieves a user with their preferences
func (ur *UserRepository) GetUserWithPreferences(userID uuid.UUID) (*models.User, error) {
	user, err := ur.GetUserByID(userID)
	if err != nil {
		log.Printf("GetUserWithPreferences: Failed to get user by ID %s: %v", userID, err)
		return nil, err
	}

	preferences, err := ur.GetUserPreferences(userID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("GetUserWithPreferences: Failed to get user preferences for user %s: %v", userID, err)
		return nil, err
	}

	// Only set preferences if we successfully retrieved them (no error)
	if err == nil {
		user.Preferences = preferences
		log.Printf("GetUserWithPreferences: Successfully retrieved preferences for user %s", userID)
	} else {
		log.Printf("GetUserWithPreferences: No preferences found for user %s (this is normal)", userID)
	}

	return user, nil
}

// UpdateUserOnlineStatus updates user's online status
func (ur *UserRepository) UpdateUserOnlineStatus(userID uuid.UUID, isOnline bool) error {
	query := `UPDATE users SET is_online = $2 WHERE id = $1`
	_, err := ur.db.Exec(query, userID, isOnline)
	return err
}

// UpdateUserLastActive updates user's last active timestamp
func (ur *UserRepository) UpdateUserLastActive(userID uuid.UUID) error {
	query := `UPDATE users SET last_active = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := ur.db.Exec(query, userID)
	return err
}

// GetUserStatistics retrieves comprehensive user statistics
func (ur *UserRepository) GetUserStatistics(userID uuid.UUID) (*models.UserStatistics, error) {
	stats := &models.UserStatistics{
		UserID:      userID,
		LastUpdated: time.Now(),
	}

	// Get basic user stats from users table
	userQuery := `
		SELECT total_points, current_streak, best_streak, total_quizzes_completed, average_score, is_test_data
		FROM users
		WHERE id = $1
	`
	err := ur.db.QueryRow(userQuery, userID).Scan(
		&stats.TotalPoints, &stats.CurrentStreak, &stats.BestStreak,
		&stats.CompletedQuizzes, &stats.AverageScore, &stats.IsTestData,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get basic user stats: %w", err)
	}

	// Get total attempts and completion rate
	attemptsQuery := `
		SELECT
			COUNT(*) as total_attempts,
			COUNT(CASE WHEN is_completed = true THEN 1 END) as completed_attempts,
			AVG(CASE WHEN is_completed = true THEN time_spent END) as avg_completion_time
		FROM quiz_attempts
		WHERE user_id = $1
	`
	var avgCompletionTime sql.NullFloat64
	var completedFromAttempts int
	err = ur.db.QueryRow(attemptsQuery, userID).Scan(
		&stats.TotalAttempts, &completedFromAttempts, &avgCompletionTime,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get attempt stats: %w", err)
	}

	// Use the completed count from attempts table as it's more accurate
	stats.CompletedQuizzes = completedFromAttempts

	// Calculate completion rate
	if stats.TotalAttempts > 0 {
		stats.CompletionRate = float64(stats.CompletedQuizzes) / float64(stats.TotalAttempts) * 100
	}

	if avgCompletionTime.Valid {
		stats.AverageCompletionTime = int(avgCompletionTime.Float64)
	}

	// Get category performance
	categoryPerformance, err := ur.getCategoryPerformance(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category performance: %w", err)
	}
	stats.CategoryPerformance = categoryPerformance

	// Get recent activity
	recentActivity, err := ur.getRecentActivity(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}
	stats.RecentActivity = recentActivity

	// Get quizzes by difficulty
	quizzesByDifficulty, err := ur.getQuizzesByDifficulty(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quizzes by difficulty: %w", err)
	}
	stats.QuizzesByDifficulty = quizzesByDifficulty

	// Get score distribution
	scoreDistribution, err := ur.getScoreDistribution(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get score distribution: %w", err)
	}
	stats.ScoreDistribution = scoreDistribution

	// Get monthly progress
	monthlyProgress, err := ur.getMonthlyProgress(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly progress: %w", err)
	}
	stats.MonthlyProgress = monthlyProgress

	return stats, nil
}

// UpdateUserStatistics updates user statistics after quiz completion
func (ur *UserRepository) UpdateUserStatistics(userID uuid.UUID, newScore float64) error {
	// Start a transaction to ensure consistency
	tx, err := ur.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// First, get current user statistics
	var currentTotalQuizzes int
	var currentAverageScore sql.NullFloat64

	query := `
		SELECT total_quizzes_completed, average_score
		FROM users
		WHERE id = $1
	`

	err = tx.QueryRow(query, userID).Scan(&currentTotalQuizzes, &currentAverageScore)
	if err != nil {
		return fmt.Errorf("failed to get current user statistics: %w", err)
	}

	// Calculate new statistics
	newTotalQuizzes := currentTotalQuizzes + 1

	// Calculate new average score
	var newAverageScore float64
	if currentAverageScore.Valid && currentTotalQuizzes > 0 {
		// Existing average exists, calculate weighted average
		totalPreviousScore := currentAverageScore.Float64 * float64(currentTotalQuizzes)
		newAverageScore = (totalPreviousScore + newScore) / float64(newTotalQuizzes)
	} else {
		// First quiz, new score becomes the average
		newAverageScore = newScore
	}

	// Update user statistics
	updateQuery := `
		UPDATE users
		SET total_quizzes_completed = $1,
		    average_score = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err = tx.Exec(updateQuery, newTotalQuizzes, newAverageScore, userID)
	if err != nil {
		return fmt.Errorf("failed to update user statistics: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// getCategoryPerformance retrieves performance metrics by category
func (ur *UserRepository) getCategoryPerformance(userID uuid.UUID) ([]models.CategoryPerformance, error) {
	query := `
		SELECT
			q.category_id,
			q.category_id as category_name,
			COUNT(CASE WHEN qa.is_completed = true THEN 1 END) as quizzes_completed,
			COALESCE(AVG(CASE WHEN qa.is_completed = true THEN qa.score END), 0) as average_score,
			COUNT(*) as total_attempts,
			COALESCE(MAX(qa.score), 0) as best_score,
			MAX(qa.completed_at) as last_attempt
		FROM quiz_attempts qa
		JOIN quizzes q ON qa.quiz_id = q.id
		WHERE qa.user_id = $1
		GROUP BY q.category_id
		ORDER BY quizzes_completed DESC
	`

	rows, err := ur.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryPerformance []models.CategoryPerformance
	for rows.Next() {
		var cp models.CategoryPerformance
		err := rows.Scan(
			&cp.CategoryID, &cp.CategoryName, &cp.QuizzesCompleted,
			&cp.AverageScore, &cp.TotalAttempts, &cp.BestScore, &cp.LastAttempt,
		)
		if err != nil {
			return nil, err
		}
		categoryPerformance = append(categoryPerformance, cp)
	}

	return categoryPerformance, nil
}

// getRecentActivity retrieves recent quiz activities
func (ur *UserRepository) getRecentActivity(userID uuid.UUID) ([]models.RecentActivityItem, error) {
	query := `
		SELECT
			qa.quiz_id,
			q.title,
			qa.score,
			q.category_id,
			qa.completed_at,
			qa.time_spent,
			q.difficulty
		FROM quiz_attempts qa
		JOIN quizzes q ON qa.quiz_id = q.id
		WHERE qa.user_id = $1 AND qa.is_completed = true
		ORDER BY qa.completed_at DESC
		LIMIT 10
	`

	rows, err := ur.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.RecentActivityItem
	for rows.Next() {
		var activity models.RecentActivityItem
		err := rows.Scan(
			&activity.QuizID, &activity.QuizTitle, &activity.Score,
			&activity.Category, &activity.CompletedAt, &activity.TimeSpent,
			&activity.Difficulty,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

// getQuizzesByDifficulty retrieves quiz completion counts by difficulty
func (ur *UserRepository) getQuizzesByDifficulty(userID uuid.UUID) (map[string]int, error) {
	query := `
		SELECT
			q.difficulty,
			COUNT(CASE WHEN qa.is_completed = true THEN 1 END) as completed_count
		FROM quiz_attempts qa
		JOIN quizzes q ON qa.quiz_id = q.id
		WHERE qa.user_id = $1
		GROUP BY q.difficulty
	`

	rows, err := ur.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quizzesByDifficulty := make(map[string]int)
	for rows.Next() {
		var difficulty string
		var count int
		err := rows.Scan(&difficulty, &count)
		if err != nil {
			return nil, err
		}
		quizzesByDifficulty[difficulty] = count
	}

	return quizzesByDifficulty, nil
}

// getScoreDistribution retrieves score distribution
func (ur *UserRepository) getScoreDistribution(userID uuid.UUID) (models.ScoreDistribution, error) {
	query := `
		SELECT
			COUNT(CASE WHEN score >= 0 AND score <= 20 THEN 1 END) as range_0_to_20,
			COUNT(CASE WHEN score > 20 AND score <= 40 THEN 1 END) as range_21_to_40,
			COUNT(CASE WHEN score > 40 AND score <= 60 THEN 1 END) as range_41_to_60,
			COUNT(CASE WHEN score > 60 AND score <= 80 THEN 1 END) as range_61_to_80,
			COUNT(CASE WHEN score > 80 AND score <= 100 THEN 1 END) as range_81_to_100
		FROM quiz_attempts
		WHERE user_id = $1 AND is_completed = true
	`

	var distribution models.ScoreDistribution
	err := ur.db.QueryRow(query, userID).Scan(
		&distribution.Range0to20, &distribution.Range21to40, &distribution.Range41to60,
		&distribution.Range61to80, &distribution.Range81to100,
	)
	if err != nil {
		return models.ScoreDistribution{}, err
	}

	return distribution, nil
}

// getMonthlyProgress retrieves monthly progress data
func (ur *UserRepository) getMonthlyProgress(userID uuid.UUID) ([]models.MonthlyProgressItem, error) {
	query := `
		SELECT
			TO_CHAR(qa.completed_at, 'YYYY-MM') as month,
			COUNT(CASE WHEN qa.is_completed = true THEN 1 END) as quizzes_completed,
			COALESCE(AVG(CASE WHEN qa.is_completed = true THEN qa.score END), 0) as average_score,
			COALESCE(SUM(CASE WHEN qa.is_completed = true THEN qa.total_points END), 0) as total_points
		FROM quiz_attempts qa
		WHERE qa.user_id = $1 AND qa.completed_at IS NOT NULL
		  AND qa.completed_at >= CURRENT_DATE - INTERVAL '12 months'
		GROUP BY TO_CHAR(qa.completed_at, 'YYYY-MM')
		ORDER BY month DESC
		LIMIT 12
	`

	rows, err := ur.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyProgress []models.MonthlyProgressItem
	for rows.Next() {
		var item models.MonthlyProgressItem
		err := rows.Scan(&item.Month, &item.QuizzesCompleted, &item.AverageScore, &item.TotalPoints)
		if err != nil {
			return nil, err
		}
		monthlyProgress = append(monthlyProgress, item)
	}

	return monthlyProgress, nil
}
