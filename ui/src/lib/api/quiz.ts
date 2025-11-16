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
    apiLogger.debug("Fetching quizzes with filters", filters);
    const response = await apiClient.get<{ data: QuizListResponse }>(API_ENDPOINTS.QUIZZES.LIST, {
      params: filters,
    }) as unknown as { data: QuizListResponse };
    apiLogger.debug("Quizzes fetched successfully", {
      count: response.data.quizzes.length,
    });
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
    apiLogger.debug("Fetching quiz", { quizId: id });
    const response = await apiClient.get<Quiz>(API_ENDPOINTS.QUIZ.GET(id)) as unknown as Quiz;
    apiLogger.debug("Quiz fetched successfully", { quizId: id });
    return response;
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
    apiLogger.debug("Fetching quiz questions", { quizId });
    const response = await apiClient.get<Question[]>(
      API_ENDPOINTS.QUIZ.QUESTIONS(quizId)
    ) as unknown as Question[];
    apiLogger.debug("Quiz questions fetched successfully", {
      quizId,
      count: response.length,
    });
    return response;
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
    apiLogger.debug("Fetching featured quizzes");
    const response = await apiClient.get<{ data: { quizzes: Quiz[] } }>(
      API_ENDPOINTS.QUIZZES.FEATURED
    ) as unknown as { data: { quizzes: Quiz[] } };
    apiLogger.debug("Featured quizzes fetched successfully", {
      count: response.data.quizzes.length,
    });
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
    apiLogger.debug("Fetching quizzes by category", { categoryId });
    const response = await apiClient.get<{ data: { quizzes: Quiz[] } }>(
      API_ENDPOINTS.QUIZZES.BY_CATEGORY(categoryId)
    ) as unknown as { data: { quizzes: Quiz[] } };
    apiLogger.debug("Category quizzes fetched successfully", {
      categoryId,
      count: response.data.quizzes.length,
    });
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
    apiLogger.debug("Fetching categories");
    const response = await apiClient.get<{ data: Category[] }>(
      API_ENDPOINTS.CATEGORIES.LIST
    ) as unknown as { data: Category[] };
    apiLogger.debug("Categories fetched successfully", {
      count: response.data.length,
    });
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
    apiLogger.debug("Searching quizzes", { query });
    const response = await apiClient.get<{ data: QuizListResponse }>(API_ENDPOINTS.QUIZZES.LIST, {
      params: { search: query },
    }) as unknown as { data: QuizListResponse };
    apiLogger.debug("Quiz search completed", {
      query,
      count: response.data.quizzes.length,
    });
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
    apiLogger.debug("Starting quiz attempt", { quizId });
    const response = await apiClient.post<QuizAttempt>(
      API_ENDPOINTS.QUIZ.START_ATTEMPT(quizId)
    ) as unknown as QuizAttempt;
    apiLogger.debug("Quiz attempt started", {
      quizId,
      attemptId: response.id,
    });
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
    apiLogger.debug("Updating quiz attempt", {
      quizId,
      attemptId,
      answersCount: answers.length,
    });
    const response = await apiClient.put<QuizAttempt>(
      API_ENDPOINTS.QUIZ.UPDATE_ATTEMPT(quizId, attemptId),
      { answers }
    ) as unknown as QuizAttempt;
    apiLogger.debug("Quiz attempt updated", { quizId, attemptId });
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
    apiLogger.debug("Submitting quiz attempt", {
      quizId,
      attemptId,
      answersCount: answers.length,
    });
    const response = await apiClient.post<QuizResults>(
      API_ENDPOINTS.QUIZ.SUBMIT_ATTEMPT(quizId, attemptId),
      { answers }
    ) as unknown as QuizResults;
    apiLogger.debug("Quiz attempt submitted", {
      quizId,
      attemptId,
      score: response.percentage,
    });
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
    apiLogger.debug("Pausing quiz session", { quizId, attemptId, pauseRequest });
    const response = await apiClient.post<{ data: SessionActionResponse }>(
      API_ENDPOINTS.QUIZ.PAUSE(quizId, attemptId),
      pauseRequest
    ) as unknown as { data: SessionActionResponse };
    apiLogger.debug("Quiz session paused", { quizId, attemptId, response: response.data });
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
    apiLogger.debug("Resuming quiz session", { quizId, attemptId });
    const response = await apiClient.post<{ data: ResumeSessionResponse }>(
      API_ENDPOINTS.QUIZ.RESUME(quizId, attemptId)
    ) as unknown as { data: ResumeSessionResponse };
    apiLogger.debug("Quiz session resumed", { quizId, attemptId, response: response.data });
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
  try {
    apiLogger.debug("Saving session progress", {
      quizId,
      attemptId,
      progressRequest,
    });
    await apiClient.put(
      API_ENDPOINTS.QUIZ.SAVE_PROGRESS(quizId, attemptId),
      progressRequest
    );
    apiLogger.debug("Session progress saved", { quizId, attemptId });
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
    apiLogger.debug("Abandoning quiz session", { quizId, attemptId });
    const response = await apiClient.delete<{ data: SessionActionResponse }>(
      API_ENDPOINTS.QUIZ.ABANDON(quizId, attemptId)
    ) as unknown as { data: SessionActionResponse };
    apiLogger.debug("Quiz session abandoned", { quizId, attemptId, response: response.data });
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
    apiLogger.debug("Fetching user active sessions", filters);
    const response = await apiClient.get<{ data: ActiveSessionsResponse }>(
      API_ENDPOINTS.USERS.ACTIVE_SESSIONS,
      { params: filters }
    ) as unknown as { data: ActiveSessionsResponse };
    apiLogger.debug("Active sessions fetched", {
      total: response.data.total,
      active_count: response.data.active_count,
      paused_count: response.data.paused_count
    });
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
    apiLogger.debug("Fetching active session for quiz", { quizId });
    const filters: SessionFilters = {
      quiz_id: quizId,
      session_state: 'active',
      page_size: 1
    };
    const response = await getUserActiveSessions(filters);
    const session = response.sessions.length > 0 ? response.sessions[0] : null;
    apiLogger.debug("Active session for quiz", { quizId, found: !!session });
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
    apiLogger.debug("Fetching user quiz attempts");
    const response = await apiClient.get<QuizAttempt[]>(
      API_ENDPOINTS.USERS.ATTEMPTS
    ) as unknown as QuizAttempt[];
    apiLogger.debug("User attempts fetched", { count: response.length });
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
    apiLogger.debug("Fetching attempt details", { attemptId });
    const response = await apiClient.get<QuizAttempt>(
      API_ENDPOINTS.USERS.ATTEMPT_DETAILS(attemptId)
    ) as unknown as QuizAttempt;
    apiLogger.debug("Attempt details fetched", { attemptId });
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
    apiLogger.debug("Fetching active sessions");
    const response = await apiClient.get<QuizAttempt[]>(
      API_ENDPOINTS.USERS.ACTIVE_SESSIONS
    ) as unknown as QuizAttempt[];
    apiLogger.debug("Active sessions fetched", { count: response.length });
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
    apiLogger.debug("Fetching active attempt for quiz", { quizId });
    const response = await apiClient.get<QuizAttempt>(
      API_ENDPOINTS.USERS.QUIZ_ATTEMPT(quizId)
    ) as unknown as QuizAttempt;
    apiLogger.debug("Active attempt fetched", { quizId, attemptId: response?.id });
    return response;
  } catch (error: any) {
    // If 404, no active attempt exists
    if (error.status === 404) {
      apiLogger.debug("No active attempt for quiz", { quizId });
      return null;
    }
    apiLogger.error("Error fetching active attempt", { quizId, error });
    throw error;
  }
}
