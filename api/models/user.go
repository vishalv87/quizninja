package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	Email                 string    `json:"email" db:"email"`
	PasswordHash          string    `json:"-" db:"password_hash"`
	Name                  string    `json:"name" db:"name"`
	Level                 string    `json:"level" db:"level"`
	TotalPoints           int       `json:"total_points" db:"total_points"`
	CurrentStreak         int       `json:"current_streak" db:"current_streak"`
	BestStreak            int       `json:"best_streak" db:"best_streak"`
	TotalQuizzesCompleted int       `json:"total_quizzes_completed" db:"total_quizzes_completed"`
	AverageScore          float64   `json:"average_score" db:"average_score"`
	IsOnline              bool      `json:"is_online" db:"is_online"`
	LastActive            time.Time `json:"last_active" db:"last_active"`
	AvatarURL             *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
	IsTestData            bool      `json:"is_test_data" db:"is_test_data"`

	// Supabase auth integration fields
	AuthMethod     string     `json:"auth_method" db:"auth_method"`           // "supabase" or "jwt"
	SupabaseID     *string    `json:"supabase_id,omitempty" db:"supabase_id"` // Supabase user ID
	LastAuthMethod string     `json:"last_auth_method" db:"last_auth_method"` // Track last successful auth
	MigratedAt     *time.Time `json:"migrated_at,omitempty" db:"migrated_at"` // When user was migrated between auth systems

	Preferences *UserPreferences `json:"preferences,omitempty"`
}

type UserPreferences struct {
	ID                    uuid.UUID              `json:"id" db:"id"`
	UserID                uuid.UUID              `json:"user_id" db:"user_id"`
	SelectedCategories    StringArray            `json:"selected_categories" db:"selected_categories"`
	DifficultyPreference  string                 `json:"difficulty_preference" db:"difficulty_preference"`
	NotificationsEnabled  bool                   `json:"notifications_enabled" db:"notifications_enabled"`
	NotificationFrequency string                 `json:"notification_frequency" db:"notification_frequency"`
	ProfileVisibility     bool                   `json:"profile_visibility" db:"profile_visibility"`
	ShowOnlineStatus      bool                   `json:"show_online_status" db:"show_online_status"`
	AllowFriendRequests   bool                   `json:"allow_friend_requests" db:"allow_friend_requests"`
	ShareActivityStatus   bool                   `json:"share_activity_status" db:"share_activity_status"`
	NotificationTypes     map[string]interface{} `json:"notification_types" db:"notification_types"`
	OnboardingCompletedAt *time.Time             `json:"onboarding_completed_at,omitempty" db:"onboarding_completed_at"`
	CreatedAt             time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at" db:"updated_at"`
	IsTestData            bool                   `json:"is_test_data" db:"is_test_data"`
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

// RegisterRequest for frontend-initiated Supabase auth
type RegisterRequest struct {
	// Supabase user data (received after frontend auth)
	SupabaseUserID string                  `json:"supabase_user_id" binding:"required"`
	Email          string                  `json:"email" binding:"required,email"`
	Name           string                  `json:"name" binding:"required"`
	Preferences    *UserPreferencesRequest `json:"preferences,omitempty"`

	// Optional: additional user metadata
	AvatarURL *string `json:"avatar_url,omitempty"`
}

type UserPreferencesRequest struct {
	SelectedCategories    []string `json:"selected_categories"`
	DifficultyPreference  string   `json:"difficulty_preference" binding:"omitempty,oneof=Easy Medium Hard"`
	NotificationsEnabled  bool     `json:"notifications_enabled"`
	NotificationFrequency string   `json:"notification_frequency" binding:"omitempty,oneof=Daily Weekly Never"`
}

// LoginRequest for frontend-initiated Supabase auth (user sync)
type LoginRequest struct {
	// Supabase user data (received after frontend auth)
	SupabaseUserID string  `json:"supabase_user_id" binding:"required"`
	Email          string  `json:"email" binding:"required,email"`
	Name           string  `json:"name" binding:"required"`
	AvatarURL      *string `json:"avatar_url,omitempty"`
}

// AuthResponse for frontend-initiated auth (returns user profile only)
type AuthResponse struct {
	User    User   `json:"user"`
	Message string `json:"message,omitempty"`
}

type UpdateProfileRequest struct {
	Name      *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// Category represents quiz topic categories
type Category struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Description string    `json:"description" db:"description"`
	IconURL     string    `json:"icon_url" db:"icon_url"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	QuizCount   int       `json:"quiz_count" db:"quiz_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsTestData  bool      `json:"is_test_data" db:"is_test_data"`
}

// DifficultyLevel represents quiz difficulty levels
type DifficultyLevel struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"display_name" db:"display_name"`
	Description string `json:"description" db:"description"`
	Order       int    `json:"order" db:"order"`
	IsActive    bool   `json:"is_active" db:"is_active"`
	IsTestData  bool   `json:"is_test_data" db:"is_test_data"`
}

// NotificationFrequency represents notification frequency options
type NotificationFrequency struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"display_name" db:"display_name"`
	Description string `json:"description" db:"description"`
	IsActive    bool   `json:"is_active" db:"is_active"`
	IsTestData  bool   `json:"is_test_data" db:"is_test_data"`
}

