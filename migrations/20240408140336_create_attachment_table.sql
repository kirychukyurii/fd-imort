-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.attachment
(
    id           serial,
    name         varchar not null,
    content_type varchar not null,
    file_size    bigint  not null,
    url          varchar not null,
    created_at   timestamp default now(),
    updated_at   timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.attachment;
-- +goose StatementEnd
