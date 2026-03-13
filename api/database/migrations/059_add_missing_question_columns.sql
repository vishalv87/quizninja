-- Migration: Add points and updated_at columns to questions table
-- Description: Adds points column for question scoring and updated_at column for tracking modifications
-- Author: Claude Code
-- Date: 2025-11-19

-- Add points column to questions table
ALTER TABLE questions
ADD COLUMN IF NOT EXISTS points INTEGER NOT NULL DEFAULT 1;

-- Add updated_at column to questions table
ALTER TABLE questions
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP;

-- Add comments to document the column purposes
COMMENT ON COLUMN questions.points IS 'Point value for this question (default: 1)';
COMMENT ON COLUMN questions.updated_at IS 'Timestamp when the question was last updated';

-- Create or replace trigger function for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for questions table
DROP TRIGGER IF EXISTS update_questions_updated_at ON questions;
CREATE TRIGGER update_questions_updated_at
    BEFORE UPDATE ON questions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
