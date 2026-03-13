-- Add retake functionality fields to quiz_attempts table
-- This migration adds fields required for tracking quiz retakes and performance comparison

-- Add retake tracking fields to quiz_attempts table
ALTER TABLE quiz_attempts
ADD COLUMN IF NOT EXISTS retake_count INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS original_attempt_id UUID REFERENCES quiz_attempts(id),
ADD COLUMN IF NOT EXISTS performance_comparison JSONB;

-- Add indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_retake_count ON quiz_attempts(retake_count);
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_original_id ON quiz_attempts(original_attempt_id);

-- Add check constraint for retake_count (max 3 retakes)
ALTER TABLE quiz_attempts
ADD CONSTRAINT check_retake_count
CHECK (retake_count >= 0 AND retake_count <= 3);

-- Update existing records to have retake_count = 0 if null
UPDATE quiz_attempts
SET retake_count = 0
WHERE retake_count IS NULL;

-- Add comment for documentation
COMMENT ON COLUMN quiz_attempts.retake_count IS 'Number of times this quiz has been retaken (0 for original attempt)';
COMMENT ON COLUMN quiz_attempts.original_attempt_id IS 'Reference to the original attempt if this is a retake';
COMMENT ON COLUMN quiz_attempts.performance_comparison IS 'JSON data comparing performance with previous attempts';