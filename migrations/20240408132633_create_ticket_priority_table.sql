-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.ticket_priority
(
    id    serial,
    name  varchar not null,
    value int     not null
);

INSERT INTO fresh.ticket_priority (name, value)
VALUES ('Low', 1)
     , ('Medium', 2)
     , ('High', 3)
     , ('Urgent', 4);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.ticket_priority;
-- +goose StatementEnd
