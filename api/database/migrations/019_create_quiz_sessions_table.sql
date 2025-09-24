-- Create quiz sessions table for handling quiz continuation functionality
-- This migration adds the ability to pause and resume quizzes

-- Quiz sessions table to track active/paused quiz states
CREATE TABLE quiz_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attempt_id UUID NOT NULL REFERENCES quiz_attempts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    current_question_index INTEGER NOT NULL DEFAULT 0,
    current_answers JSONB DEFAULT '[]'::jsonb,
    session_state VARCHAR(20) NOT NULL DEFAULT 'active',
    time_remaining INTEGER, -- seconds remaining from original time limit
    time_spent_so_far INTEGER NOT NULL DEFAULT 0, -- total time spent in seconds
    last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paused_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CHECK (session_state IN ('active', 'paused', 'completed', 'abandoned')),
    CHECK (current_question_index >= 0),
    CHECK (time_spent_so_far >= 0),
    CHECK (time_remaining IS NULL OR time_remaining >= 0),

    -- Ensure one session per attempt
    UNIQUE (attempt_id)
);

-- Add indexes for performance
CREATE INDEX idx_quiz_sessions_user_id ON quiz_sessions(user_id);
CREATE INDEX idx_quiz_sessions_attempt_id ON quiz_sessions(attempt_id);
CREATE INDEX idx_quiz_sessions_quiz_id ON quiz_sessions(quiz_id);
CREATE INDEX idx_quiz_sessions_state ON quiz_sessions(session_state);
CREATE INDEX idx_quiz_sessions_last_activity ON quiz_sessions(last_activity_at);
CREATE INDEX idx_quiz_sessions_user_active ON quiz_sessions(user_id, session_state)
    WHERE session_state IN ('active', 'paused');

-- Update quiz_attempts table to support session states
ALTER TABLE quiz_attempts ADD COLUMN IF NOT EXISTS session_data JSONB;

-- Update the status constraint to include 'paused' state
ALTER TABLE quiz_attempts DROP CONSTRAINT IF EXISTS check_quiz_attempt_status;
ALTER TABLE quiz_attempts ADD CONSTRAINT check_quiz_attempt_status
    CHECK (status IN ('started', 'paused', 'completed', 'abandoned'));

-- Add trigger to update quiz_sessions updated_at timestamp
CREATE TRIGGER update_quiz_sessions_updated_at
    BEFORE UPDATE ON quiz_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add trigger to update last_activity_at on session updates
CREATE OR REPLACE FUNCTION update_quiz_session_activity()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_activity_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_quiz_session_activity_trigger
    BEFORE UPDATE ON quiz_sessions
    FOR EACH ROW EXECUTE FUNCTION update_quiz_session_activity();

-- Function to cleanup expired sessions (older than 24 hours of inactivity)
CREATE OR REPLACE FUNCTION cleanup_expired_quiz_sessions()
RETURNS INTEGER AS $$
DECLARE
    expired_count INTEGER;
BEGIN
    -- Mark sessions as abandoned if inactive for more than 24 hours
    UPDATE quiz_sessions
    SET session_state = 'abandoned',
        updated_at = CURRENT_TIMESTAMP
    WHERE session_state IN ('active', 'paused')
    AND last_activity_at < CURRENT_TIMESTAMP - INTERVAL '24 hours';

    GET DIAGNOSTICS expired_count = ROW_COUNT;

    -- Also update corresponding quiz attempts
    UPDATE quiz_attempts
    SET status = 'abandoned',
        updated_at = CURRENT_TIMESTAMP
    WHERE id IN (
        SELECT attempt_id FROM quiz_sessions
        WHERE session_state = 'abandoned'
        AND status != 'abandoned'
    );

    RETURN expired_count;
END;
$$ LANGUAGE plpgsql;

-- Add some sample comments for documentation
COMMENT ON TABLE quiz_sessions IS 'Tracks active quiz sessions for pause/resume functionality';
COMMENT ON COLUMN quiz_sessions.current_question_index IS 'Zero-based index of current question';
COMMENT ON COLUMN quiz_sessions.current_answers IS 'JSON array of user answers so far';
COMMENT ON COLUMN quiz_sessions.session_state IS 'Current state: active, paused, completed, abandoned';
COMMENT ON COLUMN quiz_sessions.time_remaining IS 'Seconds remaining from original quiz time limit';
COMMENT ON COLUMN quiz_sessions.time_spent_so_far IS 'Total active time spent on quiz in seconds';
COMMENT ON FUNCTION cleanup_expired_quiz_sessions() IS 'Cleans up sessions inactive for 24+ hours';