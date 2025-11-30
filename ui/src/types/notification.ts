import { type NotificationType as NotificationTypeEnum } from '@/constants';

export interface Notification {
  id: string;
  user_id: string;
  type: NotificationTypeEnum;
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
  type?: NotificationTypeEnum;
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

// Re-export NotificationType and type guard from constants for backwards compatibility
export { type NotificationType, isNotificationType as isValidNotificationType } from '@/constants';

// Helper type for notifications with additional computed properties
export interface NotificationWithMeta extends Notification {
  isExpired?: boolean;
  timeSinceCreated?: string;
  actionUrl?: string;
}
