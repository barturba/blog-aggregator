package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/barturba/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/users", apiCfg.handleUsers)
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handleFeeds))
	mux.HandleFunc("GET /v1/users", apiCfg.getUsers)
	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerErr)

	srv := http.Server{
		Handler:      mux,
		Addr:         ":" + port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	fmt.Println("server started on ", port)
	err = srv.ListenAndServe()
	log.Fatal(err)

	fmt.Printf("the blog-aggregator has started\n")
}

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
