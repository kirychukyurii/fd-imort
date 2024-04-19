-- +goose Up
-- +goose StatementBegin
CREATE TABLE fresh.domain
(
    id   serial,
    name varchar
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fresh.domain;
-- +goose StatementEnd
