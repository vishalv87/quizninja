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
  PauseSessionRequest,
  SaveProgressRequest,
  SessionActionResponse,
  ResumeSessionResponse,
  ActiveSessionsResponse,
  SessionFilters,
  QuizSession,
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
    console.log("[DEBUGGG] getQuizQuestions - Raw response:", response);
    console.log("[DEBUGGG] getQuizQuestions - Response type:", typeof response);
    console.log("[DEBUGGG] getQuizQuestions - Response keys:", response ? Object.keys(response) : "null");

    // Try to extract data - handle both wrapped and unwrapped responses
    const questions = (response as any)?.data ?? response;
    console.log("[DEBUGGG] getQuizQuestions - Extracted questions:", questions);
    console.log("[DEBUGGG] getQuizQuestions - Questions is array:", Array.isArray(questions));
    console.log("[DEBUGGG] getQuizQuestions - Questions length:", Array.isArray(questions) ? questions.length : "N/A");

    if (Array.isArray(questions) && questions.length > 0) {
      console.log("[DEBUGGG] getQuizQuestions - First question:", questions[0]);
      console.log("[DEBUGGG] getQuizQuestions - First question options:", questions[0]?.options);
    }

    return questions as Question[];
  } catch (error) {
    apiLogger.error("Error fetching quiz questions", { quizId, error });
    console.log("[DEBUGGG] getQuizQuestions - ERROR:", error);
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
    const response = await apiClient.post<QuizAttempt>(
      API_ENDPOINTS.QUIZ.START_ATTEMPT(quizId)
    ) as unknown as QuizAttempt;
    return response;
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
    const response = await apiClient.put<QuizAttempt>(
      API_ENDPOINTS.QUIZ.UPDATE_ATTEMPT(quizId, attemptId),
      { answers }
    ) as unknown as QuizAttempt;
    return response;
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

    const response = await apiClient.post<QuizResults>(
      API_ENDPOINTS.QUIZ.SUBMIT_ATTEMPT(quizId, attemptId),
      {
        attemptId,
        answers: formattedAnswers,
      }
    ) as unknown as QuizResults;
    return response;
  } catch (error) {
    apiLogger.error("Error submitting quiz attempt", {
      quizId,
      attemptId,
      error,
    });
    throw error;
  }
}

// ============ QUIZ SESSION MANAGEMENT (Pause/Resume) ============

/**
 * Pause quiz session - Saves current progress and sets state to paused
 */
export async function pauseQuizSession(
  quizId: string,
  attemptId: string,
  pauseRequest: PauseSessionRequest
): Promise<SessionActionResponse> {
  try {
    const response = await apiClient.post<{ data: SessionActionResponse }>(
      API_ENDPOINTS.QUIZ.PAUSE(quizId, attemptId),
      pauseRequest
    ) as unknown as { data: SessionActionResponse };
    return response.data;
  } catch (error) {
    apiLogger.error("Error pausing quiz session", {
      quizId,
      attemptId,
      error,
    });
    throw error;
  }
}

/**
 * Resume quiz session - Restores paused session back to active
 */
export async function resumeQuizSession(
  quizId: string,
  attemptId: string
): Promise<ResumeSessionResponse> {
  try {
    const response = await apiClient.post<{ data: ResumeSessionResponse }>(
      API_ENDPOINTS.QUIZ.RESUME(quizId, attemptId)
    ) as unknown as { data: ResumeSessionResponse };
    return response.data;
  } catch (error) {
    apiLogger.error("Error resuming quiz session", {
      quizId,
      attemptId,
      error,
    });
    throw error;
  }
}

/**
 * Save quiz progress without changing session state
 */
