-- +goose Up
ALTER TABLE users ADD COLUMN useragent VARCHAR(100);

-- +goose Down
ALTER TABLE users DROP COLUMN useragent;