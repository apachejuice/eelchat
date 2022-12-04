
-- +migrate Up
CREATE TABLE user (
    id              CHAR(22)        PRIMARY KEY, -- base64'd uuid
    username        VARCHAR(100)    NOT NULL,
    discriminator   CHAR(4)         NOT NULL,
    password_hash   CHAR(16)        NOT NULL,
    registered_at   DATETIME        NOT NULL,
    last_login      DATETIME,
    email           VARCHAR(255)
);

-- +migrate Down
DROP TABLE user;
