package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/barturba/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleFeeds(w http.ResponseWriter, r *http.Request, u database.User) {

	// Decode
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      params.Name,
			Url:       params.URL,
			UserID:    u.ID,
		})

	// Respond
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}
	respondWithJSON(w, http.StatusOK, feed)
}

// func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {
// 	authorization := r.Header.Get("Authorization")
// 	if len(authorization) == 0 {
// 		respondWithError(w, http.StatusInternalServerError, "No Apikey provided")
// 		return
// 	}
// 	authorization = strings.TrimPrefix(authorization, "ApiKey ")

// 	user, err := cfg.DB.GetUser(r.Context(), authorization)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Error getting user")
// 		return

// 	}
// 	respondWithJSON(w, http.StatusOK, user)
// }
