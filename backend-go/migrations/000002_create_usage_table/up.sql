CREATE TABLE usage_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    video_name TEXT,
    prompt_text TEXT,
    summary_length INTEGER,
    processing_time INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);