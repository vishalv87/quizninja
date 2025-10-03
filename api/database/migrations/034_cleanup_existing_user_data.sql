-- Migration: Cleanup existing user data for fresh start
-- This migration safely removes all existing user data while preserving system data
-- Use this for development/testing environments only

-- ========================================
-- PHASE 1: PRE-REMOVAL VERIFICATION
-- ========================================

DO $$
DECLARE
    user_count INTEGER;
    test_data_count INTEGER;
BEGIN
    SELECT COUNT(*), COUNT(CASE WHEN is_test_data = true THEN 1 END)
    INTO user_count, test_data_count
    FROM users;

    RAISE NOTICE '=== PRE-CLEANUP VERIFICATION ===';
    RAISE NOTICE 'Total users: %, Test data users: %', user_count, test_data_count;

    IF user_count = 0 THEN
        RAISE NOTICE 'No users found in database';
    END IF;
END $$;

-- Show current users
DO $$
DECLARE
    user_record RECORD;
BEGIN
    RAISE NOTICE 'Current users in database:';
    FOR user_record IN (SELECT email, auth_method, is_test_data, created_at FROM users ORDER BY created_at) LOOP
        RAISE NOTICE 'User: % | Auth: % | Test Data: % | Created: %',
            user_record.email, user_record.auth_method, user_record.is_test_data, user_record.created_at;
    END LOOP;
END $$;

-- Count related data
DO $$
DECLARE
    count_val INTEGER;
    total_records INTEGER := 0;
BEGIN
    RAISE NOTICE 'User-related data counts before cleanup:';

    SELECT COUNT(*) INTO count_val FROM user_preferences;
    total_records := total_records + count_val;
    IF count_val > 0 THEN RAISE NOTICE '  user_preferences: %', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM quiz_attempts;
    total_records := total_records + count_val;
    IF count_val > 0 THEN RAISE NOTICE '  quiz_attempts: %', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM refresh_tokens;
    total_records := total_records + count_val;
    IF count_val > 0 THEN RAISE NOTICE '  refresh_tokens: %', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM notifications;
    total_records := total_records + count_val;
    IF count_val > 0 THEN RAISE NOTICE '  notifications: %', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM challenges;
    total_records := total_records + count_val;
    IF count_val > 0 THEN RAISE NOTICE '  challenges: %', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM friend_requests;
    total_records := total_records + count_val;
    IF count_val > 0 THEN RAISE NOTICE '  friend_requests: %', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM friendships;
    total_records := total_records + count_val;
    IF count_val > 0 THEN RAISE NOTICE '  friendships: %', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM user_achievements;
    total_records := total_records + count_val;
    IF count_val > 0 THEN RAISE NOTICE '  user_achievements: %', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM user_category_performance;
    total_records := total_records + count_val;
    IF count_val > 0 THEN RAISE NOTICE '  user_category_performance: %', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM user_quiz_favorites;
    total_records := total_records + count_val;
    IF count_val > 0 THEN RAISE NOTICE '  user_quiz_favorites: %', count_val; END IF;

    RAISE NOTICE 'Total user-related records to be removed: %', total_records;
END $$;

-- ========================================
-- PHASE 2: SAFE USER REMOVAL
-- ========================================

DO $$
DECLARE
    deleted_count INTEGER;
    quiz_count INTEGER;
BEGIN
    RAISE NOTICE '=== EXECUTING USER CLEANUP ===';

    -- First, handle quizzes created by users (set created_by to NULL)
    UPDATE quizzes SET created_by = NULL WHERE created_by IS NOT NULL;
    GET DIAGNOSTICS quiz_count = ROW_COUNT;

    IF quiz_count > 0 THEN
        RAISE NOTICE 'Updated % quizzes to remove user references', quiz_count;
    END IF;

    -- Now remove all users (CASCADE will handle all related data automatically)
    DELETE FROM users;
    GET DIAGNOSTICS deleted_count = ROW_COUNT;

    RAISE NOTICE 'Deleted % users and all related data via CASCADE', deleted_count;
END $$;

-- ========================================
-- PHASE 3: POST-REMOVAL VERIFICATION
-- ========================================

DO $$
DECLARE
    user_count INTEGER;
    count_val INTEGER;
    total_remaining INTEGER := 0;
BEGIN
    RAISE NOTICE '=== POST-CLEANUP VERIFICATION ===';

    SELECT COUNT(*) INTO user_count FROM users;
    RAISE NOTICE 'Remaining users: %', user_count;

    -- Check all user-related tables are empty
    SELECT COUNT(*) INTO count_val FROM user_preferences;
    total_remaining := total_remaining + count_val;
    IF count_val > 0 THEN RAISE WARNING 'user_preferences still has % records!', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM quiz_attempts;
    total_remaining := total_remaining + count_val;
    IF count_val > 0 THEN RAISE WARNING 'quiz_attempts still has % records!', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM refresh_tokens;
    total_remaining := total_remaining + count_val;
    IF count_val > 0 THEN RAISE WARNING 'refresh_tokens still has % records!', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM notifications;
    total_remaining := total_remaining + count_val;
    IF count_val > 0 THEN RAISE WARNING 'notifications still has % records!', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM challenges;
    total_remaining := total_remaining + count_val;
    IF count_val > 0 THEN RAISE WARNING 'challenges still has % records!', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM friend_requests;
    total_remaining := total_remaining + count_val;
    IF count_val > 0 THEN RAISE WARNING 'friend_requests still has % records!', count_val; END IF;

    SELECT COUNT(*) INTO count_val FROM friendships;
    total_remaining := total_remaining + count_val;
    IF count_val > 0 THEN RAISE WARNING 'friendships still has % records!', count_val; END IF;

    IF total_remaining = 0 THEN
        RAISE NOTICE 'SUCCESS: All user-related data has been removed';
    ELSE
        RAISE WARNING 'WARNING: % user-related records still remain', total_remaining;
    END IF;
END $$;

-- Verify system tables are preserved
DO $$
DECLARE
    count_val INTEGER;
    system_tables_count INTEGER := 0;
BEGIN
    RAISE NOTICE 'System tables preservation check:';

    SELECT COUNT(*) INTO count_val FROM achievements;
    system_tables_count := system_tables_count + count_val;
    RAISE NOTICE '  achievements: % records preserved', count_val;

    SELECT COUNT(*) INTO count_val FROM app_settings;
    system_tables_count := system_tables_count + count_val;
    RAISE NOTICE '  app_settings: % records preserved', count_val;

    SELECT COUNT(*) INTO count_val FROM difficulty_levels;
    system_tables_count := system_tables_count + count_val;
    RAISE NOTICE '  difficulty_levels: % records preserved', count_val;

    SELECT COUNT(*) INTO count_val FROM interests;
    system_tables_count := system_tables_count + count_val;
    RAISE NOTICE '  interests: % records preserved', count_val;

    SELECT COUNT(*) INTO count_val FROM quizzes;
    system_tables_count := system_tables_count + count_val;
    RAISE NOTICE '  quizzes: % records preserved', count_val;

    SELECT COUNT(*) INTO count_val FROM questions;
    system_tables_count := system_tables_count + count_val;
    RAISE NOTICE '  questions: % records preserved', count_val;

    RAISE NOTICE 'Total system records preserved: %', system_tables_count;
    RAISE NOTICE '=== DATABASE CLEANUP COMPLETED SUCCESSFULLY ===';
    RAISE NOTICE 'Ready for fresh user registration and Supabase auth testing';
END $$;