// AppSettings represents application configuration settings
type AppSettings struct {
	ID             uuid.UUID `json:"id" db:"id"`
	Key            string    `json:"key" db:"key"`
	Value          string    `json:"value" db:"value"`
	Description    string    `json:"description" db:"description"`
	Category       string    `json:"category" db:"category"`
	IsPublic       bool      `json:"is_public" db:"is_public"`
	LastModifiedBy uuid.UUID `json:"last_modified_by" db:"last_modified_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// CategoryGroup represents quiz category groups for API response
type CategoryGroup struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	DisplayName string     `json:"display_name"`
	Description string     `json:"description"`
	IconURL     string     `json:"icon_url"`
	IsActive    bool       `json:"is_active"`
	Categories  []Category `json:"categories,omitempty"`
}

// API Request/Response types
type UpdatePreferencesRequest struct {
	SelectedCategories    []string               `json:"selected_categories" binding:"dive,required"`
	DifficultyPreference  string                 `json:"difficulty_preference" binding:"required,oneof=Easy Medium Hard"`
	NotificationsEnabled  bool                   `json:"notifications_enabled"`
	NotificationFrequency string                 `json:"notification_frequency" binding:"required,oneof=Daily Weekly Never"`
	ProfileVisibility     *bool                  `json:"profile_visibility,omitempty"`
	ShowOnlineStatus      *bool                  `json:"show_online_status,omitempty"`
	AllowFriendRequests   *bool                  `json:"allow_friend_requests,omitempty"`
	ShareActivityStatus   *bool                  `json:"share_activity_status,omitempty"`
	NotificationTypes     map[string]interface{} `json:"notification_types,omitempty"`
}

type AppSettingsResponse struct {
	Settings    map[string]any `json:"settings"`
	Version     string         `json:"version"`
	LastUpdated time.Time      `json:"last_updated"`
}

// Quiz models

type Quiz struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	Title         string          `json:"title" db:"title"`
	Description   string          `json:"description" db:"description"`
	Category      string          `json:"category" db:"category_id"`
	Difficulty    string          `json:"difficulty" db:"difficulty"`
	TimeLimit     int             `json:"time_limit" db:"time_limit_minutes"` // in minutes
	QuestionCount int             `json:"question_count" db:"total_questions"`
	IsFeatured    bool            `json:"is_featured" db:"is_featured"`
	IsPublic      bool            `json:"is_public" db:"is_public"`
	CreatedBy     uuid.UUID       `json:"created_by" db:"created_by"`
	Tags          StringArray     `json:"tags" db:"tags"`
	ThumbnailURL  *string         `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
	IsTestData    bool            `json:"is_test_data" db:"is_test_data"`
	Questions     []Question      `json:"questions,omitempty"`
	Statistics    *QuizStatistics `json:"statistics,omitempty"`
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
	ID            uuid.UUID   `json:"id" db:"id"`
	QuizID        uuid.UUID   `json:"quiz_id" db:"quiz_id"`
	QuestionText  string      `json:"question_text" db:"question_text"`
	QuestionType  string      `json:"question_type" db:"question_type"` // multiple_choice, true_false, short_answer
	Options       StringArray `json:"options,omitempty" db:"options"`
	CorrectAnswer string      `json:"correct_answer" db:"correct_answer"`
	Explanation   *string     `json:"explanation,omitempty" db:"explanation"`
	Points        int         `json:"points" db:"points"`
	Order         int         `json:"order" db:"order"`
	ImageURL      *string     `json:"image_url,omitempty" db:"image_url"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	IsTestData    bool        `json:"is_test_data" db:"is_test_data"`
}

type QuizStatistics struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	QuizID            uuid.UUID  `json:"quiz_id" db:"quiz_id"`
	TotalAttempts     int        `json:"total_attempts" db:"total_attempts"`
	CompletedAttempts int        `json:"completed_attempts" db:"completed_attempts"`
	AverageScore      float64    `json:"average_score" db:"average_score"`
	AverageTime       int        `json:"average_time" db:"average_time"` // in seconds
	HighestScore      float64    `json:"highest_score" db:"highest_score"`
	LowestScore       float64    `json:"lowest_score" db:"lowest_score"`
	PopularityScore   int        `json:"popularity_score" db:"popularity_score"`
	LastAttemptAt     *time.Time `json:"last_attempt_at,omitempty" db:"last_attempt_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	IsTestData        bool       `json:"is_test_data" db:"is_test_data"`
}

