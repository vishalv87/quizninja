package repository

import (
	"quizninja-api/models"

	"github.com/google/uuid"
)

// UserRepositoryInterface defines the contract for user data operations
type UserRepositoryInterface interface {
	// User CRUD operations
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uuid.UUID) error

	// User preferences operations
	CreateUserPreferences(preferences *models.UserPreferences) error
	GetUserPreferences(userID uuid.UUID) (*models.UserPreferences, error)
	UpdateUserPreferences(preferences *models.UserPreferences) error
	DeleteUserPreferences(userID uuid.UUID) error

	// User with preferences operations
	GetUserWithPreferences(userID uuid.UUID) (*models.User, error)


	// User status operations
	UpdateUserOnlineStatus(userID uuid.UUID, isOnline bool) error
	UpdateUserLastActive(userID uuid.UUID) error

	// User statistics operations
	GetUserStatistics(userID uuid.UUID) (*models.UserStatistics, error)
	UpdateUserStatistics(userID uuid.UUID, newScore float64) error
}

// QuizRepositoryInterface defines the contract for quiz data operations
type QuizRepositoryInterface interface {
	// Quiz read operations
	GetQuizByID(id uuid.UUID) (*models.Quiz, error)
	GetQuizByIDWithQuestions(id uuid.UUID) (*models.Quiz, error)
	GetQuizByIDWithStatistics(id uuid.UUID) (*models.Quiz, error)
	GetQuizByIDWithAll(id uuid.UUID) (*models.Quiz, error)

	// Quiz list operations with filtering and pagination
	GetQuizzes(filters *models.QuizFilters) ([]models.Quiz, int, error)
	GetFeaturedQuizzes(limit int) ([]models.Quiz, error)
	GetQuizzesByCategory(category string, limit int) ([]models.Quiz, error)
	GetQuizzesByUser(userID uuid.UUID, offset, limit int) ([]models.Quiz, int, error)

	// Question read operations
	GetQuestionsByQuizID(quizID uuid.UUID) ([]models.Question, error)

	// Quiz statistics read operations
	GetQuizStatistics(quizID uuid.UUID) (*models.QuizStatistics, error)
	CreateOrUpdateQuizStatistics(quizID uuid.UUID, score float64, timeSpent int) error

	// Quiz attempt operations
	CreateQuizAttempt(attempt *models.QuizAttempt) error
	UpdateQuizAttempt(attempt *models.QuizAttempt) error
	GetQuizAttempt(id uuid.UUID) (*models.QuizAttempt, error)
	GetUserQuizAttempts(userID, quizID uuid.UUID) ([]models.QuizAttempt, error)
	GetActiveQuizAttempt(userID, quizID uuid.UUID) (*models.QuizAttempt, error)
	DeleteActiveQuizAttempt(userID, quizID uuid.UUID) error

	// Retake operations
	ValidateRetakeLimit(userID, quizID uuid.UUID) error
	GetQuizAttemptsForComparison(userID, quizID uuid.UUID) ([]models.QuizAttempt, error)
	CalculatePerformanceComparison(currentAttempt *models.QuizAttempt, previousAttempts []models.QuizAttempt) map[string]interface{}

	// Attempt history operations
	GetUserAttempts(userID uuid.UUID, filters *models.AttemptFilters) ([]models.QuizAttemptWithDetails, int, error)
	GetAttemptWithDetails(attemptID uuid.UUID) (*models.QuizAttemptWithDetails, error)

	// Favorites operations
	AddFavorite(userID, quizID uuid.UUID) error
	RemoveFavorite(userID, quizID uuid.UUID) error
	GetUserFavorites(userID uuid.UUID, page, pageSize int) ([]models.UserQuizFavorite, int, error)
	IsFavorite(userID, quizID uuid.UUID) (bool, error)
}

