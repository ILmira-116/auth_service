-- +goose Up
INSERT INTO apps (id, name, secret) VALUES (1, 'test', 'test_secret');

-- +goose Down
DELETE FROM apps WHERE id = 1;