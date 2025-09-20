package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	Email                 string     `json:"email" db:"email"`
	PasswordHash          string     `json:"-" db:"password_hash"`
	Name                  string     `json:"name" db:"name"`
	Age                   *int       `json:"age,omitempty" db:"age"`
	Level                 string     `json:"level" db:"level"`
	TotalPoints           int        `json:"total_points" db:"total_points"`
	CurrentStreak         int        `json:"current_streak" db:"current_streak"`
	BestStreak            int        `json:"best_streak" db:"best_streak"`
	TotalQuizzesCompleted int        `json:"total_quizzes_completed" db:"total_quizzes_completed"`
	AverageScore          float64    `json:"average_score" db:"average_score"`
	IsOnline              bool       `json:"is_online" db:"is_online"`
	LastActive            time.Time  `json:"last_active" db:"last_active"`
	AvatarURL             *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
	Preferences           *UserPreferences `json:"preferences,omitempty"`
}

type UserPreferences struct {
	ID                      uuid.UUID    `json:"id" db:"id"`
	UserID                  uuid.UUID    `json:"user_id" db:"user_id"`
	SelectedInterests       StringArray  `json:"selected_interests" db:"selected_interests"`
	DifficultyPreference    string       `json:"difficulty_preference" db:"difficulty_preference"`
	NotificationsEnabled    bool         `json:"notifications_enabled" db:"notifications_enabled"`
	NotificationFrequency   string       `json:"notification_frequency" db:"notification_frequency"`
	CreatedAt               time.Time    `json:"created_at" db:"created_at"`
}

// StringArray handles PostgreSQL array type
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return pq.Array(a).Value()
}

func (a *StringArray) Scan(value any) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}
	return pq.Array(a).Scan(value)
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type RegisterRequest struct {
	Email       string                     `json:"email" binding:"required,email"`
	Password    string                     `json:"password" binding:"required,min=6"`
	Name        string                     `json:"name" binding:"required"`
	Age         *int                       `json:"age,omitempty"`
	Preferences *UserPreferencesRequest    `json:"preferences,omitempty"`
}

type UserPreferencesRequest struct {
	SelectedInterests       []string `json:"selected_interests"`
	DifficultyPreference    string   `json:"difficulty_preference" binding:"omitempty,oneof=Easy Medium Hard"`
	NotificationsEnabled    bool     `json:"notifications_enabled"`
	NotificationFrequency   string   `json:"notification_frequency" binding:"omitempty,oneof=Daily Weekly Never"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileRequest struct {
	Name      *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	Age       *int    `json:"age,omitempty" binding:"omitempty,min=13,max=120"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// Interest represents quiz topic interests
type Interest struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"display_name" db:"display_name"`
	Description string `json:"description" db:"description"`
	IconURL     string `json:"icon_url" db:"icon_url"`
	IsActive    bool   `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// DifficultyLevel represents quiz difficulty levels
type DifficultyLevel struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"display_name" db:"display_name"`
	Description string `json:"description" db:"description"`
	Order       int    `json:"order" db:"order"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// NotificationFrequency represents notification frequency options
type NotificationFrequency struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"display_name" db:"display_name"`
	Description string `json:"description" db:"description"`
	IsActive    bool   `json:"is_active" db:"is_active"`
}

// AppSettings represents application configuration settings
type AppSettings struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Key              string    `json:"key" db:"key"`
	Value            string    `json:"value" db:"value"`
	Description      string    `json:"description" db:"description"`
	Category         string    `json:"category" db:"category"`
	IsPublic         bool      `json:"is_public" db:"is_public"`
	LastModifiedBy   uuid.UUID `json:"last_modified_by" db:"last_modified_by"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Category represents quiz categories for API response
type Category struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Description string      `json:"description"`
	IconURL     string      `json:"icon_url"`
	IsActive    bool        `json:"is_active"`
	Interests   []Interest  `json:"interests,omitempty"`
}

