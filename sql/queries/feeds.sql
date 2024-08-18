-- name: CreateFeed :one
INSERT into feeds (id, created_at, updated_at, name, url, user_id) 
values ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
where last_fetched_at is null order by last_fetched_at asc
LIMIT $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(), 
updated_at =  NOW()
WHERE id = $1
RETURNING *;