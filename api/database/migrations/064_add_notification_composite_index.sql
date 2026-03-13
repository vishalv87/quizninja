-- Add composite index optimizing all common notification query patterns:
--   WHERE user_id = ? AND is_deleted = false [AND is_read = ?] ORDER BY created_at DESC
-- This eliminates sequential scans and reduces query time for GetNotifications,
-- GetNotificationStats, and GetUnreadNotificationCount.
CREATE INDEX IF NOT EXISTS idx_notifications_user_undeleted_read_created
ON notifications (user_id, is_deleted, is_read, created_at DESC);
