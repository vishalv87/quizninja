-- Migration: Add validation constraints for quiz_attempts data consistency
-- This ensures that completed quiz attempts have all required fields properly set

-- Add constraint to ensure completed attempts have completed status
ALTER TABLE quiz_attempts
ADD CONSTRAINT check_completed_status
CHECK (
    (is_completed = false) OR
    (is_completed = true AND status = 'completed')
);

-- Add constraint to ensure completed attempts have percentage_score set
ALTER TABLE quiz_attempts
ADD CONSTRAINT check_completed_percentage
CHECK (
    (is_completed = false) OR
    (is_completed = true AND percentage_score IS NOT NULL)
);

-- Add constraint to ensure completed attempts have completed_at timestamp
ALTER TABLE quiz_attempts
ADD CONSTRAINT check_completed_timestamp
CHECK (
    (is_completed = false) OR
    (is_completed = true AND completed_at IS NOT NULL)
);

-- Add constraint to ensure passed flag aligns with percentage_score
-- Passing threshold is 60%
ALTER TABLE quiz_attempts
ADD CONSTRAINT check_passed_threshold
CHECK (
    (percentage_score IS NULL) OR
    (percentage_score < 60 AND passed = false) OR
    (percentage_score >= 60 AND passed = true)
);
