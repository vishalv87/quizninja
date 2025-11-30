import type { AchievementCategory, AchievementRequirementType } from '@/constants';

export interface Achievement {
  id: string;
  key: string;
  name: string;
  description: string;
  icon: string;
  category: AchievementCategory;
  points: number;
  requirement_type: AchievementRequirementType;
  requirement_value: number;
  is_secret: boolean;
  created_at: string;
}

export interface UserAchievement {
  id: string;
  user_id: string;
  achievement_id: string;
  unlocked_at: string;
  achievement: Achievement;
}

export interface AchievementProgress {
  achievement_id: string;
  achievement: Achievement;
  current_value: number;
  target_value: number;
  progress_percentage: number;
  is_unlocked: boolean;
  unlocked_at?: string;
}

export interface AchievementStats {
  total_achievements: number;
  unlocked_achievements: number;
  total_points: number;
  points_earned: number;
  completion_percentage: number;
}