// FriendsRepositoryInterface defines the contract for friends data operations
type FriendsRepositoryInterface interface {
	// Friend request operations
	SendFriendRequest(requesterID, requestedID uuid.UUID, message *string) (*models.FriendRequest, error)
	GetFriendRequest(id uuid.UUID) (*models.FriendRequest, error)
	GetFriendRequestBetweenUsers(requesterID, requestedID uuid.UUID) (*models.FriendRequest, error)
	RespondToFriendRequest(requestID uuid.UUID, status string) error
	CancelFriendRequest(requestID uuid.UUID, requesterID uuid.UUID) error
	GetPendingFriendRequests(userID uuid.UUID) ([]models.FriendRequest, error)
	GetSentFriendRequests(userID uuid.UUID) ([]models.FriendRequest, error)

	// Friendship operations
	GetFriends(userID uuid.UUID) ([]models.Friend, error)
	GetFriendship(user1ID, user2ID uuid.UUID) (*models.Friendship, error)
	RemoveFriend(userID, friendID uuid.UUID) error
	AreFriends(user1ID, user2ID uuid.UUID) (bool, error)

	// User search operations
	SearchUsers(searchQuery string, currentUserID uuid.UUID, limit, offset int) ([]models.UserSearchResult, int, error)

	// Friend notification operations
	GetFriendNotifications(userID uuid.UUID, limit, offset int) ([]models.FriendNotification, int, error)
	MarkNotificationAsRead(notificationID uuid.UUID, userID uuid.UUID) error
	MarkAllNotificationsAsRead(userID uuid.UUID) error
	GetUnreadNotificationCount(userID uuid.UUID) (int, error)
}

// ChallengesRepositoryInterface defines the contract for challenges data operations
type ChallengesRepositoryInterface interface {
	// Challenge CRUD operations
	CreateChallenge(challenge *models.Challenge) error
	GetChallengeByID(id uuid.UUID) (*models.Challenge, error)
	GetChallengeWithDetails(id uuid.UUID) (*models.ChallengeWithDetails, error)
	UpdateChallenge(challenge *models.Challenge) error
	UpdateChallengeStatus(challengeID uuid.UUID, status string) error
	UpdateChallengeScore(challengeID uuid.UUID, userID uuid.UUID, score float64) error
	DeleteChallenge(id uuid.UUID) error

	// Challenge list operations
	GetUserChallenges(userID uuid.UUID, filters *models.ChallengeFilters) ([]models.ChallengeWithDetails, int, error)
	GetPendingChallenges(userID uuid.UUID) ([]models.ChallengeWithDetails, error)
	GetActiveChallenges(userID uuid.UUID) ([]models.ChallengeWithDetails, error)
	GetCompletedChallenges(userID uuid.UUID) ([]models.ChallengeWithDetails, error)

	// Challenge status operations
	AcceptChallenge(challengeID uuid.UUID, userID uuid.UUID) error
	DeclineChallenge(challengeID uuid.UUID, userID uuid.UUID) error
	CompleteChallenge(challengeID uuid.UUID) error

	// Challenge statistics
	GetChallengeStats(userID uuid.UUID) (*models.ChallengeStatsResponse, error)

	// Challenge validation
	CanUserChallenge(challengerID, challengedID uuid.UUID) (bool, error)
	HasPendingChallenge(challengerID, challengedID uuid.UUID, quizID uuid.UUID) (bool, error)

	// Utility operations
	ExpireChallenges() error
}

// LeaderboardRepositoryInterface defines the contract for leaderboard data operations
type LeaderboardRepositoryInterface interface {
	// Leaderboard operations
	GetGlobalLeaderboard(period string, limit, offset int) ([]models.LeaderboardEntry, int, error)
	GetFriendsLeaderboard(userID uuid.UUID, period string, limit, offset int) ([]models.LeaderboardEntry, int, error)
	GetUserRank(userID uuid.UUID, period string) (*models.UserRankInfo, error)

	// User score update operations
	UpdateUserScore(userID uuid.UUID, points int, quizID uuid.UUID) error
	RecalculateUserLevel(userID uuid.UUID) error

	// Achievement operations for leaderboard
	GetUserAchievements(userID uuid.UUID) ([]string, error)
	GetUserCategoryPoints(userID uuid.UUID) (map[string]int, error)
}

// AchievementRepositoryInterface defines the contract for achievement data operations
type AchievementRepositoryInterface interface {
	// Achievement read operations
	GetAllAchievements() ([]models.Achievement, error)
	GetAchievementByKey(key string) (*models.Achievement, error)

	// User achievement operations
	GetUserAchievements(userID uuid.UUID) ([]models.UserAchievement, error)
	UnlockAchievement(userID uuid.UUID, achievementKey string) (*models.UserAchievement, error)
	HasUserAchievement(userID, achievementID uuid.UUID) (bool, error)
	HasUserAchievementByKey(userID uuid.UUID, achievementKey string) (bool, error)

	// Achievement progress operations
	GetAchievementProgress(userID uuid.UUID) ([]models.AchievementProgress, error)
}

