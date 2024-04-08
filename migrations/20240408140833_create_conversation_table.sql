-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.conversation
(
    id                  serial,
    parent_id           bigint,
    body                text,
    body_text           text,
    incoming            boolean,
    to_emails           varchar[],
    private             boolean,
    source              bigint,
    support_email       varchar,
    ticket_id           bigint,
    user_id             bigint,
    last_edited_at      timestamp,
    last_edited_user_id bigint,
    created_at          timestamp,
    updated_at          timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.conversation;
-- +goose StatementEnd
