package main

import (
	"net/http"

	"github.com/barturba/blog-aggregator/internal/database"
)

func (cfg *apiConfig) getPosts(w http.ResponseWriter, r *http.Request, u database.User) {

	// posts, err := cfg.DB.Get(r.Context(),
	// 	database.CreateFeedFollowParams{
	// 		ID:        uuid.New(),
	// 		CreatedAt: time.Time{},
	// 		UpdatedAt: time.Time{},
	// 		FeedID:    params.FeedID,
	// 		UserID:    u.ID,
	// 	})
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't create feed follow")
	// 	return
	// }

	// respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}
