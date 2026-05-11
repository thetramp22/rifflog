CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash NOT NULL TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);