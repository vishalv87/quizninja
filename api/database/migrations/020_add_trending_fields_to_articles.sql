-- Migration: Add trending fields to digest_articles table
-- Date: 2025-01-15
-- Description: Adds trending-specific fields to digest_articles for trending topics functionality

-- Add trending fields to digest_articles table
ALTER TABLE digest_articles
ADD COLUMN IF NOT EXISTS is_trending BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS trending_score DECIMAL(5,2) DEFAULT 0.0,
ADD COLUMN IF NOT EXISTS trending_rank INTEGER DEFAULT NULL;

-- Create indexes for trending queries
CREATE INDEX IF NOT EXISTS idx_digest_articles_is_trending ON digest_articles(is_trending);
CREATE INDEX IF NOT EXISTS idx_digest_articles_trending_score ON digest_articles(trending_score DESC);
CREATE INDEX IF NOT EXISTS idx_digest_articles_trending_rank ON digest_articles(trending_rank);

-- Create composite index for trending articles ordered by score
CREATE INDEX IF NOT EXISTS idx_digest_articles_trending_composite ON digest_articles(is_trending, trending_score DESC, created_at DESC) WHERE is_trending = TRUE;

-- Create function to calculate trending score based on article engagement
-- This is a simple example - in production, you'd want more sophisticated algorithms
CREATE OR REPLACE FUNCTION calculate_trending_score(
    article_id UUID,
    base_score DECIMAL DEFAULT 50.0
) RETURNS DECIMAL AS $$
DECLARE
    final_score DECIMAL := base_score;
    article_age_hours INTEGER;
    category_boost DECIMAL := 0.0;
    breaking_boost DECIMAL := 0.0;
    hot_boost DECIMAL := 0.0;
    recency_factor DECIMAL := 1.0;
BEGIN
    -- Get article age in hours
    SELECT EXTRACT(EPOCH FROM (NOW() - created_at)) / 3600
    INTO article_age_hours
    FROM digest_articles
    WHERE id = article_id;

    -- Apply recency factor (newer articles get higher scores)
    IF article_age_hours <= 2 THEN
        recency_factor := 2.0;
    ELSIF article_age_hours <= 6 THEN
        recency_factor := 1.5;
    ELSIF article_age_hours <= 12 THEN
        recency_factor := 1.2;
    ELSIF article_age_hours <= 24 THEN
        recency_factor := 1.0;
    ELSE
        recency_factor := 0.8;
    END IF;

    -- Apply category-specific boosts
    SELECT CASE
        WHEN category IN ('Technology', 'Business') THEN 10.0
        WHEN category IN ('Environment', 'Science') THEN 8.0
        WHEN category IN ('Health', 'Politics') THEN 6.0
        ELSE 5.0
    END INTO category_boost
    FROM digest_articles
    WHERE id = article_id;

    -- Apply breaking/hot article boosts
    SELECT
        CASE WHEN is_breaking THEN 20.0 ELSE 0.0 END,
        CASE WHEN is_hot THEN 15.0 ELSE 0.0 END
    INTO breaking_boost, hot_boost
    FROM digest_articles
    WHERE id = article_id;

    -- Calculate final score
    final_score := (base_score + category_boost + breaking_boost + hot_boost) * recency_factor;

    -- Cap the score at 100
    IF final_score > 100.0 THEN
        final_score := 100.0;
    END IF;

    RETURN final_score;
END;
$$ LANGUAGE plpgsql;

-- Create function to update trending rankings
CREATE OR REPLACE FUNCTION update_trending_rankings()
RETURNS VOID AS $$
DECLARE
    article_record RECORD;
    rank_counter INTEGER := 1;
BEGIN
    -- First, calculate trending scores for all articles
    FOR article_record IN
        SELECT id FROM digest_articles
        WHERE created_at >= NOW() - INTERVAL '7 days' -- Only consider articles from last 7 days
        ORDER BY created_at DESC
    LOOP
        UPDATE digest_articles
        SET trending_score = calculate_trending_score(article_record.id)
        WHERE id = article_record.id;
    END LOOP;

    -- Mark top articles as trending (top 20 by score)
    UPDATE digest_articles SET is_trending = FALSE, trending_rank = NULL;

    FOR article_record IN
        SELECT id FROM digest_articles
        WHERE trending_score > 60.0 -- Minimum score threshold
        AND created_at >= NOW() - INTERVAL '7 days'
        ORDER BY trending_score DESC
        LIMIT 20
    LOOP
        UPDATE digest_articles
        SET is_trending = TRUE, trending_rank = rank_counter
        WHERE id = article_record.id;

        rank_counter := rank_counter + 1;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Create view for trending articles
CREATE OR REPLACE VIEW trending_articles AS
SELECT
    da.*,
    d.date as digest_date,
    d.title as digest_title
FROM digest_articles da
JOIN digests d ON da.digest_id = d.id
WHERE da.is_trending = TRUE
ORDER BY da.trending_rank ASC, da.trending_score DESC, da.created_at DESC;

-- Create view for trending articles by category
CREATE OR REPLACE VIEW trending_articles_by_category AS
SELECT
    da.*,
    d.date as digest_date,
    d.title as digest_title,
    ROW_NUMBER() OVER (PARTITION BY da.category ORDER BY da.trending_score DESC) as category_rank
FROM digest_articles da
JOIN digests d ON da.digest_id = d.id
WHERE da.is_trending = TRUE
ORDER BY da.category, da.trending_score DESC;

-- Insert some sample trending data for testing
-- Update some existing articles to be trending (if any exist)
UPDATE digest_articles
SET is_trending = TRUE, trending_score = 95.5, trending_rank = 1
WHERE id IN (
    SELECT id FROM digest_articles
    WHERE is_dummy = TRUE
    AND (is_breaking = TRUE OR is_hot = TRUE)
    LIMIT 1
);

UPDATE digest_articles
SET is_trending = TRUE, trending_score = 88.2, trending_rank = 2
WHERE id IN (
    SELECT id FROM digest_articles
    WHERE is_dummy = TRUE
    AND category = 'Technology'
    AND id NOT IN (SELECT id FROM digest_articles WHERE trending_rank = 1)
    LIMIT 1
);

UPDATE digest_articles
SET is_trending = TRUE, trending_score = 82.7, trending_rank = 3
WHERE id IN (
    SELECT id FROM digest_articles
    WHERE is_dummy = TRUE
    AND category = 'Environment'
    AND id NOT IN (SELECT id FROM digest_articles WHERE trending_rank IN (1, 2))
    LIMIT 1
);

-- Create a scheduled job placeholder comment
-- In production, you would set up a CRON job or scheduled task to run:
-- SELECT update_trending_rankings();
-- This should be run every hour or few hours to keep trending data fresh

COMMENT ON FUNCTION update_trending_rankings() IS 'Should be called periodically (hourly) to refresh trending article rankings';
COMMENT ON FUNCTION calculate_trending_score(UUID, DECIMAL) IS 'Calculates trending score for an article based on various factors';