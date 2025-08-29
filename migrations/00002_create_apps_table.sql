-- +goose Up
CREATE TABLE apps (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    secret TEXT NOT NULL
);

-- +goose Down
DROP TABLE apps;
