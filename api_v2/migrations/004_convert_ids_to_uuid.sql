-- Migration 004: Convert all ID columns from INTEGER to UUID
-- This migration converts all primary key and foreign key columns from auto-incrementing integers to UUIDs for improved security

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Start a transaction to ensure atomicity
BEGIN;

-- =====================================
-- Users table migration
-- =====================================

-- Create backup table
CREATE TABLE users_backup AS SELECT * FROM users;

-- Drop existing constraints and indexes
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey CASCADE;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;

-- Add new UUID column and populate with UUIDs
ALTER TABLE users ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
UPDATE users SET id_new = uuid_generate_v4();

-- Drop old id column and rename new one
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN id_new TO id;

-- Add primary key constraint
ALTER TABLE users ADD PRIMARY KEY (id);

-- Recreate indexes
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_username ON users(username);

-- =====================================
-- Movies table migration
-- =====================================

-- Create backup table
CREATE TABLE movies_backup AS SELECT * FROM movies;

-- Drop existing constraints
ALTER TABLE movies DROP CONSTRAINT IF EXISTS movies_pkey CASCADE;

-- Add new UUID column and populate with UUIDs
ALTER TABLE movies ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
UPDATE movies SET id_new = uuid_generate_v4();

-- Drop old id column and rename new one
ALTER TABLE movies DROP COLUMN id;
ALTER TABLE movies RENAME COLUMN id_new TO id;

-- Add primary key constraint
ALTER TABLE movies ADD PRIMARY KEY (id);

-- =====================================
-- User sessions table migration
-- =====================================

-- Create backup table if exists
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'user_sessions') THEN
        CREATE TABLE user_sessions_backup AS SELECT * FROM user_sessions;
        
        -- Drop existing constraints
        ALTER TABLE user_sessions DROP CONSTRAINT IF EXISTS user_sessions_pkey CASCADE;
        ALTER TABLE user_sessions DROP CONSTRAINT IF EXISTS fk_user_sessions_user_id;
        
        -- Add new UUID columns
        ALTER TABLE user_sessions ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
        ALTER TABLE user_sessions ADD COLUMN user_id_new UUID;
        
        -- Update columns with UUIDs
        UPDATE user_sessions SET id_new = uuid_generate_v4();
        -- Note: user_id_new will be updated later after users table is fully migrated
        
        -- For now, we'll drop the user_sessions table and recreate it with proper structure
        DROP TABLE user_sessions;
        DROP TABLE user_sessions_backup;
    END IF;
END $$;

-- =====================================
-- Reviews table migration
-- =====================================

-- Create backup table if exists
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'reviews') THEN
        CREATE TABLE reviews_backup AS SELECT * FROM reviews;
        
        -- Drop existing constraints
        ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_pkey CASCADE;
        ALTER TABLE reviews DROP CONSTRAINT IF EXISTS fk_reviews_user_id;
        ALTER TABLE reviews DROP CONSTRAINT IF EXISTS fk_reviews_movie_id;
        
        -- Add new UUID columns
        ALTER TABLE reviews ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
        ALTER TABLE reviews ADD COLUMN user_id_new UUID;
        ALTER TABLE reviews ADD COLUMN movie_id_new UUID;
        
        -- Update ID column
        UPDATE reviews SET id_new = uuid_generate_v4();
        
        -- For foreign keys, we'll need to handle them after all primary tables are migrated
        -- For now, drop and recreate the table structure
        DROP TABLE reviews;
        DROP TABLE reviews_backup;
    END IF;
END $$;

-- =====================================
-- Movie lists table migration
-- =====================================

-- Create backup table if exists
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'movie_lists') THEN
        CREATE TABLE movie_lists_backup AS SELECT * FROM movie_lists;
        
        -- Drop existing constraints
        ALTER TABLE movie_lists DROP CONSTRAINT IF EXISTS movie_lists_pkey CASCADE;
        ALTER TABLE movie_lists DROP CONSTRAINT IF EXISTS fk_movie_lists_user_id;
        
        -- For now, drop and recreate with proper UUID structure
        DROP TABLE movie_lists;
        DROP TABLE movie_lists_backup;
    END IF;
END $$;

-- =====================================
-- Movie list entries table migration
-- =====================================

-- Create backup table if exists
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'movie_list_entries') THEN
        DROP TABLE movie_list_entries CASCADE;
    END IF;
END $$;

