import { z } from 'zod'

/**
 * Create Challenge Schema
 * Validates challenge creation request with opponent and quiz selection
 */
export const createChallengeSchema = z.object({
  opponent_id: z
    .string()
    .min(1, 'Opponent is required')
    .uuid('Invalid opponent ID format'),
  quiz_id: z
    .string()
    .min(1, 'Quiz is required')
    .uuid('Invalid quiz ID format'),
})

export type CreateChallengeData = z.infer<typeof createChallengeSchema>

/**
 * Link Attempt Schema
 * Validates linking a quiz attempt to a challenge
 */
export const linkAttemptSchema = z.object({
  attempt_id: z
    .string()
    .min(1, 'Attempt ID is required')
    .uuid('Invalid attempt ID format'),
})

export type LinkAttemptData = z.infer<typeof linkAttemptSchema>

/**
 * Challenge Action Schema
 * Validates challenge actions (accept, decline)
 */
export const challengeActionSchema = z.object({
  action: z.enum(['accept', 'decline'], {
    required_error: 'Action is required',
    invalid_type_error: 'Action must be either accept or decline',
  }),
})

export type ChallengeActionData = z.infer<typeof challengeActionSchema>
