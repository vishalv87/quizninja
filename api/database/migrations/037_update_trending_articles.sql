-- Migration: Update trending articles
-- Date: 2025-10-05
-- Description: Sets existing articles as trending with appropriate scores and rankings

-- Update specific articles to be trending with realistic scores
UPDATE digest_articles
SET
    is_trending = true,
    trending_score = 85.0,
    trending_rank = 1
WHERE title = 'AI Revolution Transforms Healthcare Diagnosis'
    AND is_dummy = true;

UPDATE digest_articles
SET
    is_trending = true,
    trending_score = 78.5,
    trending_rank = 2
WHERE title = 'Climate Summit 2025: Historic Agreement Reached'
    AND is_dummy = true;

UPDATE digest_articles
SET
    is_trending = true,
    trending_score = 72.0,
    trending_rank = 3
WHERE title = 'SpaceX Achieves New Milestone in Mars Mission'
    AND is_dummy = true;

UPDATE digest_articles
SET
    is_trending = true,
    trending_score = 68.5,
    trending_rank = 4
WHERE title = 'Tech Giants Report Record Quarterly Earnings'
    AND is_dummy = true;

UPDATE digest_articles
SET
    is_trending = true,
    trending_score = 65.0,
    trending_rank = 5
WHERE title = 'Breakthrough Medical Research Promises New Cancer Treatment'
    AND is_dummy = true;

-- Optional: Add a few more trending articles from historical data
UPDATE digest_articles
SET
    is_trending = true,
    trending_score = 60.0,
    trending_rank = 6
WHERE title = 'Global Markets Show Strong Recovery'
    AND is_dummy = true;