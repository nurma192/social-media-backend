-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       content VARCHAR(255) NOT NULL,
                       user_id INT REFERENCES users(id) ON DELETE CASCADE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
-- +goose StatementEnd
