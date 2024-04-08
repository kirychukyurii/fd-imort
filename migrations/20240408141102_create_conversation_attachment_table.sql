-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.conversation_attachment
(
    id              serial,
    conversation_id bigint not null,
    attachment_id   bigint not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.conversation_attachment;
-- +goose StatementEnd
