import { z } from 'zod'

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
  action: z.enum(['accept', 'decline'], {
    required_error: 'Action is required',
    invalid_type_error: 'Action must be either accept or decline',
  }),
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
