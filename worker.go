package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/barturba/blog-aggregator/internal/database"
)

func (c *apiConfig) runWorker(wg *sync.WaitGroup, maxFeeds int, workerDelay time.Duration) {
	// this worker is run within its own goroutine
	// 1. fetch n entries from the database
	// 2. perform data fetching and processing concurrently for each feed in the batch
	//   a. fetch data
	//   b. process data
	// 3. sleep for a period, workerDelay
	fmt.Printf("starting main worker\n")
	fmt.Printf("maxFeeds %v\n", maxFeeds)
	fmt.Printf("workerDelay %v\n", workerDelay)

	ctx := context.Background()
	var feedsToFetch []database.Feed
	var err error

	for {

		fmt.Printf("fetching data from database\n")
		time.Sleep(1 * time.Second)
		// TODO: GetNextFeedsToFetch should accept input for n
		feedsToFetch, err = c.DB.GetNextFeedsToFetch(ctx, int32(maxFeeds))
		if err != nil {
			log.Printf("error when getting feeds from database: %s\n", err)
			return
		}
		fmt.Printf("done fetching data from database\n")
		fmt.Printf("len(feedsToFetch): %v\n", len(feedsToFetch))
		for _, feed := range feedsToFetch {
			fmt.Printf("will fetch data for feed %s, last fetched: %v\n", feed.Name, feed.LastFetchedAt)
		}

		// for each feed fetch and process it using a goroutine

		for _, feed := range feedsToFetch {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Printf("starting data fetching worker for feed %s\n", feed.ID.String())
				data, err := c.Client.FetchRSS("https://blog.boot.dev/index.xml")
				if err != nil {
					log.Fatal(fmt.Printf("error fetching data %v\n", err))
				}
				fmt.Printf("stopped data fetching worker for feed %s\n", feed.ID.String())

				fmt.Printf("starting data processing worker for feed %s\n", feed.ID.String())
				time.Sleep(time.Second * 10)
				c.Client.ProcessRSS(data)
				fmt.Printf("stopped data processing worker for feed %s\n", feed.ID.String())
			}()
		}

		fmt.Printf("sleeping\n")
		time.Sleep(workerDelay)
	}
}
