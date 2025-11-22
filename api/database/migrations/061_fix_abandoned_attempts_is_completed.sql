-- Migration: 061_fix_abandoned_attempts_is_completed.sql
-- Description: Fix existing abandoned quiz attempts that incorrectly have is_completed = false
-- This was caused by a bug in the AbandonQuizAttempt function that set is_completed = false
-- instead of is_completed = true when abandoning attempts.
--
-- The unique constraint idx_quiz_attempts_user_quiz_active only allows one attempt per user/quiz
-- where is_completed = false. Abandoned attempts should have is_completed = true so they don't
-- block new attempts.
--
-- We need to update several check constraints to allow abandoned attempts with is_completed = true

-- Step 1: Drop existing check constraints that don't account for abandoned status
ALTER TABLE quiz_attempts DROP CONSTRAINT IF EXISTS check_completed_status;
ALTER TABLE quiz_attempts DROP CONSTRAINT IF EXISTS check_completed_percentage;
ALTER TABLE quiz_attempts DROP CONSTRAINT IF EXISTS check_completed_timestamp;

-- Step 2: Add updated check constraints that allow abandoned attempts

-- Status constraint: allow 'completed' or 'abandoned' when is_completed = true
ALTER TABLE quiz_attempts
ADD CONSTRAINT check_completed_status
CHECK (
    (is_completed = false) OR
    (is_completed = true AND status IN ('completed', 'abandoned'))
);

-- Percentage constraint: require percentage_score only for 'completed' status, not 'abandoned'
ALTER TABLE quiz_attempts
ADD CONSTRAINT check_completed_percentage
CHECK (
    (is_completed = false) OR
    (is_completed = true AND status = 'abandoned') OR
    (is_completed = true AND status = 'completed' AND percentage_score IS NOT NULL)
);

-- Timestamp constraint: require completed_at only for 'completed' status, not 'abandoned'
ALTER TABLE quiz_attempts
ADD CONSTRAINT check_completed_timestamp
CHECK (
    (is_completed = false) OR
    (is_completed = true AND status = 'abandoned') OR
    (is_completed = true AND status = 'completed' AND completed_at IS NOT NULL)
);

-- Step 3: Fix existing abandoned attempts that have is_completed = false
UPDATE quiz_attempts
SET is_completed = true, updated_at = CURRENT_TIMESTAMP
WHERE status = 'abandoned' AND is_completed = false;
