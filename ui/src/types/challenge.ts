export interface Challenge {
  id: string;
  challenger_id: string;
  opponent_id: string;
  quiz_id: string;
  status: "pending" | "accepted" | "declined" | "completed" | "expired";
  challenger_score?: number;
  opponent_score?: number;
  challenger_attempt_id?: string;
  opponent_attempt_id?: string;
  winner_id?: string;
  expires_at: string;
  created_at: string;
  updated_at: string;
  completed_at?: string;
  challenger: {
    id: string;
    full_name: string;
    avatar_url?: string;
  };
  opponent: {
    id: string;
    full_name: string;
    avatar_url?: string;
  };
  quiz: {
    id: string;
    title: string;
    difficulty: string;
    category: string;
  };
}

export interface CreateChallengeRequest {
  opponent_id: string;
  quiz_id: string;
}

export interface ChallengeStats {
  total_challenges: number;
  pending_challenges: number;
  active_challenges: number;
  completed_challenges: number;
  won_challenges: number;
  lost_challenges: number;
  win_rate: number;
}
