-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.ticket_source
(
    id    serial,
    name  varchar not null,
    value int     not null
);

INSERT INTO fresh.ticket_source (name, value)
VALUES ('Email', 1)
     , ('Portal', 2)
     , ('Phone', 3)
     , ('Chat', 7)
     , ('Feedback Widget', 9)
     , ('Outbound Email', 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fresh.ticket_source;
-- +goose StatementEnd
