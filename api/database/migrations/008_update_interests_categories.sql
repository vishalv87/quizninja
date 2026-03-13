-- Update interests table to include all categories used by frontend
-- This migration adds missing categories and updates existing ones to match frontend expectations

-- Add missing categories that are used in frontend fallback but not in DB
INSERT INTO interests (id, name, description, icon_name, color_hex) VALUES
-- Categories from frontend fallback that were missing
('politics', 'Politics & Government', 'Elections, policies, international relations', 'flag', '#FF5252'),
('business', 'Business & Economy', 'Markets, companies, economic trends', 'trending_up', '#4CAF50'),
('entertainment', 'Entertainment', 'Movies, celebrities, culture, arts', 'star', '#E91E63'),
('world', 'World News', 'International events, global affairs', 'public', '#3F51B5'),
('environment', 'Environment', 'Climate change, sustainability, nature', 'eco', '#009688')
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon_name = EXCLUDED.icon_name,
    color_hex = EXCLUDED.color_hex;

-- Update existing categories to use more consistent naming and better descriptions
UPDATE interests SET
    name = 'Science & Health',
    description = 'Research, discoveries, medical breakthroughs',
    icon_name = 'book',
    color_hex = '#2196F3'
WHERE id = 'science';

UPDATE interests SET
    name = 'Technology',
    description = 'AI, startups, gadgets, innovation',
    icon_name = 'flash_on',
    color_hex = '#9C27B0'
WHERE id = 'technology';

UPDATE interests SET
    name = 'Sports',
    description = 'Cricket, football, Olympics, tournaments',
    icon_name = 'emoji_events',
    color_hex = '#FF9800'
WHERE id = 'sports';

-- Update any quizzes that were using movies_tv category first
UPDATE quizzes SET category_id = 'entertainment' WHERE category_id = 'movies_tv';

-- Delete movies_tv if it exists, then insert entertainment
DELETE FROM interests WHERE id = 'movies_tv';

-- Insert entertainment category
INSERT INTO interests (id, name, description, icon_name, color_hex) VALUES
('entertainment', 'Entertainment', 'Movies, celebrities, culture, arts', 'star', '#E91E63')
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    icon_name = EXCLUDED.icon_name,
    color_hex = EXCLUDED.color_hex;