import { z } from 'zod'
import { friendRequestActionSchema } from '@/constants/schemas'

/**
 * Send Friend Request Schema
 */
export const sendFriendRequestSchema = z.object({
  receiver_id: z
    .string()
    .min(1, 'User ID is required')
    .uuid('Invalid user ID format'),
})

export type SendFriendRequestData = z.infer<typeof sendFriendRequestSchema>

/**
 * Friend Request Response Schema
 */
export const friendRequestResponseSchema = z.object({
  action: friendRequestActionSchema,
})

export type FriendRequestResponseData = z.infer<typeof friendRequestResponseSchema>

/**
 * Search Users Schema
 */
export const searchUsersSchema = z.object({
  query: z
    .string()
    .min(2, 'Search query must be at least 2 characters')
    .max(100, 'Search query must be less than 100 characters'),
})

export type SearchUsersData = z.infer<typeof searchUsersSchema>
