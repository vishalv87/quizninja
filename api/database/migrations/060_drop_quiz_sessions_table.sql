-- Migration to remove quiz session management functionality
-- Simplifies the system: each quiz start creates a fresh attempt, no pause/resume

-- Drop the cleanup function first (depends on quiz_sessions)
DROP FUNCTION IF EXISTS cleanup_expired_quiz_sessions();

-- Drop the activity update function and trigger
DROP TRIGGER IF EXISTS update_quiz_session_activity_trigger ON quiz_sessions;
DROP FUNCTION IF EXISTS update_quiz_session_activity();

-- Drop the updated_at trigger
DROP TRIGGER IF EXISTS update_quiz_sessions_updated_at ON quiz_sessions;

-- Drop the quiz_sessions table entirely
DROP TABLE IF EXISTS quiz_sessions;

-- Remove session_data column from quiz_attempts (no longer needed)
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS session_data;

-- Update the status constraint to remove 'paused' state
-- Status should only be: started, completed, abandoned
ALTER TABLE quiz_attempts DROP CONSTRAINT IF EXISTS check_quiz_attempt_status;
ALTER TABLE quiz_attempts ADD CONSTRAINT check_quiz_attempt_status
    CHECK (status IN ('started', 'completed', 'abandoned'));

-- Update any existing paused attempts to abandoned
UPDATE quiz_attempts
SET status = 'abandoned', updated_at = CURRENT_TIMESTAMP
WHERE status = 'paused';
