export const APP_NAME = "QuizNinja";
export const APP_DESCRIPTION = "Test Your Knowledge, Challenge Your Friends";

export const ROUTES = {
  HOME: "/",
  LOGIN: "/login",
  REGISTER: "/register",
  DASHBOARD: "/dashboard",
  QUIZZES: "/quizzes",
  QUIZ_DETAIL: (id: string) => `/quizzes/${id}`,
  QUIZ_TAKE: (id: string) => `/quizzes/${id}/take`,
  PROFILE: "/profile",
  PROFILE_EDIT: "/profile/edit",
  FRIENDS: "/friends",
  FRIEND_REQUESTS: "/friends/requests",
  CHALLENGES: "/challenges",
  CHALLENGE_DETAIL: (id: string) => `/challenges/${id}`,
  ACHIEVEMENTS: "/achievements",
  LEADERBOARD: "/leaderboard",
  DISCUSSIONS: "/discussions",
  DISCUSSION_DETAIL: (id: string) => `/discussions/${id}`,
  NOTIFICATIONS: "/notifications",
  SETTINGS: "/settings",
  ONBOARDING_WELCOME: "/welcome",
  ONBOARDING_PREFERENCES: "/preferences",
} as const;

export const DIFFICULTY_LEVELS = [
  { value: "easy", label: "Easy" },
  { value: "medium", label: "Medium" },
  { value: "hard", label: "Hard" },
] as const;

export const QUIZ_CATEGORIES = [
  { value: "science", label: "Science" },
  { value: "history", label: "History" },
  { value: "geography", label: "Geography" },
  { value: "literature", label: "Literature" },
  { value: "sports", label: "Sports" },
  { value: "entertainment", label: "Entertainment" },
  { value: "technology", label: "Technology" },
  { value: "art", label: "Art" },
] as const;

export const NOTIFICATION_TYPES = {
  FRIEND_REQUEST: "friend_request",
  FRIEND_ACCEPTED: "friend_accepted",
  CHALLENGE_RECEIVED: "challenge_received",
  CHALLENGE_ACCEPTED: "challenge_accepted",
  CHALLENGE_COMPLETED: "challenge_completed",
  ACHIEVEMENT_UNLOCKED: "achievement_unlocked",
  QUIZ_REMINDER: "quiz_reminder",
  DISCUSSION_REPLY: "discussion_reply",
  SYSTEM: "system",
} as const;

export const ACHIEVEMENT_CATEGORIES = {
  QUIZ_MASTER: "quiz_master",
  SOCIAL: "social",
  STREAK: "streak",
  KNOWLEDGE: "knowledge",
  COMPETITOR: "competitor",
} as const;

export const QUERY_KEYS = {
  USER: "user",
  PROFILE: "profile",
  QUIZZES: "quizzes",
  QUIZ: "quiz",
  QUIZ_ATTEMPTS: "quiz-attempts",
  FRIENDS: "friends",
  FRIEND_REQUESTS: "friend-requests",
  CHALLENGES: "challenges",
  CHALLENGE: "challenge",
  ACHIEVEMENTS: "achievements",
  LEADERBOARD: "leaderboard",
  NOTIFICATIONS: "notifications",
  DISCUSSIONS: "discussions",
  DISCUSSION: "discussion",
  CATEGORIES: "categories",
  PREFERENCES: "preferences",
  USER_STATS: "user-stats",
} as const;

export const STORAGE_KEYS = {
  AUTH_TOKEN: "auth_token",
  THEME: "theme",
  RECENT_QUIZZES: "recent_quizzes",
} as const;

export const API_ERROR_MESSAGES = {
  NETWORK_ERROR: "Network error. Please check your connection.",
  UNAUTHORIZED: "Please log in to continue.",
  FORBIDDEN: "You don't have permission to perform this action.",
  NOT_FOUND: "The requested resource was not found.",
  SERVER_ERROR: "Server error. Please try again later.",
  VALIDATION_ERROR: "Please check your input and try again.",
} as const;
