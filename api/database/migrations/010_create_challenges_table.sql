-- Migration: Create challenges table
-- Created: 2024-09-21

-- Create challenges table
CREATE TABLE challenges (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    challenger_id UUID NOT NULL,
    challenged_id UUID NOT NULL,
    quiz_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    challenger_score DECIMAL(5,2),
    challenged_score DECIMAL(5,2),
    message TEXT,
    expires_at TIMESTAMP,
    is_group_challenge BOOLEAN DEFAULT false,
    participant_ids UUID[],
    participant_scores JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (challenger_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (challenged_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE,
    CHECK (status IN ('pending', 'accepted', 'completed', 'expired', 'declined')),
    CHECK (challenger_id != challenged_id)
);

-- Create indexes for better performance
CREATE INDEX idx_challenges_challenger_id ON challenges(challenger_id);
CREATE INDEX idx_challenges_challenged_id ON challenges(challenged_id);
CREATE INDEX idx_challenges_quiz_id ON challenges(quiz_id);
CREATE INDEX idx_challenges_status ON challenges(status);
CREATE INDEX idx_challenges_created_at ON challenges(created_at);
CREATE INDEX idx_challenges_expires_at ON challenges(expires_at);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_challenges_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_challenges_updated_at
    BEFORE UPDATE ON challenges
    FOR EACH ROW
    EXECUTE FUNCTION update_challenges_updated_at();

-- Create trigger to automatically expire challenges
CREATE OR REPLACE FUNCTION expire_challenges()
RETURNS TRIGGER AS $$
BEGIN
    -- Auto-expire challenges that are past their expiry date
    UPDATE challenges
    SET status = 'expired', updated_at = CURRENT_TIMESTAMP
    WHERE expires_at < CURRENT_TIMESTAMP
    AND status IN ('pending', 'accepted');

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to create challenge notifications
CREATE OR REPLACE FUNCTION create_challenge_notification()
RETURNS TRIGGER AS $$
BEGIN
    -- Create notification for new challenge
    IF TG_OP = 'INSERT' THEN
        INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, created_at)
        SELECT
            NEW.challenged_id,
            'challenge_received',
            'New Challenge Received',
            u.name || ' challenged you to a quiz!',
            NEW.challenger_id,
            CURRENT_TIMESTAMP
        FROM users u WHERE u.id = NEW.challenger_id;

        RETURN NEW;
    END IF;

    -- Create notification for challenge status changes
    IF TG_OP = 'UPDATE' THEN
        -- Challenge accepted
        IF NEW.status = 'accepted' AND OLD.status = 'pending' THEN
            INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, created_at)
            SELECT
                NEW.challenger_id,
                'challenge_accepted',
                'Challenge Accepted',
                u.name || ' accepted your challenge!',
                NEW.challenged_id,
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challenged_id;
        END IF;

        -- Challenge declined
        IF NEW.status = 'declined' AND OLD.status = 'pending' THEN
            INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, created_at)
            SELECT
                NEW.challenger_id,
                'challenge_declined',
                'Challenge Declined',
                u.name || ' declined your challenge.',
                NEW.challenged_id,
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challenged_id;
        END IF;

        -- Challenge completed
        IF NEW.status = 'completed' AND OLD.status = 'accepted' THEN
            -- Determine winner and create notifications for both users
            DECLARE
                winner_id UUID;
                winner_name TEXT;
                loser_id UUID;
            BEGIN
                -- Determine winner
                IF NEW.challenger_score > NEW.challenged_score THEN
                    winner_id := NEW.challenger_id;
                    loser_id := NEW.challenged_id;
                ELSIF NEW.challenged_score > NEW.challenger_score THEN
                    winner_id := NEW.challenged_id;
                    loser_id := NEW.challenger_id;
                END IF;

                -- Create completion notification for challenger
                INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, created_at)
                SELECT
                    NEW.challenger_id,
                    'challenge_completed',
                    'Challenge Completed',
                    CASE
                        WHEN NEW.challenger_score > NEW.challenged_score THEN 'You won the challenge against ' || u.name || '!'
                        WHEN NEW.challenger_score < NEW.challenged_score THEN 'You lost the challenge against ' || u.name || '.'
                        ELSE 'Your challenge with ' || u.name || ' ended in a tie!'
                    END,
                    NEW.challenged_id,
                    CURRENT_TIMESTAMP
                FROM users u WHERE u.id = NEW.challenged_id;

                -- Create completion notification for challenged user
                INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, created_at)
                SELECT
                    NEW.challenged_id,
                    'challenge_completed',
                    'Challenge Completed',
                    CASE
                        WHEN NEW.challenged_score > NEW.challenger_score THEN 'You won the challenge against ' || u.name || '!'
                        WHEN NEW.challenged_score < NEW.challenger_score THEN 'You lost the challenge against ' || u.name || '.'
                        ELSE 'Your challenge with ' || u.name || ' ended in a tie!'
                    END,
                    NEW.challenger_id,
                    CURRENT_TIMESTAMP
                FROM users u WHERE u.id = NEW.challenger_id;
            END;
        END IF;

        RETURN NEW;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_challenge_notification
    AFTER INSERT OR UPDATE ON challenges
    FOR EACH ROW
    EXECUTE FUNCTION create_challenge_notification();

-- Update friend_notifications table to support challenge notification types
DO $$
BEGIN
    -- Check if the constraint exists and drop it
    IF EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'friend_notifications_type_check'
    ) THEN
        ALTER TABLE friend_notifications DROP CONSTRAINT friend_notifications_type_check;
    END IF;

    -- Add the new constraint with challenge notification types
    ALTER TABLE friend_notifications ADD CONSTRAINT friend_notifications_type_check
    CHECK (type IN ('friend_request', 'friend_accepted', 'friend_rejected',
                   'challenge_received', 'challenge_accepted', 'challenge_declined', 'challenge_completed'));
END $$;