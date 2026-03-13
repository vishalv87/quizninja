-- Migration to remove automatic session expiry
-- This allows users to resume paused quiz sessions indefinitely

-- Update the cleanup function to only clean up completed sessions, not active/paused ones
CREATE OR REPLACE FUNCTION cleanup_expired_quiz_sessions()
RETURNS INTEGER AS $$
DECLARE
    cleaned_count INTEGER;
BEGIN
    -- Only clean up completed sessions older than 30 days to free up space
    -- Active and paused sessions are preserved indefinitely to allow resume
    DELETE FROM quiz_sessions
    WHERE session_state = 'completed'
    AND updated_at < CURRENT_TIMESTAMP - INTERVAL '30 days';

    GET DIAGNOSTICS cleaned_count = ROW_COUNT;

    RETURN cleaned_count;
END;
$$ LANGUAGE plpgsql;

-- Update the function comment to reflect the new behavior
COMMENT ON FUNCTION cleanup_expired_quiz_sessions() IS 'Cleans up completed sessions older than 30 days. Active/paused sessions are preserved indefinitely.';