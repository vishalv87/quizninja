-- Migration: 039_rename_interests_to_categories.sql
-- Description: Rename interests table to categories and update related references
-- This migration standardizes terminology across the application

-- Rename the main interests table to categories
ALTER TABLE interests RENAME TO categories;

-- Rename the column in user_preferences table
ALTER TABLE user_preferences
  RENAME COLUMN selected_interests TO selected_categories;

-- Note: quizzes.category_id already uses correct naming
-- PostgreSQL automatically updates foreign key references when renaming the table

-- Update any indexes that reference the old table name
ALTER INDEX IF EXISTS interests_pkey RENAME TO categories_pkey;
ALTER INDEX IF EXISTS idx_interests_name RENAME TO idx_categories_name;

-- Note: Foreign key constraints will automatically work with the renamed table
-- PostgreSQL handles this transparently when renaming tables
