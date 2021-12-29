-- Delete all existing data before adding mock testdata
TRUNCATE TABLE users CASCADE;
TRUNCATE TABLE friends CASCADE;
TRUNCATE TABLE subscriptions CASCADE;
TRUNCATE TABLE user_blocks CASCADE;


INSERT INTO users(id, name, email, created_at, updated_at) VALUES
(100, 'john','john@example.com', now(), now()),
(101, 'andy','andy@example.com', now(), now()),
(102, 'common','common@example.com', now(), now()),
(103, 'lisa','lisa@example.com', now(), now()),
(104, 'kate','kate@example.com', now(), now());

INSERT INTO friends(user_id, friend_id) VALUES
(100, 102),
(101, 102),
(102, 103);

INSERT INTO user_blocks(requestor_id, target_id) VALUES (100,103);
INSERT INTO user_blocks(requestor_id, target_id) VALUES (100,104);

INSERT INTO subscriptions(subscription_requestor_id, subscription_target_id) VALUES (101,103);


