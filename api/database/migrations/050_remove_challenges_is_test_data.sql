-- Migration: Remove is_test_data field from challenges table
-- Description: Removes the is_test_data column and its index from the challenges table
-- Date: 2025-11-03

-- Drop the index first
DROP INDEX IF EXISTS idx_challenges_is_test_data;

-- Drop the column
ALTER TABLE challenges DROP COLUMN IF EXISTS is_test_data;
