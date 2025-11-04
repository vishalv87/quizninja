-- Migration: Remove is_test_data field from notifications table
-- Date: 2025-11-04
-- Description: Removes is_test_data boolean field and its index from notifications table

-- Drop the index on is_test_data column
DROP INDEX IF EXISTS idx_notifications_is_test_data;

-- Drop the is_test_data column from notifications table
ALTER TABLE notifications DROP COLUMN IF EXISTS is_test_data;

-- Add comment to document this change
COMMENT ON TABLE notifications IS 'Stores unified notifications for all notification types (friend requests, challenges, achievements, etc.). is_test_data field removed as part of cleanup.';
