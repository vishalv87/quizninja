-- Migration: Fix quiz_statistics updated_at column naming
-- Date: 2025-09-27
-- Description: Renames last_updated to updated_at in quiz_statistics table for consistency
-- This fixes the remaining migration 022 error where trigger expects updated_at but column is named last_updated

-- Rename the column from last_updated to updated_at for consistency
ALTER TABLE quiz_statistics RENAME COLUMN last_updated TO updated_at;

-- The existing trigger update_quiz_statistics_last_updated will now work correctly
-- since it calls update_updated_at_column() which sets NEW.updated_at

-- Log completion
DO $$
BEGIN
    RAISE NOTICE 'Migration 022.2 completed: Renamed quiz_statistics.last_updated to updated_at for trigger consistency';
END $$;