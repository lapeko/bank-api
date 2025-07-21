-- name: CreateTransfer :one
INSERT INTO transfers (account_from, account_to, amount)
VALUES ($1, $2, $3)
RETURNING *;