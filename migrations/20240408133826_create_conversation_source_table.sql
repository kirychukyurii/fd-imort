-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.conversation_source
(
    id    serial,
    name  varchar not null,
    value int     not null
);

INSERT INTO fresh.conversation_source (name, value)
VALUES ('Reply', 0)
     , ('Note', 2)
     , ('Created from tweets', 5)
     , ('Created from survey feedback', 6)
     , ('Created from Facebook post	', 7);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.conversation_source;
-- +goose StatementEnd
