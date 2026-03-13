-- Migration 028: Fix user_preferences unique constraint
-- This migration adds a UNIQUE constraint on user_preferences.user_id to fix the settings save issue
-- The issue was that the UPSERT operation (ON CONFLICT) wasn't working because there was no unique constraint

-- First, check if there are any duplicate user preferences that need to be cleaned up
-- We'll keep the most recent preference record for each user and delete the older ones

-- Step 1: Remove duplicate user_preferences records, keeping only the most recent one
WITH ranked_preferences AS (
    SELECT
        id,
        user_id,
        ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY updated_at DESC, created_at DESC, id DESC) as rn
    FROM user_preferences
),
duplicates_to_delete AS (
    SELECT id
    FROM ranked_preferences
    WHERE rn > 1
)
DELETE FROM user_preferences
WHERE id IN (SELECT id FROM duplicates_to_delete);

-- Step 2: Add the unique constraint on user_id
-- This will ensure that each user can only have one preferences record
ALTER TABLE user_preferences
ADD CONSTRAINT uq_user_preferences_user_id UNIQUE (user_id);

-- Step 3: Create a partial index to improve performance for onboarding queries
-- This replaces the existing similar index with a more descriptive name
DROP INDEX IF EXISTS idx_user_preferences_pending_onboarding;
CREATE INDEX idx_user_preferences_onboarding_pending
ON user_preferences(user_id)
WHERE onboarding_completed_at IS NULL;

-- Step 4: Verify the fix by checking that UpdateUserPreferences UPSERT will work
-- The ON CONFLICT (user_id) clause will now function correctly with the unique constraint

-- Add a comment explaining the fix
COMMENT ON CONSTRAINT uq_user_preferences_user_id ON user_preferences IS
'Ensures each user has only one preferences record, enabling proper UPSERT operations in UpdateUserPreferences';