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
  difficulty: "easy" | "medium" | "hard";
  question_count: number;
  time_limit_minutes?: number;
  points_per_question: number;
  total_points: number;
  is_featured: boolean;
  created_at: string;
  updated_at: string;
  attempts_count?: number;
  average_score?: number;
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
  options?: QuestionOption[];
  correct_answer?: string;
  explanation?: string;
}

export interface QuestionOption {
  id: string;
  option_text: string;
  is_correct: boolean;
}

export interface QuizAttempt {
  id: string;
  quiz_id: string;
  user_id: string;
  status: "in_progress" | "completed" | "abandoned";
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
