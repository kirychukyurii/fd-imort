-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.ticket_raw
(
    row_id       serial,
    aws_key      varchar,
    requester_id bigint,
    ticket_id    bigint,
    ticket       jsonb,
    imported_at  timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.ticket_raw;
-- +goose StatementEnd
