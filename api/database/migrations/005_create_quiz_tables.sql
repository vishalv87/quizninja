-- Create quiz system tables
-- This migration adds the core quiz functionality tables

-- Quiz content table
CREATE TABLE quizzes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    category_id VARCHAR(50) REFERENCES interests(id),
    difficulty VARCHAR(50) REFERENCES difficulty_levels(id),
    article_summary TEXT,
    total_questions INTEGER NOT NULL,
    time_limit_minutes INTEGER NOT NULL,
    points INTEGER NOT NULL,
    is_featured BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quiz questions table
CREATE TABLE questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    quiz_id UUID REFERENCES quizzes(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL, -- multipleChoice, trueFalse, fillInBlank
    options TEXT[], -- Array of options for multiple choice
    correct_answer TEXT NOT NULL,
    explanation TEXT,
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Quiz statistics table
CREATE TABLE quiz_statistics (
    quiz_id UUID PRIMARY KEY REFERENCES quizzes(id) ON DELETE CASCADE,
    total_attempts INTEGER DEFAULT 0,
    total_completions INTEGER DEFAULT 0,
    average_score DECIMAL(5,2) DEFAULT 0.0,
    average_time_seconds INTEGER DEFAULT 0,
    difficulty_rating DECIMAL(3,2) DEFAULT 0.0, -- Computed from user performance
    popularity_score INTEGER DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User quiz favorites/bookmarks table
CREATE TABLE user_quiz_favorites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    quiz_id UUID REFERENCES quizzes(id) ON DELETE CASCADE,
    favorited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, quiz_id)
);

-- Quiz ratings/reviews table
CREATE TABLE quiz_ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    quiz_id UUID REFERENCES quizzes(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating BETWEEN 1 AND 5),
    review TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, quiz_id)
);

-- Add indexes for performance
CREATE INDEX idx_quizzes_category_id ON quizzes(category_id);
CREATE INDEX idx_quizzes_difficulty ON quizzes(difficulty);
CREATE INDEX idx_quizzes_is_featured ON quizzes(is_featured);
CREATE INDEX idx_quizzes_is_active ON quizzes(is_active);
CREATE INDEX idx_quizzes_created_at ON quizzes(created_at);

CREATE INDEX idx_questions_quiz_id ON questions(quiz_id);
CREATE INDEX idx_questions_order_index ON questions(order_index);
CREATE INDEX idx_questions_question_type ON questions(question_type);

CREATE INDEX idx_quiz_statistics_total_attempts ON quiz_statistics(total_attempts);
CREATE INDEX idx_quiz_statistics_average_score ON quiz_statistics(average_score);
CREATE INDEX idx_quiz_statistics_popularity_score ON quiz_statistics(popularity_score);

CREATE INDEX idx_user_quiz_favorites_user_id ON user_quiz_favorites(user_id);
CREATE INDEX idx_user_quiz_favorites_quiz_id ON user_quiz_favorites(quiz_id);
CREATE INDEX idx_user_quiz_favorites_favorited_at ON user_quiz_favorites(favorited_at);

CREATE INDEX idx_quiz_ratings_user_id ON quiz_ratings(user_id);
CREATE INDEX idx_quiz_ratings_quiz_id ON quiz_ratings(quiz_id);
CREATE INDEX idx_quiz_ratings_rating ON quiz_ratings(rating);
CREATE INDEX idx_quiz_ratings_created_at ON quiz_ratings(created_at);

-- Add updated_at triggers for tables that need them
CREATE TRIGGER update_quizzes_updated_at
    BEFORE UPDATE ON quizzes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_quiz_statistics_last_updated
BEFORE UPDATE ON quiz_statistics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();