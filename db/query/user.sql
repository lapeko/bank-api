-- name: CreateUser :one
INSERT INTO users (full_name, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;