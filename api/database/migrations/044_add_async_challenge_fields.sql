-- Migration: Add asynchronous challenge tracking fields
-- Created: 2024-10-17
-- Description: Adds fields to support asynchronous challenge completion where users can complete quizzes independently

-- Add attempt tracking columns to challenges table
ALTER TABLE challenges
ADD COLUMN IF NOT EXISTS challenger_attempt_id UUID,
ADD COLUMN IF NOT EXISTS challengee_attempt_id UUID,
ADD COLUMN IF NOT EXISTS challenger_completed_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS challengee_completed_at TIMESTAMP;

-- Add foreign key constraints for attempt IDs
ALTER TABLE challenges
ADD CONSTRAINT fk_challenges_challenger_attempt
    FOREIGN KEY (challenger_attempt_id) REFERENCES quiz_attempts(id) ON DELETE SET NULL;

ALTER TABLE challenges
ADD CONSTRAINT fk_challenges_challengee_attempt
    FOREIGN KEY (challengee_attempt_id) REFERENCES quiz_attempts(id) ON DELETE SET NULL;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_challenges_challenger_attempt_id ON challenges(challenger_attempt_id);
CREATE INDEX IF NOT EXISTS idx_challenges_challengee_attempt_id ON challenges(challengee_attempt_id);

-- Update status constraint to include new intermediate states
ALTER TABLE challenges DROP CONSTRAINT IF EXISTS challenges_status_check;
ALTER TABLE challenges ADD CONSTRAINT challenges_status_check
CHECK (status IN ('pending', 'accepted', 'challenger_completed', 'challengee_completed', 'completed', 'expired', 'declined'));

-- Add challenge tracking columns to quiz_attempts table
ALTER TABLE quiz_attempts
ADD COLUMN IF NOT EXISTS challenge_id UUID,
ADD COLUMN IF NOT EXISTS is_challenge_attempt BOOLEAN DEFAULT false;

-- Add foreign key constraint for challenge_id
ALTER TABLE quiz_attempts
ADD CONSTRAINT fk_quiz_attempts_challenge
    FOREIGN KEY (challenge_id) REFERENCES challenges(id) ON DELETE SET NULL;

-- Create index for challenge_id in quiz_attempts
CREATE INDEX IF NOT EXISTS idx_quiz_attempts_challenge_id ON quiz_attempts(challenge_id);

-- Update the challenge notification trigger to handle new status transitions
CREATE OR REPLACE FUNCTION create_challenge_notification()
RETURNS TRIGGER AS $$
BEGIN
    -- Create notification for new challenge
    IF TG_OP = 'INSERT' THEN
        INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
        SELECT
            NEW.challenged_id,
            'challenge_received',
            'New Challenge Received',
            u.name || ' challenged you to a quiz!',
            NEW.challenger_id,
            NEW.id,
            'challenge',
            jsonb_build_object(
                'challenge_id', NEW.id,
                'challenger_name', u.name,
                'quiz_id', NEW.quiz_id,
                'message', NEW.message,
                'expires_at', NEW.expires_at
            ),
            CURRENT_TIMESTAMP
        FROM users u WHERE u.id = NEW.challenger_id;

        RETURN NEW;
    END IF;

    -- Create notification for challenge status changes
    IF TG_OP = 'UPDATE' THEN
        -- Challenge accepted
        IF NEW.status = 'accepted' AND OLD.status = 'pending' THEN
            INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
            SELECT
                NEW.challenger_id,
                'challenge_accepted',
                'Challenge Accepted',
                u.name || ' accepted your challenge!',
                NEW.challenged_id,
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'accepter_name', u.name,
                    'quiz_id', NEW.quiz_id
                ),
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challenged_id;
        END IF;

        -- Challenge declined
        IF NEW.status = 'declined' AND OLD.status = 'pending' THEN
            INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
            SELECT
                NEW.challenger_id,
                'challenge_declined',
                'Challenge Declined',
                u.name || ' declined your challenge.',
                NEW.challenged_id,
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'decliner_name', u.name,
                    'quiz_id', NEW.quiz_id
                ),
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challenged_id;
        END IF;

        -- Challenger completed (notify challengee)
        IF NEW.status = 'challenger_completed' AND OLD.status = 'accepted' THEN
            INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
            SELECT
                NEW.challenged_id,
                'challenge_progress',
                'Opponent Completed Quiz',
                u.name || ' completed the challenge quiz. Your turn!',
                NEW.challenger_id,
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'opponent_name', u.name,
                    'quiz_id', NEW.quiz_id
                ),
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challenger_id;
        END IF;

        -- Challengee completed (notify challenger)
        IF NEW.status = 'challengee_completed' AND OLD.status = 'accepted' THEN
            INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
            SELECT
                NEW.challenger_id,
                'challenge_progress',
                'Opponent Completed Quiz',
                u.name || ' completed the challenge quiz. Your turn!',
                NEW.challenged_id,
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'opponent_name', u.name,
                    'quiz_id', NEW.quiz_id
                ),
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challenged_id;
        END IF;

        -- Challenge completed (both finished) - determine winner and notify both
        IF NEW.status = 'completed' AND OLD.status IN ('challenger_completed', 'challengee_completed') THEN
            -- Create completion notification for challenger
            INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
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
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'opponent_name', u.name,
                    'quiz_id', NEW.quiz_id,
                    'your_score', NEW.challenger_score,
                    'opponent_score', NEW.challenged_score,
                    'result', CASE
                        WHEN NEW.challenger_score > NEW.challenged_score THEN 'won'
                        WHEN NEW.challenger_score < NEW.challenged_score THEN 'lost'
                        ELSE 'tie'
                    END
                ),
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challenged_id;

            -- Create completion notification for challengee
            INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
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
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'opponent_name', u.name,
                    'quiz_id', NEW.quiz_id,
                    'your_score', NEW.challenged_score,
                    'opponent_score', NEW.challenger_score,
                    'result', CASE
                        WHEN NEW.challenged_score > NEW.challenger_score THEN 'won'
                        WHEN NEW.challenged_score < NEW.challenger_score THEN 'lost'
                        ELSE 'tie'
                    END
                ),
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challenger_id;
        END IF;

        RETURN NEW;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Recreate the trigger
DROP TRIGGER IF EXISTS trigger_create_challenge_notification ON challenges;
CREATE TRIGGER trigger_create_challenge_notification
    AFTER INSERT OR UPDATE ON challenges
    FOR EACH ROW
    EXECUTE FUNCTION create_challenge_notification();

-- Update notifications table constraint to support new notification type: challenge_progress
DO $$
BEGIN
    -- Check if the constraint exists and drop it
    IF EXISTS (
        SELECT 1 FROM information_schema.check_constraints
        WHERE constraint_name = 'notifications_type_check'
    ) THEN
        ALTER TABLE notifications DROP CONSTRAINT notifications_type_check;
    END IF;

    -- Add the updated constraint with challenge_progress notification type
    ALTER TABLE notifications ADD CONSTRAINT notifications_type_check
    CHECK (type IN (
        'friend_request', 'friend_accepted', 'friend_rejected',
        'challenge_received', 'challenge_accepted', 'challenge_declined', 'challenge_progress', 'challenge_completed',
        'achievement_unlocked',
        'general', 'system_announcement'
    ));
END $$;
