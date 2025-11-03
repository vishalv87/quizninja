-- Migration: Remove is_test_data columns from preferences-related tables
-- Date: 2025-11-03

-- Drop the is_test_data column from difficulty_levels table
ALTER TABLE difficulty_levels DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from notification_frequencies table
ALTER TABLE notification_frequencies DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from user_preferences table
ALTER TABLE user_preferences DROP COLUMN IF EXISTS is_test_data;
