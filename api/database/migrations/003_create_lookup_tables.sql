-- Create lookup tables for interests, difficulty levels, and notification frequencies

-- Quiz categories/interests
CREATE TABLE interests (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon_name VARCHAR(100),
    color_hex VARCHAR(7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Difficulty levels
CREATE TABLE difficulty_levels (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon_name VARCHAR(100),
    background_color_hex VARCHAR(7)
);

-- Notification frequency options
CREATE TABLE notification_frequencies (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon_name VARCHAR(100)
);

-- Add indexes for performance
CREATE INDEX idx_interests_name ON interests(name);
CREATE INDEX idx_difficulty_levels_name ON difficulty_levels(name);
CREATE INDEX idx_notification_frequencies_name ON notification_frequencies(name);