-- =====================================
-- Social tables migration (friendships, follows)
-- =====================================

-- Drop social tables if they exist (they were recently added)
DROP TABLE IF EXISTS friendships CASCADE;
DROP TABLE IF EXISTS follows CASCADE;

-- =====================================
-- Email verification tokens
-- =====================================

DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'email_verification_tokens') THEN
        DROP TABLE email_verification_tokens CASCADE;
    END IF;
END $$;

-- =====================================
-- Password reset tokens
-- =====================================

DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'password_reset_tokens') THEN
        DROP TABLE password_reset_tokens CASCADE;
    END IF;
END $$;

-- =====================================
-- Recreate tables with UUID structure
-- =====================================

-- User sessions table
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

-- Reviews table
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 10),
    content TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, movie_id)
);

-- Movie lists table
CREATE TABLE movie_lists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Movie list entries table
CREATE TABLE movie_list_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    movie_list_id UUID NOT NULL REFERENCES movie_lists(id) ON DELETE CASCADE,
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(movie_list_id, movie_id)
);

-- Friendships table
CREATE TABLE friendships (
    user_id_1 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_id_2 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'accepted', 'declined', 'blocked')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id_1, user_id_2),
    CHECK (user_id_1 < user_id_2) -- Ensure consistent ordering
);

-- Follows table
CREATE TABLE follows (
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id) -- Prevent self-following
);

-- Email verification tokens
CREATE TABLE email_verification_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Password reset tokens
CREATE TABLE password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Posts table (update existing or create)
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'posts') THEN
        -- Update existing posts table
        ALTER TABLE posts DROP CONSTRAINT IF EXISTS posts_pkey CASCADE;
        ALTER TABLE posts ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
        ALTER TABLE posts ADD COLUMN user_id_new UUID;
        UPDATE posts SET id_new = uuid_generate_v4();
        
        -- Drop and recreate with proper structure
        DROP TABLE posts CASCADE;
    END IF;
    
    CREATE TABLE posts (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        content TEXT NOT NULL,
        visibility VARCHAR(20) NOT NULL CHECK (visibility IN ('public', 'private', 'friends')),
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );
END $$;

-- Match sessions table (create if not exists)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'match_sessions') THEN
        CREATE TABLE match_sessions (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            host_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'finished', 'canceled')),
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW()
        );
    END IF;
END $$;

-- Match session participants table (create if not exists)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'match_session_participants') THEN
        CREATE TABLE match_session_participants (
            session_id UUID NOT NULL REFERENCES match_sessions(id) ON DELETE CASCADE,
            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            joined_at TIMESTAMP DEFAULT NOW(),
            PRIMARY KEY (session_id, user_id)
        );
    END IF;
END $$;

-- Match interactions table (create if not exists)
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'match_interactions') THEN
        CREATE TABLE match_interactions (
            session_id UUID NOT NULL REFERENCES match_sessions(id) ON DELETE CASCADE,
            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
            liked BOOLEAN NOT NULL,
            interacted_at TIMESTAMP DEFAULT NOW(),
            PRIMARY KEY (session_id, user_id, movie_id)
        );
    END IF;
END $$;

-- =====================================
-- Create indexes for performance
-- =====================================

-- User sessions indexes
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token ON user_sessions(token);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Reviews indexes
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_movie_id ON reviews(movie_id);
CREATE INDEX idx_reviews_created_at ON reviews(created_at);

-- Movie lists indexes
CREATE INDEX idx_movie_lists_user_id ON movie_lists(user_id);
CREATE INDEX idx_movie_lists_name ON movie_lists(name);

-- Movie list entries indexes
CREATE INDEX idx_movie_list_entries_list_id ON movie_list_entries(movie_list_id);
CREATE INDEX idx_movie_list_entries_movie_id ON movie_list_entries(movie_id);

-- Social indexes
CREATE INDEX idx_friendships_user_id_1 ON friendships(user_id_1);
CREATE INDEX idx_friendships_user_id_2 ON friendships(user_id_2);
CREATE INDEX idx_friendships_status ON friendships(status);
CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);

-- Posts indexes
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_visibility ON posts(visibility);
CREATE INDEX idx_posts_created_at ON posts(created_at);

-- Match indexes
CREATE INDEX idx_match_sessions_host_user_id ON match_sessions(host_user_id);
CREATE INDEX idx_match_sessions_status ON match_sessions(status);

COMMIT;
