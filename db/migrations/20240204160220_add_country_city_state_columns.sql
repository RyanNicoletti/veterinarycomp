-- +goose Up
-- +goose StatementBegin
ALTER TABLE compensations
    ADD COLUMN country VARCHAR(255) DEFAULT '',
    ADD COLUMN state VARCHAR(255) DEFAULT '',
    ADD COLUMN city VARCHAR(255) DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compensations
    DROP COLUMN country,
    DROP COLUMN state,
    DROP COLUMN city;
-- +goose StatementEnd
