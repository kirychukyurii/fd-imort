-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.attachment
(
    row_id       serial,
    id           bigint,
    name         varchar not null,
    content_type varchar not null,
    file_size    bigint  not null,
    url          varchar not null,
    thumb_url    varchar,
    created_at   timestamp,
    updated_at   timestamp,
    imported_at  timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.attachment;
-- +goose StatementEnd
