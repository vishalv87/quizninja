package repository

import (
	"database/sql"
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
		INSERT INTO users (email, password_hash, name, age)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at, level, total_points, current_streak,
		          best_streak, total_quizzes_completed, average_score, is_online, last_active
	`
	err := ur.db.QueryRow(query, user.Email, user.PasswordHash, user.Name, user.Age).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Level, &user.TotalPoints,
		&user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
		&user.AverageScore, &user.IsOnline, &user.LastActive,
	)
	return err
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, age, level, total_points, current_streak,
		       best_streak, total_quizzes_completed, average_score, is_online,
		       last_active, avatar_url, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := ur.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Age, &user.Level,
		&user.TotalPoints, &user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
		&user.AverageScore, &user.IsOnline, &user.LastActive, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
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
		       last_active, avatar_url, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := ur.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Age, &user.Level,
		&user.TotalPoints, &user.CurrentStreak, &user.BestStreak, &user.TotalQuizzesCompleted,
		&user.AverageScore, &user.IsOnline, &user.LastActive, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
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
		                              notifications_enabled, notification_frequency)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := ur.db.QueryRow(query, preferences.UserID, preferences.SelectedInterests,
		preferences.DifficultyPreference, preferences.NotificationsEnabled,
		preferences.NotificationFrequency).Scan(&preferences.ID, &preferences.CreatedAt)
	return err
}

// GetUserPreferences retrieves user preferences by user ID
func (ur *UserRepository) GetUserPreferences(userID uuid.UUID) (*models.UserPreferences, error) {
	preferences := &models.UserPreferences{}
	query := `
		SELECT id, user_id, selected_interests, difficulty_preference,
		       notifications_enabled, notification_frequency, created_at
		FROM user_preferences
		WHERE user_id = $1
	`
	err := ur.db.QueryRow(query, userID).Scan(
		&preferences.ID, &preferences.UserID, &preferences.SelectedInterests,
		&preferences.DifficultyPreference, &preferences.NotificationsEnabled,
		&preferences.NotificationFrequency, &preferences.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return preferences, nil
}

// UpdateUserPreferences updates user preferences
func (ur *UserRepository) UpdateUserPreferences(preferences *models.UserPreferences) error {
	query := `
		UPDATE user_preferences SET
			selected_interests = $2, difficulty_preference = $3,
			notifications_enabled = $4, notification_frequency = $5
		WHERE user_id = $1
	`
	_, err := ur.db.Exec(query, preferences.UserID, preferences.SelectedInterests,
		preferences.DifficultyPreference, preferences.NotificationsEnabled,
		preferences.NotificationFrequency)
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
		return nil, err
	}

	preferences, err := ur.GetUserPreferences(userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err != sql.ErrNoRows {
		user.Preferences = preferences
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