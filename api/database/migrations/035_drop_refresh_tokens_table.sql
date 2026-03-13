-- Remove JWT refresh tokens table and related indexes
-- This migration removes the refresh_tokens table that was used for JWT authentication
-- Since we've moved to Supabase-only authentication, this table is no longer needed

-- Drop indexes first (PostgreSQL will automatically drop them when the table is dropped, but being explicit)
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;

-- Drop the refresh_tokens table
-- This will also drop the foreign key constraint to users(id)
DROP TABLE IF EXISTS refresh_tokens;