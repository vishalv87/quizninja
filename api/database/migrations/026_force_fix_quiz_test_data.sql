-- Migration: Force fix quiz test data flags (without WHERE conditions)
-- Date: 2025-09-27
-- Description: Forces all quiz-related data to be marked as test data
-- This addresses persistent issues where quiz data remains marked as false

-- Force update ALL quiz data to mark as test data (no WHERE condition)
UPDATE quizzes SET is_test_data = TRUE;
UPDATE questions SET is_test_data = TRUE;
UPDATE quiz_statistics SET is_test_data = TRUE;

-- Log the final counts
DO $$
DECLARE
    quiz_count INTEGER;
    question_count INTEGER;
    stats_count INTEGER;
    quiz_false_count INTEGER;
    question_false_count INTEGER;
    stats_false_count INTEGER;
BEGIN
    -- Get counts of TRUE records
    SELECT COUNT(*) INTO quiz_count FROM quizzes WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO question_count FROM questions WHERE is_test_data = TRUE;
    SELECT COUNT(*) INTO stats_count FROM quiz_statistics WHERE is_test_data = TRUE;

    -- Get counts of FALSE records (should be 0 after update)
    SELECT COUNT(*) INTO quiz_false_count FROM quizzes WHERE is_test_data = FALSE;
    SELECT COUNT(*) INTO question_false_count FROM questions WHERE is_test_data = FALSE;
    SELECT COUNT(*) INTO stats_false_count FROM quiz_statistics WHERE is_test_data = FALSE;

    -- Log the results
    RAISE NOTICE 'Migration 026 completed - Force fixed quiz test data flags:';
    RAISE NOTICE '- Quizzes with is_test_data = TRUE: % (FALSE: %)', quiz_count, quiz_false_count;
    RAISE NOTICE '- Questions with is_test_data = TRUE: % (FALSE: %)', question_count, question_false_count;
    RAISE NOTICE '- Quiz statistics with is_test_data = TRUE: % (FALSE: %)', stats_count, stats_false_count;

    IF quiz_false_count = 0 AND question_false_count = 0 AND stats_false_count = 0 THEN
        RAISE NOTICE 'SUCCESS: All quiz data is now properly marked as test data!';
    ELSE
        RAISE NOTICE 'WARNING: Some quiz data still has is_test_data = FALSE';
    END IF;
END $$;