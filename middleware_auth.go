package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/barturba/blog-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if len(authorization) == 0 {
			respondWithError(w, http.StatusInternalServerError, "No Apikey provided")
			return
		}
		fmt.Printf("middlewareAuth: got apikey: %s\n", authorization)
		authorization = strings.TrimPrefix(authorization, "ApiKey ")

		user, err := cfg.DB.GetUser(r.Context(), authorization)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error getting user")
			return

		}

		handler(w, r, user)
	}

}
