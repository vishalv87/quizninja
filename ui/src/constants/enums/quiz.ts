/**
 * Quiz-related enums and constants
 * Single source of truth for quiz difficulty, question types, and attempt status
 */

// Quiz Difficulty Levels - matches backend API values
export const QuizDifficulty = {
  BEGINNER: 'beginner',
  INTERMEDIATE: 'intermediate',
  ADVANCED: 'advanced',
} as const;

export type QuizDifficulty = typeof QuizDifficulty[keyof typeof QuizDifficulty];
export const QUIZ_DIFFICULTIES = Object.values(QuizDifficulty);

// Type guard for quiz difficulty
export function isQuizDifficulty(value: unknown): value is QuizDifficulty {
  return typeof value === 'string' && QUIZ_DIFFICULTIES.includes(value as QuizDifficulty);
}

// Question Types
export const QuestionType = {
  MULTIPLE_CHOICE: 'multiple_choice',
  TRUE_FALSE: 'true_false',
  SHORT_ANSWER: 'short_answer',
} as const;

export type QuestionType = typeof QuestionType[keyof typeof QuestionType];
export const QUESTION_TYPES = Object.values(QuestionType);

// Type guard for question type
export function isQuestionType(value: unknown): value is QuestionType {
  return typeof value === 'string' && QUESTION_TYPES.includes(value as QuestionType);
}

// Quiz Attempt Status
export const AttemptStatus = {
  IN_PROGRESS: 'in_progress',
  COMPLETED: 'completed',
  ABANDONED: 'abandoned',
} as const;

export type AttemptStatus = typeof AttemptStatus[keyof typeof AttemptStatus];
export const ATTEMPT_STATUSES = Object.values(AttemptStatus);

// Type guard for attempt status
export function isAttemptStatus(value: unknown): value is AttemptStatus {
  return typeof value === 'string' && ATTEMPT_STATUSES.includes(value as AttemptStatus);
}
