
-- name: InsertEvent :exec
INSERT INTO events (title, date, is_national_holiday)
VALUES (?, ?, ?)
ON CONFLICT(title, date) DO NOTHING;

-- name: GetEventsByMonthAndYear :many
SELECT * FROM events
WHERE strftime('%Y', date) = @year
  AND strftime('%m', date) = @month
ORDER BY date ASC;

