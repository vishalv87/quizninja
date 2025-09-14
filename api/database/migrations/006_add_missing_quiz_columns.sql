-- Add missing columns and tables to match the API models
-- This migration adds the columns and tables that are expected by the API but missing from the schema

-- Add missing columns to quizzes table
ALTER TABLE quizzes
ADD COLUMN created_by UUID REFERENCES users(id),
ADD COLUMN is_public BOOLEAN DEFAULT true,
ADD COLUMN tags TEXT[] DEFAULT '{}',
ADD COLUMN thumbnail_url VARCHAR(500);

-- Add indexes for the new columns for performance
CREATE INDEX idx_quizzes_created_by ON quizzes(created_by);
CREATE INDEX idx_quizzes_is_public ON quizzes(is_public);
CREATE INDEX idx_quizzes_thumbnail_url ON quizzes(thumbnail_url);

-- Create quiz_attempts table (missing from schema)
CREATE TABLE quiz_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quiz_id UUID REFERENCES quizzes(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    score DECIMAL(5,2) NOT NULL,
    total_points INTEGER NOT NULL,
    time_spent INTEGER NOT NULL, -- in seconds
    is_completed BOOLEAN DEFAULT false,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for quiz_attempts table
CREATE INDEX idx_quiz_attempts_quiz_id ON quiz_attempts(quiz_id);
CREATE INDEX idx_quiz_attempts_user_id ON quiz_attempts(user_id);
CREATE INDEX idx_quiz_attempts_is_completed ON quiz_attempts(is_completed);
CREATE INDEX idx_quiz_attempts_started_at ON quiz_attempts(started_at);
CREATE INDEX idx_quiz_attempts_completed_at ON quiz_attempts(completed_at);

-- Add unique constraint to prevent duplicate active attempts
CREATE UNIQUE INDEX idx_quiz_attempts_user_quiz_active
ON quiz_attempts(user_id, quiz_id)
WHERE is_completed = false;

-- Add trigger for quiz_attempts updated_at
CREATE TRIGGER update_quiz_attempts_updated_at
    BEFORE UPDATE ON quiz_attempts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Update existing records to have a default created_by (first user in the system)
-- This is safe to run even if there are no users yet
UPDATE quizzes
SET created_by = (SELECT id FROM users ORDER BY created_at LIMIT 1)
WHERE created_by IS NULL AND EXISTS (SELECT 1 FROM users LIMIT 1);