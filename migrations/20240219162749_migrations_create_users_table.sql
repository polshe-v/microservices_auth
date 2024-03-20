-- +goose Up
-- +goose StatementBegin
CREATE TYPE role AS ENUM ('UNKNOWN', 'USER', 'ADMIN');
CREATE TABLE users (
    id serial primary key,
    name text not null unique,
    role role not null,
    email text not null unique,
    password text not null,
    created_at timestamp not null default now(),
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS role;
-- +goose StatementEnd
