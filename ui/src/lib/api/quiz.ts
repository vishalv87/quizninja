import { apiClient } from "./client";
import { API_ENDPOINTS } from "./endpoints";
import type {
  Quiz,
  QuizFilters,
  QuizAttempt,
  QuizAnswer,
  QuizResults,
  Question,
  Category,
  QuizListResponse,
} from "@/types/quiz";
import { apiLogger } from "@/lib/logger";

/**
 * Quiz API Service
 * Handles all quiz-related API calls
 */

// ============ QUIZ BROWSING ============

/**
 * Get all quizzes with optional filters
 */
export async function getQuizzes(filters?: QuizFilters): Promise<Quiz[]> {
  try {
    const response = await apiClient.get<{ data: QuizListResponse }>(API_ENDPOINTS.QUIZZES.LIST, {
      params: filters,
    }) as unknown as { data: QuizListResponse };
    return response.data.quizzes;
  } catch (error) {
    apiLogger.error("Error fetching quizzes", error);
    throw error;
  }
}

/**
 * Get a single quiz by ID
 */
export async function getQuiz(id: string): Promise<Quiz> {
  try {
    const response = await apiClient.get<{ data: { quiz: Quiz } }>(
      API_ENDPOINTS.QUIZ.GET(id)
    ) as unknown as { data: { quiz: Quiz } };
    return response.data.quiz;
  } catch (error) {
    apiLogger.error("Error fetching quiz", { quizId: id, error });
    throw error;
  }
}

/**
 * Get questions for a quiz
 */
export async function getQuizQuestions(quizId: string): Promise<Question[]> {
  try {
    const response = await apiClient.get<{ data: Question[] }>(
      API_ENDPOINTS.QUIZ.QUESTIONS(quizId)
    );

    // Try to extract data - handle both wrapped and unwrapped responses
    const questions = (response as any)?.data ?? response;

    return questions as Question[];
  } catch (error) {
    apiLogger.error("Error fetching quiz questions", { quizId, error });
    throw error;
  }
}

/**
 * Get featured quizzes
 */
export async function getFeaturedQuizzes(): Promise<Quiz[]> {
  try {
    const response = await apiClient.get<{ data: { quizzes: Quiz[] } }>(
      API_ENDPOINTS.QUIZZES.FEATURED
    ) as unknown as { data: { quizzes: Quiz[] } };
    return response.data.quizzes;
  } catch (error) {
    apiLogger.error("Error fetching featured quizzes", error);
    throw error;
  }
}

/**
 * Get quizzes by category
 */
export async function getQuizzesByCategory(
  categoryId: string
): Promise<Quiz[]> {
  try {
    const response = await apiClient.get<{ data: { quizzes: Quiz[] } }>(
      API_ENDPOINTS.QUIZZES.BY_CATEGORY(categoryId)
    ) as unknown as { data: { quizzes: Quiz[] } };
    return response.data.quizzes;
  } catch (error) {
    apiLogger.error("Error fetching category quizzes", {
      categoryId,
      error,
    });
    throw error;
  }
}

/**
 * Get all categories
 */
export async function getCategories(): Promise<Category[]> {
  try {
    const response = await apiClient.get<{ data: Category[] }>(
      API_ENDPOINTS.CATEGORIES.LIST
    ) as unknown as { data: Category[] };
    return response.data;
  } catch (error) {
    apiLogger.error("Error fetching categories", error);
    throw error;
  }
}

/**
 * Search quizzes by title or description
 */
export async function searchQuizzes(query: string): Promise<Quiz[]> {
  try {
    const response = await apiClient.get<{ data: QuizListResponse }>(API_ENDPOINTS.QUIZZES.LIST, {
      params: { search: query },
    }) as unknown as { data: QuizListResponse };
    return response.data.quizzes;
  } catch (error) {
    apiLogger.error("Error searching quizzes", { query, error });
    throw error;
  }
}

// ============ QUIZ ATTEMPTS ============

/**
 * Start a new quiz attempt
 */
export async function startQuizAttempt(quizId: string): Promise<QuizAttempt> {
  try {
    const response = await apiClient.post<{ data: QuizAttempt }>(
      API_ENDPOINTS.QUIZ.START_ATTEMPT(quizId)
    ) as unknown as { data: QuizAttempt };
    return response.data;
  } catch (error) {
    apiLogger.error("Error starting quiz attempt", { quizId, error });
    throw error;
  }
}

