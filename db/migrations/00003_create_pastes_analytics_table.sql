-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pastes_analysis(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
paste_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
url TEXT ,
views INTEGER NOT NULL DEFAULT 0,

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pastes_analysis;
-- +goose StatementEnd
