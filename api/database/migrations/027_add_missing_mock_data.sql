-- Migration: Add missing mock data to restore original API behavior
-- Up Migration
-- Description: Adds missing interests and app settings that were previously hardcoded in mock implementations

-- Add missing interests that were in the original mock implementation
INSERT INTO interests (id, name, description, icon_name, color_hex, is_test_data) VALUES
('biology', 'Biology', 'Life sciences and organisms', 'leaf', '#2ECC71', true),
('chemistry', 'Chemistry', 'Chemical elements and reactions', 'flask', '#E74C3C', true),
('physics', 'Physics', 'Physical laws and phenomena', 'atom', '#3498DB', true),
('football', 'Football', 'American football trivia', 'football', '#8E44AD', true),
('basketball', 'Basketball', 'Basketball trivia and stats', 'basketball', '#E67E22', true)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon_name = EXCLUDED.icon_name,
    color_hex = EXCLUDED.color_hex,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- Add missing app settings that were in the original mock implementation
INSERT INTO app_settings (key, value, description, updated_at) VALUES
('app_name', 'QuizNinja', 'Application name displayed to users', CURRENT_TIMESTAMP),
('quiz_categories_enabled', 'true', 'Whether quiz categories feature is enabled', CURRENT_TIMESTAMP),
('leaderboard_enabled', 'true', 'Whether leaderboard feature is enabled', CURRENT_TIMESTAMP),
('achievements_enabled', 'true', 'Whether achievements feature is enabled', CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET
    value = EXCLUDED.value,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- Update max_questions_per_quiz to match original mock value (was 50, should be 20)
UPDATE app_settings
SET value = '20',
    description = 'Maximum questions allowed per quiz (updated to match original mock)',
    updated_at = CURRENT_TIMESTAMP
WHERE key = 'max_questions_per_quiz';

-- Update default_questions_per_quiz to match original mock value expectations
UPDATE app_settings
SET value = '10',
    description = 'Default number of questions per quiz',
    updated_at = CURRENT_TIMESTAMP
WHERE key = 'default_questions_per_quiz';

-- Add default_quiz_duration setting (was in original mock as default_quiz_duration: 300)
INSERT INTO app_settings (key, value, description, updated_at) VALUES
('default_quiz_duration', '300', 'Default quiz duration in seconds', CURRENT_TIMESTAMP)
ON CONFLICT (key) DO UPDATE SET
    value = EXCLUDED.value,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;