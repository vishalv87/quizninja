-- Migration: Reseed all deleted data
-- Date: 2025-10-05
-- Description: Comprehensive reseeding of all deleted seed data from previous migrations

-- ============================================================================
-- CORE LOOKUP DATA (from migration 004_seed_initial_data.sql)
-- ============================================================================

-- Insert default data for interests
INSERT INTO interests (id, name, description, icon_name, color_hex, is_test_data) VALUES
('general_knowledge', 'General Knowledge', 'Broad range of topics covering various subjects', 'lightbulb', '#FF6B6B', true),
('science', 'Science', 'Physics, chemistry, biology, and earth sciences', 'flask', '#4ECDC4', true),
('history', 'History', 'Historical events, figures, and civilizations', 'clock', '#45B7D1', true),
('sports', 'Sports', 'Various sports, athletes, and competitions', 'trophy', '#96CEB4', true),
('technology', 'Technology', 'Computing, programming, and modern tech', 'laptop', '#FFEAA7', true),
('movies_tv', 'Movies & TV', 'Films, television shows, and entertainment', 'film', '#DDA0DD', true),
('music', 'Music', 'Artists, songs, genres, and music theory', 'music', '#FFB6C1', true),
('geography', 'Geography', 'Countries, capitals, landmarks, and maps', 'globe', '#87CEEB', true),
('literature', 'Literature', 'Books, authors, poetry, and literary works', 'book', '#F4A460', true),
('art', 'Art', 'Visual arts, artists, and art movements', 'palette', '#DA70D6', true)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon_name = EXCLUDED.icon_name,
    color_hex = EXCLUDED.color_hex,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- Insert default data for difficulty levels
INSERT INTO difficulty_levels (id, name, description, icon_name, background_color_hex, is_test_data) VALUES
('beginner', 'Beginner', 'Easy questions for newcomers', 'star', '#E8F5E8', true),
('intermediate', 'Intermediate', 'Moderate difficulty questions', 'star-half', '#FFF4E6', true),
('advanced', 'Advanced', 'Challenging questions for experts', 'star-fill', '#FFE6E6', true),
('expert', 'Expert', 'Very difficult questions for masters', 'trophy', '#F0E6FF', true)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon_name = EXCLUDED.icon_name,
    background_color_hex = EXCLUDED.background_color_hex,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- Insert default data for notification frequencies
INSERT INTO notification_frequencies (id, name, description, icon_name, is_test_data) VALUES
('never', 'Never', 'No notifications', 'bell-slash', true),
('daily', 'Daily', 'Once per day', 'bell', true),
('weekly', 'Weekly', 'Once per week', 'calendar', true),
('bi_weekly', 'Bi-weekly', 'Every two weeks', 'calendar-check', true),
('monthly', 'Monthly', 'Once per month', 'calendar-month', true),
('twice_daily', 'Twice Daily', 'Morning and evening reminders', 'bell-ring', true),
('every_3_days', 'Every 3 Days', 'Notifications every three days', 'calendar-plus', true)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon_name = EXCLUDED.icon_name,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- Initial app settings configuration
INSERT INTO app_settings (key, value, description) VALUES
-- App version and metadata
('app_version', '1.0.0', 'Current application version'),
('api_version', '1.0.0', 'Current API version'),
('min_supported_version', '1.0.0', 'Minimum supported app version'),

-- Quiz configuration
('default_questions_per_quiz', '10', 'Default number of questions per quiz'),
('max_questions_per_quiz', '20', 'Maximum questions allowed per quiz (updated to match original mock)'),
('min_questions_per_quiz', '5', 'Minimum questions required per quiz'),
('quiz_time_limit_seconds', '300', 'Default time limit for quizzes in seconds (5 minutes)'),
('quiz_time_limit_enabled', 'true', 'Whether quiz time limits are enabled by default'),
('default_quiz_duration', '300', 'Default quiz duration in seconds'),

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
('app_name', 'QuizNinja', 'Application name displayed to users'),
('quiz_categories_enabled', 'true', 'Whether quiz categories feature is enabled'),
('leaderboard_enabled', 'true', 'Whether leaderboard feature is enabled'),
('achievements_enabled', 'true', 'Whether achievements feature is enabled'),

-- Maintenance
('maintenance_mode', 'false', 'Whether the app is in maintenance mode'),
('maintenance_message', 'We are currently performing scheduled maintenance. Please try again later.', 'Message shown during maintenance'),
('force_update_required', 'false', 'Whether users must update to continue using the app')

ON CONFLICT (key) DO UPDATE SET
    value = EXCLUDED.value,
    description = EXCLUDED.description,
    updated_at = CURRENT_TIMESTAMP;

