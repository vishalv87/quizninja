-- Migration: Remove is_test_data fields from leaderboard-related tables
-- Date: 2025-11-04
-- Description: Removes is_test_data boolean fields from leaderboard_snapshots, user_rank_history, and user_category_performance tables

-- Drop the is_test_data column from leaderboard_snapshots table
ALTER TABLE leaderboard_snapshots DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from user_rank_history table
ALTER TABLE user_rank_history DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from user_category_performance table
ALTER TABLE user_category_performance DROP COLUMN IF EXISTS is_test_data;

-- Add comments to document this change
COMMENT ON TABLE leaderboard_snapshots IS 'Stores historical leaderboard snapshots for tracking rank changes over time. is_test_data field removed as part of cleanup.';
COMMENT ON TABLE user_rank_history IS 'Tracks user rank changes across different time periods. is_test_data field removed as part of cleanup.';
COMMENT ON TABLE user_category_performance IS 'Stores user performance metrics by quiz category. is_test_data field removed as part of cleanup.';