// Quiz DTOs

type CreateQuizRequest struct {
	Title        string                  `json:"title" binding:"required,min=1,max=200"`
	Description  string                  `json:"description" binding:"required,min=1,max=1000"`
	Category     string                  `json:"category" binding:"required"`
	Difficulty   string                  `json:"difficulty" binding:"required,oneof=Easy Medium Hard"`
	TimeLimit    int                     `json:"time_limit" binding:"required,min=30,max=3600"`
	IsFeatured   bool                    `json:"is_featured"`
	IsPublic     bool                    `json:"is_public"`
	Tags         []string                `json:"tags"`
	ThumbnailURL *string                 `json:"thumbnail_url,omitempty"`
	Questions    []CreateQuestionRequest `json:"questions" binding:"required,dive"`
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
	Title        *string  `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Description  *string  `json:"description,omitempty" binding:"omitempty,min=1,max=1000"`
	Category     *string  `json:"category,omitempty"`
	Difficulty   *string  `json:"difficulty,omitempty" binding:"omitempty,oneof=Easy Medium Hard"`
	TimeLimit    *int     `json:"time_limit,omitempty" binding:"omitempty,min=30,max=3600"`
	IsFeatured   *bool    `json:"is_featured,omitempty"`
	IsPublic     *bool    `json:"is_public,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	ThumbnailURL *string  `json:"thumbnail_url,omitempty"`
}

type QuizListResponse struct {
	Quizzes    []QuizSummary `json:"quizzes"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

type QuizSummary struct {
	ID            uuid.UUID              `json:"id"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Category      string                 `json:"category"`
	Difficulty    string                 `json:"difficulty"`
	TimeLimit     int                    `json:"time_limit"`
	QuestionCount int                    `json:"question_count"`
	IsFeatured    bool                   `json:"is_featured"`
	Tags          []string               `json:"tags"`
	ThumbnailURL  *string                `json:"thumbnail_url,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	IsTestData    bool                   `json:"is_test_data" db:"is_test_data"`
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
	Category   string `form:"category"`   // Comma-separated category IDs for multiple category filtering
	Difficulty string `form:"difficulty" binding:"omitempty,oneof=beginner intermediate advanced"`
	Featured   *bool  `form:"featured"`
	Tags       string `form:"tags"`
	Search     string `form:"search"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=10" binding:"min=1,max=100"`
}

type QuizAttempt struct {
	ID                    uuid.UUID              `json:"id" db:"id"`
	QuizID                uuid.UUID              `json:"quiz_id" db:"quiz_id"`
	UserID                uuid.UUID              `json:"user_id" db:"user_id"`
	Answers               []AttemptAnswer        `json:"answers" db:"answers"`
	Score                 float64                `json:"score" db:"score"`
	TotalPoints           int                    `json:"total_points" db:"total_points"`
	TimeSpent             int                    `json:"time_spent" db:"time_spent"` // in seconds
	PercentageScore       float64                `json:"percentage_score" db:"percentage_score"`
	Passed                bool                   `json:"passed" db:"passed"`
	Status                string                 `json:"status" db:"status"` // started, completed, abandoned
	IsCompleted           bool                   `json:"is_completed" db:"is_completed"`
	StartedAt             time.Time              `json:"started_at" db:"started_at"`
	CompletedAt           *time.Time             `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt             time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at" db:"updated_at"`
	RetakeCount           int                    `json:"retake_count" db:"retake_count"`
	OriginalAttemptID     *uuid.UUID             `json:"original_attempt_id,omitempty" db:"original_attempt_id"`
	PerformanceComparison map[string]interface{} `json:"performance_comparison,omitempty" db:"performance_comparison"`
	IsTestData            bool                   `json:"is_test_data" db:"is_test_data"`

	// Challenge tracking fields
	ChallengeID         *uuid.UUID `json:"challenge_id,omitempty" db:"challenge_id"`
	IsChallengeAttempt  bool       `json:"is_challenge_attempt" db:"is_challenge_attempt"`
}

