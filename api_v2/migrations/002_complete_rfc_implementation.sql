-- Migration 002: Complete RFC implementation
-- Updates users table and adds all required tables per RFC-001

-- Update users table to include new fields from RFC
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_private BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS theme VARCHAR(10) DEFAULT 'light';

-- Rename avatar_url to profile_picture_url for consistency
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'users' AND column_name = 'avatar_url') THEN
        ALTER TABLE users RENAME COLUMN avatar_url TO profile_picture_url;
    END IF;
END $$;

-- Add constraints for theme
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.constraint_column_usage WHERE constraint_name = 'users_theme_check') THEN
        ALTER TABLE users ADD CONSTRAINT users_theme_check CHECK (theme IN ('light', 'dark'));
    END IF;
END $$;

-- Create user_sessions table for session management
CREATE TABLE IF NOT EXISTS user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(512) UNIQUE NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at);

-- Create movies table with TMDb integration
CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    external_api_id VARCHAR(50) UNIQUE NOT NULL, -- TMDb ID
    title VARCHAR(500) NOT NULL,
    overview TEXT,
    release_date DATE,
    poster_url TEXT,
    backdrop_url TEXT,
    genres TEXT[], -- Array of genre names
    runtime INTEGER, -- minutes
    vote_average DECIMAL(3,1),
    vote_count INTEGER,
    adult BOOLEAN DEFAULT FALSE,
    cache_expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_movies_external_api_id ON movies(external_api_id);
CREATE INDEX IF NOT EXISTS idx_movies_title ON movies USING GIN(to_tsvector('english', title));
CREATE INDEX IF NOT EXISTS idx_movies_genres ON movies USING GIN(genres);
CREATE INDEX IF NOT EXISTS idx_movies_cache_expires_at ON movies(cache_expires_at);

-- Create reviews table for movie ratings and reviews
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 10),
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, movie_id) -- One review per user per movie
);

CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_movie_id ON reviews(movie_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);
CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews(created_at);

-- Create movie lists table
CREATE TABLE IF NOT EXISTS movie_lists (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_movie_lists_user_id ON movie_lists(user_id);
CREATE INDEX IF NOT EXISTS idx_movie_lists_name ON movie_lists(name);

-- Create movie list entries table
CREATE TABLE IF NOT EXISTS movie_list_entries (
    id SERIAL PRIMARY KEY,
    movie_list_id INTEGER NOT NULL REFERENCES movie_lists(id) ON DELETE CASCADE,
    movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    position INTEGER, -- for ordered lists
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(movie_list_id, movie_id) -- One movie per list
);

CREATE INDEX IF NOT EXISTS idx_movie_list_entries_list_id ON movie_list_entries(movie_list_id);
CREATE INDEX IF NOT EXISTS idx_movie_list_entries_movie_id ON movie_list_entries(movie_id);

-- Create friendships table for friend relationships
CREATE TABLE IF NOT EXISTS friendships (
    id SERIAL PRIMARY KEY,
    user_id_1 INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_id_2 INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'accepted', 'blocked')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id_1, user_id_2),
    CHECK (user_id_1 != user_id_2)
);

CREATE INDEX IF NOT EXISTS idx_friendships_user_id_1 ON friendships(user_id_1);
CREATE INDEX IF NOT EXISTS idx_friendships_user_id_2 ON friendships(user_id_2);
CREATE INDEX IF NOT EXISTS idx_friendships_status ON friendships(status);

-- Create follows table for following relationships
CREATE TABLE IF NOT EXISTS follows (
    id SERIAL PRIMARY KEY,
    follower_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE INDEX IF NOT EXISTS idx_follows_follower_id ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_follows_following_id ON follows(following_id);

-- Create posts table for social posts
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    visibility VARCHAR(20) NOT NULL DEFAULT 'public' CHECK (visibility IN ('public', 'friends', 'private')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_visibility ON posts(visibility);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);

-- Create match sessions table for collaborative movie matching
CREATE TABLE IF NOT EXISTS match_sessions (
    id SERIAL PRIMARY KEY,
    host_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'finished', 'cancelled')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_match_sessions_host_user_id ON match_sessions(host_user_id);
CREATE INDEX IF NOT EXISTS idx_match_sessions_status ON match_sessions(status);
CREATE INDEX IF NOT EXISTS idx_match_sessions_created_at ON match_sessions(created_at);

-- Create match session participants table
CREATE TABLE IF NOT EXISTS match_session_participants (
    id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL REFERENCES match_sessions(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(session_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_match_session_participants_session_id ON match_session_participants(session_id);
CREATE INDEX IF NOT EXISTS idx_match_session_participants_user_id ON match_session_participants(user_id);

-- Create match interactions table for recording likes/dislikes
CREATE TABLE IF NOT EXISTS match_interactions (
    id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL REFERENCES match_sessions(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    liked BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(session_id, user_id, movie_id)
);

CREATE INDEX IF NOT EXISTS idx_match_interactions_session_id ON match_interactions(session_id);
CREATE INDEX IF NOT EXISTS idx_match_interactions_user_id ON match_interactions(user_id);
CREATE INDEX IF NOT EXISTS idx_match_interactions_movie_id ON match_interactions(movie_id);
CREATE INDEX IF NOT EXISTS idx_match_interactions_liked ON match_interactions(liked);

-- Create default movie lists for existing users
INSERT INTO movie_lists (user_id, name, description, is_public)
SELECT 
    u.id,
    'Want to Watch',
    'Movies I want to watch',
    FALSE
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM movie_lists ml 
    WHERE ml.user_id = u.id AND ml.name = 'Want to Watch'
);

INSERT INTO movie_lists (user_id, name, description, is_public)
SELECT 
    u.id,
    'Watched',
    'Movies I have watched',
    TRUE
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM movie_lists ml 
    WHERE ml.user_id = u.id AND ml.name = 'Watched'
);

-- Update existing users with default values for new fields
UPDATE users SET 
    password_hash = '$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$hash',  -- placeholder hash
    email_verified = TRUE,
    theme = 'light'
WHERE password_hash IS NULL;
