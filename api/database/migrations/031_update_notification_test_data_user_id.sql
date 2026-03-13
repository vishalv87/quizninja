-- Migration: Update notification test data to use correct user ID
-- Description: Updates the test notification data created in migration 030 to use the user ID
-- associated with "vishal.turnal@gmail.com" instead of the hardcoded test user IDs

-- First, let's find the user ID for vishal.turnal@gmail.com
-- This migration assumes the user exists in the database

-- Update all test notifications to use the correct user ID only if the target user exists
DO $$
DECLARE
    target_user_id UUID;
BEGIN
    -- Get the user ID for vishal.turnal@gmail.com if it exists
    SELECT id INTO target_user_id
    FROM users
    WHERE email = 'vishal.turnal@gmail.com'
    LIMIT 1;

    -- Only update if the user exists
    IF target_user_id IS NOT NULL THEN
        UPDATE notifications
        SET user_id = target_user_id
        WHERE is_test_data = true
        AND user_id IN (
            'a1b1c1d1-e1f1-4001-8001-111111111111',
            'a2b2c2d2-e2f2-4002-8002-222222222222',
            'a3b3c3d3-e3f3-4003-8003-333333333333'
        );
    ELSE
        -- If the user doesn't exist, just delete the test notifications
        DELETE FROM notifications
        WHERE is_test_data = true
        AND user_id IN (
            'a1b1c1d1-e1f1-4001-8001-111111111111',
            'a2b2c2d2-e2f2-4002-8002-222222222222',
            'a3b3c3d3-e3f3-4003-8003-333333333333'
        );
    END IF;
END $$;

-- Update related_user_id fields in test notifications to use real user IDs where appropriate
-- For friend requests and other notifications that reference other users, we'll keep them as test users
-- or set them to the main user for demonstration purposes

-- Update notifications where the related_user_id was pointing to test users
-- Set related_user_id to the main user for notifications that need a related user
DO $$
DECLARE
    target_user_id UUID;
BEGIN
    -- Get the user ID for vishal.turnal@gmail.com if it exists
    SELECT id INTO target_user_id
    FROM users
    WHERE email = 'vishal.turnal@gmail.com'
    LIMIT 1;

    -- Only update if the user exists
    IF target_user_id IS NOT NULL THEN
        UPDATE notifications
        SET related_user_id = target_user_id
        WHERE is_test_data = true
        AND related_user_id IN (
            'a2b2c2d2-e2f2-4002-8002-222222222222',
            'a3b3c3d3-e3f3-4003-8003-333333333333',
            'a4b4c4d4-e4f4-4004-8004-444444444444',
            'a5b5c5d5-e5f5-4005-8005-555555555555'
        );
    END IF;
END $$;

-- Clean up the test users that were created in migration 030 since we don't need them anymore
DELETE FROM users WHERE id IN (
    'a1b1c1d1-e1f1-4001-8001-111111111111',
    'a2b2c2d2-e2f2-4002-8002-222222222222',
    'a3b3c3d3-e3f3-4003-8003-333333333333',
    'a4b4c4d4-e4f4-4004-8004-444444444444',
    'a5b5c5d5-e5f5-4005-8005-555555555555'
) AND is_test_data = true;

-- Update the notification messages to reflect that they're for the real user
-- Update friend request messages to use generic names since related users are now the same
UPDATE notifications
SET
    message = CASE
        WHEN type = 'friend_request' AND message LIKE '%Bob Smith%' THEN 'A friend sent you a friend request!'
        WHEN type = 'friend_request' AND message LIKE '%Alice Johnson%' THEN 'A friend sent you a friend request!'
        WHEN type = 'challenge_received' AND message LIKE '%Carol Davis%' THEN 'You received a new quiz challenge!'
        WHEN type = 'challenge_completed' AND message LIKE '%David Wilson%' THEN 'You completed a challenge!'
        WHEN type = 'challenge_accepted' AND message LIKE '%Bob Smith%' THEN 'Your challenge was accepted!'
        WHEN type = 'friend_accepted' AND message LIKE '%Eva Martinez%' THEN 'Your friend request was accepted!'
        ELSE message
    END,
    data = CASE
        WHEN type = 'friend_request' THEN '{"requester_name": "Friend", "message": "Hey! Let''s be friends!"}'
        WHEN type = 'challenge_received' THEN '{"challenger_name": "Challenger", "message": "Think you can beat this quiz score?", "expires_at": "2025-10-08T00:00:00Z"}'
        WHEN type = 'challenge_completed' THEN '{"opponent_name": "Opponent", "your_score": 85, "opponent_score": 72, "result": "won"}'
        WHEN type = 'challenge_accepted' THEN '{"accepter_name": "Friend"}'
        WHEN type = 'friend_accepted' THEN '{"accepter_name": "Friend"}'
        ELSE data
    END
WHERE is_test_data = true
AND type IN ('friend_request', 'friend_accepted', 'challenge_received', 'challenge_completed', 'challenge_accepted');