// AttemptAnswer represents a user's answer to a quiz question
type AttemptAnswer struct {
	QuestionID     uuid.UUID `json:"question_id" validate:"required"`
	SelectedOption int       `json:"selected_option"` // For multiple choice
	TextAnswer     string    `json:"text_answer"`     // For text answers
	IsCorrect      bool      `json:"is_correct"`
	PointsEarned   int       `json:"points_earned"`
}

// IsRetake returns true if this attempt is a retake
func (qa *QuizAttempt) IsRetake() bool {
	return qa.OriginalAttemptID != nil
}

// CanRetake returns true if this attempt can be retaken (max 3 retakes)
func (qa *QuizAttempt) CanRetake() bool {
	return qa.RetakeCount < 3 && qa.IsCompleted
}

// GetRetakeLabel returns a label for retake attempts
func (qa *QuizAttempt) GetRetakeLabel() *string {
	if qa.RetakeCount == 0 {
		return nil
	}
	label := fmt.Sprintf("Retake #%d", qa.RetakeCount)
	return &label
}

// CreateRetakeRequest represents the request body for creating a retake attempt
type CreateRetakeRequest struct {
	OriginalAttemptID string `json:"original_attempt_id" validate:"required,uuid"`
	IsRetake          bool   `json:"is_retake"`
}

// UpdateAttemptRequest represents the request body for updating/completing a quiz attempt
type UpdateAttemptRequest struct {
	Answers   []AttemptAnswer `json:"answers" validate:"required"`
	TimeSpent int             `json:"time_spent" validate:"min=1"`
	Status    string          `json:"status" validate:"required,oneof=completed abandoned"`
}

// UserStatistics represents comprehensive user statistics
type UserStatistics struct {
	UserID                uuid.UUID             `json:"user_id"`
	TotalAttempts         int                   `json:"total_attempts"`
	CompletedQuizzes      int                   `json:"total_quizzes_completed"`
	CompletionRate        float64               `json:"completion_rate"`
	AverageScore          float64               `json:"average_score"`
	TotalPoints           int                   `json:"total_points"`
	CurrentStreak         int                   `json:"current_streak"`
	BestStreak            int                   `json:"best_streak"`
	AverageCompletionTime int                   `json:"average_completion_time"` // in seconds
	CategoryPerformance   []CategoryPerformance `json:"category_performance"`
	RecentActivity        []RecentActivityItem  `json:"recent_activity"`
	LastUpdated           time.Time             `json:"last_updated"`
	QuizzesByDifficulty   map[string]int        `json:"quizzes_by_difficulty"`
	ScoreDistribution     ScoreDistribution     `json:"score_distribution"`
	MonthlyProgress       []MonthlyProgressItem `json:"monthly_progress"`
	IsTestData            bool                  `json:"is_test_data"`
}

// CategoryPerformance represents performance metrics for a specific category
type CategoryPerformance struct {
	CategoryID       string     `json:"category_id"`
	CategoryName     string     `json:"category_name"`
	QuizzesCompleted int        `json:"quizzes_completed"`
	AverageScore     float64    `json:"average_score"`
	TotalAttempts    int        `json:"total_attempts"`
	BestScore        float64    `json:"best_score"`
	LastAttempt      *time.Time `json:"last_attempt,omitempty"`
}

// RecentActivityItem represents a recent quiz activity
type RecentActivityItem struct {
	QuizID      uuid.UUID `json:"quiz_id"`
	QuizTitle   string    `json:"quiz_title"`
	Score       float64   `json:"score"`
	Category    string    `json:"category"`
	CompletedAt time.Time `json:"completed_at"`
	TimeSpent   int       `json:"time_spent"` // in seconds
	Difficulty  string    `json:"difficulty"`
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
	QuizID     *string    `form:"quiz_id"`
	Category   string     `form:"category"`
	Difficulty string     `form:"difficulty" binding:"omitempty,oneof=Easy Medium Hard"`
	StartDate  *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate    *time.Time `form:"end_date" time_format:"2006-01-02"`
	MinScore   *float64   `form:"min_score" binding:"omitempty,min=0,max=100"`
	MaxScore   *float64   `form:"max_score" binding:"omitempty,min=0,max=100"`
	Page       int        `form:"page,default=1" binding:"min=1"`
	PageSize   int        `form:"page_size,default=10" binding:"min=1,max=100"`
	SortBy     string     `form:"sort_by,default=completed_at" binding:"omitempty,oneof=completed_at score time_spent"`
	SortOrder  string     `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
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
	IsTestData  bool       `json:"is_test_data" db:"is_test_data"`
	Requester   *User      `json:"requester,omitempty"`
	Requested   *User      `json:"requested,omitempty"`
}

// Friendship represents an accepted friendship between two users
type Friendship struct {
	ID         uuid.UUID `json:"id" db:"id"`
	User1ID    uuid.UUID `json:"user1_id" db:"user1_id"`
	User2ID    uuid.UUID `json:"user2_id" db:"user2_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	IsTestData bool      `json:"is_test_data" db:"is_test_data"`
	User1      *User     `json:"user1,omitempty"`
	User2      *User     `json:"user2,omitempty"`
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
	IsTestData      bool       `json:"is_test_data" db:"is_test_data"`
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
	IsTestData            bool      `json:"is_test_data"`
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
	IsTestData            bool      `json:"is_test_data"`
}

