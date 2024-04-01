-- +goose Up
-- +goose StatementBegin
CREATE TABLE policies (
    id serial primary key,
    endpoint text not null,
    allowed_roles role[]
);

INSERT INTO policies(endpoint, allowed_roles) VALUES
    ('/chat_v1.ChatV1/Create', ARRAY ['ADMIN']::role[]),
    ('/chat_v1.ChatV1/Delete', ARRAY ['ADMIN']::role[]),
    ('/chat_v1.ChatV1/SendMessage', ARRAY ['ADMIN', 'USER']::role[]),
    ('/chat_v1.ChatV1/Connect', ARRAY ['ADMIN', 'USER']::role[]);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS policies;
-- +goose StatementEnd
