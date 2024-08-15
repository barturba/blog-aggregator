package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/barturba/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			FeedID:    params.FeedID,
			UserID:    u.ID,
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}
