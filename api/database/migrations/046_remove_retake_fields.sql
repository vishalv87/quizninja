-- Remove retake functionality fields from quiz_attempts table
-- This migration removes fields that allowed users to retake quizzes multiple times
-- Going forward, each quiz can only be attempted once (when completed)

-- Drop check constraint for retake_count
ALTER TABLE quiz_attempts
DROP CONSTRAINT IF EXISTS check_retake_count;

-- Drop indexes for retake fields
DROP INDEX IF EXISTS idx_quiz_attempts_retake_count;
DROP INDEX IF EXISTS idx_quiz_attempts_original_id;

-- Drop retake-related columns
ALTER TABLE quiz_attempts
DROP COLUMN IF EXISTS retake_count,
DROP COLUMN IF EXISTS original_attempt_id,
DROP COLUMN IF EXISTS performance_comparison;
