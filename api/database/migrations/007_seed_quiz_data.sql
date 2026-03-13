-- Seed quiz data to populate the application with sample quizzes
-- This prevents empty API responses and gives users content to explore

-- Sample Quiz 1: General Knowledge
INSERT INTO quizzes (id, title, description, category_id, difficulty, total_questions, time_limit_minutes, points, created_by, is_featured, is_active, is_public, tags, created_at, updated_at)
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
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Questions for General Knowledge quiz
INSERT INTO questions (id, quiz_id, question_text, question_type, options, correct_answer, explanation, order_index, created_at)
VALUES
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the capital of France?', 'multipleChoice', '{"London", "Berlin", "Paris", "Madrid"}', 'Paris', 'Paris is the capital and most populous city of France.', 1, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which planet is closest to the Sun?', 'multipleChoice', '{"Venus", "Mercury", "Earth", "Mars"}', 'Mercury', 'Mercury is the smallest planet in our solar system and closest to the Sun.', 2, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'The Great Wall of China was built primarily to keep out which people?', 'multipleChoice', '{"Mongols", "Japanese", "Russians", "British"}', 'Mongols', 'The Great Wall was built to protect Chinese states from invasions by nomadic groups, particularly the Mongols.', 3, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the chemical symbol for gold?', 'multipleChoice', '{"Go", "Gd", "Au", "Ag"}', 'Au', 'Au comes from the Latin word "aurum" meaning gold.', 4, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which ocean is the largest?', 'multipleChoice', '{"Atlantic", "Pacific", "Indian", "Arctic"}', 'Pacific', 'The Pacific Ocean is the largest and deepest ocean on Earth.', 5, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Who painted the Mona Lisa?', 'multipleChoice', '{"Vincent van Gogh", "Pablo Picasso", "Leonardo da Vinci", "Michelangelo"}', 'Leonardo da Vinci', 'Leonardo da Vinci painted the Mona Lisa between 1503 and 1519.', 6, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the largest mammal in the world?', 'multipleChoice', '{"African Elephant", "Blue Whale", "Giraffe", "Polar Bear"}', 'Blue Whale', 'Blue whales are the largest animals ever known to have lived on Earth.', 7, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'In which year did World War II end?', 'multipleChoice', '{"1944", "1945", "1946", "1947"}', '1945', 'World War II ended in 1945 with the surrender of Japan in September.', 8, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the smallest country in the world?', 'multipleChoice', '{"Monaco", "Vatican City", "San Marino", "Malta"}', 'Vatican City', 'Vatican City is the smallest sovereign nation in the world by both area and population.', 9, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which element has the chemical symbol "O"?', 'multipleChoice', '{"Osmium", "Oxygen", "Ozone", "Oxide"}', 'Oxygen', 'Oxygen is a chemical element with the symbol O and atomic number 8.', 10, CURRENT_TIMESTAMP);

-- Sample Quiz 2: Technology & Programming
INSERT INTO quizzes (id, title, description, category_id, difficulty, total_questions, time_limit_minutes, points, created_by, is_featured, is_active, is_public, tags, created_at, updated_at)
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
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Questions for Technology quiz
INSERT INTO questions (id, quiz_id, question_text, question_type, options, correct_answer, explanation, order_index, created_at)
VALUES
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What does "HTML" stand for?', 'multipleChoice', '{"Hypertext Markup Language", "High-Tech Modern Language", "Home Tool Markup Language", "Hyperlink and Text Markup Language"}', 'Hypertext Markup Language', 'HTML is the standard markup language for creating web pages.', 1, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which programming language is known as the "language of the web"?', 'multipleChoice', '{"Python", "JavaScript", "Java", "C++"}', 'JavaScript', 'JavaScript is primarily known as the scripting language for web pages.', 2, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What does "CPU" stand for?', 'multipleChoice', '{"Central Processing Unit", "Computer Processing Unit", "Central Program Unit", "Computer Program Unit"}', 'Central Processing Unit', 'The CPU is the primary component of a computer that performs most processing.', 3, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which company developed the Python programming language?', 'multipleChoice', '{"Google", "Microsoft", "Guido van Rossum", "Facebook"}', 'Guido van Rossum', 'Python was created by Guido van Rossum and first released in 1991.', 4, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What does "API" stand for?', 'multipleChoice', '{"Application Programming Interface", "Advanced Programming Interface", "Application Program Integration", "Advanced Program Integration"}', 'Application Programming Interface', 'An API defines how software components should interact.', 5, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which of these is NOT a database management system?', 'multipleChoice', '{"MySQL", "PostgreSQL", "MongoDB", "Photoshop"}', 'Photoshop', 'Photoshop is image editing software, not a database management system.', 6, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What does "Git" primarily help with?', 'multipleChoice', '{"Image editing", "Version control", "Database management", "Web hosting"}', 'Version control', 'Git is a distributed version control system for tracking changes in source code.', 7, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which protocol is used for secure web browsing?', 'multipleChoice', '{"HTTP", "HTTPS", "FTP", "SMTP"}', 'HTTPS', 'HTTPS provides encrypted communication between web browsers and servers.', 8, CURRENT_TIMESTAMP);

