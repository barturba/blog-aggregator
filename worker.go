package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func (c *apiConfig) runWorker(wg *sync.WaitGroup, maxFeeds int, workerDelay time.Duration) {
	// this worker is run within its own goroutine
	// 1. fetch n entries from the database
	// 2. perform data fetching and processing concurrently
	//   a. fetch data
	//   b. process data
	// 3. sleep for a period
	fmt.Printf("starting main worker\n")
	fmt.Printf("maxFeeds %v\n", maxFeeds)
	fmt.Printf("workerDelay %v\n", workerDelay)
	ctx := context.Background()
	// data := rssapi.Rss{}
	// var err error

	for {

		fmt.Printf("fetching data from database\n")
		time.Sleep(1)
		// TODO: GetNextFeedsToFetch should accept input for n
		_, err := c.DB.GetNextFeedsToFetch(ctx)
		if err != nil {
			log.Printf("error when getting feeds from database: %s\n", err)
			return
		}
		fmt.Printf("done fetching data from database\n")

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("starting data fetching worker\n")
			time.Sleep(time.Second * 20)
			// data, err = c.Client.FetchRSS("https://blog.boot.dev/index.xml")
			// if err != nil {
			// 	log.Fatal(fmt.Printf("error fetching data %v\n", err))
			// }
			// fmt.Printf("got feeds: %s\n", data)
			fmt.Printf("stopped data fetching worker\n")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("starting data processing worker\n")
			time.Sleep(time.Second * 10)
			// c.Client.ProcessRSS(data)
			// fmt.Printf("got feeds: %s\n", data)
			fmt.Printf("stopped data processing worker\n")
		}()
		fmt.Printf("sleeping")
		time.Sleep(workerDelay)
	}
}
