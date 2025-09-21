-- Add privacy settings to user_preferences table
-- Migration: 015_add_privacy_settings

-- Add privacy settings columns to user_preferences table
ALTER TABLE user_preferences ADD COLUMN IF NOT EXISTS profile_visibility BOOLEAN DEFAULT true;
ALTER TABLE user_preferences ADD COLUMN IF NOT EXISTS show_online_status BOOLEAN DEFAULT true;
ALTER TABLE user_preferences ADD COLUMN IF NOT EXISTS allow_friend_requests BOOLEAN DEFAULT true;
ALTER TABLE user_preferences ADD COLUMN IF NOT EXISTS share_activity_status BOOLEAN DEFAULT true;

-- Add notification type preferences as JSON
ALTER TABLE user_preferences ADD COLUMN IF NOT EXISTS notification_types JSONB DEFAULT '{
  "quiz_reminders": true,
  "friend_activity": true,
  "challenges": true,
  "achievements": true,
  "leaderboard_updates": false,
  "system_announcements": true
}'::jsonb;

-- Update existing rows to have default values
UPDATE user_preferences
SET
  profile_visibility = true,
  show_online_status = true,
  allow_friend_requests = true,
  share_activity_status = true
WHERE profile_visibility IS NULL;

-- Add updated_at column for tracking preference changes
ALTER TABLE user_preferences ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Create trigger to update updated_at automatically
CREATE OR REPLACE FUNCTION update_user_preferences_updated_at()
RETURNS trigger AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Drop trigger if it exists and recreate it
DROP TRIGGER IF EXISTS trigger_update_user_preferences_updated_at ON user_preferences;
CREATE TRIGGER trigger_update_user_preferences_updated_at
    BEFORE UPDATE ON user_preferences
    FOR EACH ROW
    EXECUTE FUNCTION update_user_preferences_updated_at();

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_user_preferences_updated_at ON user_preferences(updated_at);