// NotificationRepositoryInterface defines the contract for notification data operations
type NotificationRepositoryInterface interface {
	// Notification CRUD operations
	CreateNotification(notification *models.CreateNotificationRequest) (*models.Notification, error)
	GetNotifications(userID uuid.UUID, filters *models.NotificationFilters) ([]models.Notification, int, error)
	GetNotificationByID(notificationID uuid.UUID, userID uuid.UUID) (*models.Notification, error)
	DeleteNotification(notificationID uuid.UUID, userID uuid.UUID) error // Soft delete
	HardDeleteNotification(notificationID uuid.UUID, userID uuid.UUID) error // Hard delete (admin)
	RestoreNotification(notificationID uuid.UUID, userID uuid.UUID) error // Restore soft deleted

	// Notification read status operations
	MarkNotificationAsRead(notificationID uuid.UUID, userID uuid.UUID) error
	MarkNotificationAsUnread(notificationID uuid.UUID, userID uuid.UUID) error
	MarkAllNotificationsAsRead(userID uuid.UUID) error
	GetUnreadNotificationCount(userID uuid.UUID) (int, error)

	// Notification statistics and management
	GetNotificationStats(userID uuid.UUID) (*models.NotificationStatsResponse, error)
	CleanupExpiredNotifications() error

	// Backward compatibility for friend notifications
	GetFriendNotifications(userID uuid.UUID, limit, offset int) ([]models.FriendNotificationCompat, int, error)
	GetFriendUnreadNotificationCount(userID uuid.UUID) (int, error)
}

// QuizSessionRepositoryInterface defines the contract for quiz session data operations
type QuizSessionRepositoryInterface interface {
	// Session CRUD operations
	CreateSession(session *models.QuizSession) error
	UpdateSession(session *models.QuizSession) error
	GetSessionByID(id uuid.UUID) (*models.QuizSession, error)
	GetSessionByAttemptID(attemptID uuid.UUID) (*models.QuizSession, error)
	DeleteSession(id uuid.UUID) error

	// Session state operations
	PauseSession(attemptID uuid.UUID, pauseData *models.PauseSessionRequest) error
	ResumeSession(attemptID uuid.UUID) error
	AbandonSession(attemptID uuid.UUID) error
	CompleteSession(attemptID uuid.UUID) error

	// Session queries
	GetActiveSession(userID, quizID uuid.UUID) (*models.QuizSession, error)
	GetUserActiveSessions(userID uuid.UUID, filters *models.SessionFilters) ([]models.QuizSessionWithDetails, int, error)
	GetSessionWithDetails(sessionID uuid.UUID) (*models.QuizSessionWithDetails, error)

	// Session management
	SaveSessionProgress(sessionID uuid.UUID, updateData *models.UpdateQuizSessionRequest) error
	UpdateSessionActivity(sessionID uuid.UUID) error
	CleanupExpiredSessions() (int, error)

	// Session validation
	HasActiveSession(userID, quizID uuid.UUID) (bool, error)
	CanResumeSession(attemptID uuid.UUID, userID uuid.UUID) (bool, error)
}

// Repository aggregates all repository interfaces
type Repository struct {
	User         UserRepositoryInterface
	Quiz         QuizRepositoryInterface
	QuizSession  QuizSessionRepositoryInterface
	Friends      FriendsRepositoryInterface
	Challenges   ChallengesRepositoryInterface
	Leaderboard  LeaderboardRepositoryInterface
	Achievement  AchievementRepositoryInterface
	Notification NotificationRepositoryInterface
	Discussion   DiscussionRepositoryInterface
}

// NewRepository creates a new repository instance
func NewRepository() *Repository {
	return &Repository{
		User:         NewUserRepository(),
		Quiz:         NewQuizRepository(),
		QuizSession:  NewQuizSessionRepository(),
		Friends:      NewFriendsRepository(),
		Challenges:   NewChallengesRepository(),
		Leaderboard:  NewLeaderboardRepository(),
		Achievement:  NewAchievementRepository(),
		Notification: NewNotificationRepository(),
		Discussion:   NewDiscussionRepository(),
	}
}
