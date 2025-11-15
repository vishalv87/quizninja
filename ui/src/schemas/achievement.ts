import { z } from 'zod'

/**
 * Achievement Filter Schema
 * Validates achievement filtering options
 */
export const achievementFilterSchema = z.object({
  category: z.string().optional(),
  unlocked: z.boolean().optional(),
  search: z.string().optional(),
})

export type AchievementFilterData = z.infer<typeof achievementFilterSchema>

/**
 * Achievement Category Schema
 * Validates achievement category
 */
export const achievementCategorySchema = z.object({
  category: z
    .string()
    .min(1, 'Category is required')
    .max(50, 'Category must be less than 50 characters'),
})

export type AchievementCategoryData = z.infer<typeof achievementCategorySchema>
