-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts
(
    id         SERIAL PRIMARY KEY,
    content    VARCHAR(255) NOT NULL,
    user_id    INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE postImages
(
    id       SERIAL PRIMARY KEY,
    post_id   INT NOT NULL,
    image_url TEXT NOT NULL,
    CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
-- +goose StatementEnd
