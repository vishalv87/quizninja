-- Migration: Add soft delete capability to notifications
-- This migration adds soft delete fields to the notifications table
-- allowing notifications to be hidden from users while preserving data

-- Add soft delete fields to notifications table
ALTER TABLE notifications
ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN deleted_at TIMESTAMP NULL;

-- Create indexes for better performance on soft delete queries
CREATE INDEX idx_notifications_is_deleted ON notifications(is_deleted);
CREATE INDEX idx_notifications_deleted_at ON notifications(deleted_at);

-- Create composite index for user_id and is_deleted for efficient user notification queries
CREATE INDEX idx_notifications_user_id_is_deleted ON notifications(user_id, is_deleted);

-- Update existing data to ensure all current notifications are marked as not deleted
UPDATE notifications SET is_deleted = FALSE WHERE is_deleted IS NULL;

-- Add comment for documentation
COMMENT ON COLUMN notifications.is_deleted IS 'Soft delete flag - when TRUE, notification is hidden from user but preserved in database';
COMMENT ON COLUMN notifications.deleted_at IS 'Timestamp when notification was soft deleted, NULL if not deleted';

-- Create function to handle soft delete of notifications
CREATE OR REPLACE FUNCTION soft_delete_notification(
    p_notification_id UUID,
    p_user_id UUID
)
RETURNS BOOLEAN AS $$
DECLARE
    affected_rows INTEGER;
BEGIN
    UPDATE notifications
    SET is_deleted = TRUE, deleted_at = CURRENT_TIMESTAMP
    WHERE id = p_notification_id
      AND user_id = p_user_id
      AND is_deleted = FALSE;

    GET DIAGNOSTICS affected_rows = ROW_COUNT;

    IF affected_rows > 0 THEN
        RETURN TRUE;
    ELSE
        RETURN FALSE;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create function to restore soft deleted notification (for potential undo functionality)
CREATE OR REPLACE FUNCTION restore_notification(
    p_notification_id UUID,
    p_user_id UUID
)
RETURNS BOOLEAN AS $$
DECLARE
    affected_rows INTEGER;
BEGIN
    UPDATE notifications
    SET is_deleted = FALSE, deleted_at = NULL
    WHERE id = p_notification_id
      AND user_id = p_user_id
      AND is_deleted = TRUE;

    GET DIAGNOSTICS affected_rows = ROW_COUNT;

    IF affected_rows > 0 THEN
        RETURN TRUE;
    ELSE
        RETURN FALSE;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create function for hard delete cleanup (admin use)
CREATE OR REPLACE FUNCTION cleanup_old_deleted_notifications(
    p_days_old INTEGER DEFAULT 30
)
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM notifications
    WHERE is_deleted = TRUE
      AND deleted_at < (CURRENT_TIMESTAMP - INTERVAL '1 day' * p_days_old);

    GET DIAGNOSTICS deleted_count = ROW_COUNT;

    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Add constraint to ensure deleted_at is set when is_deleted is true
ALTER TABLE notifications
ADD CONSTRAINT check_deleted_at_when_deleted
CHECK (
    (is_deleted = FALSE AND deleted_at IS NULL) OR
    (is_deleted = TRUE AND deleted_at IS NOT NULL)
);