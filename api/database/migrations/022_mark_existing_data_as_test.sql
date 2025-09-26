-- Migration: Mark existing data as test data
-- Date: 2025-09-27
-- Description: Updates all existing records to mark them as test/sample data
-- This is Step 1.2 of the Sample Data to API Migration Plan
-- WARNING: This marks ALL existing data as test data. Run only if you want to treat current data as sample/test data.

-- Core Content Tables - Mark existing data as test data
UPDATE quizzes SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE questions SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE quiz_attempts SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE quiz_statistics SET is_test_data = TRUE WHERE is_test_data = FALSE;

-- User & Social Tables - Mark existing data as test data
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

-- System Tables - Mark existing data as test data
UPDATE friend_notifications SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE interests SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE difficulty_levels SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE notification_frequencies SET is_test_data = TRUE WHERE is_test_data = FALSE;

-- Digest Tables - Mark existing data as test data
-- Note: We'll mark both is_test_data and is_dummy for digest tables to maintain consistency
UPDATE digests SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE digest_articles SET is_test_data = TRUE WHERE is_test_data = FALSE;

-- Additional Tables - Mark existing data as test data
UPDATE quiz_ratings SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE user_quiz_favorites SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE user_rank_history SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE quiz_sessions SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE discussion_likes SET is_test_data = TRUE WHERE is_test_data = FALSE;
UPDATE discussion_reply_likes SET is_test_data = TRUE WHERE is_test_data = FALSE;

-- Optional: Log the counts of updated records for verification
DO $$
DECLARE
    quiz_count INTEGER;
    user_count INTEGER;
    question_count INTEGER;
    attempt_count INTEGER;
BEGIN
    -- Get counts of updated records
    SELECT COUNT(*) INTO quiz_count FROM quizzes WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO user_count FROM users WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO question_count FROM questions WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO attempt_count FROM quiz_attempts WHERE is_test_data = TRUE;

    -- Log the results (these will appear in migration output)
    RAISE NOTICE 'Migration 022 completed:';
    RAISE NOTICE '- Quizzes marked as test data: %', quiz_count;
    RAISE NOTICE '- Users marked as test data: %', user_count;
    RAISE NOTICE '- Questions marked as test data: %', question_count;
    RAISE NOTICE '- Quiz attempts marked as test data: %', attempt_count;
    RAISE NOTICE 'All existing data has been marked as test data. New production data will default to is_test_data = FALSE.';
END $$;