package main

import (
	"fmt"
	"log"
	"os"

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
	databaseURL := os.Getenv("DATABASE_URL")
	fmt.Printf("PORT: %v\n", port)
	fmt.Printf("DATABASE_URL: %v\n", databaseURL)
	fmt.Printf("the blog-aggregator has started\n")
}
