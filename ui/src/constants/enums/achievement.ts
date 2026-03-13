/**
 * Achievement-related enums and constants
 * Single source of truth for achievement categories and requirement types
 */

// Achievement Categories
export const AchievementCategory = {
  QUIZ_MASTER: 'quiz_master',
  SOCIAL: 'social',
  STREAK: 'streak',
  KNOWLEDGE: 'knowledge',
  COMPETITOR: 'competitor',
} as const;

export type AchievementCategory = typeof AchievementCategory[keyof typeof AchievementCategory];
export const ACHIEVEMENT_CATEGORIES = Object.values(AchievementCategory);

// Type guard for achievement category
export function isAchievementCategory(value: unknown): value is AchievementCategory {
  return typeof value === 'string' && ACHIEVEMENT_CATEGORIES.includes(value as AchievementCategory);
}

// Achievement Requirement Types
export const AchievementRequirementType = {
  QUIZZES_COMPLETED: 'quizzes_completed',
  TOTAL_POINTS: 'total_points',
  ACCURACY_PERCENTAGE: 'accuracy_percentage',
  STREAK_REACHED: 'streak_reached',
  FRIENDS_ADDED: 'friends_added',
  DISCUSSIONS_STARTED: 'discussions_started',
} as const;

export type AchievementRequirementType = typeof AchievementRequirementType[keyof typeof AchievementRequirementType];
export const ACHIEVEMENT_REQUIREMENT_TYPES = Object.values(AchievementRequirementType);

// Type guard for achievement requirement type
export function isAchievementRequirementType(value: unknown): value is AchievementRequirementType {
  return typeof value === 'string' && ACHIEVEMENT_REQUIREMENT_TYPES.includes(value as AchievementRequirementType);
}