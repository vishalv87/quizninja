-- Migration: Remove is_test_data fields from friends-related tables
-- Date: 2025-11-03
-- Description: Removes is_test_data boolean fields from friend_requests and friendships tables

-- Drop the is_test_data column from friend_requests table
ALTER TABLE friend_requests DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from friendships table
ALTER TABLE friendships DROP COLUMN IF EXISTS is_test_data;

-- Add comments to document this change
COMMENT ON TABLE friend_requests IS 'Stores friend requests between users. is_test_data field removed as part of cleanup.';
COMMENT ON TABLE friendships IS 'Stores accepted friendships between users. is_test_data field removed as part of cleanup.';
