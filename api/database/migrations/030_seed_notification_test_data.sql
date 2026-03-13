-- Migration: Seed notification test data
-- Description: Creates test users and comprehensive notification data to verify database-to-UI flow

-- Create test users for realistic notifications (if they don't already exist)
INSERT INTO users (id, email, name, password_hash, level, total_points, current_streak, is_test_data) VALUES
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'alice.johnson@testuser.com',
    'Alice Johnson',
    '$2a$12$test.hash.for.alice.johnson.test.user.placeholder',
    'Intermediate',
    1250,
    5,
    true
),
(
    'a2b2c2d2-e2f2-4002-8002-222222222222',
    'bob.smith@testuser.com',
    'Bob Smith',
    '$2a$12$test.hash.for.bob.smith.test.user.placeholder',
    'Advanced',
    2340,
    12,
    true
),
(
    'a3b3c3d3-e3f3-4003-8003-333333333333',
    'carol.davis@testuser.com',
    'Carol Davis',
    '$2a$12$test.hash.for.carol.davis.test.user.placeholder',
    'Expert',
    3890,
    8,
    true
),
(
    'a4b4c4d4-e4f4-4004-8004-444444444444',
    'david.wilson@testuser.com',
    'David Wilson',
    '$2a$12$test.hash.for.david.wilson.test.user.placeholder',
    'Beginner',
    650,
    3,
    true
),
(
    'a5b5c5d5-e5f5-4005-8005-555555555555',
    'eva.martinez@testuser.com',
    'Eva Martinez',
    '$2a$12$test.hash.for.eva.martinez.test.user.placeholder',
    'Advanced',
    2890,
    7,
    true
)
ON CONFLICT (id) DO UPDATE SET
    email = EXCLUDED.email,
    name = EXCLUDED.name,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- Create comprehensive notification test data covering all types
INSERT INTO notifications (user_id, type, title, message, data, related_user_id, related_entity_id, related_entity_type, is_read, created_at, read_at, is_test_data) VALUES

-- UNREAD NOTIFICATIONS (Today)
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'friend_request',
    'New Friend Request',
    'Bob Smith sent you a friend request!',
    '{"requester_name": "Bob Smith", "message": "Hey! I saw your amazing quiz scores. Let''s be friends!"}',
    'a2b2c2d2-e2f2-4002-8002-222222222222',
    uuid_generate_v4(),
    'friend_request',
    false,
    CURRENT_TIMESTAMP - INTERVAL '2 hours',
    NULL,
    true
),
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'challenge_received',
    'New Challenge Received',
    'Carol Davis challenged you to a quiz!',
    '{"challenger_name": "Carol Davis", "message": "Think you can beat my Science quiz score?", "expires_at": "2025-10-08T00:00:00Z"}',
    'a3b3c3d3-e3f3-4003-8003-333333333333',
    uuid_generate_v4(),
    'challenge',
    false,
    CURRENT_TIMESTAMP - INTERVAL '4 hours',
    NULL,
    true
),
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'achievement_unlocked',
    'Achievement Unlocked!',
    'You unlocked: Quiz Master',
    '{"achievement_key": "quiz_master", "achievement_title": "Quiz Master", "achievement_description": "Complete 10 quizzes with 90% or higher score", "points_awarded": 100}',
    NULL,
    uuid_generate_v4(),
    'achievement',
    false,
    CURRENT_TIMESTAMP - INTERVAL '6 hours',
    NULL,
    true
),

-- READ NOTIFICATIONS (Today)
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'general',
    'New Quizzes Available',
    '5 new quizzes have been added in your favorite categories',
    '{"new_quizzes_count": 5, "categories": ["Science", "History", "Technology"]}',
    NULL,
    NULL,
    NULL,
    true,
    CURRENT_TIMESTAMP - INTERVAL '8 hours',
    CURRENT_TIMESTAMP - INTERVAL '7 hours',
    true
),

