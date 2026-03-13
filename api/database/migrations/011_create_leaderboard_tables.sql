-- Migration: Create leaderboard and achievements tables
-- This migration adds support for leaderboard functionality, achievements, and user rankings

-- Achievements table to store available achievements
CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(100) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    icon VARCHAR(100),
    color VARCHAR(50) DEFAULT '#FFD700',
    points_reward INTEGER DEFAULT 0,
    category VARCHAR(50) NOT NULL DEFAULT 'general',
    is_rare BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User achievements to track which achievements users have unlocked
CREATE TABLE user_achievements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID REFERENCES achievements(id) ON DELETE CASCADE,
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    points_awarded INTEGER DEFAULT 0,
    UNIQUE(user_id, achievement_id)
);

-- User category performance to track points by quiz category
CREATE TABLE user_category_performance (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    category_id VARCHAR(100) NOT NULL,
    category_name VARCHAR(255) NOT NULL,
    total_points INTEGER DEFAULT 0,
    quizzes_completed INTEGER DEFAULT 0,
    average_score DECIMAL(5,2) DEFAULT 0.0,
    best_score DECIMAL(5,2) DEFAULT 0.0,
    last_attempt_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, category_id)
);

-- Leaderboard snapshots for historical tracking (optional for future use)
CREATE TABLE leaderboard_snapshots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    snapshot_date DATE NOT NULL,
    period_type VARCHAR(20) NOT NULL, -- 'daily', 'weekly', 'monthly', 'yearly'
    rank_position INTEGER NOT NULL,
    total_points INTEGER NOT NULL,
    quizzes_completed INTEGER NOT NULL,
    average_score DECIMAL(5,2) NOT NULL,
    current_streak INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, snapshot_date, period_type)
);

-- User rank history to track rank changes over time
CREATE TABLE user_rank_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    period_type VARCHAR(20) NOT NULL, -- 'daily', 'weekly', 'monthly', 'alltime'
    old_rank INTEGER,
    new_rank INTEGER NOT NULL,
    rank_change INTEGER NOT NULL DEFAULT 0, -- positive for improvement, negative for decline
    points_at_change INTEGER NOT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_user_achievements_user_id ON user_achievements(user_id);
CREATE INDEX idx_user_achievements_achievement_id ON user_achievements(achievement_id);
CREATE INDEX idx_user_achievements_unlocked_at ON user_achievements(unlocked_at);

CREATE INDEX idx_user_category_performance_user_id ON user_category_performance(user_id);
CREATE INDEX idx_user_category_performance_category_id ON user_category_performance(category_id);
CREATE INDEX idx_user_category_performance_total_points ON user_category_performance(total_points DESC);

CREATE INDEX idx_leaderboard_snapshots_user_id ON leaderboard_snapshots(user_id);
CREATE INDEX idx_leaderboard_snapshots_date_period ON leaderboard_snapshots(snapshot_date, period_type);
CREATE INDEX idx_leaderboard_snapshots_rank ON leaderboard_snapshots(rank_position);

CREATE INDEX idx_user_rank_history_user_id ON user_rank_history(user_id);
CREATE INDEX idx_user_rank_history_period ON user_rank_history(period_type);
CREATE INDEX idx_user_rank_history_changed_at ON user_rank_history(changed_at);

-- Additional indexes on users table for leaderboard queries
CREATE INDEX idx_users_total_points ON users(total_points DESC);
CREATE INDEX idx_users_average_score ON users(average_score DESC);
CREATE INDEX idx_users_current_streak ON users(current_streak DESC);
CREATE INDEX idx_users_total_quizzes_completed ON users(total_quizzes_completed DESC);
CREATE INDEX idx_users_level ON users(level);

-- Triggers for updated_at columns
CREATE TRIGGER update_achievements_updated_at
    BEFORE UPDATE ON achievements
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_category_performance_updated_at
    BEFORE UPDATE ON user_category_performance
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert some initial achievements
INSERT INTO achievements (key, title, description, icon, color, points_reward, category, is_rare) VALUES
('first_win', 'First Victory', 'Complete your first quiz successfully', 'trophy', '#FFD700', 50, 'quiz', false),
('week_warrior', 'Week Warrior', 'Maintain a 7-day streak', 'fire', '#FF6347', 100, 'streak', false),
('perfect_score', 'Perfect Score', 'Score 100% on any quiz', 'star', '#9932CC', 150, 'score', true),
('social_butterfly', 'Social Butterfly', 'Add 5 friends', 'people', '#FF69B4', 75, 'social', false),
('quiz_master', 'Quiz Master', 'Complete 100 quizzes', 'crown', '#FFD700', 500, 'quiz', true),
('streak_legend', 'Streak Legend', 'Maintain a 30-day streak', 'fire', '#FF4500', 300, 'streak', true),
('tech_genius', 'Tech Genius', 'Score above 90% in 10 Technology quizzes', 'code', '#00CED1', 200, 'category', false),
('sports_expert', 'Sports Expert', 'Score above 90% in 10 Sports quizzes', 'sports', '#32CD32', 200, 'category', false),
('rising_star', 'Rising Star', 'Climb 10 positions in the leaderboard', 'trending-up', '#FFD700', 100, 'leaderboard', false),
('speed_demon', 'Speed Demon', 'Complete 5 quizzes in under 2 minutes each', 'zap', '#FF6347', 150, 'speed', false);

-- Function to update user statistics after quiz completion
CREATE OR REPLACE FUNCTION update_user_leaderboard_stats()
RETURNS TRIGGER AS $$
BEGIN
    -- This function will be called when quiz_attempts are updated/inserted
    -- For now, it's a placeholder for future automatic stat updates
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to calculate user level based on total points
CREATE OR REPLACE FUNCTION calculate_user_level(points INTEGER)
RETURNS VARCHAR(50) AS $$
BEGIN
    CASE
        WHEN points >= 5000 THEN RETURN 'Master';
        WHEN points >= 3000 THEN RETURN 'Expert';
        WHEN points >= 1500 THEN RETURN 'Advanced';
        WHEN points >= 500 THEN RETURN 'Intermediate';
        ELSE RETURN 'Beginner';
    END CASE;
END;
$$ LANGUAGE plpgsql;

-- Function to update user level when points change
CREATE OR REPLACE FUNCTION update_user_level()
RETURNS TRIGGER AS $$
BEGIN
    -- Update level based on new total_points
    NEW.level = calculate_user_level(NEW.total_points);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update user level when points change
CREATE TRIGGER update_users_level_on_points_change
    BEFORE UPDATE OF total_points ON users
    FOR EACH ROW
    WHEN (OLD.total_points IS DISTINCT FROM NEW.total_points)
    EXECUTE FUNCTION update_user_level();