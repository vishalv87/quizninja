-- Migration: Remove is_test_data fields from achievements and user_achievements tables
-- Created: 2025-11-02
-- Description: Removes the is_test_data column from achievements and user_achievements tables
--              and drops the associated index on achievements table

-- Drop the index first (indexes must be dropped before columns)
DROP INDEX IF EXISTS idx_achievements_is_test_data;

-- Remove is_test_data column from achievements table
ALTER TABLE achievements DROP COLUMN IF EXISTS is_test_data;

-- Remove is_test_data column from user_achievements table
ALTER TABLE user_achievements DROP COLUMN IF EXISTS is_test_data;
