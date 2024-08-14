package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load .env file")
	}
	fmt.Printf("PORT: %v\n", os.Getenv("PORT"))
	fmt.Printf("the blog-aggregator has started\n")
}
