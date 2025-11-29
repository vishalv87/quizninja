export interface UserPreferences {
  user_id: string;
  preferred_categories: string[];
  preferred_difficulty: string;
  notification_frequency: string;
  email_notifications: boolean;
  onboarding_completed: boolean;
  theme: "light" | "dark" | "system";
  // Privacy settings
  profile_visibility?: "public" | "friends_only" | "private";
  show_achievements?: boolean;
  show_stats?: boolean;
  allow_friend_requests?: boolean;
  created_at: string;
  updated_at: string;
}

export interface UserStats {
  user_id: string;
  total_quizzes_taken: number;
  total_quizzes_completed: number;
  total_points: number;
  average_score: number;
  total_time_spent_minutes: number;
  current_streak: number;
  longest_streak: number;
  achievements_unlocked: number;
  rank: number;
}

export interface Friend {
  id: string;
  name: string;
  email: string;
  avatar_url?: string;
  level: number;
  total_points: number;
  current_streak: number;
  best_streak: number;
  total_quizzes_completed: number;
  average_score: number;
  is_online: boolean;
  last_active: string;
  friends_since: string;
}

export interface FriendRequest {
  id: string;
  requester_id: string;
  requested_id: string;
  status: "pending" | "accepted" | "rejected";
  message?: string;
  requester: {
    id: string;
    full_name: string;
    avatar_url?: string;
  };
  requested: {
    id: string;
    full_name: string;
    avatar_url?: string;
  };
  created_at: string;
  responded_at?: string;
}

export interface UserProfile {
  id: string;
  user_id: string;
  name?: string;
  email: string;
  full_name?: string;
  avatar_url?: string;
  bio?: string;
  created_at: string;
  updated_at: string;
  // Stats (privacy-aware, may be null if hidden)
  stats?: UserStats | null;
  // Privacy settings
  preferences?: Pick<UserPreferences, 'profile_visibility' | 'show_achievements' | 'show_stats'>;
  // Friendship status
  is_friend?: boolean;
  friend_request_status?: 'none' | 'pending_sent' | 'pending_received' | 'friends';
}
