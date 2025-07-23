-- name: CreateAccount :one
INSERT INTO accounts (user_id, currency)
VALUES ($1, $2)
RETURNING *;

-- name: GetAccountById :one
SELECT a.id, a.user_id, u.full_name, u.email, a.currency, a.balance, a.created_at
FROM accounts as a
JOIN users as u
ON a.user_id = u.id
WHERE a.id = $1;

-- name: GetAccountByIdForUpdate :one
SELECT *
FROM accounts
WHERE id = $1
FOR UPDATE;

-- name: GetTwoAccountsByIdForUpdate :many
SELECT *
FROM accounts
WHERE id IN ($1, $2)
ORDER BY id
FOR UPDATE;

-- name: ListAccounts :many
SELECT a.id, a.user_id, u.full_name, u.email, a.currency, a.balance, a.created_at
FROM accounts as a
JOIN users as u
ON a.user_id = u.id
LIMIT $1
OFFSET $2;

-- name: GetTotalAccountsCount :one
SELECT COUNT(*) AS total_count FROM accounts;

-- name: OffsetAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(delta)
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :one
DELETE FROM accounts
WHERE id = $1
RETURNING *;
