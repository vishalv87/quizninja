-- Migration: Add image_url column to questions table
-- Description: Adds support for displaying images in quiz questions (diagrams, charts, photos, etc.)
-- Author: Claude Code
-- Date: 2025-11-16

-- Add image_url column to questions table
ALTER TABLE questions
ADD COLUMN IF NOT EXISTS image_url VARCHAR(500);

-- Add comment to document the column purpose
COMMENT ON COLUMN questions.image_url IS 'Optional URL to an image for the question (diagrams, charts, photos, etc.)';
