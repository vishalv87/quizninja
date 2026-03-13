import { z } from "zod";
import { quizDifficultySchema } from "@/constants/schemas";

/**
 * Quiz Validation Schemas
 * Using Zod for runtime type validation
 */

// ============ QUIZ FILTERS ============

export const quizFiltersSchema = z.object({
  category: z.string().optional(),
  difficulty: quizDifficultySchema.optional(),
  search: z.string().min(1).max(100).optional(),
  is_featured: z.boolean().optional(),
  limit: z.number().min(1).max(100).optional().default(20),
  offset: z.number().min(0).optional().default(0),
});

export type QuizFiltersInput = z.infer<typeof quizFiltersSchema>;

// ============ QUIZ SEARCH ============

export const quizSearchSchema = z.object({
  query: z
    .string()
    .min(1, "Search query must be at least 1 character")
    .max(100, "Search query is too long")
    .trim(),
});

export type QuizSearchInput = z.infer<typeof quizSearchSchema>;

// ============ QUIZ ANSWER ============

export const quizAnswerSchema = z.object({
  question_id: z.string().uuid("Invalid question ID"),
  selected_answer: z.string().min(1, "Answer cannot be empty"),
  time_spent_seconds: z.number().min(0).optional(),
});

export const quizAnswersSchema = z.array(quizAnswerSchema);

export type QuizAnswerInput = z.infer<typeof quizAnswerSchema>;

// ============ QUIZ ATTEMPT SUBMISSION ============

export const submitQuizAttemptSchema = z.object({
  answers: quizAnswersSchema,
  time_spent_seconds: z.number().min(0).optional(),
});

export type SubmitQuizAttemptInput = z.infer<typeof submitQuizAttemptSchema>;

// ============ VALIDATORS ============

/**
 * Validate quiz filters
 */
export function validateQuizFilters(filters: unknown) {
  return quizFiltersSchema.safeParse(filters);
}

/**
 * Validate search query
 */
export function validateSearchQuery(query: unknown) {
  return quizSearchSchema.safeParse(query);
}

/**
 * Validate quiz answer
 */
export function validateQuizAnswer(answer: unknown) {
  return quizAnswerSchema.safeParse(answer);
}

/**
 * Validate quiz attempt submission
 */
export function validateQuizSubmission(data: unknown) {
  return submitQuizAttemptSchema.safeParse(data);
}