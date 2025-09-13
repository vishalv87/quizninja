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
		INSERT INTO users (email, password_hash, name)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at, level, total_points
	`
	err := ur.db.QueryRow(query, user.Email, user.PasswordHash, user.Name).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Level, &user.TotalPoints,
	)
	return err
}

func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, age, level, total_points, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := ur.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Age, &user.Level, &user.TotalPoints, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, name, age, level, total_points, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := ur.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Age, &user.Level, &user.TotalPoints, &user.CreatedAt, &user.UpdatedAt,
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