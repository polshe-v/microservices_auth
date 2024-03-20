-- +goose Up
-- +goose StatementBegin
CREATE TABLE hmac_keys (
    id serial primary key,
    timestamp timestamp not null default now(),
    key text not null,
    value text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS hmac_keys;
-- +goose StatementEnd