// API Request/Response types
type UpdatePreferencesRequest struct {
	SelectedInterests       []string `json:"selected_interests" binding:"dive,required"`
	DifficultyPreference    string   `json:"difficulty_preference" binding:"required,oneof=Easy Medium Hard"`
	NotificationsEnabled    bool     `json:"notifications_enabled"`
	NotificationFrequency   string   `json:"notification_frequency" binding:"required,oneof=Daily Weekly Never"`
}

type AppSettingsResponse struct {
	Settings map[string]any `json:"settings"`
	Version  string         `json:"version"`
	LastUpdated time.Time   `json:"last_updated"`
}

// Quiz models

type Quiz struct {
	ID                uuid.UUID    `json:"id" db:"id"`
	Title             string       `json:"title" db:"title"`
	Description       string       `json:"description" db:"description"`
	Category          string       `json:"category" db:"category_id"`
	Difficulty        string       `json:"difficulty" db:"difficulty"`
	TimeLimit         int          `json:"time_limit" db:"time_limit_minutes"` // in minutes
	QuestionCount     int          `json:"question_count" db:"total_questions"`
	IsFeatured        bool         `json:"is_featured" db:"is_featured"`
	IsPublic          bool         `json:"is_public" db:"is_public"`
	CreatedBy         uuid.UUID    `json:"created_by" db:"created_by"`
	Tags              StringArray  `json:"tags" db:"tags"`
	ThumbnailURL      *string      `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	CreatedAt         time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at" db:"updated_at"`
	Questions         []Question   `json:"questions,omitempty"`
	Statistics        *QuizStatistics `json:"statistics,omitempty"`
}

// MarshalJSON converts TimeLimit from minutes to seconds for JSON response
func (q Quiz) MarshalJSON() ([]byte, error) {
	type Alias Quiz
	return json.Marshal(&struct {
		TimeLimit int `json:"time_limit"` // Convert minutes to seconds
		*Alias
	}{
		TimeLimit: q.TimeLimit * 60, // Convert minutes to seconds
		Alias:     (*Alias)(&q),
	})
}

