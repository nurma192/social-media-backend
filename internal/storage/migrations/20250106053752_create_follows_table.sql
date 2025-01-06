-- +goose Up
-- +goose StatementBegin
CREATE TABLE follows
(
    follower_id  INT REFERENCES users (id) ON DELETE CASCADE,
    following_id INT REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (follower_id, following_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS follows;
-- +goose StatementEnd
