-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.conversation
(
    row_id                 serial,
    id                     bigint,
    parent_id              bigint,
    body                   text,
    body_text              text,
    incoming               boolean,
    to_emails              varchar[],
    category               bigint,
    from_email             varchar,
    cc_emails              varchar[],
    bcc_emails             varchar[],
    private                boolean,
    source                 bigint,
    source_additional_info text,
    support_email          varchar,
    ticket_id              bigint,
    cloud_files            jsonb,
    association_type       bigint,
    email_failure_count    bigint,
    thread_id              bigint,
    thread_message_id      bigint,
    auto_response          boolean,
    automation_id          bigint,
    automation_type_id     bigint,
    outgoing_failures      jsonb,
    user_id                bigint,
    last_edited_at         timestamp,
    last_edited_user_id    bigint,
    attachment_ids         bigint[],
    created_at             timestamp,
    updated_at             timestamp,
    imported_at            timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.conversation;
-- +goose StatementEnd
