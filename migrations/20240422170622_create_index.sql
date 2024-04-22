-- +goose Up
-- +goose StatementBegin
CREATE INDEX ticket_raw_aws_key_idx ON fresh.ticket_raw USING btree (domain_id, aws_key);
CREATE INDEX ticket_created_at_idx ON fresh.ticket USING btree (domain_id, created_at DESC);
CREATE INDEX ticket_id_idx ON fresh.ticket USING btree (domain_id, id);
CREATE INDEX conversation_id_idx ON fresh.conversation USING btree (domain_id, id);
CREATE INDEX conversation_ticket_id_idx ON fresh.conversation USING btree (domain_id, ticket_id);
CREATE INDEX attachment_id_idx ON fresh.attachment USING btree (domain_id, id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX fresh.attachment_id_idx;
DROP INDEX fresh.conversation_ticket_id_idx;
DROP INDEX fresh.conversation_id_idx;
DROP INDEX fresh.ticket_id_idx;
DROP INDEX fresh.ticket_created_at_idx;
DROP INDEX fresh.ticket_raw_aws_key_idx;
-- +goose StatementEnd
