-- Migration: Remove remaining is_test_data columns and indexes
-- Description: Removes is_test_data column from quiz_ratings and users tables, and associated index

-- Remove is_test_data from quiz_ratings table
ALTER TABLE quiz_ratings DROP COLUMN IF EXISTS is_test_data;

-- Remove index on users.is_test_data
DROP INDEX IF EXISTS idx_users_is_test_data;

-- Remove is_test_data from users table
ALTER TABLE users DROP COLUMN IF EXISTS is_test_data;
