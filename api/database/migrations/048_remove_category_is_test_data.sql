-- Migration: Remove is_test_data column from categories table
-- Date: 2025-11-03

-- Drop the is_test_data column from categories table
ALTER TABLE categories DROP COLUMN IF EXISTS is_test_data;
