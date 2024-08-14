package main

import (
	"net/http"

	"github.com/barturba/blog-aggregator/internal/auth"
	"github.com/barturba/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), authorization)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error getting user")
			return
		}

		handler(w, r, user)
	}

}
