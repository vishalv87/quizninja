export interface Notification {
  id: string;
  user_id: string;
  type: string;
  title: string;
  message: string;
  data?: Record<string, any>;
  is_read: boolean;
  created_at: string;
  expires_at?: string;
}

export interface NotificationStats {
  total_notifications: number;
  unread_notifications: number;
  read_notifications: number;
}

export interface NotificationFilter {
  type?: string;
  is_read?: boolean;
  limit?: number;
  offset?: number;
}

export interface NotificationListResponse {
  notifications: Notification[];
  total: number;
  page: number;
  limit: number;
  hasMore: boolean;
}

// Notification type constants for type-safe filtering
export type NotificationType =
  | 'friend_request'
  | 'friend_accepted'
  | 'achievement_unlocked'
  | 'quiz_reminder'
  | 'discussion_reply'
  | 'system';

// Type guard to check if a notification type is valid
export function isValidNotificationType(type: string): type is NotificationType {
  return [
    'friend_request',
    'friend_accepted',
    'achievement_unlocked',
    'quiz_reminder',
    'discussion_reply',
    'system',
  ].includes(type);
}

// Helper type for notifications with additional computed properties
export interface NotificationWithMeta extends Notification {
  isExpired?: boolean;
  timeSinceCreated?: string;
  actionUrl?: string;
}