// FriendNotificationsResponse represents the response for friend notifications
type FriendNotificationsResponse struct {
	Notifications []FriendNotification `json:"notifications"`
	UnreadCount   int                  `json:"unread_count"`
	Total         int                  `json:"total"`
}

// Challenge-related models

// Challenge represents a quiz challenge between users
type Challenge struct {
	ID                uuid.UUID              `json:"id" db:"id"`
	ChallengerID      uuid.UUID              `json:"challenger_id" db:"challenger_id"`
	ChallengeeID      uuid.UUID              `json:"challengee_id" db:"challengee_id"`
	QuizID            uuid.UUID              `json:"quiz_id" db:"quiz_id"`
	Status            string                 `json:"status" db:"status"`
	ChallengerScore   *float64               `json:"challenger_score,omitempty" db:"challenger_score"`
	ChallengeeScore   *float64               `json:"challengee_score,omitempty" db:"challengee_score"`
	Message           *string                `json:"message,omitempty" db:"message"`
	ExpiresAt         *time.Time             `json:"expires_at,omitempty" db:"expires_at"`
	IsGroupChallenge  bool                   `json:"is_group_challenge" db:"is_group_challenge"`
	ParticipantIDs    []uuid.UUID            `json:"participant_ids,omitempty" db:"participant_ids"`
	ParticipantScores map[string]interface{} `json:"participant_scores,omitempty" db:"participant_scores"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`
	IsTestData        bool                   `json:"is_test_data" db:"is_test_data"`

	// Asynchronous challenge tracking fields
	ChallengerAttemptID   *uuid.UUID `json:"challenger_attempt_id,omitempty" db:"challenger_attempt_id"`
	ChallengeeAttemptID   *uuid.UUID `json:"challengee_attempt_id,omitempty" db:"challengee_attempt_id"`
	ChallengerCompletedAt *time.Time `json:"challenger_completed_at,omitempty" db:"challenger_completed_at"`
	ChallengeeCompletedAt *time.Time `json:"challengee_completed_at,omitempty" db:"challengee_completed_at"`

	// Relationships
	Challenger        *User                  `json:"challenger,omitempty"`
	Challengee        *User                  `json:"challengee,omitempty"`
	Quiz              *QuizSummary           `json:"quiz,omitempty"`
}

// ChallengeWithDetails represents a challenge with full user and quiz details
type ChallengeWithDetails struct {
	Challenge
	ChallengerName   string `json:"challenger_name"`
	ChallengerAvatar string `json:"challenger_avatar"`
	ChallengeeName   string `json:"challengee_name"`
	ChallengeeAvatar string `json:"challengee_avatar"`
	QuizTitle        string `json:"quiz_title"`
	QuizCategory     string `json:"quiz_category"`
}

// Challenge DTOs

// CreateChallengeRequest represents the request to create a new challenge
type CreateChallengeRequest struct {
	ChallengeeUserID uuid.UUID   `json:"challengee_user_id" binding:"required"`
	QuizID           uuid.UUID   `json:"quiz_id" binding:"required"`
	Message          *string     `json:"message,omitempty"`
	ExpiresAt        *time.Time  `json:"expires_at,omitempty"`
	IsGroupChallenge bool        `json:"is_group_challenge"`
	ParticipantIDs   []uuid.UUID `json:"participant_ids,omitempty"`
}

// AcceptChallengeRequest represents the request to accept a challenge
type AcceptChallengeRequest struct {
	// No additional fields needed - just the challenge ID from URL
}

// UpdateChallengeScoreRequest represents the request to update challenge scores
type UpdateChallengeScoreRequest struct {
	UserScore float64 `json:"user_score" binding:"required,min=0,max=100"`
}

