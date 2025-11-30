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
  name: string;
  avatar?: string;
  points: number;
  quizzes_completed: number;
  average_score: number;
  current_streak: number;
  level: string;
  is_current_user: boolean;
  is_friend: boolean;
  last_active: string;
  achievements: string[];
  category_points: Record<string, number>;
}

export interface Category {
  id: string;
  name: string;
  description?: string;
  icon?: string;
  quiz_count?: number;
}