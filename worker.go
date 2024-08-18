package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/barturba/blog-aggregator/internal/database"
)

func runWorker(db *database.Queries, maxFeeds int, workerDelay time.Duration) {
	var wg sync.WaitGroup
	ctx := context.Background()
	var feedsToFetch []database.Feed
	var err error

	for {

		time.Sleep(1 * time.Second)
		feedsToFetch, err = db.GetNextFeedsToFetch(ctx, int32(maxFeeds))
		if err != nil {
			log.Printf("error when getting feeds from database: %s\n", err)
			return
		}

		for _, feed := range feedsToFetch {
			wg.Add(1)
			go func() {
				defer wg.Done()
				data, err := fetchRSS(feed.Url)
				if err != nil {
					log.Fatal(fmt.Printf("error fetching data %v\n", err))
				}

				for _, item := range data.Channel.Item {
					fmt.Printf("found post: %v\n", item.Title)
				}
				log.Printf("Feed %s collected, %v posts found", feed.Name, len(data.Channel.Item))
			}()
		}

		fmt.Printf("sleeping\n")
		time.Sleep(workerDelay)
	}
}

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func fetchRSS(url string) (Rss, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Rss{}, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return Rss{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Rss{}, errors.New("location not found")
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Rss{}, err
	}

	rssResp := Rss{}
	err = xml.Unmarshal(dat, &rssResp)
	if err != nil {
		return Rss{}, err
	}

	return rssResp, nil
}

func processRSS(data Rss) {
	for _, item := range data.Channel.Item {
		fmt.Printf("got item: %v\n", item.Title)
	}
}