// ChallengeFilters represents filters for challenge queries
type ChallengeFilters struct {
	Status    string     `form:"status" binding:"omitempty,oneof=pending accepted completed expired declined"`
	QuizID    *uuid.UUID `form:"quiz_id"`
	UserType  string     `form:"user_type" binding:"omitempty,oneof=challenger challenged all"`
	StartDate *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time `form:"end_date" time_format:"2006-01-02"`
	Page      int        `form:"page,default=1" binding:"min=1"`
	PageSize  int        `form:"page_size,default=10" binding:"min=1,max=100"`
	SortBy    string     `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at expires_at status"`
	SortOrder string     `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// ChallengeListResponse represents the response for challenge list
type ChallengeListResponse struct {
	Challenges []ChallengeWithDetails `json:"challenges"`
	Total      int                    `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

// ChallengeDetailResponse represents the response for a single challenge detail
type ChallengeDetailResponse struct {
	Challenge ChallengeWithDetails `json:"challenge"`
}

// ChallengeStatsResponse represents challenge statistics for a user
type ChallengeStatsResponse struct {
	TotalChallenges     int     `json:"total_challenges"`
	PendingChallenges   int     `json:"pending_challenges"`
	ActiveChallenges    int     `json:"active_challenges"`
	CompletedChallenges int     `json:"completed_challenges"`
	WonChallenges       int     `json:"won_challenges"`
	LostChallenges      int     `json:"lost_challenges"`
	WinRate             float64 `json:"win_rate"`
	AverageScore        float64 `json:"average_score"`
}

// Leaderboard-related models

// LeaderboardEntry represents a user's position in the leaderboard
type LeaderboardEntry struct {
	UserID           uuid.UUID      `json:"user_id" db:"user_id"`
	Name             string         `json:"name" db:"name"`
	Avatar           *string        `json:"avatar" db:"avatar_url"`
	Rank             int            `json:"rank"`
	Points           int            `json:"points" db:"total_points"`
	QuizzesCompleted int            `json:"quizzes_completed" db:"total_quizzes_completed"`
	AverageScore     float64        `json:"average_score" db:"average_score"`
	CurrentStreak    int            `json:"current_streak" db:"current_streak"`
	Level            string         `json:"level" db:"level"`
	IsCurrentUser    bool           `json:"is_current_user"`
	IsFriend         bool           `json:"is_friend"`
	LastActive       time.Time      `json:"last_active" db:"last_active"`
	Achievements     []string       `json:"achievements"`
	CategoryPoints   map[string]int `json:"category_points"`
}

// LeaderboardFilters represents filters for leaderboard queries
type LeaderboardFilters struct {
	Period      string `form:"period,default=alltime" binding:"omitempty,oneof=today week month alltime"`
	FriendsOnly bool   `form:"friends_only,default=false"`
	Limit       int    `form:"limit,default=50" binding:"min=1,max=100"`
	Offset      int    `form:"offset,default=0" binding:"min=0"`
}

// LeaderboardResponse represents the API response for leaderboard
type LeaderboardResponse struct {
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	UserRank    *UserRankInfo      `json:"user_rank,omitempty"`
	Total       int                `json:"total"`
	Period      string             `json:"period"`
	FriendsOnly bool               `json:"friends_only"`
}

// UserRankInfo represents the current user's rank information
type UserRankInfo struct {
	UserID       uuid.UUID `json:"user_id"`
	Rank         int       `json:"rank"`
	Points       int       `json:"points"`
	PointsToNext int       `json:"points_to_next"`
	RankChange   int       `json:"rank_change"` // +/- change from previous period
}

// Achievement-related models

// Achievement represents an achievement that users can unlock
type Achievement struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Key          string    `json:"key" db:"key"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	Icon         *string   `json:"icon,omitempty" db:"icon"`
	Color        string    `json:"color" db:"color"`
	PointsReward int       `json:"points_reward" db:"points_reward"`
	Category     string    `json:"category" db:"category"`
	IsRare       bool      `json:"is_rare" db:"is_rare"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	IsTestData   bool      `json:"is_test_data" db:"is_test_data"`
}

// UserAchievement represents a user's unlocked achievement
type UserAchievement struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	UserID        uuid.UUID    `json:"user_id" db:"user_id"`
	AchievementID uuid.UUID    `json:"achievement_id" db:"achievement_id"`
	UnlockedAt    time.Time    `json:"unlocked_at" db:"unlocked_at"`
	PointsAwarded int          `json:"points_awarded" db:"points_awarded"`
	IsTestData    bool         `json:"is_test_data" db:"is_test_data"`
	Achievement   *Achievement `json:"achievement,omitempty"`
}

// Achievement DTOs

// AchievementListResponse represents the response for user achievements
type AchievementListResponse struct {
	Achievements []UserAchievement `json:"achievements"`
	Total        int               `json:"total"`
}

