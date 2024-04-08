-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.ticket_status
(
    id    serial,
    name  varchar not null,
    value int     not null
);

INSERT INTO fresh.ticket_status (name, value)
VALUES ('Open', 2)
     , ('Pending', 3)
     , ('Resolved', 4)
     , ('Closed', 5);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.ticket_status;
-- +goose StatementEnd
