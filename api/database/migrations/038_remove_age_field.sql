-- Migration: Remove age field from users table
-- Created: 2025-10-06
-- Description: Removes the optional age column from users table as it's not needed for the application

-- Drop the age column from users table
ALTER TABLE users DROP COLUMN IF EXISTS age;

-- Note: This is a non-reversible migration. If age data needs to be preserved,
-- ensure a backup is created before running this migration.