-- Sample Quiz 3: Science & Nature
INSERT INTO quizzes (id, title, description, category_id, difficulty, total_questions, time_limit_minutes, points, created_by, is_featured, is_active, is_public, tags, created_at, updated_at)
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
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

-- Questions for Science quiz
INSERT INTO questions (id, quiz_id, question_text, question_type, options, correct_answer, explanation, order_index, created_at)
VALUES
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the speed of light in a vacuum?', 'multipleChoice', '{"300,000 km/s", "299,792,458 m/s", "150,000 km/s", "186,000 miles/s"}', '299,792,458 m/s', 'The speed of light in a vacuum is exactly 299,792,458 meters per second.', 1, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which gas makes up about 78% of Earth''s atmosphere?', 'multipleChoice', '{"Oxygen", "Carbon Dioxide", "Nitrogen", "Argon"}', 'Nitrogen', 'Nitrogen makes up approximately 78% of Earth''s atmosphere by volume.', 2, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the powerhouse of the cell?', 'multipleChoice', '{"Nucleus", "Mitochondria", "Ribosome", "Chloroplast"}', 'Mitochondria', 'Mitochondria are often called the powerhouse of the cell because they generate ATP.', 3, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the hardest natural substance on Earth?', 'multipleChoice', '{"Gold", "Iron", "Diamond", "Granite"}', 'Diamond', 'Diamond is the hardest known natural material on the Mohs scale.', 4, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'How many chambers does a human heart have?', 'multipleChoice', '{"2", "3", "4", "5"}', '4', 'The human heart has four chambers: two atria and two ventricles.', 5, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the most abundant element in the universe?', 'multipleChoice', '{"Oxygen", "Carbon", "Hydrogen", "Helium"}', 'Hydrogen', 'Hydrogen is the most abundant element in the observable universe.', 6, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'At what temperature does water boil at sea level?', 'multipleChoice', '{"90°C", "95°C", "100°C", "105°C"}', '100°C', 'Water boils at 100°C (212°F) at sea level atmospheric pressure.', 7, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Which planet is known as the Red Planet?', 'multipleChoice', '{"Venus", "Mars", "Jupiter", "Saturn"}', 'Mars', 'Mars appears red due to iron oxide (rust) on its surface.', 8, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What type of animal is a Komodo dragon?', 'multipleChoice', '{"Snake", "Lizard", "Crocodile", "Turtle"}', 'Lizard', 'The Komodo dragon is the largest living species of lizard.', 9, CURRENT_TIMESTAMP),
(uuid_generate_v4(), 'b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'What is the chemical formula for water?', 'multipleChoice', '{"H2O", "CO2", "NaCl", "CH4"}', 'H2O', 'Water is composed of two hydrogen atoms and one oxygen atom (H2O).', 10, CURRENT_TIMESTAMP);

-- Create initial statistics for each quiz
INSERT INTO quiz_statistics (quiz_id, total_attempts, total_completions, average_score, average_time_seconds, difficulty_rating, popularity_score, last_updated)
VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 0, 0, 0.0, 0, 0.0, 0, CURRENT_TIMESTAMP),
('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 0, 0, 0.0, 0, 0.0, 0, CURRENT_TIMESTAMP),
('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 0, 0, 0.0, 0, 0.0, 0, CURRENT_TIMESTAMP);