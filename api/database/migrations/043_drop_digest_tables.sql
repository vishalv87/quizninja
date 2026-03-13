-- Migration: Drop digest tables and related objects
-- Date: 2025-01-15
-- Description: Removes all digest-related database objects including tables, views, functions, and triggers

-- Drop triggers first
DROP TRIGGER IF EXISTS trigger_update_digest_article_count_insert ON digest_articles;
DROP TRIGGER IF EXISTS trigger_update_digest_article_count_delete ON digest_articles;
DROP TRIGGER IF EXISTS trigger_update_digests_updated_at ON digests;

-- Drop views
DROP VIEW IF EXISTS todays_digest;
DROP VIEW IF EXISTS digest_with_stats;

-- Drop functions
DROP FUNCTION IF EXISTS update_digest_article_count();
DROP FUNCTION IF EXISTS update_digest_updated_at_column();
DROP FUNCTION IF EXISTS get_or_create_todays_digest();
DROP FUNCTION IF EXISTS get_digest_categories();
DROP FUNCTION IF EXISTS update_trending_rankings();

-- Drop tables (CASCADE will handle foreign key constraints)
DROP TABLE IF EXISTS digest_articles CASCADE;
DROP TABLE IF EXISTS digests CASCADE;
