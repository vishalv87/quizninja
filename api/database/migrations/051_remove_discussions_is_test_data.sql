-- Migration: Remove is_test_data field from discussion-related tables
-- Description: Removes the is_test_data column and its index from discussions, discussion_replies, discussion_likes, and discussion_reply_likes tables

-- Drop the index on is_test_data from discussions table
DROP INDEX IF EXISTS idx_discussions_is_test_data;

-- Drop the is_test_data column from discussions table
ALTER TABLE discussions DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from discussion_replies table
ALTER TABLE discussion_replies DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from discussion_likes table
ALTER TABLE discussion_likes DROP COLUMN IF EXISTS is_test_data;

-- Drop the is_test_data column from discussion_reply_likes table
ALTER TABLE discussion_reply_likes DROP COLUMN IF EXISTS is_test_data;
