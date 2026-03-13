import { z } from 'zod'
import {
  quizDifficultySchema,
  themeSchema,
  profileVisibilitySchema,
  notificationFrequencySchema,
} from '@/constants/schemas'
import { QUIZ_DIFFICULTIES } from '@/constants'

// Extended difficulty schema that includes 'all' option for preferences
const preferredDifficultySchema = z.enum(
  ['beginner', 'intermediate', 'advanced', 'all'] as const,
  {
    errorMap: () => ({ message: 'Please select a valid difficulty level' }),
  }
)

/**
 * User Preferences Validation Schema
 * For updating user preferences including categories, difficulty, notifications, privacy, etc.
 */
export const preferencesSchema = z.object({
  preferred_categories: z
    .array(z.string())
    .min(1, 'Please select at least one category')
    .max(10, 'You can select up to 10 categories')
    .optional(),
  preferred_difficulty: preferredDifficultySchema.optional(),
  notification_frequency: notificationFrequencySchema.optional(),
  email_notifications: z
    .boolean()
    .optional(),
  theme: themeSchema.optional(),
  // Privacy settings
  profile_visibility: profileVisibilitySchema.optional(),
  show_achievements: z
    .boolean()
    .optional(),
  show_stats: z
    .boolean()
    .optional(),
  allow_friend_requests: z
    .boolean()
    .optional(),
})

export type PreferencesFormData = z.infer<typeof preferencesSchema>

/**
 * Category Preference Schema (standalone)
 */
export const categoryPreferenceSchema = z
  .array(z.string())
  .min(1, 'Please select at least one category')
  .max(10, 'You can select up to 10 categories')

/**
 * Difficulty Preference Schema (standalone)
 * Includes 'all' option for preferences
 */
export const difficultyPreferenceSchema = preferredDifficultySchema

/**
 * Notification Settings Schema
 */
export const notificationSettingsSchema = z.object({
  email_notifications: z.boolean(),
  notification_frequency: notificationFrequencySchema,
})

export type NotificationSettingsFormData = z.infer<typeof notificationSettingsSchema>

/**
 * Privacy Settings Schema
 */
export const privacySettingsSchema = z.object({
  profile_visibility: profileVisibilitySchema.optional(),
  show_achievements: z.boolean().optional(),
  show_stats: z.boolean().optional(),
  allow_friend_requests: z.boolean().optional(),
})

export type PrivacySettingsFormData = z.infer<typeof privacySettingsSchema>

/**
 * Theme Preference Schema (standalone)
 */
export const themePreferenceSchema = themeSchema
