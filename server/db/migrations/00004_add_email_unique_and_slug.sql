-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email);

ALTER TABLE pastes ADD COLUMN IF NOT EXISTS slug TEXT;
UPDATE pastes SET slug = SUBSTRING(url FROM '[^/]+$') WHERE slug IS NULL;
ALTER TABLE pastes ALTER COLUMN slug SET NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_pastes_slug ON pastes(slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_pastes_slug;
ALTER TABLE pastes DROP COLUMN IF EXISTS slug;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_unique;
-- +goose StatementEnd
