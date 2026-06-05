-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'buyer';

INSERT INTO users (username, password_hash, password_salt, role)
VALUES (
    'seller1',
    '$2a$12$pryu/2vk1sJPT1bepJfbreG8uQyHhORFHUCsgCOYXAZzJfoH6WJyy',
    '0000000000000000000000000000000000000000000000000000000000000000',
    'seller'
) ON CONFLICT (username) DO NOTHING;

-- +goose Down
DELETE FROM users WHERE username = 'seller1';
ALTER TABLE users DROP COLUMN IF EXISTS role;
