-- +goose Up
-- +goose StatementBegin
CREATE TABLE likes
(
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL,
    post_id    INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_like FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_post_like FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS likes;
-- +goose StatementEnd