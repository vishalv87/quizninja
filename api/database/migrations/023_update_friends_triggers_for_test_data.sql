-- Migration: Update friends triggers to set is_test_data = true
-- Date: 2025-09-27
-- Description: Updates the existing friends triggers to include is_test_data = true when creating friendships and notifications

-- Drop and recreate the trigger function to include is_test_data = true
DROP TRIGGER IF EXISTS trigger_create_friendship_on_accept ON friend_requests;
DROP FUNCTION IF EXISTS create_friendship_on_accept();

CREATE OR REPLACE FUNCTION create_friendship_on_accept()
RETURNS TRIGGER AS $$
BEGIN
    -- Only create friendship if status changed to 'accepted'
    IF NEW.status = 'accepted' AND OLD.status != 'accepted' THEN
        -- Ensure user1_id < user2_id for consistent ordering
        INSERT INTO friendships (user1_id, user2_id, created_at, is_test_data)
        VALUES (
            LEAST(NEW.requester_id, NEW.requested_id),
            GREATEST(NEW.requester_id, NEW.requested_id),
            CURRENT_TIMESTAMP,
            true
        )
        ON CONFLICT (user1_id, user2_id) DO NOTHING;

        -- Update responded_at timestamp
        NEW.responded_at = CURRENT_TIMESTAMP;

        -- Create notification for requester
        INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, friend_request_id, created_at, is_test_data)
        SELECT
            NEW.requester_id,
            'friend_accepted',
            'Friend Request Accepted',
            u.name || ' accepted your friend request!',
            NEW.requested_id,
            NEW.id,
            CURRENT_TIMESTAMP,
            true
        FROM users u WHERE u.id = NEW.requested_id;
    END IF;

    -- Create notification for rejection
    IF NEW.status = 'rejected' AND OLD.status != 'rejected' THEN
        NEW.responded_at = CURRENT_TIMESTAMP;

        INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, friend_request_id, created_at, is_test_data)
        SELECT
            NEW.requester_id,
            'friend_rejected',
            'Friend Request Declined',
            u.name || ' declined your friend request.',
            NEW.requested_id,
            NEW.id,
            CURRENT_TIMESTAMP,
            true
        FROM users u WHERE u.id = NEW.requested_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_friendship_on_accept
    BEFORE UPDATE ON friend_requests
    FOR EACH ROW
    EXECUTE FUNCTION create_friendship_on_accept();

-- Drop and recreate the friend request notification trigger to include is_test_data = true
DROP TRIGGER IF EXISTS trigger_create_friend_request_notification ON friend_requests;
DROP FUNCTION IF EXISTS create_friend_request_notification();

CREATE OR REPLACE FUNCTION create_friend_request_notification()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, friend_request_id, created_at, is_test_data)
    SELECT
        NEW.requested_id,
        'friend_request',
        'New Friend Request',
        u.name || ' sent you a friend request!',
        NEW.requester_id,
        NEW.id,
        CURRENT_TIMESTAMP,
        true
    FROM users u WHERE u.id = NEW.requester_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_friend_request_notification
    AFTER INSERT ON friend_requests
    FOR EACH ROW
    EXECUTE FUNCTION create_friend_request_notification();