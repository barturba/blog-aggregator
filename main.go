package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/barturba/blog-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

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
	mux.HandleFunc("GET /v1/feeds", apiCfg.getFeeds)

	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handleFeedFollows))
	mux.HandleFunc("DELETE /v1/feed_follows/", apiCfg.middlewareAuth(apiCfg.deleteFeedFollows))
	mux.HandleFunc("GET /v1/feed_follows/", apiCfg.middlewareAuth(apiCfg.getFeedFollows))

	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.getUsers))
	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerErr)

	srv := http.Server{
		Handler:      mux,
		Addr:         ":" + port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	const collectionConcurrency = 10
	const collectionTimeInterval = time.Minute
	go runWorker(apiCfg.DB, collectionConcurrency, collectionTimeInterval)

	fmt.Println("server started on ", port)
	err = srv.ListenAndServe()
	log.Fatal(err)

	fmt.Printf("the blog-aggregator has started\n")
}
