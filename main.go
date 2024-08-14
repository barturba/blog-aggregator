package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/barturba/blog-aggregator/internal/database"
	"github.com/google/uuid"
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
	// check to make suer each of these env vars are set
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	databaseURL := os.Getenv("DATABASE_URL")
	fmt.Printf("PORT: %v\n", port)
	fmt.Printf("DATABASE_URL: %v\n", databaseURL)

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	cfg := apiConfig{
		DB: dbQueries,
	}
	fmt.Printf("dbQueries: %v\n", dbQueries)

	m := http.NewServeMux()

	m.HandleFunc("POST /v1/users", cfg.handleUsers)

	srv := http.Server{
		Handler:      m,
		Addr:         port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	fmt.Println("server started on ", port)
	err = srv.ListenAndServe()
	log.Fatal(err)

	fmt.Printf("the blog-aggregator has started\n")
}

func (cfg *apiConfig) handleUsers(w http.ResponseWriter, r *http.Request) {

	// Decode
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a new user
	user, err := cfg.DB.CreateUser(r.Context(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      params.Name,
		})
	if err != nil {
		log.Printf("Couldn't create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}
