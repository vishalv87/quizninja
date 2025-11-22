package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"quizninja-api/database"
	"quizninja-api/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
		INSERT INTO users (email, password_hash, name, auth_method, supabase_id, last_auth_method)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at, level, total_points, current_streak,
		          best_streak, total_quizzes_completed, average_score, is_online, last_active,
		          auth_method, supabase_id, last_auth_method, migrated_at
	`

	// Set default values if not provided
	if user.AuthMethod == "" {
		user.AuthMethod = "jwt"
	}
	if user.LastAuthMethod == "" {
		user.LastAuthMethod = user.AuthMethod
	}

	err := ur.db.QueryRow(query, user.Email, user.PasswordHash, user.Name,
		user.AuthMethod, user.SupabaseID, user.LastAuthMethod).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Level, &user.TotalPoints,
		&user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
		&user.AverageScore, &user.IsOnline, &user.LastActive,
		&user.AuthMethod, &user.SupabaseID, &user.LastAuthMethod, &user.MigratedAt,
	)
	return err
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, level, total_points, current_streak,
		       best_streak, total_quizzes_completed, average_score, is_online,
		       last_active, avatar_url, created_at, updated_at,
		       auth_method, supabase_id, last_auth_method, migrated_at
		FROM users
		WHERE email = $1
	`
	err := ur.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Level,
		&user.TotalPoints, &user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
		&user.AverageScore, &user.IsOnline, &user.LastActive, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
		&user.AuthMethod, &user.SupabaseID, &user.LastAuthMethod, &user.MigratedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, level, total_points, current_streak,
		       best_streak, total_quizzes_completed, average_score, is_online,
		       last_active, avatar_url, created_at, updated_at,
		       auth_method, supabase_id, last_auth_method, migrated_at
		FROM users
		WHERE id = $1
	`
	err := ur.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Level,
		&user.TotalPoints, &user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
		&user.AverageScore, &user.IsOnline, &user.LastActive, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
		&user.AuthMethod, &user.SupabaseID, &user.LastAuthMethod, &user.MigratedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates an existing user
func (ur *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users SET
			email = $2, name = $3, level = $4, total_points = $5,
			current_streak = $6, best_streak = $7, total_quizzes_completed = $8,
			average_score = $9, is_online = $10, last_active = $11, avatar_url = $12,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := ur.db.Exec(query, user.ID, user.Email, user.Name, user.Level,
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
		INSERT INTO user_preferences (
			user_id, selected_categories, difficulty_preference, notifications_enabled,
			notification_frequency, profile_visibility, show_online_status,
			allow_friend_requests, share_activity_status, notification_types,
			onboarding_completed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	// Handle JSONB field properly - marshal to JSON bytes
	var notificationTypesData map[string]interface{}
	if len(preferences.NotificationTypes) > 0 {
		notificationTypesData = preferences.NotificationTypes
	} else {
		notificationTypesData = map[string]interface{}{
			"challenges":           true,
			"achievements":         true,
			"quiz_reminders":       true,
			"friend_activity":      true,
			"leaderboard_updates":  false,
			"system_announcements": true,
		}
	}

	notificationTypesJSON, err := json.Marshal(notificationTypesData)
	if err != nil {
		return fmt.Errorf("failed to marshal notification_types: %w", err)
	}

	err = ur.db.QueryRow(query,
		preferences.UserID, preferences.SelectedCategories, preferences.DifficultyPreference,
		preferences.NotificationsEnabled, preferences.NotificationFrequency,
		preferences.ProfileVisibility, preferences.ShowOnlineStatus,
		preferences.AllowFriendRequests, preferences.ShareActivityStatus,
		string(notificationTypesJSON), preferences.OnboardingCompletedAt).Scan(&preferences.ID, &preferences.CreatedAt, &preferences.UpdatedAt)
	return err
}

// GetUserPreferences retrieves user preferences by user ID
func (ur *UserRepository) GetUserPreferences(userID uuid.UUID) (*models.UserPreferences, error) {
	preferences := &models.UserPreferences{}
	query := `
		SELECT id, user_id, selected_categories, difficulty_preference,
		       notifications_enabled, notification_frequency, profile_visibility,
		       show_online_status, allow_friend_requests, share_activity_status,
		       notification_types, onboarding_completed_at, created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`

	var selectedCategoriesSlice []string
	var notificationTypesJSON string
	err := ur.db.QueryRow(query, userID).Scan(
		&preferences.ID, &preferences.UserID, pq.Array(&selectedCategoriesSlice),
		&preferences.DifficultyPreference, &preferences.NotificationsEnabled,
		&preferences.NotificationFrequency, &preferences.ProfileVisibility,
		&preferences.ShowOnlineStatus, &preferences.AllowFriendRequests,
		&preferences.ShareActivityStatus, &notificationTypesJSON,
		&preferences.OnboardingCompletedAt, &preferences.CreatedAt,
		&preferences.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Convert []string to StringArray
	preferences.SelectedCategories = models.StringArray(selectedCategoriesSlice)

	// Unmarshal the JSON for notification_types
	if notificationTypesJSON != "" {
		err = json.Unmarshal([]byte(notificationTypesJSON), &preferences.NotificationTypes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal notification_types: %w", err)
		}
	}

	return preferences, nil
}

// UpdateUserPreferences updates or inserts user preferences (UPSERT)
func (ur *UserRepository) UpdateUserPreferences(preferences *models.UserPreferences) error {
	query := `
		INSERT INTO user_preferences (
			user_id, selected_categories, difficulty_preference,
			notifications_enabled, notification_frequency,
			profile_visibility, show_online_status,
			allow_friend_requests, share_activity_status,
			notification_types, onboarding_completed_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
		)
		ON CONFLICT (user_id) DO UPDATE SET
			selected_categories = EXCLUDED.selected_categories,
			difficulty_preference = EXCLUDED.difficulty_preference,
			notifications_enabled = EXCLUDED.notifications_enabled,
			notification_frequency = EXCLUDED.notification_frequency,
			profile_visibility = EXCLUDED.profile_visibility,
			show_online_status = EXCLUDED.show_online_status,
			allow_friend_requests = EXCLUDED.allow_friend_requests,
			share_activity_status = EXCLUDED.share_activity_status,
			notification_types = EXCLUDED.notification_types,
			onboarding_completed_at = EXCLUDED.onboarding_completed_at,
			updated_at = CURRENT_TIMESTAMP
	`
	// Handle JSONB field properly - marshal to JSON bytes
	var notificationTypesData map[string]interface{}
	if len(preferences.NotificationTypes) > 0 {
		notificationTypesData = preferences.NotificationTypes
	} else {
		// Use default notification types if nil or empty
		notificationTypesData = map[string]interface{}{
			"challenges":           true,
			"achievements":         true,
			"quiz_reminders":       true,
			"friend_activity":      true,
			"leaderboard_updates":  false,
			"system_announcements": true,
		}
	}

	notificationTypesJSON, err := json.Marshal(notificationTypesData)
	if err != nil {
		return fmt.Errorf("failed to marshal notification_types: %w", err)
	}

	_, err = ur.db.Exec(query,
		preferences.UserID,
		preferences.SelectedCategories,
		preferences.DifficultyPreference,
		preferences.NotificationsEnabled,
		preferences.NotificationFrequency,
		preferences.ProfileVisibility,
		preferences.ShowOnlineStatus,
		preferences.AllowFriendRequests,
		preferences.ShareActivityStatus,
		string(notificationTypesJSON),
		preferences.OnboardingCompletedAt)
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
	query := `UPDATE users SET is_online = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := ur.db.Exec(query, userID, isOnline)
	return err
}