// AchievementNotification represents a notification for a newly unlocked achievement
type AchievementNotification struct {
	AchievementID uuid.UUID `json:"achievement_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Icon          *string   `json:"icon,omitempty"`
	Color         string    `json:"color"`
	PointsAwarded int       `json:"points_awarded"`
	IsRare        bool      `json:"is_rare"`
}

// AchievementProgress represents progress towards an achievement
type AchievementProgress struct {
	AchievementID uuid.UUID `json:"achievement_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Icon          *string   `json:"icon,omitempty"`
	Color         string    `json:"color"`
	Category      string    `json:"category"`
	IsRare        bool      `json:"is_rare"`
	CurrentValue  int       `json:"current_value"`
	TargetValue   int       `json:"target_value"`
	Progress      float64   `json:"progress"` // Percentage (0-100)
	IsUnlocked    bool      `json:"is_unlocked"`
}

// AchievementProgressResponse represents the response for achievement progress
type AchievementProgressResponse struct {
	Progress []AchievementProgress `json:"progress"`
	Total    int                   `json:"total"`
}

// Favorites-related models

// UserQuizFavorite represents a user's favorite quiz
type UserQuizFavorite struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	UserID      uuid.UUID    `json:"user_id" db:"user_id"`
	QuizID      uuid.UUID    `json:"quiz_id" db:"quiz_id"`
	FavoritedAt time.Time    `json:"favorited_at" db:"favorited_at"`
	IsTestData  bool         `json:"is_test_data" db:"is_test_data"`
	Quiz        *QuizSummary `json:"quiz,omitempty"`
}

// FavoritesListResponse represents the response for user favorites
type FavoritesListResponse struct {
	Favorites  []UserQuizFavorite `json:"favorites"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// AddFavoriteRequest represents the request to add a quiz to favorites
type AddFavoriteRequest struct {
	QuizID uuid.UUID `json:"quiz_id" binding:"required"`
}

// Quiz Session-related models

// QuizSession represents an active or paused quiz session for continuation functionality
type QuizSession struct {
	ID                   uuid.UUID       `json:"id" db:"id"`
	AttemptID            uuid.UUID       `json:"attempt_id" db:"attempt_id"`
	UserID               uuid.UUID       `json:"user_id" db:"user_id"`
	QuizID               uuid.UUID       `json:"quiz_id" db:"quiz_id"`
	CurrentQuestionIndex int             `json:"current_question_index" db:"current_question_index"`
	CurrentAnswers       []AttemptAnswer `json:"current_answers" db:"current_answers"`
	SessionState         string          `json:"session_state" db:"session_state"`             // active, paused, completed, abandoned
	TimeRemaining        *int            `json:"time_remaining,omitempty" db:"time_remaining"` // seconds
	TimeSpentSoFar       int             `json:"time_spent_so_far" db:"time_spent_so_far"`     // seconds
	LastActivityAt       time.Time       `json:"last_activity_at" db:"last_activity_at"`
	PausedAt             *time.Time      `json:"paused_at,omitempty" db:"paused_at"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
	IsTestData           bool            `json:"is_test_data" db:"is_test_data"`
	Quiz                 *QuizSummary    `json:"quiz,omitempty"`
	Attempt              *QuizAttempt    `json:"attempt,omitempty"`
}

// IsActive returns true if the session is currently active
func (qs *QuizSession) IsActive() bool {
	return qs.SessionState == "active"
}

// IsPaused returns true if the session is currently paused
func (qs *QuizSession) IsPaused() bool {
	return qs.SessionState == "paused"
}

// CanResume returns true if the session can be resumed
func (qs *QuizSession) CanResume() bool {
	return qs.SessionState == "paused" || qs.SessionState == "active"
}

// GetProgress returns the completion progress as a percentage (0-100)
func (qs *QuizSession) GetProgress() float64 {
	if qs.Quiz == nil || qs.Quiz.QuestionCount == 0 {
		return 0.0
	}
	return (float64(qs.CurrentQuestionIndex) / float64(qs.Quiz.QuestionCount)) * 100
}

// GetQuestionsRemaining returns the number of questions left to complete
func (qs *QuizSession) GetQuestionsRemaining() int {
	if qs.Quiz == nil {
		return 0
	}
	return qs.Quiz.QuestionCount - qs.CurrentQuestionIndex
}

