-- Add onboarding completion tracking
-- Migration: 016_add_onboarding_completion_tracking

-- Add onboarding completion timestamp to user_preferences table
ALTER TABLE user_preferences ADD COLUMN IF NOT EXISTS onboarding_completed_at TIMESTAMP NULL;

-- Add index for onboarding status queries
CREATE INDEX IF NOT EXISTS idx_user_preferences_onboarding_completed ON user_preferences(onboarding_completed_at);

-- Add index for filtering users who haven't completed onboarding
CREATE INDEX IF NOT EXISTS idx_user_preferences_pending_onboarding ON user_preferences(user_id) WHERE onboarding_completed_at IS NULL;

-- Update existing user preferences to mark them as completed (since they were created before separation)
UPDATE user_preferences
SET onboarding_completed_at = created_at
WHERE onboarding_completed_at IS NULL
  AND selected_interests IS NOT NULL
  AND difficulty_preference IS NOT NULL;