-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(100) NOT NULL UNIQUE,
    username      VARCHAR(100) NOT NULL UNIQUE,
    password      VARCHAR(255) NOT NULL,
    firstname     VARCHAR(100) NOT NULL,
    lastname      VARCHAR(100) NOT NULL,
    avatar_url    VARCHAR(255),
    date_of_birth DATE,
    bio           TEXT,
    verified      BOOLEAN   DEFAULT FALSE,
    location      VARCHAR(100),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