// GetTimeRemainingFormatted returns the remaining time in a human-readable format
func (qs *QuizSession) GetTimeRemainingFormatted() string {
	if qs.TimeRemaining == nil {
		return "No time limit"
	}

	minutes := *qs.TimeRemaining / 60
	seconds := *qs.TimeRemaining % 60

	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// HasExpired returns true if the session has been inactive for too long
func (qs *QuizSession) HasExpired() bool {
	// Sessions no longer expire automatically - users can resume anytime
	return false
}

// QuizSessionWithDetails represents a quiz session with full quiz and attempt details
type QuizSessionWithDetails struct {
	QuizSession
	QuizTitle         string  `json:"quiz_title"`
	QuizCategory      string  `json:"quiz_category"`
	QuizDifficulty    string  `json:"quiz_difficulty"`
	QuizThumbnail     *string `json:"quiz_thumbnail,omitempty"`
	TotalQuestions    int     `json:"total_questions"`
	OriginalTimeLimit int     `json:"original_time_limit"` // in seconds
	QuizIsTestData    bool    `json:"quiz_is_test_data" db:"quiz_is_test_data"`
}

// Quiz Session DTOs

// CreateQuizSessionRequest represents the request to create a new quiz session
type CreateQuizSessionRequest struct {
	AttemptID            uuid.UUID `json:"attempt_id" binding:"required"`
	CurrentQuestionIndex int       `json:"current_question_index" binding:"min=0"`
	TimeRemaining        *int      `json:"time_remaining,omitempty" binding:"omitempty,min=0"`
}

// UpdateQuizSessionRequest represents the request to update a quiz session
type UpdateQuizSessionRequest struct {
	CurrentQuestionIndex int             `json:"current_question_index" binding:"min=0"`
	CurrentAnswers       []AttemptAnswer `json:"current_answers"`
	TimeSpentSoFar       int             `json:"time_spent_so_far" binding:"min=0"`
	TimeRemaining        *int            `json:"time_remaining,omitempty" binding:"omitempty,min=0"`
}

// PauseSessionRequest represents the request to pause a quiz session
type PauseSessionRequest struct {
	CurrentQuestionIndex int             `json:"current_question_index" binding:"min=0"`
	CurrentAnswers       []AttemptAnswer `json:"current_answers"`
	TimeSpentSoFar       int             `json:"time_spent_so_far" binding:"min=0"`
	TimeRemaining        *int            `json:"time_remaining,omitempty" binding:"omitempty,min=0"`
}

// ResumeSessionRequest represents the request to resume a quiz session
type ResumeSessionRequest struct {
	// No additional fields needed - session state is loaded from database
}

// ActiveSessionsResponse represents the response for user's active/paused sessions
type ActiveSessionsResponse struct {
	Sessions    []QuizSessionWithDetails `json:"sessions"`
	Total       int                      `json:"total"`
	ActiveCount int                      `json:"active_count"`
	PausedCount int                      `json:"paused_count"`
}

// SessionFilters represents filters for quiz session queries
type SessionFilters struct {
	SessionState string     `form:"session_state" binding:"omitempty,oneof=active paused completed abandoned"`
	QuizID       *uuid.UUID `form:"quiz_id"`
	Category     string     `form:"category"`
	Difficulty   string     `form:"difficulty" binding:"omitempty,oneof=Easy Medium Hard"`
	StartDate    *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate      *time.Time `form:"end_date" time_format:"2006-01-02"`
	Page         int        `form:"page,default=1" binding:"min=1"`
	PageSize     int        `form:"page_size,default=10" binding:"min=1,max=100"`
	SortBy       string     `form:"sort_by,default=last_activity_at" binding:"omitempty,oneof=last_activity_at created_at"`
	SortOrder    string     `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}

// SessionActionResponse represents the response for session actions (pause/resume/abandon)
type SessionActionResponse struct {
	SessionID     uuid.UUID `json:"session_id"`
	Action        string    `json:"action"` // paused, resumed, abandoned
	SessionState  string    `json:"session_state"`
	Message       string    `json:"message"`
	TimeRemaining *int      `json:"time_remaining,omitempty"`
	Progress      float64   `json:"progress"`
}

// Onboarding-related models

// OnboardingCompleteRequest represents the request to complete onboarding
type OnboardingCompleteRequest struct {
	SelectedCategories    []string `json:"selected_categories" binding:"required,dive,required"`
	DifficultyPreference  string   `json:"difficulty_preference" binding:"required,oneof=Easy Medium Hard"`
	NotificationsEnabled  bool     `json:"notifications_enabled"`
	NotificationFrequency string   `json:"notification_frequency" binding:"required,oneof=Daily Weekly Never"`
}

// OnboardingStatusResponse represents the response for onboarding status
type OnboardingStatusResponse struct {
	IsCompleted bool             `json:"is_completed"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
	Preferences *UserPreferences `json:"preferences,omitempty"`
}
