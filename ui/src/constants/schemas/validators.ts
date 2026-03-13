/**
 * Zod schema validators derived from enums
 * Single source of truth - these should be used in all Zod schemas
 */

import { z } from 'zod';
import {
  QUIZ_DIFFICULTIES,
  QUESTION_TYPES,
  ATTEMPT_STATUSES,
  type QuizDifficulty,
  type QuestionType,
  type AttemptStatus,
} from '../enums/quiz';
import {
  THEMES,
  PROFILE_VISIBILITIES,
  FRIEND_REQUEST_STATUSES,
  FRIENDSHIP_STATUSES,
  NOTIFICATION_FREQUENCIES,
  type Theme,
  type ProfileVisibility,
  type FriendRequestStatus,
  type FriendshipStatus,
  type NotificationFrequency,
} from '../enums/user';
import {
  NOTIFICATION_TYPES,
  NOTIFICATION_ACTIONS,
  type NotificationType,
  type NotificationAction,
} from '../enums/notification';
import {
  DISCUSSION_SORTS,
  DISCUSSION_TYPES,
  type DiscussionSort,
  type DiscussionType,
} from '../enums/discussion';
import {
  ACHIEVEMENT_CATEGORIES,
  ACHIEVEMENT_REQUIREMENT_TYPES,
  type AchievementCategory,
  type AchievementRequirementType,
} from '../enums/achievement';
import {
  LEADERBOARD_PERIODS,
  LEADERBOARD_SORT_FIELDS,
  type LeaderboardPeriod,
  type LeaderboardSortField,
} from '../enums/leaderboard';
import {
  SORT_ORDERS,
  type SortOrder,
} from '../enums/common';
import {
  FRIEND_REQUEST_ACTIONS,
  type FriendRequestAction,
} from '../enums/user';

// Quiz-related validators
export const quizDifficultySchema = z.enum(
  QUIZ_DIFFICULTIES as [QuizDifficulty, ...QuizDifficulty[]]
);

export const questionTypeSchema = z.enum(
  QUESTION_TYPES as [QuestionType, ...QuestionType[]]
);

export const attemptStatusSchema = z.enum(
  ATTEMPT_STATUSES as [AttemptStatus, ...AttemptStatus[]]
);

// User-related validators
export const themeSchema = z.enum(
  THEMES as [Theme, ...Theme[]]
);

export const profileVisibilitySchema = z.enum(
  PROFILE_VISIBILITIES as [ProfileVisibility, ...ProfileVisibility[]]
);

export const friendRequestStatusSchema = z.enum(
  FRIEND_REQUEST_STATUSES as [FriendRequestStatus, ...FriendRequestStatus[]]
);

export const friendshipStatusSchema = z.enum(
  FRIENDSHIP_STATUSES as [FriendshipStatus, ...FriendshipStatus[]]
);

export const notificationFrequencySchema = z.enum(
  NOTIFICATION_FREQUENCIES as [NotificationFrequency, ...NotificationFrequency[]]
);

// Notification validators
export const notificationTypeSchema = z.enum(
  NOTIFICATION_TYPES as [NotificationType, ...NotificationType[]]
);

// Discussion validators
export const discussionSortSchema = z.enum(
  DISCUSSION_SORTS as [DiscussionSort, ...DiscussionSort[]]
);

// Extended validators with custom error messages
export const quizDifficultySchemaWithError = z.enum(
  QUIZ_DIFFICULTIES as [QuizDifficulty, ...QuizDifficulty[]],
  {
    errorMap: () => ({ message: 'Please select a valid difficulty level' }),
  }
);

export const themeSchemaWithError = z.enum(
  THEMES as [Theme, ...Theme[]],
  {
    errorMap: () => ({ message: 'Please select a valid theme' }),
  }
);

export const profileVisibilitySchemaWithError = z.enum(
  PROFILE_VISIBILITIES as [ProfileVisibility, ...ProfileVisibility[]],
  {
    errorMap: () => ({ message: 'Please select a valid visibility setting' }),
  }
);

// Leaderboard validators
export const leaderboardPeriodSchema = z.enum(
  LEADERBOARD_PERIODS as [LeaderboardPeriod, ...LeaderboardPeriod[]]
);

export const leaderboardSortFieldSchema = z.enum(
  LEADERBOARD_SORT_FIELDS as [LeaderboardSortField, ...LeaderboardSortField[]]
);

// Common validators
export const sortOrderSchema = z.enum(
  SORT_ORDERS as [SortOrder, ...SortOrder[]]
);

// Notification action validator
export const notificationActionSchema = z.enum(
  NOTIFICATION_ACTIONS as [NotificationAction, ...NotificationAction[]]
);

// Friend request action validator
export const friendRequestActionSchema = z.enum(
  FRIEND_REQUEST_ACTIONS as [FriendRequestAction, ...FriendRequestAction[]]
);

// Achievement validators
export const achievementCategorySchema = z.enum(
  ACHIEVEMENT_CATEGORIES as [AchievementCategory, ...AchievementCategory[]]
);

export const achievementRequirementTypeSchema = z.enum(
  ACHIEVEMENT_REQUIREMENT_TYPES as [AchievementRequirementType, ...AchievementRequirementType[]]
);

// Discussion type validator
export const discussionTypeSchema = z.enum(
  DISCUSSION_TYPES as [DiscussionType, ...DiscussionType[]]
);