// UnmarshalJSON converts TimeLimit from seconds to minutes when receiving JSON
func (q *Quiz) UnmarshalJSON(data []byte) error {
	type Alias Quiz
	aux := &struct {
		TimeLimit int `json:"time_limit"` // Expect seconds from JSON
		*Alias
	}{
		Alias: (*Alias)(q),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	q.TimeLimit = aux.TimeLimit / 60 // Convert seconds to minutes for storage
	return nil
}

type Question struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	QuizID        uuid.UUID    `json:"quiz_id" db:"quiz_id"`
	QuestionText  string       `json:"question_text" db:"question_text"`
	QuestionType  string       `json:"question_type" db:"question_type"` // multiple_choice, true_false, short_answer
	Options       StringArray  `json:"options,omitempty" db:"options"`
	CorrectAnswer string       `json:"correct_answer" db:"correct_answer"`
	Explanation   *string      `json:"explanation,omitempty" db:"explanation"`
	Points        int          `json:"points" db:"points"`
	Order         int          `json:"order" db:"order"`
	ImageURL      *string      `json:"image_url,omitempty" db:"image_url"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

type QuizStatistics struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	QuizID             uuid.UUID `json:"quiz_id" db:"quiz_id"`
	TotalAttempts      int       `json:"total_attempts" db:"total_attempts"`
	CompletedAttempts  int       `json:"completed_attempts" db:"completed_attempts"`
	AverageScore       float64   `json:"average_score" db:"average_score"`
	AverageTime        int       `json:"average_time" db:"average_time"` // in seconds
	HighestScore       float64   `json:"highest_score" db:"highest_score"`
	LowestScore        float64   `json:"lowest_score" db:"lowest_score"`
	LastAttemptAt      *time.Time `json:"last_attempt_at,omitempty" db:"last_attempt_at"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// Quiz DTOs

type CreateQuizRequest struct {
	Title         string                   `json:"title" binding:"required,min=1,max=200"`
	Description   string                   `json:"description" binding:"required,min=1,max=1000"`
	Category      string                   `json:"category" binding:"required"`
	Difficulty    string                   `json:"difficulty" binding:"required,oneof=Easy Medium Hard"`
	TimeLimit     int                      `json:"time_limit" binding:"required,min=30,max=3600"`
	IsFeatured    bool                     `json:"is_featured"`
	IsPublic      bool                     `json:"is_public"`
	Tags          []string                 `json:"tags"`
	ThumbnailURL  *string                  `json:"thumbnail_url,omitempty"`
	Questions     []CreateQuestionRequest  `json:"questions" binding:"required,dive"`
}

type CreateQuestionRequest struct {
	QuestionText  string   `json:"question_text" binding:"required,min=1,max=500"`
	QuestionType  string   `json:"question_type" binding:"required,oneof=multiple_choice true_false short_answer"`
	Options       []string `json:"options,omitempty"`
	CorrectAnswer string   `json:"correct_answer" binding:"required"`
	Explanation   *string  `json:"explanation,omitempty"`
	Points        int      `json:"points" binding:"required,min=1,max=100"`
	Order         int      `json:"order" binding:"required,min=1"`
	ImageURL      *string  `json:"image_url,omitempty"`
}

type UpdateQuizRequest struct {
	Title         *string  `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Description   *string  `json:"description,omitempty" binding:"omitempty,min=1,max=1000"`
	Category      *string  `json:"category,omitempty"`
	Difficulty    *string  `json:"difficulty,omitempty" binding:"omitempty,oneof=Easy Medium Hard"`
	TimeLimit     *int     `json:"time_limit,omitempty" binding:"omitempty,min=30,max=3600"`
	IsFeatured    *bool    `json:"is_featured,omitempty"`
	IsPublic      *bool    `json:"is_public,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	ThumbnailURL  *string  `json:"thumbnail_url,omitempty"`
}

type QuizListResponse struct {
	Quizzes    []QuizSummary `json:"quizzes"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

type QuizSummary struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Difficulty    string    `json:"difficulty"`
	TimeLimit     int       `json:"time_limit"`
	QuestionCount int       `json:"question_count"`
	IsFeatured    bool      `json:"is_featured"`
	Tags          []string  `json:"tags"`
	ThumbnailURL  *string   `json:"thumbnail_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Statistics    *QuizStatisticsSummary `json:"statistics,omitempty"`
}

type QuizStatisticsSummary struct {
	TotalAttempts int     `json:"total_attempts"`
	AverageScore  float64 `json:"average_score"`
	AverageTime   int     `json:"average_time"`
}

type QuizDetailResponse struct {
	Quiz Quiz `json:"quiz"`
}

type QuizFilters struct {
	Category   string `form:"category"`
	Difficulty string `form:"difficulty" binding:"omitempty,oneof=Easy Medium Hard"`
	Featured   *bool  `form:"featured"`
	Tags       string `form:"tags"`
	Search     string `form:"search"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=10" binding:"min=1,max=100"`
}

type QuizAttempt struct {
	ID          uuid.UUID `json:"id" db:"id"`
	QuizID      uuid.UUID `json:"quiz_id" db:"quiz_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Score       float64   `json:"score" db:"score"`
	TotalPoints int       `json:"total_points" db:"total_points"`
	TimeSpent   int       `json:"time_spent" db:"time_spent"` // in seconds
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
	StartedAt   time.Time `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserStatistics represents comprehensive user statistics
type UserStatistics struct {
	UserID                 uuid.UUID                 `json:"user_id"`
	TotalAttempts          int                       `json:"total_attempts"`
	CompletedQuizzes       int                       `json:"completed_quizzes"`
	CompletionRate         float64                   `json:"completion_rate"`
	AverageScore           float64                   `json:"average_score"`
	TotalPoints            int                       `json:"total_points"`
	CurrentStreak          int                       `json:"current_streak"`
	BestStreak             int                       `json:"best_streak"`
	AverageCompletionTime  int                       `json:"average_completion_time"` // in seconds
	CategoryPerformance    []CategoryPerformance     `json:"category_performance"`
	RecentActivity         []RecentActivityItem      `json:"recent_activity"`
	LastUpdated            time.Time                 `json:"last_updated"`
	QuizzesByDifficulty    map[string]int            `json:"quizzes_by_difficulty"`
	ScoreDistribution      ScoreDistribution         `json:"score_distribution"`
	MonthlyProgress        []MonthlyProgressItem     `json:"monthly_progress"`
}

// CategoryPerformance represents performance metrics for a specific category
type CategoryPerformance struct {
	CategoryID       string  `json:"category_id"`
	CategoryName     string  `json:"category_name"`
	QuizzesCompleted int     `json:"quizzes_completed"`
	AverageScore     float64 `json:"average_score"`
	TotalAttempts    int     `json:"total_attempts"`
	BestScore        float64 `json:"best_score"`
	LastAttempt      *time.Time `json:"last_attempt,omitempty"`
}

// RecentActivityItem represents a recent quiz activity
type RecentActivityItem struct {
	QuizID       uuid.UUID  `json:"quiz_id"`
	QuizTitle    string     `json:"quiz_title"`
	Score        float64    `json:"score"`
	Category     string     `json:"category"`
	CompletedAt  time.Time  `json:"completed_at"`
	TimeSpent    int        `json:"time_spent"` // in seconds
	Difficulty   string     `json:"difficulty"`
}

// ScoreDistribution represents the distribution of scores
type ScoreDistribution struct {
	Range0to20   int `json:"range_0_to_20"`
	Range21to40  int `json:"range_21_to_40"`
	Range41to60  int `json:"range_41_to_60"`
	Range61to80  int `json:"range_61_to_80"`
	Range81to100 int `json:"range_81_to_100"`
}

// MonthlyProgressItem represents progress for a specific month
type MonthlyProgressItem struct {
	Month            string  `json:"month"`
	QuizzesCompleted int     `json:"quizzes_completed"`
	AverageScore     float64 `json:"average_score"`
	TotalPoints      int     `json:"total_points"`
}

// UserStatisticsResponse represents the API response for user statistics
type UserStatisticsResponse struct {
	Statistics UserStatistics `json:"statistics"`
	Message    string         `json:"message,omitempty"`
}

// AttemptFilters represents filters for user attempt history
type AttemptFilters struct {
	QuizID      *string    `form:"quiz_id"`
	Category    string     `form:"category"`
	Difficulty  string     `form:"difficulty" binding:"omitempty,oneof=Easy Medium Hard"`
	StartDate   *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate     *time.Time `form:"end_date" time_format:"2006-01-02"`
	MinScore    *float64   `form:"min_score" binding:"omitempty,min=0,max=100"`
	MaxScore    *float64   `form:"max_score" binding:"omitempty,min=0,max=100"`
	Page        int        `form:"page,default=1" binding:"min=1"`
	PageSize    int        `form:"page_size,default=10" binding:"min=1,max=100"`
	SortBy      string     `form:"sort_by,default=completed_at" binding:"omitempty,oneof=completed_at score time_spent"`
	SortOrder   string     `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// QuizAttemptWithDetails represents a quiz attempt with quiz details for API responses
type QuizAttemptWithDetails struct {
	QuizAttempt
	Quiz QuizSummary `json:"quiz"`
}

// AttemptHistoryResponse represents the API response for user attempt history
type AttemptHistoryResponse struct {
	Attempts   []QuizAttemptWithDetails `json:"attempts"`
	Total      int                      `json:"total"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPages int                      `json:"total_pages"`
}

// AttemptDetailResponse represents the API response for a single attempt detail
type AttemptDetailResponse struct {
	Attempt QuizAttemptWithDetails `json:"attempt"`
}

// Friends-related models

// FriendRequest represents a friend request
type FriendRequest struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	RequesterID uuid.UUID  `json:"requester_id" db:"requester_id"`
	RequestedID uuid.UUID  `json:"requested_id" db:"requested_id"`
	Status      string     `json:"status" db:"status"`
	Message     *string    `json:"message,omitempty" db:"message"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	RespondedAt *time.Time `json:"responded_at,omitempty" db:"responded_at"`
	Requester   *User      `json:"requester,omitempty"`
	Requested   *User      `json:"requested,omitempty"`
}

// Friendship represents an accepted friendship between two users
type Friendship struct {
	ID        uuid.UUID `json:"id" db:"id"`
	User1ID   uuid.UUID `json:"user1_id" db:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id" db:"user2_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	User1     *User     `json:"user1,omitempty"`
	User2     *User     `json:"user2,omitempty"`
}

// FriendNotification represents a notification related to friends
type FriendNotification struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	Type            string     `json:"type" db:"type"`
	Title           string     `json:"title" db:"title"`
	Message         *string    `json:"message,omitempty" db:"message"`
	RelatedUserID   *uuid.UUID `json:"related_user_id,omitempty" db:"related_user_id"`
	FriendRequestID *uuid.UUID `json:"friend_request_id,omitempty" db:"friend_request_id"`
	IsRead          bool       `json:"is_read" db:"is_read"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	ReadAt          *time.Time `json:"read_at,omitempty" db:"read_at"`
	RelatedUser     *User      `json:"related_user,omitempty"`
}

// Friend represents a user's friend with additional metadata
type Friend struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	AvatarURL             *string   `json:"avatar_url,omitempty"`
	Level                 string    `json:"level"`
	TotalPoints           int       `json:"total_points"`
	CurrentStreak         int       `json:"current_streak"`
	BestStreak            int       `json:"best_streak"`
	TotalQuizzesCompleted int       `json:"total_quizzes_completed"`
	AverageScore          float64   `json:"average_score"`
	IsOnline              bool      `json:"is_online"`
	LastActive            time.Time `json:"last_active"`
	FriendsSince          time.Time `json:"friends_since"`
}

// Friends DTOs

// SendFriendRequestRequest represents the request to send a friend request
type SendFriendRequestRequest struct {
	RequestedUserID uuid.UUID `json:"requested_user_id" binding:"required"`
	Message         *string   `json:"message,omitempty"`
}

// RespondToFriendRequestRequest represents the request to respond to a friend request
type RespondToFriendRequestRequest struct {
	Status string `json:"status" binding:"required,oneof=accepted rejected"`
}

// FriendRequestsResponse represents the response for friend requests
type FriendRequestsResponse struct {
	PendingRequests []FriendRequest `json:"pending_requests"`
	SentRequests    []FriendRequest `json:"sent_requests"`
	Total           int             `json:"total"`
}

// FriendsListResponse represents the response for friends list
type FriendsListResponse struct {
	Friends []Friend `json:"friends"`
	Total   int      `json:"total"`
}

// UserSearchResponse represents the response for user search
type UserSearchResponse struct {
	Users []UserSearchResult `json:"users"`
	Total int                `json:"total"`
}

// UserSearchResult represents a user in search results
type UserSearchResult struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	AvatarURL             *string   `json:"avatar_url,omitempty"`
	Level                 string    `json:"level"`
	TotalPoints           int       `json:"total_points"`
	TotalQuizzesCompleted int       `json:"total_quizzes_completed"`
	AverageScore          float64   `json:"average_score"`
	IsOnline              bool      `json:"is_online"`
	IsFriend              bool      `json:"is_friend"`
	HasPendingRequest     bool      `json:"has_pending_request"`
	RequestSentByMe       bool      `json:"request_sent_by_me"`
}

// FriendNotificationsResponse represents the response for friend notifications
type FriendNotificationsResponse struct {
	Notifications []FriendNotification `json:"notifications"`
	UnreadCount   int                  `json:"unread_count"`
	Total         int                  `json:"total"`
}