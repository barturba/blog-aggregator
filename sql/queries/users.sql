-- name: CreateUser :one
INSERT into users (id, created_at, updated_at, name) 
values ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByAPIKey :one
SELECT * FROM users WHERE apikey = $1;