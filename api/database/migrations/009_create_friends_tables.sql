-- Migration: Create friends system tables
-- Created: 2024-09-20

-- Create friend_requests table
CREATE TABLE friend_requests (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    requester_id UUID NOT NULL,
    requested_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    responded_at TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (requester_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (requested_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (status IN ('pending', 'accepted', 'rejected', 'cancelled')),
    CHECK (requester_id != requested_id),
    UNIQUE (requester_id, requested_id)
);

-- Create friendships table (for accepted friend relationships)
CREATE TABLE friendships (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user1_id UUID NOT NULL,
    user2_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (user1_id != user2_id),
    CHECK (user1_id < user2_id), -- Ensure consistent ordering to prevent duplicates
    UNIQUE (user1_id, user2_id)
);

-- Create friend_notifications table
CREATE TABLE friend_notifications (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT,
    related_user_id UUID,
    friend_request_id UUID,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (related_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (friend_request_id) REFERENCES friend_requests(id) ON DELETE CASCADE,
    CHECK (type IN ('friend_request', 'friend_accepted', 'friend_rejected'))
);

-- Create indexes for better performance
CREATE INDEX idx_friend_requests_requester_id ON friend_requests(requester_id);
CREATE INDEX idx_friend_requests_requested_id ON friend_requests(requested_id);
CREATE INDEX idx_friend_requests_status ON friend_requests(status);
CREATE INDEX idx_friend_requests_created_at ON friend_requests(created_at);

CREATE INDEX idx_friendships_user1_id ON friendships(user1_id);
CREATE INDEX idx_friendships_user2_id ON friendships(user2_id);
CREATE INDEX idx_friendships_created_at ON friendships(created_at);

CREATE INDEX idx_friend_notifications_user_id ON friend_notifications(user_id);
CREATE INDEX idx_friend_notifications_type ON friend_notifications(type);
CREATE INDEX idx_friend_notifications_is_read ON friend_notifications(is_read);
CREATE INDEX idx_friend_notifications_created_at ON friend_notifications(created_at);
CREATE INDEX idx_friend_notifications_related_user_id ON friend_notifications(related_user_id);

-- Create trigger to automatically create friendship when friend request is accepted
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
        INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, friend_request_id, created_at)
        SELECT
            NEW.requester_id,
            'friend_accepted',
            'Friend Request Accepted',
            u.name || ' accepted your friend request!',
            NEW.requested_id,
            NEW.id,
            CURRENT_TIMESTAMP
        FROM users u WHERE u.id = NEW.requested_id;
    END IF;

    -- Create notification for rejection
    IF NEW.status = 'rejected' AND OLD.status != 'rejected' THEN
        NEW.responded_at = CURRENT_TIMESTAMP;

        INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, friend_request_id, created_at)
        SELECT
            NEW.requester_id,
            'friend_rejected',
            'Friend Request Declined',
            u.name || ' declined your friend request.',
            NEW.requested_id,
            NEW.id,
            CURRENT_TIMESTAMP
        FROM users u WHERE u.id = NEW.requested_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_friendship_on_accept
    BEFORE UPDATE ON friend_requests
    FOR EACH ROW
    EXECUTE FUNCTION create_friendship_on_accept();

-- Create trigger to create notification when friend request is created
CREATE OR REPLACE FUNCTION create_friend_request_notification()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO friend_notifications (user_id, type, title, message, related_user_id, friend_request_id, created_at)
    SELECT
        NEW.requested_id,
        'friend_request',
        'New Friend Request',
        u.name || ' sent you a friend request!',
        NEW.requester_id,
        NEW.id,
        CURRENT_TIMESTAMP
    FROM users u WHERE u.id = NEW.requester_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_friend_request_notification
    AFTER INSERT ON friend_requests
    FOR EACH ROW
    EXECUTE FUNCTION create_friend_request_notification();

-- Create trigger to cleanup friendship when friend request is cancelled/rejected after being accepted
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