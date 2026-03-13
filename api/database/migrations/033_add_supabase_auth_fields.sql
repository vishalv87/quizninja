-- Migration: Add Supabase authentication support fields to users table
-- This migration adds fields to support dual authentication (Supabase + JWT fallback)
-- and tracks user migration between authentication systems

-- Add Supabase authentication fields to users table
ALTER TABLE users
ADD COLUMN auth_method VARCHAR(20) DEFAULT 'jwt' NOT NULL CHECK (auth_method IN ('jwt', 'supabase')),
ADD COLUMN supabase_id UUID NULL,
ADD COLUMN last_auth_method VARCHAR(20) DEFAULT 'jwt' NOT NULL CHECK (last_auth_method IN ('jwt', 'supabase')),
ADD COLUMN migrated_at TIMESTAMP NULL;

-- Create indexes for efficient authentication queries
CREATE INDEX idx_users_auth_method ON users(auth_method);
CREATE INDEX idx_users_supabase_id ON users(supabase_id);
CREATE INDEX idx_users_last_auth_method ON users(last_auth_method);
CREATE INDEX idx_users_migrated_at ON users(migrated_at);

-- Create unique index for supabase_id to prevent duplicates (only when not null)
CREATE UNIQUE INDEX idx_users_supabase_id_unique
ON users(supabase_id)
WHERE supabase_id IS NOT NULL;

-- Add comments for documentation
COMMENT ON COLUMN users.auth_method IS 'Primary authentication method: jwt (custom) or supabase';
COMMENT ON COLUMN users.supabase_id IS 'Supabase user ID when user is linked to Supabase Auth';
COMMENT ON COLUMN users.last_auth_method IS 'Last successful authentication method used';
COMMENT ON COLUMN users.migrated_at IS 'Timestamp when user was migrated between auth systems';

-- Initialize existing users with JWT authentication method
UPDATE users
SET
    auth_method = 'jwt',
    last_auth_method = 'jwt'
WHERE auth_method IS NULL OR last_auth_method IS NULL;

-- Create function to migrate user from JWT to Supabase
CREATE OR REPLACE FUNCTION migrate_user_to_supabase(
    p_user_id UUID,
    p_supabase_id UUID
)
RETURNS BOOLEAN AS $$
DECLARE
    affected_rows INTEGER;
BEGIN
    -- Check if supabase_id is already in use
    IF EXISTS (SELECT 1 FROM users WHERE supabase_id = p_supabase_id) THEN
        RAISE EXCEPTION 'Supabase ID already linked to another user';
    END IF;

    UPDATE users
    SET
        auth_method = 'supabase',
        supabase_id = p_supabase_id,
        last_auth_method = 'supabase',
        migrated_at = CURRENT_TIMESTAMP
    WHERE id = p_user_id;

    GET DIAGNOSTICS affected_rows = ROW_COUNT;

    IF affected_rows > 0 THEN
        RETURN TRUE;
    ELSE
        RETURN FALSE;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create function to migrate user from Supabase back to JWT
CREATE OR REPLACE FUNCTION migrate_user_to_jwt(
    p_user_id UUID
)
RETURNS BOOLEAN AS $$
DECLARE
    affected_rows INTEGER;
BEGIN
    UPDATE users
    SET
        auth_method = 'jwt',
        supabase_id = NULL,
        last_auth_method = 'jwt',
        migrated_at = CURRENT_TIMESTAMP
    WHERE id = p_user_id;

    GET DIAGNOSTICS affected_rows = ROW_COUNT;

    IF affected_rows > 0 THEN
        RETURN TRUE;
    ELSE
        RETURN FALSE;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create function to update last auth method
CREATE OR REPLACE FUNCTION update_user_last_auth_method(
    p_user_id UUID,
    p_auth_method VARCHAR(20)
)
RETURNS BOOLEAN AS $$
DECLARE
    affected_rows INTEGER;
BEGIN
    -- Validate auth method
    IF p_auth_method NOT IN ('jwt', 'supabase') THEN
        RAISE EXCEPTION 'Invalid auth method: %. Must be jwt or supabase', p_auth_method;
    END IF;

    UPDATE users
    SET
        last_auth_method = p_auth_method,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_user_id;

    GET DIAGNOSTICS affected_rows = ROW_COUNT;

    IF affected_rows > 0 THEN
        RETURN TRUE;
    ELSE
        RETURN FALSE;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create function to find user by Supabase ID
CREATE OR REPLACE FUNCTION get_user_by_supabase_id(
    p_supabase_id UUID
)
RETURNS TABLE (
    id UUID,
    email VARCHAR,
    name VARCHAR,
    auth_method VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        u.id,
        u.email,
        u.name,
        u.auth_method,
        u.created_at,
        u.updated_at
    FROM users u
    WHERE u.supabase_id = p_supabase_id;
END;
$$ LANGUAGE plpgsql;

-- Create view for authentication statistics
CREATE OR REPLACE VIEW auth_method_stats AS
SELECT
    auth_method,
    COUNT(*) as user_count,
    COUNT(CASE WHEN migrated_at IS NOT NULL THEN 1 END) as migrated_count,
    MIN(created_at) as first_user_created,
    MAX(last_active) as last_activity
FROM users
GROUP BY auth_method;

-- Add constraint to ensure migrated_at is set when auth_method changes from default
-- (This will be enforced by application logic rather than database constraint for flexibility)

-- Create trigger to automatically update last_active when authentication occurs
CREATE OR REPLACE FUNCTION update_last_active_on_auth()
RETURNS TRIGGER AS $$
BEGIN
    -- Only update if last_auth_method actually changed
    IF OLD.last_auth_method IS DISTINCT FROM NEW.last_auth_method THEN
        NEW.last_active = CURRENT_TIMESTAMP;
        NEW.updated_at = CURRENT_TIMESTAMP;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_last_active_on_auth
    BEFORE UPDATE OF last_auth_method ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_last_active_on_auth();