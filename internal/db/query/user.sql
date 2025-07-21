-- name: CreateUser :one
INSERT INTO users (full_name, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT *
FROM users
LIMIT $2
OFFSET $1;

-- name: UpdateUserFullName :one
UPDATE users
SET full_name = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users
SET email = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;