-- Migration to add provider tracking and sync timestamp
-- Date: 2025-11-08

-- Add provider column to track data source
ALTER TABLE movies 
ADD COLUMN IF NOT EXISTS provider VARCHAR(20) DEFAULT 'omdb' CHECK (provider IN ('omdb', 'tmdb', 'internal'));

-- Add last_sync_at to track when data was last updated from external API
ALTER TABLE movies 
ADD COLUMN IF NOT EXISTS last_sync_at TIMESTAMP WITH TIME ZONE;

-- Create index for provider lookups
CREATE INDEX IF NOT EXISTS idx_movies_provider ON movies(provider);

-- Create index for sync tracking
CREATE INDEX IF NOT EXISTS idx_movies_last_sync_at ON movies(last_sync_at);

-- Update existing movies to have provider 'internal' if they don't have one
UPDATE movies SET provider = 'internal' WHERE provider IS NULL;
