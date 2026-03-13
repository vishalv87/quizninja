-- Add title column to discussions table
-- This addresses the mismatch between UI expectations and backend implementation

-- Add title column (nullable initially for existing data)
ALTER TABLE discussions
ADD COLUMN title VARCHAR(200);

-- Backfill existing records with truncated content as title
UPDATE discussions
SET title = CASE
    WHEN LENGTH(content) > 50 THEN SUBSTRING(content FROM 1 FOR 47) || '...'
    ELSE content
END
WHERE title IS NULL;

-- Make title NOT NULL after backfill
ALTER TABLE discussions
ALTER COLUMN title SET NOT NULL;

-- Add index on title for search performance
CREATE INDEX idx_discussions_title ON discussions(title);
