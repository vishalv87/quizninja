import { z } from "zod";

/**
 * Create Discussion Schema
 * Validates discussion creation request
 */
export const createDiscussionSchema = z.object({
  quiz_id: z
    .string()
    .min(1, "Quiz is required")
    .uuid("Invalid quiz ID format"),
  title: z
    .string()
    .min(3, "Title must be at least 3 characters")
    .max(200, "Title must not exceed 200 characters"),
  content: z
    .string()
    .min(10, "Content must be at least 10 characters")
    .max(5000, "Content must not exceed 5000 characters"),
});

export type CreateDiscussionData = z.infer<typeof createDiscussionSchema>;

/**
 * Update Discussion Schema
 * Validates discussion update request
 */
export const updateDiscussionSchema = z.object({
  title: z
    .string()
    .min(3, "Title must be at least 3 characters")
    .max(200, "Title must not exceed 200 characters")
    .optional(),
  content: z
    .string()
    .min(10, "Content must be at least 10 characters")
    .max(5000, "Content must not exceed 5000 characters")
    .optional(),
});

export type UpdateDiscussionData = z.infer<typeof updateDiscussionSchema>;

/**
 * Create Discussion Reply Schema
 * Validates reply creation request
 */
export const createDiscussionReplySchema = z.object({
  content: z
    .string()
    .min(1, "Reply content is required")
    .max(2000, "Reply must not exceed 2000 characters"),
});

export type CreateDiscussionReplyData = z.infer<
  typeof createDiscussionReplySchema
>;

/**
 * Discussion Filters Schema
 * Validates discussion list filtering parameters
 */
export const discussionFiltersSchema = z.object({
  quiz_id: z.string().uuid("Invalid quiz ID format").optional(),
  sort: z.enum(["recent", "popular"]).optional(),
  limit: z.number().min(1).max(100).optional(),
  offset: z.number().min(0).optional(),
});

export type DiscussionFiltersData = z.infer<typeof discussionFiltersSchema>;
