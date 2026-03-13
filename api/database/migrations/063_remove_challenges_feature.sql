-- Migration: Remove Challenges Feature
-- This migration completely removes the challenges feature from the database

-- Step 1: Remove challenge-related notification data
DELETE FROM notifications WHERE type IN (
  'challenge_received',
  'challenge_accepted',
  'challenge_declined',
  'challenge_progress',
  'challenge_completed'
);

-- Step 2: Remove challenge-related fields from quiz_attempts
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS challenge_id;
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS is_challenge_attempt;

-- Step 3: Drop challenges table (will cascade delete related data)
DROP TABLE IF EXISTS challenges CASCADE;