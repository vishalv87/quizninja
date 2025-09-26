-- Migration: Add missing updated_at columns to tables that need them
-- Date: 2025-09-27
-- Description: Adds updated_at timestamp columns to tables that have UPDATE triggers but were missing the column
-- This fixes the migration 022 error: "record new has no field updated_at"

-- Add updated_at column to difficulty_levels table
ALTER TABLE difficulty_levels ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Add updated_at column to interests table
ALTER TABLE interests ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Add updated_at column to notification_frequencies table
ALTER TABLE notification_frequencies ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Add UPDATE triggers for these tables to automatically update the timestamp
-- Note: The triggers may already exist but will now work correctly with the new columns

-- Trigger for difficulty_levels (create if not exists)
DROP TRIGGER IF EXISTS update_difficulty_levels_updated_at ON difficulty_levels;
CREATE TRIGGER update_difficulty_levels_updated_at
    BEFORE UPDATE ON difficulty_levels
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for interests (create if not exists)
DROP TRIGGER IF EXISTS update_interests_updated_at ON interests;
CREATE TRIGGER update_interests_updated_at
    BEFORE UPDATE ON interests
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for notification_frequencies (create if not exists)
DROP TRIGGER IF EXISTS update_notification_frequencies_updated_at ON notification_frequencies;
CREATE TRIGGER update_notification_frequencies_updated_at
    BEFORE UPDATE ON notification_frequencies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Log completion
DO $$
BEGIN
    RAISE NOTICE 'Migration 022.1 completed: Added updated_at columns to difficulty_levels, interests, and notification_frequencies tables';
END $$;