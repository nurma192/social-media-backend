-- +goose Up
-- +goose StatementBegin
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
DROP TABLE IF EXISTS postImages;
-- +goose StatementEnd