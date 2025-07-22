-- name: CreateEntry :one
INSERT INTO entries (account_id, amount)
VALUES ($1, $2)
RETURNING *;

-- name: ListEntries :many
SELECT *
FROM entries
LIMIT $1
OFFSET $2;

-- name: GetEntryById :one
SELECT *
FROM entries
WHERE id = $1;

-- name: ListEntriesByAccount :many
SELECT *
FROM entries
WHERE account_id = $1
LIMIT $2
OFFSET $3;

-- name: GetTotalEntriesCount :one
SELECT COUNT(*) as total_count
FROM entries;

-- name: GetTotalEntriesCountByAccount :one
SELECT COUNT(*) as total_count
FROM entries
WHERE account_id = $1;
