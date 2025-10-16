-- Initial database schema for CineVerse
-- Creates users table

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(30) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    bio TEXT,
    avatar_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

INSERT INTO users (username, email, display_name, bio) VALUES 
    ('cinefilo', 'cinefilo@exemplo.com', 'Cineasta Apaixonado', 'Adoro filmes clássicos e independentes'),
    ('moviefan', 'fan@exemplo.com', 'Movie Fan', 'Sempre em busca do próximo blockbuster')
ON CONFLICT (username) DO NOTHING;
