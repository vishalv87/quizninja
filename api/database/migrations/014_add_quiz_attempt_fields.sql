-- Add missing fields to quiz_attempts table for enhanced attempt management
-- This migration adds fields required for storing answers, percentage scores, pass status, and attempt status

-- Add missing columns to quiz_attempts table
ALTER TABLE quiz_attempts
ADD COLUMN IF NOT EXISTS answers JSONB,
ADD COLUMN IF NOT EXISTS percentage_score DECIMAL(5,2),
ADD COLUMN IF NOT EXISTS passed BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'started';

-- Add check constraint for status field
ALTER TABLE quiz_attempts
ADD CONSTRAINT check_quiz_attempt_status
CHECK (status IN ('started', 'completed', 'abandoned'));

-- Add index for status column for performance
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_status ON quiz_attempts(status);

-- Add index for passed column for performance
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_passed ON quiz_attempts(passed);

-- Update existing records to have proper status based on is_completed
UPDATE quiz_attempts
SET status = CASE
    WHEN is_completed = true THEN 'completed'
    ELSE 'started'
END
WHERE status IS NULL OR status = 'started';

-- Update existing completed attempts to calculate percentage_score if missing
UPDATE quiz_attempts
SET percentage_score = CASE
    WHEN total_points > 0 THEN ROUND((score / total_points * 100), 2)
    ELSE 0.0
END,
passed = CASE
    WHEN total_points > 0 AND (score / total_points * 100) >= 60 THEN true
    ELSE false
END
WHERE is_completed = true AND percentage_score IS NULL;