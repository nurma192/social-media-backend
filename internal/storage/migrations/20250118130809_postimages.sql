-- +goose Up
-- +goose StatementBegin
CREATE TABLE postImages
(
    id       SERIAL PRIMARY KEY,
    postId   INT NOT NULL,
    imageURL TEXT NOT NULL,
    CONSTRAINT fk_post FOREIGN KEY (postId) REFERENCES posts (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS postImages;
-- +goose StatementEnd