-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pastes_analytics(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
pasteID UUID NOT NULL UNIQUE  REFERENCES pastes(id) ON DELETE CASCADE,
url TEXT ,
views INTEGER NOT NULL DEFAULT 0,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pastes_analytics;
-- +goose StatementEnd
