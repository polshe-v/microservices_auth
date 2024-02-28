-- +goose Up
-- +goose StatementBegin
CREATE TABLE transaction_log (
    id serial primary key,
    timestamp timestamp not null default now(),
    log text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_log;
-- +goose StatementEnd
