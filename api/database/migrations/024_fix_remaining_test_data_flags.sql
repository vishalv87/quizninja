-- Migration: Fix remaining test data flags
-- Date: 2025-09-27
-- Description: Updates any remaining records that still have is_test_data = false
-- This ensures all existing data is properly marked as test data

-- Core Content Tables - Ensure all existing data is marked as test data
UPDATE quizzes SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE questions SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE quiz_attempts SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE quiz_statistics SET is_test_data = TRUE WHERE is_test_data = FALSE;

-- User & Social Tables - Ensure all existing data is marked as test data
UPDATE users SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE user_preferences SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE user_category_performance SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE friendships SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE friend_requests SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE challenges SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE discussions SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE discussion_replies SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE leaderboard_snapshots SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE achievements SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE user_achievements SET is_test_data = TRUE WHERE is_test_data = FALSE;

-- System Tables - Ensure all existing data is marked as test data
UPDATE friend_notifications SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE interests SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE difficulty_levels SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE notification_frequencies SET is_test_data = TRUE WHERE is_test_data = FALSE;

-- Digest Tables - Ensure all existing data is marked as test data
UPDATE digests SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE digest_articles SET is_test_data = TRUE WHERE is_test_data = FALSE;

-- Additional Tables - Ensure all existing data is marked as test data
UPDATE quiz_ratings SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE user_quiz_favorites SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE user_rank_history SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE quiz_sessions SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE discussion_likes SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE discussion_reply_likes SET is_test_data = TRUE WHERE is_test_data = FALSE;

-- Log the counts of updated records for verification
DO $$
DECLARE
    quiz_count INTEGER;
    user_count INTEGER;
    question_count INTEGER;
    attempt_count INTEGER;
    achievement_count INTEGER;
    total_updates INTEGER := 0;
BEGIN
    -- Get counts of records marked as test data
    SELECT COUNT(*) INTO quiz_count FROM quizzes WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO user_count FROM users WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO question_count FROM questions WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO attempt_count FROM quiz_attempts WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO achievement_count FROM achievements WHERE is_test_data = TRUE;

    -- Log the results
    RAISE NOTICE 'Migration 024 completed - Final counts of test data:';
    RAISE NOTICE '- Quizzes marked as test data: %', quiz_count;
    RAISE NOTICE '- Users marked as test data: %', user_count;
    RAISE NOTICE '- Questions marked as test data: %', question_count;
    RAISE NOTICE '- Quiz attempts marked as test data: %', attempt_count;
    RAISE NOTICE '- Achievements marked as test data: %', achievement_count;
    RAISE NOTICE 'All existing data should now be marked as test data.';
    RAISE NOTICE 'New data created via handlers will automatically have is_test_data = TRUE.';
END $$;