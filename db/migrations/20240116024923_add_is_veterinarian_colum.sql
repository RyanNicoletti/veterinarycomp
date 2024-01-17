-- +goose Up
-- +goose StatementBegin
ALTER TABLE compensations ADD COLUMN is_veterinarian BOOLEAN DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compensations DROP COLUMN is_veterinarian;
-- +goose StatementEnd
