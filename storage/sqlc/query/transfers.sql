-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id LIMIT $1 OFFSET $2;

-- name: ListTransfersBySender :many
SELECT * FROM transfers
WHERE account_from = $1
ORDER BY id LIMIT $2 OFFSET $3;

-- name: ListTransfersByReceiver :many
SELECT * FROM transfers
WHERE account_to = $1
ORDER BY id LIMIT $2 OFFSET $3;

-- name: ListTransfersBySenderAndReceiver :many
SELECT * FROM transfers
WHERE account_to = $1 AND account_from = $2
ORDER BY id LIMIT $3 OFFSET $4;

-- name: CreateTransfer :one
INSERT INTO transfers (account_from, account_to, amount)
VALUES ($1, $2, $3)
RETURNING *;
