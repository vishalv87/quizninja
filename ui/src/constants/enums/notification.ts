/**
 * Notification-related enums and constants
 * Single source of truth for notification types
 */

// Notification Types
export const NotificationType = {
  FRIEND_REQUEST: 'friend_request',
  FRIEND_ACCEPTED: 'friend_accepted',
  ACHIEVEMENT_UNLOCKED: 'achievement_unlocked',
  QUIZ_REMINDER: 'quiz_reminder',
  DISCUSSION_REPLY: 'discussion_reply',
  SYSTEM: 'system',
} as const;

export type NotificationType = typeof NotificationType[keyof typeof NotificationType];
export const NOTIFICATION_TYPES = Object.values(NotificationType);

// Type guard for notification type
export function isNotificationType(value: unknown): value is NotificationType {
  return typeof value === 'string' && NOTIFICATION_TYPES.includes(value as NotificationType);
}

// Notification Actions (mark as read/unread)
export const NotificationAction = {
  READ: 'read',
  UNREAD: 'unread',
} as const;

export type NotificationAction = typeof NotificationAction[keyof typeof NotificationAction];
export const NOTIFICATION_ACTIONS = Object.values(NotificationAction);

// Type guard for notification action
export function isNotificationAction(value: unknown): value is NotificationAction {
  return typeof value === 'string' && NOTIFICATION_ACTIONS.includes(value as NotificationAction);
}