-- ============================================================================
-- ADDITIONAL INTERESTS (from migration 027_add_missing_mock_data.sql)
-- ============================================================================

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

-- ============================================================================
-- SAMPLE QUIZ CONTENT (from migration 007_seed_quiz_data.sql)
-- ============================================================================

-- Sample Quiz 1: General Knowledge
INSERT INTO quizzes (id, title, description, category_id, difficulty, total_questions, time_limit_minutes, points, created_by, is_featured, is_active, is_public, tags, is_test_data, created_at, updated_at)
VALUES (
    'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'General Knowledge Basics',
    'Test your basic knowledge across various topics including geography, history, science, and culture.',
    'general_knowledge',
    'beginner',
    10,
    15,
    100,
    (SELECT id FROM users ORDER BY created_at LIMIT 1),
    true,
    true,
    true,
    '{"trivia", "basics", "general"}',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    category_id = EXCLUDED.category_id,
    difficulty = EXCLUDED.difficulty,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- Questions for General Knowledge quiz
INSERT INTO questions (id, quiz_id, question_text, question_type, options, correct_answer, explanation, order_index, is_test_data, created_at)
VALUES
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the capital of France?', 'multipleChoice', '{"London", "Berlin", "Paris", "Madrid"}', 'Paris', 'Paris is the capital and most populous city of France.', 1, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which planet is closest to the Sun?', 'multipleChoice', '{"Venus", "Mercury", "Earth", "Mars"}', 'Mercury', 'Mercury is the smallest planet in our solar system and closest to the Sun.', 2, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'The Great Wall of China was built primarily to keep out which people?', 'multipleChoice', '{"Mongols", "Japanese", "Russians", "British"}', 'Mongols', 'The Great Wall was built to protect Chinese states from invasions by nomadic groups, particularly the Mongols.', 3, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the chemical symbol for gold?', 'multipleChoice', '{"Go", "Gd", "Au", "Ag"}', 'Au', 'Au comes from the Latin word "aurum" meaning gold.', 4, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which ocean is the largest?', 'multipleChoice', '{"Atlantic", "Pacific", "Indian", "Arctic"}', 'Pacific', 'The Pacific Ocean is the largest and deepest ocean on Earth.', 5, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Who painted the Mona Lisa?', 'multipleChoice', '{"Vincent van Gogh", "Pablo Picasso", "Leonardo da Vinci", "Michelangelo"}', 'Leonardo da Vinci', 'Leonardo da Vinci painted the Mona Lisa between 1503 and 1519.', 6, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the largest mammal in the world?', 'multipleChoice', '{"African Elephant", "Blue Whale", "Giraffe", "Polar Bear"}', 'Blue Whale', 'Blue whales are the largest animals ever known to have lived on Earth.', 7, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'In which year did World War II end?', 'multipleChoice', '{"1944", "1945", "1946", "1947"}', '1945', 'World War II ended in 1945 with the surrender of Japan in September.', 8, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the smallest country in the world?', 'multipleChoice', '{"Monaco", "Vatican City", "San Marino", "Malta"}', 'Vatican City', 'Vatican City is the smallest sovereign nation in the world by both area and population.', 9, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which element has the chemical symbol "O"?', 'multipleChoice', '{"Osmium", "Oxygen", "Ozone", "Oxide"}', 'Oxygen', 'Oxygen is a chemical element with the symbol O and atomic number 8.', 10, true, CURRENT_TIMESTAMP);

