CREATE TABLE users (
    id uuid NOT NULL,
    email varchar,
    username varchar NOT NULL UNIQUE,
    password varchar NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE user_subs (
    sub_id uuid NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id),
    subbed_to_user_id uuid NOT NULL REFERENCES users(id),
    PRIMARY KEY(sub_id)
);

CREATE TABLE posts (
    post_id uuid NOT NULL,
    posted_at timestamptz NOT NULL,
    poster_id uuid NOT NULL REFERENCES users(id),
    body varchar,
    membership_locked boolean NOT NULL,
    membership_tier integer,
    image_ref varchar,
    PRIMARY KEY(post_id)
);

CREATE TABLE tokens (
    id uuid NOT NULL,
    token varchar NOT NULL,
    user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(id)
);

CREATE TABLE memberships (
    id uuid NOT NULL,
    owner_id uuid NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE tiers (
    id uuid NOT NULL,
    name varchar NOT NULL,
    price bigint NOT NULL,
    rewards varchar NOT NULL,
    membership_id uuid NOT NULL REFERENCES memberships(id) ON DELETE CASCADE,
    PRIMARY KEY(id)
);

CREATE TABLE members (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tier_id uuid NOT NULL REFERENCES tiers(id) ON DELETE CASCADE
);

INSERT INTO users (id, username, password) VALUES ('c03867f8-0f7c-4aef-8ff6-16ab6aa24215', '3xwr', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8');
INSERT INTO users (id, username, password) VALUES ('0622dea2-ee79-4aa9-8560-b3ba5a09fa26', 'admin', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8');
INSERT INTO users (id, username, password) VALUES ('abd4528d-cd53-4366-83a7-1a12739904f5', 'atd_mf', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8');
INSERT INTO user_subs (sub_id, user_id, subbed_to_user_id) VALUES ('e7afd1e4-ce3b-11ec-9d64-0242ac120002', 'c03867f8-0f7c-4aef-8ff6-16ab6aa24215', '0622dea2-ee79-4aa9-8560-b3ba5a09fa26');
INSERT INTO posts (post_id, posted_at, poster_id, body, membership_locked) VALUES ('3994378e-a378-4d9d-a367-a9104f2a3c43','2014-04-04 06:00:00','0622dea2-ee79-4aa9-8560-b3ba5a09fa26','first ever Subby post!', false);
INSERT INTO posts (post_id, posted_at, poster_id, body, membership_locked) VALUES ('3994378e-a378-4d9d-a367-a9104fff3c43','2014-04-04 07:00:00','0622dea2-ee79-4aa9-8560-b3ba5a09fa26','second Subby post!', false);
INSERT INTO memberships (id, owner_id) VALUES ('096c791f-f42b-4fa6-a303-0046e6c09b15','abd4528d-cd53-4366-83a7-1a12739904f5');
INSERT INTO tiers (id, name, price, rewards, membership_id) VALUES ('540722b9-9bd9-4fe4-896b-fc7985cdc6bb', 'Biggest Fan', 2000, 'Shoutout at the end of the video','096c791f-f42b-4fa6-a303-0046e6c09b15');
INSERT INTO members (user_id, tier_id) VALUES ('c03867f8-0f7c-4aef-8ff6-16ab6aa24215', '540722b9-9bd9-4fe4-896b-fc7985cdc6bb');