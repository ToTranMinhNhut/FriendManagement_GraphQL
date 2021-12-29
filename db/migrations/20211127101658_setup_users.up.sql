-- Setup database for friends managment and dummy data.

BEGIN;

-- Setup users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);
CREATE INDEX email_on_users ON users(email);

-- Setup friends table
CREATE TABLE friends (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users NOT NULL,
    friend_id INTEGER REFERENCES users NOT NULL,
    CONSTRAINT constraint_friends_pkey UNIQUE (user_id, friend_id)
);

-- Setup subscriptions table
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    subscription_requestor_id INTEGER REFERENCES users NOT NULL,
    subscription_target_id INTEGER REFERENCES users NOT NULL,
    CONSTRAINT constraint_subscriptions_pkey UNIQUE (subscription_requestor_id, subscription_target_id)
);

-- Setup user_blocks table
CREATE TABLE user_blocks (
    id SERIAL PRIMARY KEY,
    requestor_id INTEGER REFERENCES users NOT NULL,
    target_id INTEGER REFERENCES users NOT NULL,
    CONSTRAINT constraint_user_blocks_pkey UNIQUE (requestor_id, target_id)
);

-- dummy data for testing
INSERT INTO users(name, email, created_at, updated_at) VALUES
('john','john@example.com', now(), now()),
('andy','andy@example.com', now(), now()),
('common','common@example.com', now(), now()),
('lisa','lisa@example.com', now(), now()),
('kate','kate@example.com', now(), now());

COMMIT;
