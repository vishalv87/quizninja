-- Migration: Fix quiz test data flags
-- Date: 2025-09-27
-- Description: Ensures all quiz-related data is properly marked as test data
-- This addresses cases where quiz data was seeded after the initial test data migration

-- Force update all quiz data to mark as test data (without WHERE condition)
-- This ensures all quiz data is marked as test data regardless of current value
UPDATE quizzes SET is_test_data = TRUE;
UPDATE questions SET is_test_data = TRUE;
UPDATE quiz_statistics SET is_test_data = TRUE;

-- Log the updates for verification
DO $$
DECLARE
    quiz_count INTEGER;
    question_count INTEGER;
    stats_count INTEGER;
BEGIN
    -- Get counts of records with is_test_data = TRUE
    SELECT COUNT(*) INTO quiz_count FROM quizzes WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO question_count FROM questions WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO stats_count FROM quiz_statistics WHERE is_test_data = TRUE;

    -- Log the results
    RAISE NOTICE 'Migration 025 completed - Fixed quiz test data flags:';
    RAISE NOTICE '- Quizzes with is_test_data = TRUE: %', quiz_count;
    RAISE NOTICE '- Questions with is_test_data = TRUE: %', question_count;
    RAISE NOTICE '- Quiz statistics with is_test_data = TRUE: %', stats_count;
    RAISE NOTICE 'All quiz data is now properly marked as test data.';
END $$;