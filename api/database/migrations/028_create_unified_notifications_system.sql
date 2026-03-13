-- Migration: Create unified notifications system
-- This migration replaces the friend_notifications table with a unified notifications table
-- that supports all notification types: friend, challenge, achievement, general

-- Drop existing friend_notifications table and related triggers
DROP TRIGGER IF EXISTS trigger_create_friend_request_notification ON friend_requests;
DROP TRIGGER IF EXISTS trigger_create_friendship_on_accept ON friend_requests;
DROP TRIGGER IF EXISTS trigger_cleanup_friendship_on_status_change ON friend_requests;
DROP TRIGGER IF EXISTS trigger_create_challenge_notification ON challenges;
DROP FUNCTION IF EXISTS create_friend_request_notification();
DROP FUNCTION IF EXISTS create_friendship_on_accept();
DROP FUNCTION IF EXISTS cleanup_friendship_on_status_change();
DROP FUNCTION IF EXISTS create_challenge_notification();
DROP TABLE IF EXISTS friend_notifications;

-- Create unified notifications table
CREATE TABLE notifications (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT,
    data JSONB DEFAULT '{}',
    related_user_id UUID,
    related_entity_id UUID,
    related_entity_type VARCHAR(50),
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP,
    expires_at TIMESTAMP,
    is_test_data BOOLEAN DEFAULT false,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (related_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (type IN (
        'friend_request', 'friend_accepted', 'friend_rejected',
        'challenge_received', 'challenge_accepted', 'challenge_declined', 'challenge_completed',
        'achievement_unlocked',
        'general', 'system_announcement'
    ))
);

-- Create indexes for better performance
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
CREATE INDEX idx_notifications_related_user_id ON notifications(related_user_id);
CREATE INDEX idx_notifications_related_entity ON notifications(related_entity_id, related_entity_type);
CREATE INDEX idx_notifications_is_test_data ON notifications(is_test_data);

-- Create trigger to update updated_at timestamp (if we add this column later)
CREATE OR REPLACE FUNCTION update_notifications_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    -- This function is prepared for future use if we add updated_at column
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create enhanced friend request notification trigger
CREATE OR REPLACE FUNCTION create_friend_request_notification()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
    SELECT
        NEW.requested_id,
        'friend_request',
        'New Friend Request',
        u.name || ' sent you a friend request!',
        NEW.requester_id,
        NEW.id,
        'friend_request',
        jsonb_build_object(
            'friend_request_id', NEW.id,
            'requester_name', u.name,
            'message', NEW.message
        ),
        CURRENT_TIMESTAMP
    FROM users u WHERE u.id = NEW.requester_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create enhanced friendship acceptance/rejection notification trigger
CREATE OR REPLACE FUNCTION create_friendship_on_accept()
RETURNS TRIGGER AS $$
BEGIN
    -- Only create friendship if status changed to 'accepted'
    IF NEW.status = 'accepted' AND OLD.status != 'accepted' THEN
        -- Ensure user1_id < user2_id for consistent ordering
        INSERT INTO friendships (user1_id, user2_id, created_at)
        VALUES (
            LEAST(NEW.requester_id, NEW.requested_id),
            GREATEST(NEW.requester_id, NEW.requested_id),
            CURRENT_TIMESTAMP
        )
        ON CONFLICT (user1_id, user2_id) DO NOTHING;

        -- Update responded_at timestamp
        NEW.responded_at = CURRENT_TIMESTAMP;

        -- Create notification for requester
        INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
        SELECT
            NEW.requester_id,
            'friend_accepted',
            'Friend Request Accepted',
            u.name || ' accepted your friend request!',
            NEW.requested_id,
            NEW.id,
            'friend_request',
            jsonb_build_object(
                'friend_request_id', NEW.id,
                'accepter_name', u.name
            ),
            CURRENT_TIMESTAMP
        FROM users u WHERE u.id = NEW.requested_id;
    END IF;

    -- Create notification for rejection
    IF NEW.status = 'rejected' AND OLD.status != 'rejected' THEN
        NEW.responded_at = CURRENT_TIMESTAMP;

        INSERT INTO notifications (user_id, type, title, message, related_user_id, related_entity_id, related_entity_type, data, created_at)
        SELECT
            NEW.requester_id,
            'friend_rejected',
            'Friend Request Declined',
            u.name || ' declined your friend request.',
            NEW.requested_id,
            NEW.id,
            'friend_request',
            jsonb_build_object(
                'friend_request_id', NEW.id,
                'decliner_name', u.name
            ),
            CURRENT_TIMESTAMP
        FROM users u WHERE u.id = NEW.requested_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create enhanced challenge notification trigger
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

        -- Challenge completed
        IF NEW.status = 'completed' AND OLD.status = 'accepted' THEN
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

            -- Create completion notification for challenged user
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

-- Create achievement notification function
CREATE OR REPLACE FUNCTION create_achievement_notification(
    p_user_id UUID,
    p_achievement_id UUID,
    p_achievement_key VARCHAR,
    p_achievement_title VARCHAR,
    p_achievement_description TEXT,
    p_points_awarded INTEGER
)
RETURNS UUID AS $$
DECLARE
    notification_id UUID;
BEGIN
    INSERT INTO notifications (
        user_id,
        type,
        title,
        message,
        related_entity_id,
        related_entity_type,
        data,
        created_at
    )
    VALUES (
        p_user_id,
        'achievement_unlocked',
        'Achievement Unlocked!',
        'You unlocked: ' || p_achievement_title,
        p_achievement_id,
        'achievement',
        jsonb_build_object(
            'achievement_id', p_achievement_id,
            'achievement_key', p_achievement_key,
            'achievement_title', p_achievement_title,
            'achievement_description', p_achievement_description,
            'points_awarded', p_points_awarded
        ),
        CURRENT_TIMESTAMP
    )
    RETURNING id INTO notification_id;

    RETURN notification_id;
END;
$$ LANGUAGE plpgsql;

-- Recreate triggers
CREATE TRIGGER trigger_create_friend_request_notification
    AFTER INSERT ON friend_requests
    FOR EACH ROW
    EXECUTE FUNCTION create_friend_request_notification();

CREATE TRIGGER trigger_create_friendship_on_accept
    BEFORE UPDATE ON friend_requests
    FOR EACH ROW
    EXECUTE FUNCTION create_friendship_on_accept();

CREATE TRIGGER trigger_create_challenge_notification
    AFTER INSERT OR UPDATE ON challenges
    FOR EACH ROW
    EXECUTE FUNCTION create_challenge_notification();

-- Restore the cleanup friendship trigger that was removed in the friend table creation
CREATE OR REPLACE FUNCTION cleanup_friendship_on_status_change()
RETURNS TRIGGER AS $$
BEGIN
    -- If a friend request is cancelled or deleted, remove the friendship
    IF (OLD.status = 'accepted' AND NEW.status IN ('cancelled', 'rejected')) OR TG_OP = 'DELETE' THEN
        DELETE FROM friendships
        WHERE (user1_id = LEAST(OLD.requester_id, OLD.requested_id)
               AND user2_id = GREATEST(OLD.requester_id, OLD.requested_id));
    END IF;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_cleanup_friendship_on_status_change
    BEFORE UPDATE OR DELETE ON friend_requests
    FOR EACH ROW
    EXECUTE FUNCTION cleanup_friendship_on_status_change();