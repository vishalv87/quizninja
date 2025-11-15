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
  challenges_won: number;
  challenges_lost: number;
  rank: number;
}

export interface Friend {
  id: string;
  user_id: string;
  friend_user_id: string;
  friend: {
    id: string;
    full_name: string;
    avatar_url?: string;
    email?: string;
  };
  created_at: string;
}

export interface FriendRequest {
  id: string;
  sender_id: string;
  receiver_id: string;
  status: "pending" | "accepted" | "declined";
  sender: {
    id: string;
    full_name: string;
    avatar_url?: string;
  };
  created_at: string;
  updated_at: string;
}
