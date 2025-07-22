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
WHERE account_from = $1 OR account_to = $1
LIMIT $2
OFFSET $3;
