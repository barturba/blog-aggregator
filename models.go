package main

import (
	"time"

	"github.com/barturba/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Apikey    string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Apikey:    user.Apikey,
	}
}

type Feed struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name"`
	Url           string    `json:"url"`
	UserID        uuid.UUID `json:"user_id"`
	LastFetchedAt time.Time `json:"last_fetched_at"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: feed.LastFetchedAt.Time,
	}
}

func databaseFeedsToFeeds(feeds []database.Feed) []Feed {
	var items []Feed
	for _, feed := range feeds {
		items = append(items, databaseFeedToFeed(feed))
	}
	return items
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedFollowToFeedFollow(feedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
	}
}

func databaseFeedFollowsToFeedFollows(feedFollows []database.FeedFollow) []FeedFollow {
	var items []FeedFollow
	for _, feedFollow := range feedFollows {
		items = append(items, databaseFeedFollowToFeedFollow(feedFollow))
	}
	return items
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databaseGetPostByUserRowToPost(row database.GetPostsByUserRow) Post {
	return Post{
		ID:          row.ID,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Title:       row.Title,
		Url:         row.Url,
		Description: row.Description.String,
		PublishedAt: row.PublishedAt,
		FeedID:      row.FeedID,
	}
}

func databaseGetPostsByUserRowToPosts(posts []database.GetPostsByUserRow) []Post {
	var items []Post
	for _, post := range posts {
		items = append(items, databaseGetPostByUserRowToPost(post))
	}
	return items
}
