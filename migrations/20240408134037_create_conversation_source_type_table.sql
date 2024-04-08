-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.conversation_source_type
(
    id    serial,
    name  varchar not null,
    value int     not null
);

INSERT INTO fresh.conversation_source_type (name, value)
VALUES ('Created from Forwarded Email', 8)
     , ('Created from Phone', 9)
     , ('E-Commerce', 11);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.conversation_source_type;
-- +goose StatementEnd
