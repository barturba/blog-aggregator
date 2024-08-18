package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/barturba/blog-aggregator/internal/database"
)

func (cfg *apiConfig) getPosts(w http.ResponseWriter, r *http.Request, u database.User) {
	limit := 10

	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("error parsing limit param")
			limit = 10
		}
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(),
		database.GetPostsByUserParams{
			UserID: u.ID,
			Limit:  int32(limit)})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts by user")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseGetPostsByUserRowToPosts(posts))
}
