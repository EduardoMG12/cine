-- Migration: 003_social_features.sql
-- Create tables for friendship and follow functionality

-- Friendships table for friend requests and friend relationships
CREATE TABLE IF NOT EXISTS friendships (
    user_id_1 INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_id_2 INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status friendship_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id_1, user_id_2),
    CONSTRAINT friendship_different_users CHECK (user_id_1 != user_id_2),
    CONSTRAINT friendship_ordered_users CHECK (user_id_1 < user_id_2)
);

-- Create friendship status enum if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'friendship_status') THEN
        CREATE TYPE friendship_status AS ENUM ('pending', 'accepted', 'declined', 'blocked');
    END IF;
END
$$;

-- Follows table for follow relationships (asymmetric)
CREATE TABLE IF NOT EXISTS follows (
    follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CONSTRAINT follow_different_users CHECK (follower_id != following_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_friendships_user_id_1 ON friendships(user_id_1);
CREATE INDEX IF NOT EXISTS idx_friendships_user_id_2 ON friendships(user_id_2);
CREATE INDEX IF NOT EXISTS idx_friendships_status ON friendships(status);
CREATE INDEX IF NOT EXISTS idx_friendships_created_at ON friendships(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_follows_follower_id ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_follows_following_id ON follows(following_id);
CREATE INDEX IF NOT EXISTS idx_follows_created_at ON follows(created_at DESC);

-- Function to ensure friendship user order (user_id_1 < user_id_2)
CREATE OR REPLACE FUNCTION ensure_friendship_order()
RETURNS TRIGGER AS $$
BEGIN
    -- Ensure user_id_1 is always less than user_id_2 for consistency
    IF NEW.user_id_1 > NEW.user_id_2 THEN
        -- Swap the user IDs
        DECLARE
            temp_id INTEGER;
        BEGIN
            temp_id := NEW.user_id_1;
            NEW.user_id_1 := NEW.user_id_2;
            NEW.user_id_2 := temp_id;
        END;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to enforce friendship user order
DROP TRIGGER IF EXISTS trigger_ensure_friendship_order ON friendships;
CREATE TRIGGER trigger_ensure_friendship_order
    BEFORE INSERT OR UPDATE ON friendships
    FOR EACH ROW
    EXECUTE FUNCTION ensure_friendship_order();

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to auto-update updated_at on friendships
DROP TRIGGER IF EXISTS trigger_friendships_updated_at ON friendships;
CREATE TRIGGER trigger_friendships_updated_at
    BEFORE UPDATE ON friendships
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
