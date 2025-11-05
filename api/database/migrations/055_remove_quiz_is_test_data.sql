-- Migration: Remove is_test_data columns from quiz-related tables
-- Date: 2025-11-04

-- Drop the is_test_data column from quizzes table
ALTER TABLE quizzes DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from questions table
ALTER TABLE questions DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from quiz_statistics table
ALTER TABLE quiz_statistics DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from quiz_attempts table
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from quiz_sessions table
ALTER TABLE quiz_sessions DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from user_quiz_favorites table
ALTER TABLE user_quiz_favorites DROP COLUMN IF EXISTS is_test_data;
