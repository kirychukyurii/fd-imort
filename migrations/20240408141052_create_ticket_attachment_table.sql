-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.ticket_attachment
(
    id            serial,
    ticket_id     bigint not null,
    attachment_id bigint not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.ticket_attachment;
-- +goose StatementEnd
