-- Migration: Add is_test_data field to all relevant tables
-- Date: 2025-09-27
-- Description: Adds is_test_data boolean field to distinguish between test/sample data and real production data
-- This is Step 1.1 of the Sample Data to API Migration Plan

-- Core Content Tables
ALTER TABLE quizzes ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE questions ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE quiz_attempts ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE quiz_statistics ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;

-- User & Social Tables
ALTER TABLE users ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE user_preferences ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE user_category_performance ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE friendships ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE friend_requests ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE challenges ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE discussions ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE discussion_replies ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE leaderboard_snapshots ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE achievements ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE user_achievements ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;

-- System Tables
ALTER TABLE friend_notifications ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE interests ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE difficulty_levels ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE notification_frequencies ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;

-- Digest Tables (Note: digest_articles already has is_dummy field, but adding is_test_data for consistency)
ALTER TABLE digests ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE digest_articles ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;

-- Additional Tables
ALTER TABLE quiz_ratings ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE user_quiz_favorites ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE user_rank_history ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE quiz_sessions ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE discussion_likes ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;
ALTER TABLE discussion_reply_likes ADD COLUMN is_test_data BOOLEAN DEFAULT FALSE;

-- Create indexes for performance optimization on key tables that will be frequently filtered
CREATE INDEX idx_quizzes_is_test_data ON quizzes(is_test_data);
CREATE INDEX idx_questions_is_test_data ON questions(is_test_data);
CREATE INDEX idx_users_is_test_data ON users(is_test_data);
CREATE INDEX idx_quiz_attempts_is_test_data ON quiz_attempts(is_test_data);
CREATE INDEX idx_discussions_is_test_data ON discussions(is_test_data);
CREATE INDEX idx_challenges_is_test_data ON challenges(is_test_data);
CREATE INDEX idx_achievements_is_test_data ON achievements(is_test_data);
CREATE INDEX idx_digests_is_test_data ON digests(is_test_data);
CREATE INDEX idx_digest_articles_is_test_data ON digest_articles(is_test_data);

-- Add comments to document the purpose of this field
COMMENT ON COLUMN quizzes.is_test_data IS 'Identifies whether this quiz is test/sample data (true) or real production data (false)';
COMMENT ON COLUMN questions.is_test_data IS 'Identifies whether this question is test/sample data (true) or real production data (false)';
COMMENT ON COLUMN users.is_test_data IS 'Identifies whether this user is test/sample data (true) or real production data (false)';
COMMENT ON COLUMN quiz_attempts.is_test_data IS 'Identifies whether this quiz attempt is test/sample data (true) or real production data (false)';
COMMENT ON COLUMN discussions.is_test_data IS 'Identifies whether this discussion is test/sample data (true) or real production data (false)';
COMMENT ON COLUMN challenges.is_test_data IS 'Identifies whether this challenge is test/sample data (true) or real production data (false)';
COMMENT ON COLUMN achievements.is_test_data IS 'Identifies whether this achievement is test/sample data (true) or real production data (false)';
COMMENT ON COLUMN digests.is_test_data IS 'Identifies whether this digest is test/sample data (true) or real production data (false)';
COMMENT ON COLUMN digest_articles.is_test_data IS 'Identifies whether this article is test/sample data (true) or real production data (false)';