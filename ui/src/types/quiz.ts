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
  status: "in_progress" | "completed" | "abandoned" | "paused";
  score?: number;
  total_questions: number;
  correct_answers: number;
  started_at: string;
  completed_at?: string;
  time_spent_seconds?: number;
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

// Quiz Session Types for Pause/Resume functionality
export type SessionState = 'active' | 'paused' | 'completed' | 'abandoned';

export interface AttemptAnswer {
  question_id: string;
  selected_answer: string;
  selected_option_index?: number; // Index in the options array (used for submission)
  is_correct?: boolean;
  points_earned?: number;
}

export interface QuizSession {
  id: string;
  attempt_id: string;
  user_id: string;
  quiz_id: string;
  current_question_index: number;
  current_answers: AttemptAnswer[];
  session_state: SessionState;
  time_remaining?: number; // in seconds
  time_spent_so_far: number; // in seconds
  last_activity_at: string;
  paused_at?: string;
  created_at: string;
  updated_at: string;
  quiz?: Quiz; // Optional quiz details
}

export interface QuizSessionWithDetails extends QuizSession {
  quiz_title: string;
  quiz_category: string;
  quiz_difficulty: string;
  total_questions: number;
  original_time_limit?: number; // in seconds
  progress: number; // percentage 0-100
}

export interface PauseSessionRequest {
  current_question_index: number;
  current_answers: AttemptAnswer[];
  time_spent_so_far: number;
  time_remaining?: number;
}

export interface SaveProgressRequest {
  current_question_index: number;
  current_answers: AttemptAnswer[];
  time_spent_so_far: number;
  time_remaining?: number;
}

export interface SessionActionResponse {
  session_id: string;
  action: string;
  session_state: SessionState;
  message: string;
  time_remaining?: number;
  progress?: number;
}

export interface ResumeSessionResponse {
  session_id: string;
  action: string;
  session_state: SessionState;
  message: string;
  quiz: Quiz & { questions: Question[] };
  current_question_index: number;
  current_answers: AttemptAnswer[];
  time_remaining?: number;
  time_spent_so_far: number;
  progress?: number;
}

export interface ActiveSessionsResponse {
  sessions: QuizSessionWithDetails[];
  total: number;
  active_count: number;
  paused_count: number;
}

export interface SessionFilters {
  session_state?: SessionState;
  quiz_id?: string;
  category?: string;
  difficulty?: string;
  page?: number;
  page_size?: number;
  sort_by?: 'last_activity_at' | 'created_at';
  sort_order?: 'asc' | 'desc';
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
