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
	fmt.Printf("starting main worker. maxFeeds %v, workerDelay %v\n", maxFeeds, workerDelay)

	ctx := context.Background()
	var feedsToFetch []database.Feed
	var err error

	for {

		time.Sleep(1 * time.Second)
		feedsToFetch, err = c.DB.GetNextFeedsToFetch(ctx, int32(maxFeeds))
		if err != nil {
			log.Printf("error when getting feeds from database: %s\n", err)
			return
		}
		for _, feed := range feedsToFetch {
			fmt.Printf("will fetch data for feed %s, last fetched: %v\n", feed.Name, feed.LastFetchedAt)
		}

		for _, feed := range feedsToFetch {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Printf("starting data fetching worker for feed %s\n", feed.ID.String())
				data, err := c.Client.FetchRSS(feed.Url)
				if err != nil {
					log.Fatal(fmt.Printf("error fetching data %v\n", err))
				}
				fmt.Printf("stopped data fetching worker for feed %s\n", feed.ID.String())

				fmt.Printf("starting data processing worker for feed %s\n", feed.ID.String())
				c.Client.ProcessRSS(data)
				fmt.Printf("stopped data processing worker for feed %s\n", feed.ID.String())
			}()
		}

		fmt.Printf("sleeping\n")
		time.Sleep(workerDelay)
	}
}
