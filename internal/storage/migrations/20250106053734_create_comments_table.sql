-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments (
                          id         SERIAL PRIMARY KEY,
                          content    TEXT NOT NULL,
                          user_id    INT NOT NULL,
                          post_id    INT NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT fk_user_comment FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
                          CONSTRAINT fk_post_comment FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comments;
-- +goose StatementEnd