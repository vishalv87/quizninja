-- Migration: Create digest tables
-- Date: 2025-01-15
-- Description: Creates tables for daily digest functionality including digests and digest_articles

-- Create extension for UUID generation if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Digests table - stores daily digest information
CREATE TABLE IF NOT EXISTS digests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date DATE NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    summary TEXT,
    article_count INTEGER DEFAULT 0,
    is_dummy BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Articles table - stores individual articles in each digest
CREATE TABLE IF NOT EXISTS digest_articles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    digest_id UUID NOT NULL REFERENCES digests(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    summary TEXT NOT NULL,
    source VARCHAR(255),
    author VARCHAR(255),
    published_at TIMESTAMP WITH TIME ZONE,
    category VARCHAR(100) NOT NULL,
    image_url VARCHAR(500),
    external_url VARCHAR(500),
    read_time_minutes INTEGER,
    is_breaking BOOLEAN DEFAULT false,
    is_hot BOOLEAN DEFAULT false,
    is_dummy BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_digests_date ON digests(date DESC);
CREATE INDEX IF NOT EXISTS idx_digests_created_at ON digests(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_digests_is_dummy ON digests(is_dummy);

CREATE INDEX IF NOT EXISTS idx_digest_articles_digest_id ON digest_articles(digest_id);
CREATE INDEX IF NOT EXISTS idx_digest_articles_category ON digest_articles(category);
CREATE INDEX IF NOT EXISTS idx_digest_articles_published_at ON digest_articles(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_digest_articles_is_breaking ON digest_articles(is_breaking);
CREATE INDEX IF NOT EXISTS idx_digest_articles_is_hot ON digest_articles(is_hot);
CREATE INDEX IF NOT EXISTS idx_digest_articles_is_dummy ON digest_articles(is_dummy);
CREATE INDEX IF NOT EXISTS idx_digest_articles_created_at ON digest_articles(created_at DESC);

-- Create trigger to automatically update article_count in digests table
CREATE OR REPLACE FUNCTION update_digest_article_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE digests
        SET article_count = (
            SELECT COUNT(*)
            FROM digest_articles
            WHERE digest_id = NEW.digest_id
        ),
        updated_at = NOW()
        WHERE id = NEW.digest_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE digests
        SET article_count = (
            SELECT COUNT(*)
            FROM digest_articles
            WHERE digest_id = OLD.digest_id
        ),
        updated_at = NOW()
        WHERE id = OLD.digest_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers
DROP TRIGGER IF EXISTS trigger_update_digest_article_count_insert ON digest_articles;
CREATE TRIGGER trigger_update_digest_article_count_insert
    AFTER INSERT ON digest_articles
    FOR EACH ROW
    EXECUTE FUNCTION update_digest_article_count();

DROP TRIGGER IF EXISTS trigger_update_digest_article_count_delete ON digest_articles;
CREATE TRIGGER trigger_update_digest_article_count_delete
    AFTER DELETE ON digest_articles
    FOR EACH ROW
    EXECUTE FUNCTION update_digest_article_count();

-- Create trigger to automatically update updated_at column
CREATE OR REPLACE FUNCTION update_digest_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_digests_updated_at ON digests;
CREATE TRIGGER trigger_update_digests_updated_at
    BEFORE UPDATE ON digests
    FOR EACH ROW
    EXECUTE FUNCTION update_digest_updated_at_column();

-- Views for easier querying

-- View to get digest with article statistics
CREATE OR REPLACE VIEW digest_with_stats AS
SELECT
    d.*,
    COUNT(da.id) as total_articles,
    COUNT(CASE WHEN da.is_breaking = true THEN 1 END) as breaking_articles,
    COUNT(CASE WHEN da.is_hot = true THEN 1 END) as hot_articles,
    STRING_AGG(DISTINCT da.category, ', ' ORDER BY da.category) as categories,
    MAX(da.published_at) as latest_article_published_at
FROM digests d
LEFT JOIN digest_articles da ON d.id = da.digest_id
GROUP BY d.id, d.date, d.title, d.summary, d.article_count, d.is_dummy, d.created_at, d.updated_at;

-- View to get today's digest with articles
CREATE OR REPLACE VIEW todays_digest AS
SELECT
    d.*,
    da.id as article_id,
    da.title as article_title,
    da.content as article_content,
    da.summary as article_summary,
    da.source as article_source,
    da.author as article_author,
    da.published_at as article_published_at,
    da.category as article_category,
    da.image_url as article_image_url,
    da.external_url as article_external_url,
    da.read_time_minutes as article_read_time_minutes,
    da.is_breaking as article_is_breaking,
    da.is_hot as article_is_hot,
    da.is_dummy as article_is_dummy,
    da.created_at as article_created_at
FROM digests d
LEFT JOIN digest_articles da ON d.id = da.digest_id
WHERE d.date = CURRENT_DATE
ORDER BY
    CASE WHEN da.is_breaking = true THEN 0 ELSE 1 END,
    CASE WHEN da.is_hot = true THEN 0 ELSE 1 END,
    da.published_at DESC NULLS LAST;

-- Function to get or create today's digest
CREATE OR REPLACE FUNCTION get_or_create_todays_digest()
RETURNS UUID AS $$
DECLARE
    digest_id UUID;
BEGIN
    -- Try to find today's digest
    SELECT id INTO digest_id
    FROM digests
    WHERE date = CURRENT_DATE;

    -- If not found, create it
    IF digest_id IS NULL THEN
        INSERT INTO digests (date, title, summary, is_dummy)
        VALUES (
            CURRENT_DATE,
            'Daily News Digest',
            'Stay informed with the most important news from around the world.',
            false
        )
        RETURNING id INTO digest_id;
    END IF;

    RETURN digest_id;
END;
$$ LANGUAGE plpgsql;

-- Function to get digest categories
CREATE OR REPLACE FUNCTION get_digest_categories()
RETURNS TABLE(category VARCHAR) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT da.category
    FROM digest_articles da
    WHERE da.is_dummy = false
    ORDER BY da.category;
END;
$$ LANGUAGE plpgsql;