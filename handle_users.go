package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/barturba/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      params.Name,
		})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if len(authorization) == 0 {
		respondWithError(w, http.StatusInternalServerError, "No Apikey provided")
		return
	}
	authorization = strings.TrimPrefix(authorization, "ApiKey ")

	user, err := cfg.DB.GetUserByAPIKey(r.Context(), authorization)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting user")
		return

	}
	respondWithJSON(w, http.StatusOK, user)
}
