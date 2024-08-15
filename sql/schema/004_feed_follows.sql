
-- +goose Up
CREATE TABLE feed_follows(
    id uuid PRIMARY KEY,
    created_at timestamp not null,
    updated_at timestamp not null,
    feed_id uuid not null,
    user_id uuid not null,
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;