-- Sample Quiz 2: Technology & Programming
INSERT INTO quizzes (id, title, description, category_id, difficulty, total_questions, time_limit_minutes, points, created_by, is_featured, is_active, is_public, tags, is_test_data, created_at, updated_at)
VALUES (
    'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Tech & Programming Fundamentals',
    'Test your knowledge of technology concepts, programming languages, and computer science basics.',
    'technology',
    'intermediate',
    8,
    12,
    120,
    (SELECT id FROM users ORDER BY created_at LIMIT 1),
    true,
    true,
    true,
    '{"programming", "technology", "computers"}',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    category_id = EXCLUDED.category_id,
    difficulty = EXCLUDED.difficulty,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- Questions for Technology quiz
INSERT INTO questions (id, quiz_id, question_text, question_type, options, correct_answer, explanation, order_index, is_test_data, created_at)
VALUES
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What does "HTML" stand for?', 'multipleChoice', '{"Hypertext Markup Language", "High-Tech Modern Language", "Home Tool Markup Language", "Hyperlink and Text Markup Language"}', 'Hypertext Markup Language', 'HTML is the standard markup language for creating web pages.', 1, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which programming language is known as the "language of the web"?', 'multipleChoice', '{"Python", "JavaScript", "Java", "C++"}', 'JavaScript', 'JavaScript is primarily known as the scripting language for web pages.', 2, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What does "CPU" stand for?', 'multipleChoice', '{"Central Processing Unit", "Computer Processing Unit", "Central Program Unit", "Computer Program Unit"}', 'Central Processing Unit', 'The CPU is the primary component of a computer that performs most processing.', 3, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which company developed the Python programming language?', 'multipleChoice', '{"Google", "Microsoft", "Guido van Rossum", "Facebook"}', 'Guido van Rossum', 'Python was created by Guido van Rossum and first released in 1991.', 4, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What does "API" stand for?', 'multipleChoice', '{"Application Programming Interface", "Advanced Programming Interface", "Application Program Integration", "Advanced Program Integration"}', 'Application Programming Interface', 'An API defines how software components should interact.', 5, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which of these is NOT a database management system?', 'multipleChoice', '{"MySQL", "PostgreSQL", "MongoDB", "Photoshop"}', 'Photoshop', 'Photoshop is image editing software, not a database management system.', 6, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What does "Git" primarily help with?', 'multipleChoice', '{"Image editing", "Version control", "Database management", "Web hosting"}', 'Version control', 'Git is a distributed version control system for tracking changes in source code.', 7, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which protocol is used for secure web browsing?', 'multipleChoice', '{"HTTP", "HTTPS", "FTP", "SMTP"}', 'HTTPS', 'HTTPS provides encrypted communication between web browsers and servers.', 8, true, CURRENT_TIMESTAMP);

-- Sample Quiz 3: Science & Nature
INSERT INTO quizzes (id, title, description, category_id, difficulty, total_questions, time_limit_minutes, points, created_by, is_featured, is_active, is_public, tags, is_test_data, created_at, updated_at)
VALUES (
    'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'Science & Nature Quiz',
    'Explore the wonders of science, biology, physics, and the natural world around us.',
    'science',
    'intermediate',
    10,
    15,
    150,
    (SELECT id FROM users ORDER BY created_at LIMIT 1),
    false,
    true,
    true,
    '{"science", "nature", "biology", "physics"}',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    category_id = EXCLUDED.category_id,
    difficulty = EXCLUDED.difficulty,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- Questions for Science quiz
INSERT INTO questions (id, quiz_id, question_text, question_type, options, correct_answer, explanation, order_index, is_test_data, created_at)
VALUES
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the speed of light in a vacuum?', 'multipleChoice', '{"300,000 km/s", "299,792,458 m/s", "150,000 km/s", "186,000 miles/s"}', '299,792,458 m/s', 'The speed of light in a vacuum is exactly 299,792,458 meters per second.', 1, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which gas makes up about 78% of Earth''s atmosphere?', 'multipleChoice', '{"Oxygen", "Carbon Dioxide", "Nitrogen", "Argon"}', 'Nitrogen', 'Nitrogen makes up approximately 78% of Earth''s atmosphere by volume.', 2, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the powerhouse of the cell?', 'multipleChoice', '{"Nucleus", "Mitochondria", "Ribosome", "Chloroplast"}', 'Mitochondria', 'Mitochondria are often called the powerhouse of the cell because they generate ATP.', 3, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the hardest natural substance on Earth?', 'multipleChoice', '{"Gold", "Iron", "Diamond", "Granite"}', 'Diamond', 'Diamond is the hardest known natural material on the Mohs scale.', 4, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'How many chambers does a human heart have?', 'multipleChoice', '{"2", "3", "4", "5"}', '4', 'The human heart has four chambers: two atria and two ventricles.', 5, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the most abundant element in the universe?', 'multipleChoice', '{"Oxygen", "Carbon", "Hydrogen", "Helium"}', 'Hydrogen', 'Hydrogen is the most abundant element in the observable universe.', 6, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'At what temperature does water boil at sea level?', 'multipleChoice', '{"90°C", "95°C", "100°C", "105°C"}', '100°C', 'Water boils at 100°C (212°F) at sea level atmospheric pressure.', 7, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which planet is known as the Red Planet?', 'multipleChoice', '{"Venus", "Mars", "Jupiter", "Saturn"}', 'Mars', 'Mars appears red due to iron oxide (rust) on its surface.', 8, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What type of animal is a Komodo dragon?', 'multipleChoice', '{"Snake", "Lizard", "Crocodile", "Turtle"}', 'Lizard', 'The Komodo dragon is the largest living species of lizard.', 9, true, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the chemical formula for water?', 'multipleChoice', '{"H2O", "CO2", "NaCl", "CH4"}', 'H2O', 'Water is composed of two hydrogen atoms and one oxygen atom (H2O).', 10, true, CURRENT_TIMESTAMP);

