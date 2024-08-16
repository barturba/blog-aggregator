-- name: CreateFeedFollow :one
INSERT into feed_follows(id, created_at, updated_at, feed_id, user_id) 
values ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollows :many
SELECT * FROM feed_follows;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows
WHERE id = $1;