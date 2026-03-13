/**
 * Human-readable display labels for enums
 * Used in UI dropdowns, forms, and display text
 */

import type { QuizDifficulty } from './enums/quiz';
import type { NotificationType } from './enums/notification';
import type { Theme, ProfileVisibility, NotificationFrequency } from './enums/user';
import type { DiscussionSort, DiscussionType } from './enums/discussion';
import type { AchievementCategory, AchievementRequirementType } from './enums/achievement';

// Quiz Difficulty Labels
export const DIFFICULTY_LABELS: Record<QuizDifficulty, string> = {
  beginner: 'Beginner',
  intermediate: 'Intermediate',
  advanced: 'Advanced',
};

// Difficulty options for dropdowns (value/label pairs)
export const DIFFICULTY_OPTIONS = [
  { value: 'beginner', label: 'Beginner' },
  { value: 'intermediate', label: 'Intermediate' },
  { value: 'advanced', label: 'Advanced' },
] as const;

// Notification Type Labels
export const NOTIFICATION_TYPE_LABELS: Record<NotificationType, string> = {
  friend_request: 'Friend Request',
  friend_accepted: 'Friend Accepted',
  achievement_unlocked: 'Achievement Unlocked',
  quiz_reminder: 'Quiz Reminder',
  discussion_reply: 'Discussion Reply',
  system: 'System',
};

// Theme Labels
export const THEME_LABELS: Record<Theme, string> = {
  light: 'Light',
  dark: 'Dark',
  system: 'System',
};

// Theme options for dropdowns
export const THEME_OPTIONS = [
  { value: 'light', label: 'Light' },
  { value: 'dark', label: 'Dark' },
  { value: 'system', label: 'System' },
] as const;

// Profile Visibility Labels
export const PROFILE_VISIBILITY_LABELS: Record<ProfileVisibility, string> = {
  public: 'Public',
  friends_only: 'Friends Only',
  private: 'Private',
};

// Profile visibility options for dropdowns
export const PROFILE_VISIBILITY_OPTIONS = [
  { value: 'public', label: 'Public' },
  { value: 'friends_only', label: 'Friends Only' },
  { value: 'private', label: 'Private' },
] as const;

// Notification Frequency Labels
export const NOTIFICATION_FREQUENCY_LABELS: Record<NotificationFrequency, string> = {
  instant: 'Instant',
  daily: 'Daily',
  weekly: 'Weekly',
  never: 'Never',
};

// Notification frequency options for dropdowns
export const NOTIFICATION_FREQUENCY_OPTIONS = [
  { value: 'instant', label: 'Instant' },
  { value: 'daily', label: 'Daily' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'never', label: 'Never' },
] as const;

// Discussion Sort Labels
export const DISCUSSION_SORT_LABELS: Record<DiscussionSort, string> = {
  recent: 'Most Recent',
  popular: 'Most Popular',
};

// Discussion sort options for dropdowns
export const DISCUSSION_SORT_OPTIONS = [
  { value: 'recent', label: 'Most Recent' },
  { value: 'popular', label: 'Most Popular' },
] as const;

// Discussion Type Labels
export const DISCUSSION_TYPE_LABELS: Record<DiscussionType, string> = {
  question: 'Question',
  general: 'General',
  bug_report: 'Bug Report',
  feature_request: 'Feature Request',
  discussion: 'Discussion',
};

// Discussion type options for dropdowns
export const DISCUSSION_TYPE_OPTIONS = [
  { value: 'question', label: 'Question' },
  { value: 'general', label: 'General' },
  { value: 'bug_report', label: 'Bug Report' },
  { value: 'feature_request', label: 'Feature Request' },
  { value: 'discussion', label: 'Discussion' },
] as const;

// Achievement Category Labels
export const ACHIEVEMENT_CATEGORY_LABELS: Record<AchievementCategory, string> = {
  quiz_master: 'Quiz Master',
  social: 'Social',
  streak: 'Streak',
  knowledge: 'Knowledge',
  competitor: 'Competitor',
};

// Achievement category options for dropdowns
export const ACHIEVEMENT_CATEGORY_OPTIONS = [
  { value: 'quiz_master', label: 'Quiz Master' },
  { value: 'social', label: 'Social' },
  { value: 'streak', label: 'Streak' },
  { value: 'knowledge', label: 'Knowledge' },
  { value: 'competitor', label: 'Competitor' },
] as const;

// Achievement Requirement Type Labels
export const ACHIEVEMENT_REQUIREMENT_TYPE_LABELS: Record<AchievementRequirementType, string> = {
  quizzes_completed: 'Quizzes Completed',
  total_points: 'Total Points',
  accuracy_percentage: 'Accuracy Percentage',
  streak_reached: 'Streak Reached',
  friends_added: 'Friends Added',
  discussions_started: 'Discussions Started',
};

// Achievement requirement type options for dropdowns
export const ACHIEVEMENT_REQUIREMENT_TYPE_OPTIONS = [
  { value: 'quizzes_completed', label: 'Quizzes Completed' },
  { value: 'total_points', label: 'Total Points' },
  { value: 'accuracy_percentage', label: 'Accuracy Percentage' },
  { value: 'streak_reached', label: 'Streak Reached' },
  { value: 'friends_added', label: 'Friends Added' },
  { value: 'discussions_started', label: 'Discussions Started' },
] as const;
