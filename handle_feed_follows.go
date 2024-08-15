package main

import (
	"encoding/json"
	"net/http"

	"github.com/barturba/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handleFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		FeedID string `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(),
		database.CreateFeedFollowsParams{})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}
