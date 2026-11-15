-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS password_recovery (
  "user_id"     uuid        PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  "code_hash"   varchar(64) NOT NULL,
  "attempts"    smallint    NOT NULL DEFAULT 0,
  "createat"    timestamp   NOT NULL DEFAULT now(),
  "expires_at"  timestamp   NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS password_recovery;
-- +goose StatementEnd
