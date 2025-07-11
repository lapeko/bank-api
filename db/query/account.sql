-- name: CreateAccount :one
INSERT INTO accounts (user_id, currency, balance)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountById :one
SELECT *
FROM accounts
WHERE id = $1;

-- name: ListAccounts :many
SELECT *
FROM accounts
LIMIT $1
OFFSET $2;
