export const APP_NAME = "QuizNinja";
export const APP_DESCRIPTION = "Test Your Knowledge, Compete with Friends";

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
  ACHIEVEMENTS: "/achievements",
  LEADERBOARD: "/leaderboard",
  DISCUSSIONS: "/discussions",
  DISCUSSION_DETAIL: (id: string) => `/discussions/${id}`,
  NOTIFICATIONS: "/notifications",
  SETTINGS: "/settings",
  ONBOARDING_PREFERENCES: "/preferences",
} as const;

// Re-export from centralized constants
export { DIFFICULTY_OPTIONS as DIFFICULTY_LEVELS } from '@/constants';
export { NotificationType as NOTIFICATION_TYPE_ENUM } from '@/constants';

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

// NOTIFICATION_TYPES is now re-exported from @/constants
// Use NotificationType from '@/constants' for type-safe usage
export { NotificationType as NOTIFICATION_TYPES } from '@/constants';

// ACHIEVEMENT_CATEGORIES is now re-exported from @/constants
// Use AchievementCategory from '@/constants' for type-safe usage
export { AchievementCategory as ACHIEVEMENT_CATEGORIES } from '@/constants';

export const QUERY_KEYS = {
  USER: "user",
  PROFILE: "profile",
  QUIZZES: "quizzes",
  QUIZ: "quiz",
  QUIZ_ATTEMPTS: "quiz-attempts",
  FRIENDS: "friends",
  FRIEND_REQUESTS: "friend-requests",
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
