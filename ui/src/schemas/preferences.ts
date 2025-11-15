import { z } from 'zod'

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
  preferred_difficulty: z
    .enum(['easy', 'medium', 'hard', 'all'], {
      errorMap: () => ({ message: 'Please select a valid difficulty level' }),
    })
    .optional(),
  notification_frequency: z
    .enum(['instant', 'daily', 'weekly', 'never'], {
      errorMap: () => ({ message: 'Please select a valid notification frequency' }),
    })
    .optional(),
  email_notifications: z
    .boolean()
    .optional(),
  theme: z
    .enum(['light', 'dark', 'system'], {
      errorMap: () => ({ message: 'Please select a valid theme' }),
    })
    .optional(),
  // Privacy settings
  profile_visibility: z
    .enum(['public', 'friends_only', 'private'], {
      errorMap: () => ({ message: 'Please select a valid visibility level' }),
    })
    .optional(),
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
 */
export const difficultyPreferenceSchema = z
  .enum(['easy', 'medium', 'hard', 'all'], {
    errorMap: () => ({ message: 'Please select a valid difficulty level' }),
  })

/**
 * Notification Settings Schema
 */
export const notificationSettingsSchema = z.object({
  email_notifications: z.boolean(),
  notification_frequency: z.enum(['instant', 'daily', 'weekly', 'never']),
})

export type NotificationSettingsFormData = z.infer<typeof notificationSettingsSchema>

/**
 * Privacy Settings Schema
 */
export const privacySettingsSchema = z.object({
  profile_visibility: z.enum(['public', 'friends_only', 'private']).optional(),
  show_achievements: z.boolean().optional(),
  show_stats: z.boolean().optional(),
  allow_friend_requests: z.boolean().optional(),
})

export type PrivacySettingsFormData = z.infer<typeof privacySettingsSchema>

/**
 * Theme Preference Schema (standalone)
 */
export const themePreferenceSchema = z.enum(['light', 'dark', 'system'], {
  errorMap: () => ({ message: 'Please select a valid theme' }),
})
