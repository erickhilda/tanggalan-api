
-- name: InsertEvent :exec
INSERT INTO events (title, date, is_national_holiday)
VALUES (?, ?, ?)
ON CONFLICT(title, date) DO NOTHING;

-- name: GetEventsByMonthAndYear :many
SELECT * FROM events
WHERE strftime('%Y', date) = ? AND strftime('%m', date) = ?
ORDER BY date ASC;

-- name: GetAllEvents :many
SELECT * FROM events ORDER BY date ASC;

