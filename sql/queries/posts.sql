-- name: CreatePost :one

INSERT into posts(id, created_at, updated_at, title, url, description, published_at, feed_id ) 
values ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- -- name: GetFeedFollows :many
-- SELECT * FROM feed_follows where user_id = $1;

-- -- name: DeleteFeedFollows :exec
-- DELETE FROM feed_follows
-- WHERE id = $1;