package repository

import (
	"database/sql"
	"fmt"
	"time"

	"quizninja-api/database"
	"quizninja-api/models"

	"github.com/google/uuid"
)

type LeaderboardRepository struct {
	db *sql.DB
}

func NewLeaderboardRepository() *LeaderboardRepository {
	return &LeaderboardRepository{
		db: database.DB,
	}
}

// GetGlobalLeaderboard retrieves the global leaderboard for a given period
func (lr *LeaderboardRepository) GetGlobalLeaderboard(period string, limit, offset int) ([]models.LeaderboardEntry, int, error) {
	var entries []models.LeaderboardEntry
	var totalCount int

	// Count total users with points > 0
	countQuery := `SELECT COUNT(*) FROM users WHERE total_points > 0`
	err := lr.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return []models.LeaderboardEntry{}, 0, fmt.Errorf("failed to count leaderboard entries: %w", err)
	}

	// Get leaderboard entries ordered by total_points
	baseQuery := `
		SELECT
			u.id,
			u.name,
			u.avatar_url,
			u.total_points,
			u.total_quizzes_completed,
			u.average_score,
			u.current_streak,
			u.level,
			u.last_active
		FROM users u
		WHERE u.total_points > 0
		ORDER BY u.total_points DESC, u.average_score DESC, u.current_streak DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := lr.db.Query(baseQuery, limit, offset)
	if err != nil {
		return []models.LeaderboardEntry{}, 0, fmt.Errorf("failed to query global leaderboard: %w", err)
	}
	defer rows.Close()

	rank := offset + 1
	for rows.Next() {
		var entry models.LeaderboardEntry
		var avatarURL sql.NullString
		var lastActive sql.NullTime

		err := rows.Scan(
			&entry.UserID,
			&entry.Name,
			&avatarURL,
			&entry.Points,
			&entry.QuizzesCompleted,
			&entry.AverageScore,
			&entry.CurrentStreak,
			&entry.Level,
			&lastActive,
		)
		if err != nil {
			return []models.LeaderboardEntry{}, 0, fmt.Errorf("failed to scan leaderboard entry: %w", err)
		}

		if avatarURL.Valid {
			entry.Avatar = &avatarURL.String
		}

		if lastActive.Valid {
			entry.LastActive = lastActive.Time
		} else {
			entry.LastActive = time.Now()
		}

		entry.Rank = rank
		rank++

		// Default values for global leaderboard (no user context)
		entry.IsCurrentUser = false
		entry.IsFriend = false

		// Get user achievements
		achievements, err := lr.GetUserAchievements(entry.UserID)
		if err != nil {
			achievements = []string{}
		}
		entry.Achievements = achievements

		// Get user category points
		categoryPoints, err := lr.GetUserCategoryPoints(entry.UserID)
		if err != nil {
			categoryPoints = make(map[string]int)
		}
		entry.CategoryPoints = categoryPoints

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return []models.LeaderboardEntry{}, 0, fmt.Errorf("error iterating leaderboard rows: %w", err)
	}

	// Ensure entries is never nil
	if entries == nil {
		entries = []models.LeaderboardEntry{}
	}

	return entries, totalCount, nil
}

// GetFriendsLeaderboard retrieves the leaderboard filtered to friends only
func (lr *LeaderboardRepository) GetFriendsLeaderboard(userID uuid.UUID, period string, limit, offset int) ([]models.LeaderboardEntry, int, error) {
	var entries []models.LeaderboardEntry
	var totalCount int

	// Build query that includes friends and current user
	baseQuery := `
		SELECT
			u.id,
			u.name,
			u.avatar_url,
			u.total_points,
			u.total_quizzes_completed,
			u.average_score,
			u.current_streak,
			u.level,
			u.last_active,
			CASE WHEN u.id = $1 THEN true ELSE false END as is_current_user,
			CASE WHEN f.user1_id IS NOT NULL OR f.user2_id IS NOT NULL THEN true ELSE false END as is_friend
		FROM users u
		LEFT JOIN friendships f ON (f.user1_id = $1 AND f.user2_id = u.id) OR (f.user2_id = $1 AND f.user1_id = u.id)
		WHERE u.total_points > 0
		AND (u.id = $1 OR f.user1_id IS NOT NULL OR f.user2_id IS NOT NULL)
	`

	// Add time filtering based on period
	timeFilter := lr.getTimeFilterForPeriod(period)
	if timeFilter != "" {
		baseQuery += " AND " + timeFilter
	}

	// Count total matching records
	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") as friends_leaderboard_count"
	err := lr.db.QueryRow(countQuery, userID).Scan(&totalCount)
	if err != nil {
		return []models.LeaderboardEntry{}, 0, fmt.Errorf("failed to count friends leaderboard entries: %w", err)
	}

	// Add ordering and pagination
	baseQuery += " ORDER BY u.total_points DESC, u.average_score DESC, u.current_streak DESC"
	baseQuery += " LIMIT $2 OFFSET $3"

	rows, err := lr.db.Query(baseQuery, userID, limit, offset)
	if err != nil {
		return []models.LeaderboardEntry{}, 0, fmt.Errorf("failed to query friends leaderboard: %w", err)
	}
	defer rows.Close()

	rank := offset + 1
	for rows.Next() {
		var entry models.LeaderboardEntry
		var avatarURL sql.NullString

		err := rows.Scan(
			&entry.UserID,
			&entry.Name,
			&avatarURL,
			&entry.Points,
			&entry.QuizzesCompleted,
			&entry.AverageScore,
			&entry.CurrentStreak,
			&entry.Level,
			&entry.LastActive,
			&entry.IsCurrentUser,
			&entry.IsFriend,
		)
		if err != nil {
			return []models.LeaderboardEntry{}, 0, fmt.Errorf("failed to scan friends leaderboard entry: %w", err)
		}

		if avatarURL.Valid {
			entry.Avatar = &avatarURL.String
		}

		entry.Rank = rank
		rank++

		// Get user achievements
		achievements, err := lr.GetUserAchievements(entry.UserID)
		if err != nil {
			achievements = []string{}
		}
		entry.Achievements = achievements

		// Get user category points
		categoryPoints, err := lr.GetUserCategoryPoints(entry.UserID)
		if err != nil {
			categoryPoints = make(map[string]int)
		}
		entry.CategoryPoints = categoryPoints

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return []models.LeaderboardEntry{}, 0, fmt.Errorf("error iterating friends leaderboard rows: %w", err)
	}

	// Ensure entries is never nil
	if entries == nil {
		entries = []models.LeaderboardEntry{}
	}

	return entries, totalCount, nil
}

// GetUserRank retrieves the current user's rank information
func (lr *LeaderboardRepository) GetUserRank(userID uuid.UUID, period string) (*models.UserRankInfo, error) {
	// Get user's profile and stats
	var userPoints int
	var quizzesCompleted int
	var fullName string
	var avatarURL sql.NullString

	userQuery := `
		SELECT total_points, total_quizzes_completed, name, avatar_url
		FROM users
		WHERE id = $1
	`
	err := lr.db.QueryRow(userQuery, userID).Scan(&userPoints, &quizzesCompleted, &fullName, &avatarURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user data: %w", err)
	}

	// Get user's achievements count
	var achievementsUnlocked int
	achievementsQuery := `
		SELECT COUNT(*)
		FROM user_achievements
		WHERE user_id = $1
	`
	err = lr.db.QueryRow(achievementsQuery, userID).Scan(&achievementsUnlocked)
	if err != nil {
		// Don't fail if achievements query fails, just set to 0
		achievementsUnlocked = 0
	}

	// Calculate user's rank
	rankQuery := `
		SELECT COUNT(*) + 1 as rank
		FROM users u
		WHERE u.total_points > $1
		AND u.total_points > 0
	`

	// Add time filtering based on period
	timeFilter := lr.getTimeFilterForPeriod(period)
	if timeFilter != "" {
		rankQuery += " AND " + timeFilter
	}

	var rank int
	err = lr.db.QueryRow(rankQuery, userPoints).Scan(&rank)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate user rank: %w", err)
	}

	// Get total users count (include current user even if they have 0 points)
	totalUsersQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE total_points > 0 OR id = $1
	`
	if timeFilter != "" {
		totalUsersQuery += " AND " + timeFilter
	}

	var totalUsers int
	err = lr.db.QueryRow(totalUsersQuery, userID).Scan(&totalUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}

	// Get points needed to reach next rank
	nextRankQuery := `
		SELECT MIN(total_points) as next_points
		FROM users u
		WHERE u.total_points > $1
		AND u.total_points > 0
	`

	if timeFilter != "" {
		nextRankQuery += " AND " + timeFilter
	}

	var nextPoints sql.NullInt64
	err = lr.db.QueryRow(nextRankQuery, userPoints).Scan(&nextPoints)
	if err != nil {
		return nil, fmt.Errorf("failed to get next rank points: %w", err)
	}

	pointsToNext := 0
	if nextPoints.Valid {
		pointsToNext = int(nextPoints.Int64) - userPoints
	}

	// For now, set rank change to 0 (would need historical data to calculate)
	rankChange := 0

	// Build user info
	userInfo := models.UserInfo{
		ID:       userID.String(),
		FullName: fullName,
	}
	if avatarURL.Valid {
		userInfo.AvatarURL = &avatarURL.String
	}

	return &models.UserRankInfo{
		Rank:                 rank,
		TotalUsers:           totalUsers,
		User:                 userInfo,
		TotalPoints:          userPoints,
		QuizzesCompleted:     quizzesCompleted,
		AchievementsUnlocked: achievementsUnlocked,
		PointsToNext:         pointsToNext,
		RankChange:           rankChange,
	}, nil
}

