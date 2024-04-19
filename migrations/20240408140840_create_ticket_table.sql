-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.ticket
(
    row_id                 serial,
    aws_key                varchar,
    id                     bigint,
    requester_name         varchar,
    parent_id              bigint,
    archived               boolean,
    meta                   jsonb,
    name                   varchar,
    cc_emails              varchar[],
    ticket_cc_emails       varchar[],
    company_id             bigint,
    custom_fields          jsonb,
    deleted                boolean,
    description            text,
    description_text       text,
    due_by                 timestamp,
    email                  varchar,
    email_config_id        bigint,
    facebook_id            varchar,
    fr_due_by              timestamp,
    fr_escalated           boolean,
    nr_due_by              timestamp,
    nr_escalated           boolean,
    fwd_emails             varchar[],
    group_id               bigint,
    is_escalated           boolean,
    phone                  varchar,
    priority               bigint,
    product_id             bigint,
    reply_cc_emails        varchar[],
    requester_id           bigint,
    responder_id           bigint,
    source                 bigint,
    spam                   boolean,
    status                 bigint,
    subject                varchar,
    tags                   varchar[],
    to_emails              varchar[],
    twitter_id             varchar,
    type                   varchar,
    association_type       bigint,
    source_additional_info text,
    support_email          varchar,
    form_id                bigint,
    attachment_ids         bigint[],
    created_at             timestamp,
    updated_at             timestamp,
    imported_at            timestamp default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.ticket;
-- +goose StatementEnd
