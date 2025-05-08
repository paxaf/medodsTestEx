-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    guid UUID PRIMARY KEY,
    refresh_token VARCHAR(100) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;