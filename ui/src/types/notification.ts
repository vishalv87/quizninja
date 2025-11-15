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
