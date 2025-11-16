/**
 * Common API response types
 */

export interface APIResponse<T = any> {
  data: T;
  message?: string;
  error?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
  has_more: boolean;
}

/**
 * Backend response type for attempt history
 * Matches the AttemptHistoryResponse struct from the Go backend
 */
export interface AttemptHistoryResponse<T = any> {
  attempts: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface APIError {
  error: string;
  message: string;
  status: number;
  details?: Record<string, any>;
}

export interface LeaderboardEntry {
  rank: number;
  user_id: string;
  user: {
    id: string;
    full_name: string;
    avatar_url?: string;
  };
  total_points: number;
  quizzes_completed: number;
  achievements_unlocked: number;
}

export interface Category {
  id: string;
  name: string;
  description?: string;
  icon?: string;
  quiz_count?: number;
}