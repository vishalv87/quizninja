-- Migration to remove all database triggers
-- Simplifies the system: all logic (updated_at, notifications, friendships, levels)
-- will be handled in application code instead of database triggers

-- ============================================================
-- 1. DROP ALL TIMESTAMP (updated_at) TRIGGERS
-- ============================================================

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_app_settings_updated_at ON app_settings;
DROP TRIGGER IF EXISTS update_achievements_updated_at ON achievements;
DROP TRIGGER IF EXISTS update_interests_updated_at ON categories;
DROP TRIGGER IF EXISTS update_difficulty_levels_updated_at ON difficulty_levels;
DROP TRIGGER IF EXISTS update_notification_frequencies_updated_at ON notification_frequencies;
DROP TRIGGER IF EXISTS update_quizzes_updated_at ON quizzes;
DROP TRIGGER IF EXISTS update_questions_updated_at ON questions;
DROP TRIGGER IF EXISTS update_quiz_statistics_last_updated ON quiz_statistics;
DROP TRIGGER IF EXISTS update_quiz_attempts_updated_at ON quiz_attempts;
DROP TRIGGER IF EXISTS update_user_category_performance_updated_at ON user_category_performance;
DROP TRIGGER IF EXISTS update_discussions_updated_at ON discussions;
DROP TRIGGER IF EXISTS update_discussion_replies_updated_at ON discussion_replies;

-- ============================================================
-- 2. DROP FRIEND SYSTEM TRIGGERS
-- ============================================================

DROP TRIGGER IF EXISTS trigger_create_friendship_on_accept ON friend_requests;
DROP TRIGGER IF EXISTS trigger_create_friend_request_notification ON friend_requests;
DROP TRIGGER IF EXISTS trigger_cleanup_friendship_on_status_change ON friend_requests;

-- ============================================================
-- 3. DROP CHALLENGE TRIGGERS
-- ============================================================

DROP TRIGGER IF EXISTS trigger_update_challenges_updated_at ON challenges;
DROP TRIGGER IF EXISTS trigger_create_challenge_notification ON challenges;

-- ============================================================
-- 4. DROP OTHER TRIGGERS
-- ============================================================

DROP TRIGGER IF EXISTS trigger_update_user_preferences_updated_at ON user_preferences;
DROP TRIGGER IF EXISTS trigger_sync_quiz_points ON quizzes;
DROP TRIGGER IF EXISTS update_users_level_on_points_change ON users;
DROP TRIGGER IF EXISTS trigger_update_last_active_on_auth ON users;

-- ============================================================
-- 5. DROP ALL ASSOCIATED FUNCTIONS
-- ============================================================

-- Timestamp function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- User level functions
DROP FUNCTION IF EXISTS update_user_level();
DROP FUNCTION IF EXISTS calculate_user_level(INTEGER);

-- Friend system functions
DROP FUNCTION IF EXISTS create_friendship_on_accept();
DROP FUNCTION IF EXISTS create_friend_request_notification();
DROP FUNCTION IF EXISTS cleanup_friendship_on_status_change();

-- Challenge functions
DROP FUNCTION IF EXISTS update_challenges_updated_at();
DROP FUNCTION IF EXISTS create_challenge_notification();

-- User preferences function
DROP FUNCTION IF EXISTS update_user_preferences_updated_at();

-- Quiz sync function
DROP FUNCTION IF EXISTS sync_quiz_points_with_questions();

-- Auth tracking function
DROP FUNCTION IF EXISTS update_last_active_on_auth();
