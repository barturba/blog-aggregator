
-- +goose Up
CREATE TABLE posts(
    id uuid PRIMARY KEY,
    created_at timestamp not null,
    updated_at timestamp not null,
    title text not null,
    url text unique not null,
    description text unique,
    published_at timestamp not null,
    feed_id uuid not null,
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;