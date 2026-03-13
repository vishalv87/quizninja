-- Sync quiz points with total questions
-- This migration ensures that quiz points always equals the number of questions (1 point per question)

-- Step 1: Update all existing quizzes to set points = total_questions
UPDATE quizzes
SET points = total_questions
WHERE points != total_questions;

-- Step 2: Create a trigger function to automatically sync points when total_questions changes
CREATE OR REPLACE FUNCTION sync_quiz_points_with_questions()
RETURNS TRIGGER AS $$
BEGIN
    -- Automatically set points to match total_questions
    NEW.points := NEW.total_questions;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Step 3: Create trigger that fires before insert or update on quizzes table
CREATE TRIGGER trigger_sync_quiz_points
    BEFORE INSERT OR UPDATE OF total_questions, points ON quizzes
    FOR EACH ROW
    EXECUTE FUNCTION sync_quiz_points_with_questions();

-- Step 4: Add CHECK constraint to enforce the rule at database level
ALTER TABLE quizzes
ADD CONSTRAINT check_points_equals_questions
CHECK (points = total_questions);
