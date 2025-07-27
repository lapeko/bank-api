-- name: CreateTransfer :one
INSERT INTO transfers (account_from, account_to, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListTransfers :many
SELECT *
FROM transfers
LIMIT $1
OFFSET $2;

-- name: ListTransfersByAccount :many
SELECT *
FROM transfers
WHERE account_from = sqlc.arg(account_id) OR account_to = sqlc.arg(account_id)
LIMIT $1
OFFSET $2;

-- name: GetTotalTransfersCount :one
SELECT COUNT(*) as total_count
FROM transfers;

-- name: GetTotalTransfersCountByAccount :one
SELECT COUNT(*) as total_count
FROM transfers
WHERE account_from = sqlc.arg(account_id) OR account_to = sqlc.arg(account_id);