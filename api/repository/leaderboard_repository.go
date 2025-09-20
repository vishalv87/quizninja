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

	// Since we don't have real user data yet, let's return sample data for now
	// In production, this would query the actual database
	sampleEntries := lr.getSampleLeaderboardData(period)

	// Apply pagination
	start := offset
	if start > len(sampleEntries) {
		start = len(sampleEntries)
	}

	end := start + limit
	if end > len(sampleEntries) {
		end = len(sampleEntries)
	}

	if start < len(sampleEntries) {
		entries = sampleEntries[start:end]
	} else {
		entries = []models.LeaderboardEntry{} // Return empty slice, not nil
	}

	// Ensure entries is never nil
	if entries == nil {
		entries = []models.LeaderboardEntry{}
	}
	totalCount = len(sampleEntries)

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
	// Get user's current points
	var userPoints int
	userQuery := "SELECT total_points FROM users WHERE id = $1"
	err := lr.db.QueryRow(userQuery, userID).Scan(&userPoints)
	if err != nil {
		return nil, fmt.Errorf("failed to get user points: %w", err)
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

	return &models.UserRankInfo{
		UserID:       userID,
		Rank:         rank,
		Points:       userPoints,
		PointsToNext: pointsToNext,
		RankChange:   rankChange,
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

	return nil
}

// RecalculateUserLevel updates user's level based on their total points
func (lr *LeaderboardRepository) RecalculateUserLevel(userID uuid.UUID) error {
	// The level is automatically updated by the database trigger when total_points changes
	// This method is kept for compatibility and future enhancements
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

// getSampleLeaderboardData returns sample leaderboard data for testing
func (lr *LeaderboardRepository) getSampleLeaderboardData(period string) []models.LeaderboardEntry {
	// In a real implementation, this would filter data based on the period
	_ = period
	sampleEntries := []models.LeaderboardEntry{
		{
			UserID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			Name:              "Alex Johnson",
			Avatar:            stringPtr("https://i.pravatar.cc/150?img=1"),
			Rank:              1,
			Points:            2850,
			QuizzesCompleted:  45,
			AverageScore:      87.5,
			CurrentStreak:     12,
			Level:             "15",
			LastActive:        time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
			IsCurrentUser:     false,
			IsFriend:          true,
			Achievements:      []string{"Quiz Master", "Speed Demon", "Perfect Score"},
			CategoryPoints:    map[string]int{"Science": 950, "History": 800, "Sports": 600, "Math": 500},
		},
		{
			UserID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			Name:              "Sarah Chen",
			Avatar:            stringPtr("https://i.pravatar.cc/150?img=2"),
			Rank:              2,
			Points:            2720,
			QuizzesCompleted:  38,
			AverageScore:      89.2,
			CurrentStreak:     8,
			Level:             "14",
			LastActive:        time.Date(2024, 1, 15, 12, 15, 0, 0, time.UTC),
			IsCurrentUser:     false,
			IsFriend:          true,
			Achievements:      []string{"Perfectionist", "Science Whiz"},
			CategoryPoints:    map[string]int{"Science": 1200, "Math": 800, "History": 720},
		},
		{
			UserID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
			Name:              "Mike Rodriguez",
			Avatar:            stringPtr("https://i.pravatar.cc/150?img=3"),
			Rank:              3,
			Points:            2450,
			QuizzesCompleted:  32,
			AverageScore:      82.1,
			CurrentStreak:     5,
			Level:             "13",
			LastActive:        time.Date(2024, 1, 15, 10, 45, 0, 0, time.UTC),
			IsCurrentUser:     true,
			IsFriend:          false,
			Achievements:      []string{"Consistency King", "Sports Expert"},
			CategoryPoints:    map[string]int{"Sports": 1100, "History": 700, "Science": 650},
		},
		{
			UserID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
			Name:              "Emma Wilson",
			Avatar:            stringPtr("https://i.pravatar.cc/150?img=4"),
			Rank:              4,
			Points:            2180,
			QuizzesCompleted:  29,
			AverageScore:      85.6,
			CurrentStreak:     3,
			Level:             "12",
			LastActive:        time.Date(2024, 1, 14, 18, 20, 0, 0, time.UTC),
			IsCurrentUser:     false,
			IsFriend:          false,
			Achievements:      []string{"History Buff", "Rising Star"},
			CategoryPoints:    map[string]int{"History": 900, "Math": 680, "Science": 600},
		},
		{
			UserID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
			Name:              "David Kim",
			Avatar:            stringPtr("https://i.pravatar.cc/150?img=5"),
			Rank:              5,
			Points:            1950,
			QuizzesCompleted:  26,
			AverageScore:      78.9,
			CurrentStreak:     7,
			Level:             "11",
			LastActive:        time.Date(2024, 1, 14, 16, 30, 0, 0, time.UTC),
			IsCurrentUser:     false,
			IsFriend:          true,
			Achievements:      []string{"Dedicated Player", "Math Genius"},
			CategoryPoints:    map[string]int{"Math": 850, "Science": 550, "Sports": 550},
		},
	}

	return sampleEntries
}

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
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

