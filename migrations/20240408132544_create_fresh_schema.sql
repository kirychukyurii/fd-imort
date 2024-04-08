-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA fresh;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA fresh;
-- +goose StatementEnd
