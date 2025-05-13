-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE events (
  id                  INTEGER PRIMARY KEY AUTOINCREMENT,
  title               TEXT    NOT NULL,
  date                TEXT    NOT NULL,
  is_national_holiday BOOLEAN NOT NULL DEFAULT FALSE,
  UNIQUE(title, date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE events;
-- +goose StatementEnd
