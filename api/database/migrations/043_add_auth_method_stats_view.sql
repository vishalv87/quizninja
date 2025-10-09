-- Migration: Add auth_method_stats view
-- Date: 2025-01-15
-- Description: Creates a view to track authentication method statistics and migration status

-- Create view to get authentication method statistics
CREATE OR REPLACE VIEW auth_method_stats AS
SELECT
    auth_method,
    COUNT(*) AS user_count,
    COUNT(CASE WHEN migrated_at IS NOT NULL THEN 1 END) AS migrated_count,
    MIN(created_at) AS first_user_created,
    MAX(last_active) AS last_activity
FROM users
GROUP BY auth_method;

-- Add helpful comment
COMMENT ON VIEW auth_method_stats IS 'Aggregates user statistics by authentication method (jwt vs supabase) and tracks migration status';
