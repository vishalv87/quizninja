export interface Category {
  id: string;
  name: string;
  display_name: string;
  description: string;
  icon_url: string;
  is_active: boolean;
  quiz_count: number;
  created_at: string;
  updated_at: string;
}

export interface Quiz {
  id: string;
  title: string;
  description: string;
  category: string;
  difficulty: "beginner" | "intermediate" | "advanced";
  question_count: number;
  time_limit?: number; // Changed from time_limit_minutes to match backend
  points: number; // Changed from points_per_question to match backend
  is_featured: boolean;
  created_at: string;
  updated_at: string;
  created_by?: string; // UUID of user who created the quiz
  tags?: string[]; // Array of tags for categorization
  thumbnail_url?: string; // Optional thumbnail image for the quiz
  attempts_count?: number;
  average_score?: number;
  article_summary?: string;
  questions?: Question[];
  user_best_score?: number;
  user_has_attempted?: boolean;
  average_rating?: number;
  total_ratings?: number;
  statistics?: QuizStatistics;
}

export interface QuizStatistics {
  total_attempts: number;
  completed_attempts: number;
  average_score: number;
  average_time: number;
  completion_rate: number;
}

export interface QuizListResponse {
  quizzes: Quiz[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface Question {
  id: string;
  quiz_id: string;
  question_text: string;
  question_type: "multiple_choice" | "true_false" | "short_answer";
  points: number;
  order_index: number;
  options?: string[]; // Array of option strings (e.g., ["Paris", "London", "Berlin"])
  correct_answer?: string;
  explanation?: string;
  image_url?: string; // Optional image URL for questions (diagrams, charts, photos)
}

export interface QuizAttempt {
  id: string;
  quiz_id: string;
  user_id: string;
  status: "in_progress" | "completed" | "abandoned";
  score: number; // Number of correct answers
  total_points: number; // Total number of questions
  started_at: string;
  completed_at?: string;
  time_spent?: number; // Time spent in seconds
  percentage_score?: number; // Pre-calculated percentage from backend
  passed?: boolean; // Whether the user passed (>= 60%)
  answers: QuizAnswer[];
}

export interface QuizAnswer {
  question_id: string;
  selected_answer: string;
  selected_option_index?: number; // Index in the options array (used for submission)
  is_correct: boolean;
  points_earned: number;
  time_spent_seconds?: number;
}

export interface QuizFilters {
  category?: string;
  difficulty?: string;
  search?: string;
  is_featured?: boolean;
  limit?: number;
  offset?: number;
}

export interface QuizResults {
  attempt: QuizAttempt;
  quiz: Quiz;
  percentage: number;
  passed: boolean;
  rank?: number;
}

// Answer type for quiz taking (before submission/grading)
export interface AttemptAnswer {
  question_id: string;
  selected_answer: string;
  selected_option_index?: number; // Index in the options array (used for submission)
  is_correct?: boolean;
  points_earned?: number;
}

// Rating Types
export interface QuizRating {
  id: string;
  quiz_id: string;
  user_id: string;
  user_name?: string;
  rating: number; // 1-5
  review?: string;
  created_at: string;
  updated_at?: string;
}

export interface RatingListResponse {
  ratings: QuizRating[];
  average_rating: number;
  total_ratings: number;
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface AverageRatingResponse {
  quiz_id: string;
  average_rating: number;
  total_ratings: number;
}

export interface CreateRatingRequest {
  rating: number;
  review?: string;
}

export interface UpdateRatingRequest {
  rating?: number;
  review?: string;
}
