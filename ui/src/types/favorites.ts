/**
 * Type definitions for favorites API responses
 * Matches the backend API structure from quizninja-api
 */

export interface QuizSummary {
  id: string;
  title: string;
  description: string;
  category: string;
  difficulty: string;
  time_limit: number;
  question_count: number;
  points: number;
  is_featured: boolean;
  tags: string[];
  thumbnail_url?: string;
  created_at: string;
}

export interface UserQuizFavorite {
  id: string;
  user_id: string;
  quiz_id: string;
  favorited_at: string;
  quiz: QuizSummary;
}

export interface FavoritesListResponse {
  favorites: UserQuizFavorite[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

/**
 * API response wrapper
 * The backend wraps all responses in a "data" field
 */
export interface FavoritesApiResponse {
  data: FavoritesListResponse;
}
