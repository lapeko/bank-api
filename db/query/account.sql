-- name: CreateAccount :one
INSERT INTO accounts (user_id, currency, balance)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccountById :one
SELECT *
FROM accounts
WHERE id = $1;

-- name: GetAccountsByIdForUpdate :many
SELECT *
FROM accounts
WHERE id IN ($1, $2)
FOR UPDATE;

-- name: ListAccounts :many
SELECT *
FROM accounts
LIMIT $1
OFFSET $2;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: OffsetBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(delta)
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;
