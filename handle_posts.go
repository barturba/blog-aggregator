package main

import (
	"net/http"

	"github.com/barturba/blog-aggregator/internal/database"
)

func (cfg *apiConfig) getPosts(w http.ResponseWriter, r *http.Request, u database.User) {

	posts, err := cfg.DB.GetPostsByUser(r.Context(),
		database.GetPostsByUserParams{
			UserID: u.ID,
			Limit:  10})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get posts by user")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseGetPostsByUserRowToPosts(posts))
}
