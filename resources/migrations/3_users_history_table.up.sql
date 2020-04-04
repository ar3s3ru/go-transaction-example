CREATE TABLE users_history (
    user_history_id SERIAL PRIMARY KEY,
    user_id         INTEGER     NOT NULL REFERENCES users(user_id),
    action          VARCHAR(32) NOT NULL,
    state           JSONB       NOT NULL,
    created_at      TIMESTAMP DEFAULT current_timestamp
);