-- UNREAD NOTIFICATIONS (Yesterday)
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'friend_accepted',
    'Friend Request Accepted',
    'Eva Martinez accepted your friend request!',
    '{"accepter_name": "Eva Martinez"}',
    'a5b5c5d5-e5f5-4005-8005-555555555555',
    uuid_generate_v4(),
    'friend_request',
    false,
    CURRENT_TIMESTAMP - INTERVAL '1 day 3 hours',
    NULL,
    true
),
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'challenge_completed',
    'Challenge Completed',
    'You won the challenge against David Wilson!',
    '{"opponent_name": "David Wilson", "your_score": 85, "opponent_score": 72, "result": "won"}',
    'a4b4c4d4-e4f4-4004-8004-444444444444',
    uuid_generate_v4(),
    'challenge',
    false,
    CURRENT_TIMESTAMP - INTERVAL '1 day 8 hours',
    NULL,
    true
),

-- READ NOTIFICATIONS (Yesterday)
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'system_announcement',
    'App Update Available',
    'Version 2.1.0 is now available with new features and improvements',
    '{"version": "2.1.0", "update_url": "https://play.google.com/store", "features": ["Improved UI", "Bug fixes", "New quiz categories"]}',
    NULL,
    NULL,
    NULL,
    true,
    CURRENT_TIMESTAMP - INTERVAL '1 day 12 hours',
    CURRENT_TIMESTAMP - INTERVAL '1 day 10 hours',
    true
),

-- READ NOTIFICATIONS (2 days ago)
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'achievement_unlocked',
    'Achievement Unlocked!',
    'You unlocked: Streak Master',
    '{"achievement_key": "streak_master", "achievement_title": "Streak Master", "achievement_description": "Maintain a 10-day quiz streak", "points_awarded": 150}',
    NULL,
    uuid_generate_v4(),
    'achievement',
    true,
    CURRENT_TIMESTAMP - INTERVAL '2 days 5 hours',
    CURRENT_TIMESTAMP - INTERVAL '2 days 4 hours',
    true
),
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'challenge_accepted',
    'Challenge Accepted',
    'Bob Smith accepted your challenge!',
    '{"accepter_name": "Bob Smith"}',
    'a2b2c2d2-e2f2-4002-8002-222222222222',
    uuid_generate_v4(),
    'challenge',
    true,
    CURRENT_TIMESTAMP - INTERVAL '2 days 10 hours',
    CURRENT_TIMESTAMP - INTERVAL '2 days 9 hours',
    true
),

-- READ NOTIFICATIONS (3 days ago)
(
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'general',
    'Weekly Digest Available',
    'Your weekly quiz performance summary is ready',
    '{"total_quizzes": 7, "average_score": 87.5, "best_category": "Science", "improvement": "+12%"}',
    NULL,
    NULL,
    NULL,
    true,
    CURRENT_TIMESTAMP - INTERVAL '3 days 2 hours',
    CURRENT_TIMESTAMP - INTERVAL '3 days 1 hour',
    true
),

-- Additional notifications for other test users to make data more realistic
(
    'a2b2c2d2-e2f2-4002-8002-222222222222',
    'friend_request',
    'New Friend Request',
    'Alice Johnson sent you a friend request!',
    '{"requester_name": "Alice Johnson", "message": "Great job on the History quiz!"}',
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    uuid_generate_v4(),
    'friend_request',
    false,
    CURRENT_TIMESTAMP - INTERVAL '1 hour',
    NULL,
    true
),
(
    'a3b3c3d3-e3f3-4003-8003-333333333333',
    'achievement_unlocked',
    'Achievement Unlocked!',
    'You unlocked: Science Genius',
    '{"achievement_key": "science_genius", "achievement_title": "Science Genius", "achievement_description": "Score 100% on 5 Science quizzes", "points_awarded": 200}',
    NULL,
    uuid_generate_v4(),
    'achievement',
    false,
    CURRENT_TIMESTAMP - INTERVAL '3 hours',
    NULL,
    true
);

-- Insert a few mock achievements to make the achievement notifications valid
INSERT INTO achievements (key, title, description, icon, color, points_reward, category, is_test_data) VALUES
(
    'quiz_master',
    'Quiz Master',
    'Complete 10 quizzes with 90% or higher score',
    'trophy',
    '#FFD700',
    100,
    'quiz_performance',
    true
),
(
    'streak_master',
    'Streak Master',
    'Maintain a 10-day quiz streak',
    'fire',
    '#FF6B35',
    150,
    'consistency',
    true
),
(
    'science_genius',
    'Science Genius',
    'Score 100% on 5 Science quizzes',
    'atom',
    '#4ECDC4',
    200,
    'subject_mastery',
    true
)
ON CONFLICT (key) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;