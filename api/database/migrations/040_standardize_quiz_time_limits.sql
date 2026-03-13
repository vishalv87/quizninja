-- Migration: Standardize all quiz time limits to 5 minutes
-- Purpose: Ensure all quizzes have a consistent 5-minute time limit
-- Note: The backend converts time_limit_minutes to seconds via MarshalJSON (5 * 60 = 300 seconds)

-- Update all existing quizzes to have 5-minute time limit
UPDATE quizzes
SET time_limit_minutes = 5,
    updated_at = NOW()
WHERE time_limit_minutes != 5;
