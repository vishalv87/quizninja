-- Migration: Rename 'challenged' fields to 'challengee' for consistency
-- Created: 2025-10-17
-- Description: Renames challenged_id and challenged_score to challengee_id and challengee_score across the challenges table and all related triggers/functions

-- Rename columns in challenges table
ALTER TABLE challenges RENAME COLUMN challenged_id TO challengee_id;
ALTER TABLE challenges RENAME COLUMN challenged_score TO challengee_score;

-- Drop and recreate the foreign key constraint with new column name
ALTER TABLE challenges DROP CONSTRAINT IF EXISTS challenges_challenged_id_fkey;
ALTER TABLE challenges ADD CONSTRAINT challenges_challengee_id_fkey
    FOREIGN KEY (challengee_id) REFERENCES users(id) ON DELETE CASCADE;

-- Update the challenge notification trigger to use new column names
CREATE OR REPLACE FUNCTION create_challenge_notification()
RETURNS TRIGGER AS $$
BEGIN
    -- Create notification for new challenge
    IF TG_OP = 'INSERT' THEN
        INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
        SELECT
            NEW.challengee_id,
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
                NEW.challengee_id,
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'accepter_name', u.name,
                    'quiz_id', NEW.quiz_id
                ),
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challengee_id;
        END IF;

        -- Challenge declined
        IF NEW.status = 'declined' AND OLD.status = 'pending' THEN
            INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
            SELECT
                NEW.challenger_id,
                'challenge_declined',
                'Challenge Declined',
                u.name || ' declined your challenge.',
                NEW.challengee_id,
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'decliner_name', u.name,
                    'quiz_id', NEW.quiz_id
                ),
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challengee_id;
        END IF;

        -- Challenger completed (notify challengee)
        IF NEW.status = 'challenger_completed' AND OLD.status = 'accepted' THEN
            INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
            SELECT
                NEW.challengee_id,
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
                NEW.challengee_id,
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'opponent_name', u.name,
                    'quiz_id', NEW.quiz_id
                ),
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challengee_id;
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
                    WHEN NEW.challenger_score > NEW.challengee_score THEN 'You won the challenge against ' || u.name || '!'
                    WHEN NEW.challenger_score < NEW.challengee_score THEN 'You lost the challenge against ' || u.name || '.'
                    ELSE 'Your challenge with ' || u.name || ' ended in a tie!'
                END,
                NEW.challengee_id,
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'opponent_name', u.name,
                    'quiz_id', NEW.quiz_id,
                    'your_score', NEW.challenger_score,
                    'opponent_score', NEW.challengee_score,
                    'result', CASE
                        WHEN NEW.challenger_score > NEW.challengee_score THEN 'won'
                        WHEN NEW.challenger_score < NEW.challengee_score THEN 'lost'
                        ELSE 'tie'
                    END
                ),
                CURRENT_TIMESTAMP
            FROM users u WHERE u.id = NEW.challengee_id;

            -- Create completion notification for challengee
            INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
            SELECT
                NEW.challengee_id,
                'challenge_completed',
                'Challenge Completed',
                CASE
                    WHEN NEW.challengee_score > NEW.challenger_score THEN 'You won the challenge against ' || u.name || '!'
                    WHEN NEW.challengee_score < NEW.challenger_score THEN 'You lost the challenge against ' || u.name || '.'
                    ELSE 'Your challenge with ' || u.name || ' ended in a tie!'
                END,
                NEW.challenger_id,
                NEW.id,
                'challenge',
                jsonb_build_object(
                    'challenge_id', NEW.id,
                    'opponent_name', u.name,
                    'quiz_id', NEW.quiz_id,
                    'your_score', NEW.challengee_score,
                    'opponent_score', NEW.challenger_score,
                    'result', CASE
                        WHEN NEW.challengee_score > NEW.challenger_score THEN 'won'
                        WHEN NEW.challengee_score < NEW.challenger_score THEN 'lost'
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
