-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcyrpto";
CREATE TABLE IF NOT EXISTS pastes(
id UUID PRIMARY KEY gen_random_uuid(),
user_id UUID REFERENCES users(id),
title TEXT NOT NULL,
is_private BOOLEAN NOT NULL DEFAULT false,
content TEXT NOT NULL,
password TEXT NOT NULL,
language TEXT NOT NULL,
URL TEXT NOT NULL,
expires_at TIMESTAMPTZ,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pastes;
-- +goose StatementEnd