// UpdateUserLastActive updates user's last active timestamp
func (ur *UserRepository) UpdateUserLastActive(userID uuid.UUID) error {
	query := `UPDATE users SET last_active = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
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
		SELECT total_points, current_streak, best_streak, total_quizzes_completed, average_score
		FROM users
		WHERE id = $1
	`
	err := ur.db.QueryRow(userQuery, userID).Scan(
		&stats.TotalPoints, &stats.CurrentStreak, &stats.BestStreak,
		&stats.CompletedQuizzes, &stats.AverageScore,
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

	// First, get current user statistics including streak values
	var currentTotalQuizzes int
	var currentAverageScore sql.NullFloat64
	var currentStreak int
	var bestStreak int

	query := `
		SELECT total_quizzes_completed, average_score, current_streak, best_streak
		FROM users
		WHERE id = $1
	`

	err = tx.QueryRow(query, userID).Scan(&currentTotalQuizzes, &currentAverageScore, &currentStreak, &bestStreak)
	if err != nil {
		return fmt.Errorf("failed to get current user statistics: %w", err)
	}

	// Get the last quiz completion date (before this current quiz)
	var lastCompletedAt sql.NullTime
	lastQuizQuery := `
		SELECT MAX(completed_at)
		FROM quiz_attempts
		WHERE user_id = $1 AND is_completed = true AND completed_at < CURRENT_TIMESTAMP
	`
	err = tx.QueryRow(lastQuizQuery, userID).Scan(&lastCompletedAt)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get last quiz completion date: %w", err)
	}

	// Calculate new streak
	newStreak := currentStreak
	now := time.Now()

	if !lastCompletedAt.Valid {
		// This is the user's first quiz completion
		newStreak = 1
	} else {
		// Calculate days since last quiz
		lastDate := lastCompletedAt.Time.Truncate(24 * time.Hour)
		currentDate := now.Truncate(24 * time.Hour)
		daysDiff := int(currentDate.Sub(lastDate).Hours() / 24)

		if daysDiff == 0 {
			// Completed another quiz on the same day - streak continues
			newStreak = currentStreak
		} else if daysDiff == 1 {
			// Completed quiz on consecutive day - increment streak
			newStreak = currentStreak + 1
		} else {
			// More than 1 day gap - streak breaks, reset to 1
			newStreak = 1
		}
	}

	// Update best streak if current streak is higher
	newBestStreak := bestStreak
	if newStreak > bestStreak {
		newBestStreak = newStreak
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

	// Update user statistics including streaks
	updateQuery := `
		UPDATE users
		SET total_quizzes_completed = $1,
		    average_score = $2,
		    current_streak = $3,
		    best_streak = $4,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	_, err = tx.Exec(updateQuery, newTotalQuizzes, newAverageScore, newStreak, newBestStreak, userID)
	if err != nil {
		return fmt.Errorf("failed to update user statistics: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Log streak update for debugging
	log.Printf("Updated streaks for user %s: current_streak=%d, best_streak=%d (days_since_last=%v)",
		userID, newStreak, newBestStreak,
		func() string {
			if lastCompletedAt.Valid {
				return fmt.Sprintf("%d", int(now.Sub(lastCompletedAt.Time).Hours()/24))
			}
			return "first_quiz"
		}())

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

// Supabase-related user management methods

// GetUserBySupabaseID retrieves a user by their Supabase ID
func (ur *UserRepository) GetUserBySupabaseID(supabaseID string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, level, total_points, current_streak,
		       best_streak, total_quizzes_completed, average_score, is_online,
		       last_active, avatar_url, created_at, updated_at,
		       auth_method, supabase_id, last_auth_method, migrated_at
		FROM users
		WHERE supabase_id = $1
	`
	err := ur.db.QueryRow(query, supabaseID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Level,
		&user.TotalPoints, &user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
		&user.AverageScore, &user.IsOnline, &user.LastActive, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
		&user.AuthMethod, &user.SupabaseID, &user.LastAuthMethod, &user.MigratedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUserWithSupabaseID creates a user with Supabase authentication
func (ur *UserRepository) CreateUserWithSupabaseID(user *models.User, supabaseID string) error {
	user.AuthMethod = "supabase"
	user.LastAuthMethod = "supabase"
	user.SupabaseID = &supabaseID

	return ur.CreateUser(user)
}

// LinkUserToSupabase links an existing JWT user to Supabase
func (ur *UserRepository) LinkUserToSupabase(userID uuid.UUID, supabaseID string) error {
	query := `
		UPDATE users
		SET auth_method = 'supabase',
		    supabase_id = $2,
		    last_auth_method = 'supabase',
		    migrated_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := ur.db.Exec(query, userID, supabaseID)
	return err
}

// UnlinkUserFromSupabase removes Supabase link and reverts to JWT
func (ur *UserRepository) UnlinkUserFromSupabase(userID uuid.UUID) error {
	query := `
		UPDATE users
		SET auth_method = 'jwt',
		    supabase_id = NULL,
		    last_auth_method = 'jwt',
		    migrated_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := ur.db.Exec(query, userID)
	return err
}

// UpdateUserAuthMethod updates the last authentication method used
func (ur *UserRepository) UpdateUserAuthMethod(userID uuid.UUID, authMethod string) error {
	query := `
		UPDATE users
		SET last_auth_method = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := ur.db.Exec(query, userID, authMethod)
	return err
}

// GetUsersByAuthMethod retrieves users by their authentication method
func (ur *UserRepository) GetUsersByAuthMethod(authMethod string, limit int) ([]*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, level, total_points, current_streak,
		       best_streak, total_quizzes_completed, average_score, is_online,
		       last_active, avatar_url, created_at, updated_at,
		       auth_method, supabase_id, last_auth_method, migrated_at
		FROM users
		WHERE auth_method = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := ur.db.Query(query, authMethod, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Level,
			&user.TotalPoints, &user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
			&user.AverageScore, &user.IsOnline, &user.LastActive, &user.AvatarURL,
			&user.CreatedAt, &user.UpdatedAt,
			&user.AuthMethod, &user.SupabaseID, &user.LastAuthMethod, &user.MigratedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetAuthMethodStats returns statistics about authentication methods
func (ur *UserRepository) GetAuthMethodStats() (map[string]int, error) {
	query := `
		SELECT auth_method, COUNT(*) as count
		FROM users
		GROUP BY auth_method
	`

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var authMethod string
		var count int
		err := rows.Scan(&authMethod, &count)
		if err != nil {
			return nil, err
		}
		stats[authMethod] = count
	}

	return stats, nil
}

// GetMigrationCandidates returns users who could be migrated to Supabase
func (ur *UserRepository) GetMigrationCandidates(limit int) ([]*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, level, total_points, current_streak,
		       best_streak, total_quizzes_completed, average_score, is_online,
		       last_active, avatar_url, created_at, updated_at,
		       auth_method, supabase_id, last_auth_method, migrated_at
		FROM users
		WHERE auth_method = 'jwt' AND supabase_id IS NULL
		ORDER BY last_active DESC
		LIMIT $1
	`

	rows, err := ur.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Level,
			&user.TotalPoints, &user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
			&user.AverageScore, &user.IsOnline, &user.LastActive, &user.AvatarURL,
			&user.CreatedAt, &user.UpdatedAt,
			&user.AuthMethod, &user.SupabaseID, &user.LastAuthMethod, &user.MigratedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