// UpdateUserScore updates user's total points after completing a quiz
func (lr *LeaderboardRepository) UpdateUserScore(userID uuid.UUID, points int, quizID uuid.UUID) error {
	query := `
		UPDATE users
		SET total_points = total_points + $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err := lr.db.Exec(query, points, userID)
	if err != nil {
		return fmt.Errorf("failed to update user score: %w", err)
	}

	// Recalculate user level after points update
	if err := lr.RecalculateUserLevel(userID); err != nil {
		return fmt.Errorf("failed to recalculate user level: %w", err)
	}

	return nil
}

// calculateUserLevel determines the level based on total points
// Mirrors the logic from the removed database trigger
func calculateUserLevel(points int) string {
	switch {
	case points < 100:
		return "Beginner"
	case points < 300:
		return "Novice"
	case points < 600:
		return "Intermediate"
	case points < 1000:
		return "Advanced"
	case points < 1500:
		return "Expert"
	case points < 2000:
		return "Master"
	default:
		return "Legend"
	}
}

// RecalculateUserLevel updates user's level based on their total points
func (lr *LeaderboardRepository) RecalculateUserLevel(userID uuid.UUID) error {
	// Get current total points
	var totalPoints int
	getPointsQuery := `SELECT total_points FROM users WHERE id = $1`
	err := lr.db.QueryRow(getPointsQuery, userID).Scan(&totalPoints)
	if err != nil {
		return fmt.Errorf("failed to get user points: %w", err)
	}

	// Calculate new level
	newLevel := calculateUserLevel(totalPoints)

	// Update level
	updateQuery := `UPDATE users SET level = $1, updated_at = NOW() WHERE id = $2`
	_, err = lr.db.Exec(updateQuery, newLevel, userID)
	if err != nil {
		return fmt.Errorf("failed to update user level: %w", err)
	}

	return nil
}

// GetUserAchievements retrieves user's achievements for leaderboard display
func (lr *LeaderboardRepository) GetUserAchievements(userID uuid.UUID) ([]string, error) {
	query := `
		SELECT a.title
		FROM user_achievements ua
		JOIN achievements a ON ua.achievement_id = a.id
		WHERE ua.user_id = $1
		ORDER BY ua.unlocked_at DESC
	`

	rows, err := lr.db.Query(query, userID)
	if err != nil {
		return []string{}, fmt.Errorf("failed to query user achievements: %w", err)
	}
	defer rows.Close()

	var achievements []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return []string{}, fmt.Errorf("failed to scan achievement: %w", err)
		}
		achievements = append(achievements, title)
	}

	if err = rows.Err(); err != nil {
		return []string{}, fmt.Errorf("error iterating achievements: %w", err)
	}

	// Ensure achievements is never nil
	if achievements == nil {
		achievements = []string{}
	}

	return achievements, nil
}

// GetUserCategoryPoints retrieves user's points breakdown by category
func (lr *LeaderboardRepository) GetUserCategoryPoints(userID uuid.UUID) (map[string]int, error) {
	query := `
		SELECT category_name, total_points
		FROM user_category_performance
		WHERE user_id = $1
		ORDER BY total_points DESC
	`

	rows, err := lr.db.Query(query, userID)
	if err != nil {
		return make(map[string]int), fmt.Errorf("failed to query user category points: %w", err)
	}
	defer rows.Close()

	categoryPoints := make(map[string]int)
	for rows.Next() {
		var categoryName string
		var points int
		if err := rows.Scan(&categoryName, &points); err != nil {
			return make(map[string]int), fmt.Errorf("failed to scan category points: %w", err)
		}
		categoryPoints[categoryName] = points
	}

	if err = rows.Err(); err != nil {
		return make(map[string]int), fmt.Errorf("error iterating category points: %w", err)
	}

	return categoryPoints, nil
}

// getTimeFilterForPeriod returns the appropriate WHERE clause for time filtering
func (lr *LeaderboardRepository) getTimeFilterForPeriod(period string) string {
	switch period {
	case "today":
		return "u.created_at >= CURRENT_DATE"
	case "week":
		return "u.created_at >= CURRENT_DATE - INTERVAL '7 days'"
	case "month":
		return "u.created_at >= CURRENT_DATE - INTERVAL '30 days'"
	case "alltime":
		return "" // No time filter for all-time leaderboard
	default:
		return "" // Default to all-time if invalid period
	}
}
