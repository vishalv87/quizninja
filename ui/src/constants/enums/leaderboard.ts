/**
 * Leaderboard-related enums and constants
 * Single source of truth for leaderboard periods and sort options
 */

// Leaderboard Time Periods
export const LeaderboardPeriod = {
  ALL_TIME: 'all-time',
  MONTHLY: 'monthly',
  WEEKLY: 'weekly',
} as const;

export type LeaderboardPeriod = typeof LeaderboardPeriod[keyof typeof LeaderboardPeriod];
export const LEADERBOARD_PERIODS = Object.values(LeaderboardPeriod);

// Type guard for leaderboard period
export function isLeaderboardPeriod(value: unknown): value is LeaderboardPeriod {
  return typeof value === 'string' && LEADERBOARD_PERIODS.includes(value as LeaderboardPeriod);
}

// Leaderboard Sort Fields
export const LeaderboardSortField = {
  POINTS: 'points',
  QUIZZES: 'quizzes',
  ACHIEVEMENTS: 'achievements',
} as const;

export type LeaderboardSortField = typeof LeaderboardSortField[keyof typeof LeaderboardSortField];
export const LEADERBOARD_SORT_FIELDS = Object.values(LeaderboardSortField);

// Type guard for leaderboard sort field
export function isLeaderboardSortField(value: unknown): value is LeaderboardSortField {
  return typeof value === 'string' && LEADERBOARD_SORT_FIELDS.includes(value as LeaderboardSortField);
}
