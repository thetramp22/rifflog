CREATE TABLE practice_sessions (
    id SERIAL PRIMARY KEY,
    skill_id INTEGER NOT NULL REFERENCES skills(id),
    duration_minutes INTEGER NOT NULL CHECK (duration_minutes > 0),
    notes TEXT NOT NULL,
    practiced_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);