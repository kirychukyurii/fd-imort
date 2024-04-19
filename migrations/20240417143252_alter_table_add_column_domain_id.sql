-- +goose Up
-- +goose StatementBegin
ALTER TABLE fresh.ticket
    ADD COLUMN domain_id bigint;

ALTER TABLE fresh.ticket_raw
    ADD COLUMN domain_id bigint;

ALTER TABLE fresh.conversation
    ADD COLUMN domain_id bigint;

ALTER TABLE fresh.attachment
    ADD COLUMN domain_id bigint;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE fresh.ticket
    DROP COLUMN domain_id;

ALTER TABLE fresh.ticket_raw
    DROP COLUMN domain_id;

ALTER TABLE fresh.conversation
    DROP COLUMN domain_id;

ALTER TABLE fresh.attachment
    DROP COLUMN domain_id;
-- +goose StatementEnd
