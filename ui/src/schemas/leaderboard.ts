import { z } from 'zod'
import {
  leaderboardPeriodSchema,
  leaderboardSortFieldSchema,
  sortOrderSchema,
} from '@/constants/schemas'

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
  period: leaderboardPeriodSchema.optional().default('all-time'),
})

export type LeaderboardFilterData = z.infer<typeof leaderboardFilterSchema>

/**
 * Leaderboard Sort Schema
 * Validates leaderboard sorting options
 */
export const leaderboardSortSchema = z.object({
  sortBy: leaderboardSortFieldSchema.optional().default('points'),
  order: sortOrderSchema.optional().default('desc'),
})

export type LeaderboardSortData = z.infer<typeof leaderboardSortSchema>