-- Create initial statistics for each quiz
INSERT INTO quiz_statistics (quiz_id, total_attempts, total_completions, average_score, average_time_seconds, difficulty_rating, popularity_score, is_test_data, updated_at)
VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 0, 0, 0.0, 0, 0.0, 0, true, CURRENT_TIMESTAMP),
('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 0, 0, 0.0, 0, 0.0, 0, true, CURRENT_TIMESTAMP),
('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 0, 0, 0.0, 0, 0.0, 0, true, CURRENT_TIMESTAMP)
ON CONFLICT (quiz_id) DO UPDATE SET
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- ============================================================================
-- ACHIEVEMENT LOOKUP DATA (from migration 030_seed_notification_test_data.sql)
-- ============================================================================

-- Insert achievement lookup data
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

-- ============================================================================
-- DIGEST CONTENT DATA (from migration 018_seed_digest_dummy_data.sql)
-- ============================================================================

-- Create today's digest
INSERT INTO digests (date, title, summary, is_dummy, is_test_data) VALUES
(CURRENT_DATE, 'Daily News Digest', 'Stay informed with the most important news from around the world.', true, true)
ON CONFLICT (date) DO UPDATE SET
    title = EXCLUDED.title,
    summary = EXCLUDED.summary,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- Get today's digest ID and insert dummy articles
DO $$
DECLARE
    today_digest_id UUID;
BEGIN
    SELECT id INTO today_digest_id FROM digests WHERE date = CURRENT_DATE;

    -- Insert dummy articles for today's digest
    INSERT INTO digest_articles (
        digest_id, title, content, summary, source, author, published_at,
        category, image_url, external_url, read_time_minutes,
        is_breaking, is_hot, is_dummy, is_test_data
    ) VALUES
    (
        today_digest_id,
        'AI Revolution Transforms Healthcare Diagnosis',
        'Revolutionary artificial intelligence algorithms are transforming the way medical professionals diagnose diseases, with new systems achieving 95% accuracy rates in detecting various conditions. These advanced AI models have been trained on millions of medical cases and can now identify patterns that human doctors might miss.

The breakthrough technology combines machine learning with advanced imaging techniques to provide real-time analysis of medical scans, blood tests, and patient symptoms. Leading hospitals worldwide are already implementing these systems, reporting significant improvements in early disease detection.

Dr. Sarah Chen, Chief Medical Officer at the Stanford AI Lab, explains: "This technology doesn''t replace doctors but enhances their capabilities. It''s like having a highly experienced specialist available 24/7 to provide a second opinion."

The AI diagnostic tools have shown particular promise in oncology, cardiology, and neurology. Early trials suggest that cancer detection rates have improved by 40% when AI assistance is used alongside traditional diagnostic methods.

While the technology promises revolutionary changes in healthcare, experts emphasize the importance of maintaining human oversight and ensuring patient privacy in AI-driven medical systems.',
        'New AI algorithms are revolutionizing medical diagnosis with 95% accuracy rates, enhancing doctor capabilities and improving early disease detection.',
        'TechHealth News',
        'Dr. Michael Rodriguez',
        NOW() - INTERVAL '2 hours',
        'Technology',
        'https://example.com/ai-healthcare.jpg',
        'https://techhealth.com/ai-diagnosis-breakthrough',
        4,
        true,
        false,
        true,
        true
    ),
    (
        today_digest_id,
        'Climate Summit 2025: Historic Agreement Reached',
        'World leaders at the Climate Summit 2025 have reached a groundbreaking agreement that commits all participating nations to unprecedented climate action. The agreement, signed by representatives from 195 countries, sets binding targets for carbon emission reductions and renewable energy adoption.

The historic pact includes a commitment to achieve net-zero emissions by 2035, a timeline that is five years ahead of previous international agreements. Additionally, developed nations have pledged $500 billion in climate aid to support developing countries in their transition to clean energy.

Key provisions of the agreement include:
• Mandatory 60% reduction in carbon emissions by 2030
• Complete phase-out of coal power by 2032
• Investment of $2 trillion in renewable energy infrastructure
• Protection of 50% of Earth''s land and oceans by 2030

UN Secretary-General Maria Santos called the agreement "a turning point in human history" and emphasized that immediate action is crucial for implementation.

Environmental groups have praised the ambitious targets while noting that success will depend on rigorous monitoring and enforcement mechanisms.',
        'World leaders commit to unprecedented climate action with binding agreements for net-zero emissions by 2035.',
        'Global Climate Network',
        'Elena Petrov',
        NOW() - INTERVAL '4 hours',
        'Environment',
        'https://example.com/climate-summit.jpg',
        'https://climatenews.org/summit-2025-agreement',
        5,
        false,
        true,
        true,
        true
    ),
    (
        today_digest_id,
        'Tech Giants Report Record Quarterly Earnings',
        'Major technology companies have announced record-breaking quarterly earnings, with Apple, Google, and Microsoft exceeding analyst expectations by significant margins. The strong performance comes amid continued growth in cloud computing, AI services, and digital transformation initiatives.

Apple reported quarterly revenue of $125 billion, driven by strong iPhone sales and services growth. The company''s Services division, which includes the App Store and Apple Pay, reached an all-time high of $25 billion in revenue.

Google''s parent company Alphabet posted revenues of $95 billion, with Google Cloud growing 35% year-over-year. The company''s AI investments are showing strong returns, particularly in search and advertising technologies.

Microsoft continued its dominance in enterprise software with $78 billion in quarterly revenue. Azure cloud services grew 42%, solidifying Microsoft''s position as a leader in business cloud solutions.

The exceptional performance has led to increased investor confidence in the technology sector, with tech stocks reaching new highs. Analysts predict continued growth as businesses worldwide accelerate their digital transformation efforts.',
        'Apple, Google, and Microsoft exceed expectations with strong quarterly results driven by cloud computing and AI growth.',
        'Financial Times Tech',
        'Robert Kim',
        NOW() - INTERVAL '6 hours',
        'Business',
        'https://example.com/tech-earnings.jpg',
        'https://fintech.com/q4-tech-earnings-record',
        3,
        false,
        false,
        true,
        true
    );

END $$;

-- Create a few historical digests for testing pagination
INSERT INTO digests (date, title, summary, is_dummy, is_test_data) VALUES
(CURRENT_DATE - INTERVAL '1 day', 'Yesterday''s News Digest', 'The most important stories from yesterday.', true, true),
(CURRENT_DATE - INTERVAL '2 days', 'Weekend News Digest', 'Key developments from the weekend.', true, true),
(CURRENT_DATE - INTERVAL '3 days', 'Weekly News Roundup', 'Major stories from the past week.', true, true)
ON CONFLICT (date) DO UPDATE SET
    title = EXCLUDED.title,
    summary = EXCLUDED.summary,
    is_test_data = EXCLUDED.is_test_data,
    updated_at = CURRENT_TIMESTAMP;

-- Add a few articles to historical digests
DO $$
DECLARE
    yesterday_digest_id UUID;
    weekend_digest_id UUID;
BEGIN
    SELECT id INTO yesterday_digest_id FROM digests WHERE date = CURRENT_DATE - INTERVAL '1 day';
    SELECT id INTO weekend_digest_id FROM digests WHERE date = CURRENT_DATE - INTERVAL '2 days';

    -- Yesterday's articles
    INSERT INTO digest_articles (
        digest_id, title, content, summary, source, author, published_at,
        category, read_time_minutes, is_dummy, is_test_data
    ) VALUES
    (
        yesterday_digest_id,
        'Global Markets Show Strong Recovery',
        'Financial markets worldwide have shown remarkable resilience following recent economic uncertainties...',
        'Stock markets rally as investors show confidence in global economic recovery.',
        'Financial News',
        'Sarah Johnson',
        (CURRENT_DATE - INTERVAL '1 day') + TIME '09:00:00',
        'Business',
        3,
        true,
        true
    ),
    (
        yesterday_digest_id,
        'New Archaeological Discovery in Egypt',
        'Archaeologists have uncovered a previously unknown tomb in the Valley of the Kings...',
        'Ancient tomb discovered in Egypt provides new insights into pharaonic burial practices.',
        'History Today',
        'Dr. Ahmed Hassan',
        (CURRENT_DATE - INTERVAL '1 day') + TIME '14:30:00',
        'Science',
        4,
        true,
        true
    );

    -- Weekend articles
    INSERT INTO digest_articles (
        digest_id, title, content, summary, source, author, published_at,
        category, read_time_minutes, is_dummy, is_test_data
    ) VALUES
    (
        weekend_digest_id,
        'International Sports Championship Results',
        'The weekend saw exciting developments in international sports competitions...',
        'Weekend sports highlights from major international championships.',
        'Sports Weekly',
        'Mike Thompson',
        (CURRENT_DATE - INTERVAL '2 days') + TIME '16:00:00',
        'Sports',
        2,
        true,
        true
    );

END $$;