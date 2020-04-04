CREATE TABLE users (
    user_id    SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    age        INT          NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp
);

CREATE UNIQUE INDEX unique_users_name ON users (name);

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