/**
 * Update quiz attempt with answers
 */
export async function updateQuizAttempt(
  quizId: string,
  attemptId: string,
  answers: QuizAnswer[]
): Promise<QuizAttempt> {
  try {
    const response = await apiClient.put<{ data: QuizAttempt }>(
      API_ENDPOINTS.QUIZ.UPDATE_ATTEMPT(quizId, attemptId),
      { answers }
    ) as unknown as { data: QuizAttempt };
    return response.data;
  } catch (error) {
    apiLogger.error("Error updating quiz attempt", {
      quizId,
      attemptId,
      error,
    });
    throw error;
  }
}

/**
 * Submit quiz attempt for scoring
 */
export async function submitQuizAttempt(
  quizId: string,
  attemptId: string,
  answers: QuizAnswer[]
): Promise<QuizResults> {
  try {
    // Transform answers to match backend expected format (camelCase)
    const formattedAnswers = answers.map((answer) => ({
      questionId: answer.question_id,
      selectedOption: answer.selected_answer,
      selectedOptionIndex: answer.selected_option_index,
    }));

    const response = await apiClient.post<{ data: QuizResults }>(
      API_ENDPOINTS.QUIZ.SUBMIT_ATTEMPT(quizId, attemptId),
      {
        attemptId,
        answers: formattedAnswers,
      }
    ) as unknown as { data: QuizResults };
    return response.data;
  } catch (error) {
    apiLogger.error("Error submitting quiz attempt", {
      quizId,
      attemptId,
      error,
    });
    throw error;
  }
}

// ============ QUIZ SESSION MANAGEMENT ============

/**
 * Abandon quiz session - Marks attempt as abandoned
 */
export async function abandonQuizSession(
  quizId: string,
  attemptId: string
): Promise<void> {
  try {
    await apiClient.delete(
      API_ENDPOINTS.QUIZ.ABANDON(quizId, attemptId)
    );
  } catch (error) {
    apiLogger.error("Error abandoning quiz session", {
      quizId,
      attemptId,
      error,
    });
    throw error;
  }
}

// ============ USER QUIZ DATA ============

/**
 * Get user's quiz attempt history
 */
export async function getUserAttempts(): Promise<QuizAttempt[]> {
  try {
    const response = await apiClient.get<QuizAttempt[]>(
      API_ENDPOINTS.USERS.ATTEMPTS
    ) as unknown as QuizAttempt[];
    return response;
  } catch (error) {
    apiLogger.error("Error fetching user attempts", error);
    throw error;
  }
}

/**
 * Get specific attempt details
 */
export async function getAttemptDetails(
  attemptId: string
): Promise<QuizAttempt> {
  try {
    const response = await apiClient.get<{ data: { attempt: QuizAttempt } }>(
      API_ENDPOINTS.USERS.ATTEMPT_DETAILS(attemptId)
    ) as unknown as { data: { attempt: QuizAttempt } };
    return response.data.attempt;
  } catch (error) {
    apiLogger.error("Error fetching attempt details", { attemptId, error });
    throw error;
  }
}

/**
 * Get active attempt for a specific quiz (if exists)
 */
export async function getActiveAttemptForQuiz(
  quizId: string
): Promise<QuizAttempt | null> {
  try {
    const response = await apiClient.get<{ data: { attempt: QuizAttempt } }>(
      API_ENDPOINTS.USERS.QUIZ_ATTEMPT(quizId)
    ) as unknown as { data: { attempt: QuizAttempt } };
    return response.data.attempt;
  } catch (error: any) {
    // If 404, no active attempt exists
    if (error.status === 404) {
      return null;
    }
    apiLogger.error("Error fetching active attempt", { quizId, error });
    throw error;
  }
}

/**
 * Get the latest completed (non-abandoned) attempt for a specific quiz
 */
export async function getLatestCompletedAttempt(
  quizId: string
): Promise<QuizAttempt | null> {
  try {
    const response = await apiClient.get<{ data: { attempt: QuizAttempt } }>(
      API_ENDPOINTS.USERS.QUIZ_COMPLETED_ATTEMPT(quizId)
    ) as unknown as { data: { attempt: QuizAttempt } };
    return response.data.attempt;
  } catch (error: any) {
    // If 404, no completed attempt exists
    if (error.status === 404) {
      return null;
    }
    apiLogger.error("Error fetching latest completed attempt", { quizId, error });
    throw error;
  }
}
