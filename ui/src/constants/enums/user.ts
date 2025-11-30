/**
 * User-related enums and constants
 * Single source of truth for theme, profile visibility, friend request status, and notification frequency
 */

// Theme Preferences
export const Theme = {
  LIGHT: 'light',
  DARK: 'dark',
  SYSTEM: 'system',
} as const;

export type Theme = typeof Theme[keyof typeof Theme];
export const THEMES = Object.values(Theme);

// Type guard for theme
export function isTheme(value: unknown): value is Theme {
  return typeof value === 'string' && THEMES.includes(value as Theme);
}

// Profile Visibility Settings
export const ProfileVisibility = {
  PUBLIC: 'public',
  FRIENDS_ONLY: 'friends_only',
  PRIVATE: 'private',
} as const;

export type ProfileVisibility = typeof ProfileVisibility[keyof typeof ProfileVisibility];
export const PROFILE_VISIBILITIES = Object.values(ProfileVisibility);

// Type guard for profile visibility
export function isProfileVisibility(value: unknown): value is ProfileVisibility {
  return typeof value === 'string' && PROFILE_VISIBILITIES.includes(value as ProfileVisibility);
}

// Friend Request Status
export const FriendRequestStatus = {
  PENDING: 'pending',
  ACCEPTED: 'accepted',
  REJECTED: 'rejected',
} as const;

export type FriendRequestStatus = typeof FriendRequestStatus[keyof typeof FriendRequestStatus];
export const FRIEND_REQUEST_STATUSES = Object.values(FriendRequestStatus);

// Type guard for friend request status
export function isFriendRequestStatus(value: unknown): value is FriendRequestStatus {
  return typeof value === 'string' && FRIEND_REQUEST_STATUSES.includes(value as FriendRequestStatus);
}

// Friendship Status (for UserProfile)
export const FriendshipStatus = {
  NONE: 'none',
  PENDING_SENT: 'pending_sent',
  PENDING_RECEIVED: 'pending_received',
  FRIENDS: 'friends',
} as const;

export type FriendshipStatus = typeof FriendshipStatus[keyof typeof FriendshipStatus];
export const FRIENDSHIP_STATUSES = Object.values(FriendshipStatus);

// Type guard for friendship status
export function isFriendshipStatus(value: unknown): value is FriendshipStatus {
  return typeof value === 'string' && FRIENDSHIP_STATUSES.includes(value as FriendshipStatus);
}

// Notification Frequency
export const NotificationFrequency = {
  INSTANT: 'instant',
  DAILY: 'daily',
  WEEKLY: 'weekly',
  NEVER: 'never',
} as const;

export type NotificationFrequency = typeof NotificationFrequency[keyof typeof NotificationFrequency];
export const NOTIFICATION_FREQUENCIES = Object.values(NotificationFrequency);

// Type guard for notification frequency
export function isNotificationFrequency(value: unknown): value is NotificationFrequency {
  return typeof value === 'string' && NOTIFICATION_FREQUENCIES.includes(value as NotificationFrequency);
}

// Friend Request Actions (accept/decline)
export const FriendRequestAction = {
  ACCEPT: 'accept',
  DECLINE: 'decline',
} as const;

export type FriendRequestAction = typeof FriendRequestAction[keyof typeof FriendRequestAction];
export const FRIEND_REQUEST_ACTIONS = Object.values(FriendRequestAction);

// Type guard for friend request action
export function isFriendRequestAction(value: unknown): value is FriendRequestAction {
  return typeof value === 'string' && FRIEND_REQUEST_ACTIONS.includes(value as FriendRequestAction);
}
