-- Additional seed data for QuizNinja application (only new data not in previous migrations)

-- Insert default data for interests
INSERT INTO interests (id, name, description, icon_name, color_hex) VALUES
('general_knowledge', 'General Knowledge', 'Broad range of topics covering various subjects', 'lightbulb', '#FF6B6B'),
('science', 'Science', 'Physics, chemistry, biology, and earth sciences', 'flask', '#4ECDC4'),
('history', 'History', 'Historical events, figures, and civilizations', 'clock', '#45B7D1'),
('sports', 'Sports', 'Various sports, athletes, and competitions', 'trophy', '#96CEB4'),
('technology', 'Technology', 'Computing, programming, and modern tech', 'laptop', '#FFEAA7'),
('movies_tv', 'Movies & TV', 'Films, television shows, and entertainment', 'film', '#DDA0DD'),
('music', 'Music', 'Artists, songs, genres, and music theory', 'music', '#FFB6C1'),
('geography', 'Geography', 'Countries, capitals, landmarks, and maps', 'globe', '#87CEEB'),
('literature', 'Literature', 'Books, authors, poetry, and literary works', 'book', '#F4A460'),
('art', 'Art', 'Visual arts, artists, and art movements', 'palette', '#DA70D6');

-- Insert default data for difficulty levels
INSERT INTO difficulty_levels (id, name, description, icon_name, background_color_hex) VALUES
('beginner', 'Beginner', 'Easy questions for newcomers', 'star', '#E8F5E8'),
('intermediate', 'Intermediate', 'Moderate difficulty questions', 'star-half', '#FFF4E6'),
('advanced', 'Advanced', 'Challenging questions for experts', 'star-fill', '#FFE6E6'),
('expert', 'Expert', 'Very difficult questions for masters', 'trophy', '#F0E6FF');

-- Insert default data for notification frequencies
INSERT INTO notification_frequencies (id, name, description, icon_name) VALUES
('never', 'Never', 'No notifications', 'bell-slash'),
('daily', 'Daily', 'Once per day', 'bell'),
('weekly', 'Weekly', 'Once per week', 'calendar'),
('bi_weekly', 'Bi-weekly', 'Every two weeks', 'calendar-check'),
('monthly', 'Monthly', 'Once per month', 'calendar-month');

-- Additional notification frequencies (beyond the 5 already in migration 003)
INSERT INTO notification_frequencies (id, name, description, icon_name) VALUES
('twice_daily', 'Twice Daily', 'Morning and evening reminders', 'bell-ring'),
('every_3_days', 'Every 3 Days', 'Notifications every three days', 'calendar-plus')
ON CONFLICT (id) DO NOTHING;


-- Initial app settings configuration (app_settings table already exists from migration 001)
INSERT INTO app_settings (key, value, description) VALUES
-- App version and metadata
('app_version', '1.0.0', 'Current application version'),
('api_version', '1.0.0', 'Current API version'),
('min_supported_version', '1.0.0', 'Minimum supported app version'),

-- Quiz configuration
('default_questions_per_quiz', '10', 'Default number of questions per quiz'),
('max_questions_per_quiz', '50', 'Maximum questions allowed per quiz'),
('min_questions_per_quiz', '5', 'Minimum questions required per quiz'),
('quiz_time_limit_seconds', '300', 'Default time limit for quizzes in seconds (5 minutes)'),
('quiz_time_limit_enabled', 'true', 'Whether quiz time limits are enabled by default'),

-- Scoring system
('points_per_correct_answer', '10', 'Points awarded for each correct answer'),
('bonus_points_streak_multiplier', '1.5', 'Multiplier for streak bonuses'),
('streak_bonus_threshold', '5', 'Number of correct answers in a row to trigger streak bonus'),
('perfect_quiz_bonus', '50', 'Extra points for completing a quiz with 100% accuracy'),

-- User progression
('level_up_points_threshold', '1000', 'Points needed to advance to the next level'),
('level_up_points_increment', '500', 'Additional points needed for each subsequent level'),
('max_user_level', '100', 'Maximum user level achievable'),

-- Feature flags
('dark_mode_enabled', 'true', 'Whether dark mode is available'),
('offline_mode_enabled', 'true', 'Whether offline quiz mode is available'),

-- Maintenance
('maintenance_mode', 'false', 'Whether the app is in maintenance mode'),
('maintenance_message', 'We are currently performing scheduled maintenance. Please try again later.', 'Message shown during maintenance'),
('force_update_required', 'false', 'Whether users must update to continue using the app')

ON CONFLICT (key) DO UPDATE SET
    value = EXCLUDED.value,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;