export async function saveSessionProgress(
  quizId: string,
  attemptId: string,
  progressRequest: SaveProgressRequest
): Promise<void> {
  // Guard against undefined attemptId to prevent invalid API calls
  if (!attemptId) {
    apiLogger.warn("Skipping save progress - attemptId is undefined", { quizId });
    return;
  }

  try {
    await apiClient.put(
      API_ENDPOINTS.QUIZ.SAVE_PROGRESS(quizId, attemptId),
      progressRequest
    );
  } catch (error) {
    apiLogger.error("Error saving session progress", {
      quizId,
      attemptId,
      error,
    });
    throw error;
  }
}

/**
 * Abandon quiz session - Marks attempt as abandoned
 */
export async function abandonQuizSession(
  quizId: string,
  attemptId: string
): Promise<SessionActionResponse> {
  try {
    const response = await apiClient.delete<{ data: SessionActionResponse }>(
      API_ENDPOINTS.QUIZ.ABANDON(quizId, attemptId)
    ) as unknown as { data: SessionActionResponse };
    return response.data;
  } catch (error) {
    apiLogger.error("Error abandoning quiz session", {
      quizId,
      attemptId,
      error,
    });
    throw error;
  }
}

/**
 * Get user's active quiz sessions with optional filters
 */
export async function getUserActiveSessions(
  filters?: SessionFilters
): Promise<ActiveSessionsResponse> {
  try {
    // Build query parameters manually to avoid Axios array serialization issues
    const params = new URLSearchParams();
    if (filters?.quiz_id) params.append('quiz_id', filters.quiz_id);
    if (filters?.session_state) params.append('session_state', filters.session_state);
    if (filters?.category) params.append('category', filters.category);
    if (filters?.difficulty) params.append('difficulty', filters.difficulty);
    if (filters?.page) params.append('page', filters.page.toString());
    if (filters?.page_size) params.append('page_size', filters.page_size.toString());
    if (filters?.sort_by) params.append('sort_by', filters.sort_by);
    if (filters?.sort_order) params.append('sort_order', filters.sort_order);

    const url = `${API_ENDPOINTS.USERS.ACTIVE_SESSIONS}${params.toString() ? '?' + params.toString() : ''}`;

    const response = await apiClient.get<{ data: ActiveSessionsResponse }>(url) as unknown as { data: ActiveSessionsResponse };
    return response.data;
  } catch (error) {
    apiLogger.error("Error fetching active sessions", error);
    throw error;
  }
}

/**
 * Get active session for a specific quiz (if exists)
 */
export async function getQuizActiveSession(
  quizId: string
): Promise<QuizSession | null> {
  try {
    const filters: SessionFilters = {
      quiz_id: quizId,
      session_state: 'active',
      page_size: 1
    };
    const response = await getUserActiveSessions(filters);
    const session = response.sessions.length > 0 ? response.sessions[0] : null;
    return session;
  } catch (error) {
    apiLogger.error("Error fetching active session for quiz", { quizId, error });
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
    const response = await apiClient.get<QuizAttempt>(
      API_ENDPOINTS.USERS.ATTEMPT_DETAILS(attemptId)
    ) as unknown as QuizAttempt;
    return response;
  } catch (error) {
    apiLogger.error("Error fetching attempt details", { attemptId, error });
    throw error;
  }
}

/**
 * Get user's active quiz sessions
 */
export async function getActiveSessions(): Promise<QuizAttempt[]> {
  try {
    const response = await apiClient.get<QuizAttempt[]>(
      API_ENDPOINTS.USERS.ACTIVE_SESSIONS
    ) as unknown as QuizAttempt[];
    return response;
  } catch (error) {
    apiLogger.error("Error fetching active sessions", error);
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
    const response = await apiClient.get<QuizAttempt>(
      API_ENDPOINTS.USERS.QUIZ_ATTEMPT(quizId)
    ) as unknown as QuizAttempt;
    return response;
  } catch (error: any) {
    // If 404, no active attempt exists
    if (error.status === 404) {
      return null;
    }
    apiLogger.error("Error fetching active attempt", { quizId, error });
    throw error;
  }
}
