-- Migration: Enable Row Level Security (RLS) on all tables
-- Date: 2025-01-15
-- Description: Enables RLS on all public tables to match Supabase production configuration
--              Note: No specific policies are defined yet - this migration only enables RLS
--              Access will be controlled by service role or future policy definitions

-- Enable RLS on user-related tables
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_preferences ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_achievements ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_category_performance ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_quiz_favorites ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_rank_history ENABLE ROW LEVEL SECURITY;

-- Enable RLS on quiz-related tables
ALTER TABLE quizzes ENABLE ROW LEVEL SECURITY;
ALTER TABLE questions ENABLE ROW LEVEL SECURITY;
ALTER TABLE quiz_attempts ENABLE ROW LEVEL SECURITY;
ALTER TABLE quiz_sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE quiz_statistics ENABLE ROW LEVEL SECURITY;
ALTER TABLE quiz_ratings ENABLE ROW LEVEL SECURITY;

-- Enable RLS on social/interaction tables
ALTER TABLE friend_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE friendships ENABLE ROW LEVEL SECURITY;
ALTER TABLE challenges ENABLE ROW LEVEL SECURITY;
ALTER TABLE notifications ENABLE ROW LEVEL SECURITY;

-- Enable RLS on discussion tables
ALTER TABLE discussions ENABLE ROW LEVEL SECURITY;
ALTER TABLE discussion_replies ENABLE ROW LEVEL SECURITY;
ALTER TABLE discussion_likes ENABLE ROW LEVEL SECURITY;
ALTER TABLE discussion_reply_likes ENABLE ROW LEVEL SECURITY;

-- Enable RLS on digest tables
ALTER TABLE digests ENABLE ROW LEVEL SECURITY;
ALTER TABLE digest_articles ENABLE ROW LEVEL SECURITY;

-- Enable RLS on lookup/reference tables
ALTER TABLE categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE difficulty_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE achievements ENABLE ROW LEVEL SECURITY;
ALTER TABLE notification_frequencies ENABLE ROW LEVEL SECURITY;
ALTER TABLE app_settings ENABLE ROW LEVEL SECURITY;

-- Enable RLS on leaderboard tables
ALTER TABLE leaderboard_snapshots ENABLE ROW LEVEL SECURITY;

-- Note: The migrations table is intentionally excluded from RLS
-- as it needs to be accessible for migration tracking

-- Important: With RLS enabled and no policies defined, all access is blocked
-- except when using the service role key which bypasses RLS.
-- Future migrations should add specific policies as needed.

COMMENT ON TABLE users IS 'RLS enabled - access controlled by service role or future policies';
COMMENT ON TABLE quizzes IS 'RLS enabled - access controlled by service role or future policies';
COMMENT ON TABLE questions IS 'RLS enabled - access controlled by service role or future policies';
