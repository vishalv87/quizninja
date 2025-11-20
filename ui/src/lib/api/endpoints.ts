/**
 * API Endpoint Constants
 * Based on backend routes from /api/v1
 */

export const API_ENDPOINTS = {
  // Health & Info
  HEALTH: "/health",
  PING: "/ping",

  // Authentication
  AUTH: {
    REGISTER: "/auth/register",
    LOGIN: "/auth/login",
    LOGOUT: "/auth/logout",
  },

  // Profile
  PROFILE: {
    GET: "/profile",
    UPDATE: "/profile",
  },

  // Users
  USERS: {
    PROFILE: (userId: string) => `/users/${userId}`,
    PREFERENCES: {
      GET: "/users/preferences",
      UPDATE: "/users/preferences",
    },
    ONBOARDING: {
      COMPLETE: "/users/onboarding/complete",
      STATUS: "/users/onboarding/status",
    },
    QUIZZES: "/users/quizzes",
    QUIZ_ATTEMPT: (quizId: string) => `/users/quizzes/${quizId}/attempt`,
    STATS: "/users/stats",
    USER_STATS: (userId: string) => `/users/${userId}/stats`,
    ATTEMPTS: "/users/attempts",
    USER_ATTEMPTS: (userId: string) => `/users/${userId}/attempts`,
    ATTEMPT_DETAILS: (attemptId: string) => `/users/attempts/${attemptId}`,
    ACHIEVEMENTS: "/users/achievements",
    ACHIEVEMENTS_BY_USER: (userId: string) => `/users/${userId}/achievements`,
  },

  // Quizzes (Public)
  QUIZZES: {
    LIST: "/quizzes",
    FEATURED: "/quizzes/featured",
    BY_CATEGORY: (category: string) => `/quizzes/category/${category}`,
    CATEGORIES: "/quizzes/categories",
  },

  // Quizzes (Protected)
  QUIZ: {
    GET: (id: string) => `/quizzes/${id}`,
    QUESTIONS: (id: string) => `/quizzes/${id}/questions`,
    START_ATTEMPT: (id: string) => `/quizzes/${id}/attempts`,
    SUBMIT_ATTEMPT: (id: string, attemptId: string) => `/quizzes/${id}/attempts/${attemptId}/submit`,
    UPDATE_ATTEMPT: (id: string, attemptId: string) => `/quizzes/${id}/attempts/${attemptId}`,
    PAUSE: (id: string, attemptId: string) => `/quizzes/${id}/attempts/${attemptId}/pause`,
    RESUME: (id: string, attemptId: string) => `/quizzes/${id}/attempts/${attemptId}/resume`,
    SAVE_PROGRESS: (id: string, attemptId: string) => `/quizzes/${id}/attempts/${attemptId}/save-progress`,
    ABANDON: (id: string, attemptId: string) => `/quizzes/${id}/attempts/${attemptId}/abandon`,
  },

  // Categories
  CATEGORIES: {
    LIST: "/categories",
    GROUPS: "/categories",
  },

  // Config
  CONFIG: {
    APP_SETTINGS: "/config/app-settings",
  },

  // Preferences (Public)
  PREFERENCES: {
    CATEGORIES: "/preferences/categories",
    DIFFICULTY_LEVELS: "/preferences/difficulty-levels",
    NOTIFICATION_FREQUENCIES: "/preferences/notification-frequencies",
  },

  // Friends
  FRIENDS: {
    LIST: "/friends",
    REQUESTS: {
      SEND: "/friends/requests",
      LIST: "/friends/requests",
      RESPOND: (id: string) => `/friends/requests/${id}`,
      CANCEL: (id: string) => `/friends/requests/${id}`,
    },
    REMOVE: (id: string) => `/friends/${id}`,
    SEARCH: "/friends/search",
    NOTIFICATIONS: "/friends/notifications",
    NOTIFICATION_READ: (id: string) => `/friends/notifications/${id}/read`,
    NOTIFICATIONS_READ_ALL: "/friends/notifications/read-all",
  },

  // Challenges
  CHALLENGES: {
    CREATE: "/challenges",
    LIST: "/challenges",
    STATS: "/challenges/stats",
    PENDING: "/challenges/pending",
    ACTIVE: "/challenges/active",
    COMPLETED: "/challenges/completed",
    GET: (id: string) => `/challenges/${id}`,
    ACCEPT: (id: string) => `/challenges/${id}/accept`,
    DECLINE: (id: string) => `/challenges/${id}/decline`,
    CANCEL: (id: string) => `/challenges/${id}/cancel`,
    UPDATE_SCORE: (id: string) => `/challenges/${id}/score`,
    LINK_ATTEMPT: (id: string) => `/challenges/${id}/link-attempt`,
    COMPLETE: (id: string) => `/challenges/${id}/complete`,
    EXPIRE: "/challenges/expire",
  },

  // Leaderboard
  LEADERBOARD: {
    GET: "/leaderboard",
    STATS: "/leaderboard/stats",
    RANK: "/leaderboard/rank",
    UPDATE_SCORE: "/leaderboard/score",
    WITH_ACHIEVEMENTS: "/leaderboard/achievements",
  },

  // Achievements
  ACHIEVEMENTS: {
    LIST: "/achievements",
    PROGRESS: "/achievements/progress",
    STATS: "/achievements/stats",
    CHECK: "/achievements/check",
    BY_CATEGORY: (category: string) => `/achievements/category/${category}`,
    UNLOCK: (key: string) => `/achievements/unlock/${key}`,
  },

  // Favorites
  FAVORITES: {
    ADD: "/favorites",
    REMOVE: (quizId: string) => `/favorites/${quizId}`,
    LIST: "/favorites",
    CHECK: (quizId: string) => `/favorites/check/${quizId}`,
  },

  // Discussions
  DISCUSSIONS: {
    LIST: "/discussions",
    CREATE: "/discussions",
    STATS: "/discussions/stats",
    GET: (id: string) => `/discussions/${id}`,
    UPDATE: (id: string) => `/discussions/${id}`,
    DELETE: (id: string) => `/discussions/${id}`,
    LIKE: (id: string) => `/discussions/${id}/like`,
    REPLIES: {
      LIST: (id: string) => `/discussions/${id}/replies`,
      CREATE: (id: string) => `/discussions/${id}/replies`,
      UPDATE: (replyId: string) => `/discussions/replies/${replyId}`,
      DELETE: (replyId: string) => `/discussions/replies/${replyId}`,
      LIKE: (replyId: string) => `/discussions/replies/${replyId}/like`,
    },
  },

  // Notifications
  NOTIFICATIONS: {
    LIST: "/notifications",
    STATS: "/notifications/stats",
    GET: (id: string) => `/notifications/${id}`,
    MARK_READ: (id: string) => `/notifications/${id}/read`,
    MARK_UNREAD: (id: string) => `/notifications/${id}/unread`,
    MARK_ALL_READ: "/notifications/read-all",
    DELETE: (id: string) => `/notifications/${id}`,
    CREATE: "/notifications",
    CLEANUP: "/notifications/cleanup",
  },

  // Admin
  ADMIN: {
    CACHE: {
      CLEAR_APP_SETTINGS: "/admin/cache/app-settings",
    },
  },
} as const;