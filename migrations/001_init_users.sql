-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    guid UUID PRIMARY KEY,
    refreshtoken VARCHAR(100) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;