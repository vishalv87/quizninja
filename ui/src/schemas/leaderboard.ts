import { z } from 'zod'

/**
 * Leaderboard Filter Schema
 * Validates leaderboard filtering and sorting options
 */
export const leaderboardFilterSchema = z.object({
  limit: z
    .number()
    .min(1, 'Limit must be at least 1')
    .max(100, 'Limit must be at most 100')
    .optional()
    .default(50),
  period: z
    .enum(['all-time', 'monthly', 'weekly'], {
      invalid_type_error: 'Period must be all-time, monthly, or weekly',
    })
    .optional()
    .default('all-time'),
})

export type LeaderboardFilterData = z.infer<typeof leaderboardFilterSchema>

/**
 * Leaderboard Sort Schema
 * Validates leaderboard sorting options
 */
export const leaderboardSortSchema = z.object({
  sortBy: z
    .enum(['points', 'quizzes', 'achievements'], {
      invalid_type_error: 'Sort by must be points, quizzes, or achievements',
    })
    .optional()
    .default('points'),
  order: z
    .enum(['asc', 'desc'], {
      invalid_type_error: 'Order must be asc or desc',
    })
    .optional()
    .default('desc'),
})

export type LeaderboardSortData = z.infer<typeof leaderboardSortSchema>