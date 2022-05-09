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
    poster_id uuid NOT NULL REFERENCES users(id),
    body varchar,
    paywall_locked boolean NOT NULL,
    paywall_tier integer,
    image_ref varchar,
    PRIMARY KEY(post_id)
);

CREATE TABLE tokens (
    id uuid NOT NULL,
    token varchar NOT NULL,
    user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(id)
);

INSERT INTO users (id, username, password) VALUES ('c03867f8-0f7c-4aef-8ff6-16ab6aa24215', '3xwr', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8');
INSERT INTO users (id, username, password) VALUES ('0622dea2-ee79-4aa9-8560-b3ba5a09fa26', 'admin', '5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8');
INSERT INTO user_subs (sub_id, user_id, subbed_to_user_id) VALUES ('e7afd1e4-ce3b-11ec-9d64-0242ac120002', 'c03867f8-0f7c-4aef-8ff6-16ab6aa24215', '0622dea2-ee79-4aa9-8560-b3ba5a09fa26');
INSERT INTO posts (post_id, poster_id, body, paywall_locked) VALUES ('3994378e-a378-4d9d-a367-a9104f2a3c43','0622dea2-ee79-4aa9-8560-b3ba5a09fa26','first ever Subby post!', false);
INSERT INTO posts (post_id, poster_id, body, paywall_locked) VALUES ('3994378e-a378-4d9d-a367-a9104fff3c43','0622dea2-ee79-4aa9-8560-b3ba5a09fa26','second Subby post!', false);