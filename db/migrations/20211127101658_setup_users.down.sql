-- Reverses the corresponding up script

BEGIN;

DROP TABLE friends;
DROP TABLE subscriptions;
DROP TABLE user_blocks;
DROP TABLE users CASCADE;

COMMIT;
