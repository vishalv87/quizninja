import { z } from 'zod'
import { notificationTypeSchema, notificationActionSchema as notificationActionEnumSchema } from '@/constants/schemas'

/**
 * Notification Filter Schema
 * Validates notification filtering and pagination parameters
 */
export const notificationFilterSchema = z.object({
  type: notificationTypeSchema.optional(),
  is_read: z.boolean().optional(),
  limit: z.number().min(1).max(100).optional(),
  offset: z.number().min(0).optional(),
})

export type NotificationFilterData = z.infer<typeof notificationFilterSchema>

/**
 * Create Notification Schema (for admin/system use)
 * Validates notification creation request
 */
export const createNotificationSchema = z.object({
  user_id: z
    .string()
    .min(1, 'User ID is required')
    .uuid('Invalid user ID format'),
  type: notificationTypeSchema,
  title: z.string().min(1, 'Title is required').max(200, 'Title too long'),
  message: z.string().min(1, 'Message is required').max(500, 'Message too long'),
  data: z.record(z.any()).optional(),
  expires_at: z.string().datetime().optional(),
})

export type CreateNotificationData = z.infer<typeof createNotificationSchema>

/**
 * Notification Action Schema
 * Validates notification actions (mark as read/unread)
 */
export const notificationActionSchema = z.object({
  action: notificationActionEnumSchema,
})

export type NotificationActionData = z.infer<typeof notificationActionSchema>
