package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/barturba/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func runWorker(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	ticker := time.NewTicker(timeBetweenRequest)
	ctx := context.Background()
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(ctx, int32(concurrency))
		if err != nil {
			log.Printf("error when getting feeds from database: %s\n", err)
			return
		}
		log.Printf("found %v feeds from fetch!", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}

}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("couldn't mark feed %s as fetched: %v", feed.Name, err)
	}

	feedData, err := fetchRSS(feed.Url)
	if err != nil {
		log.Printf("couldn't fetch feed %s as fetched: %v", feed.Name, err)
	}
	for _, item := range feedData.Channel.Item {
		publishedAt, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			log.Printf("error parsing published time %s: %v", item.PubDate, err)
			publishedAt = time.Now()
			log.Printf("set published_at to %s", publishedAt.String())
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			pqErr := err.(*pq.Error)
			// If the error is something other than a duplicate warning, then print the error to the log.
			if pqErr.Code.Class() != "23" {
				log.Printf("couldn't save feed %s: %v", feed.Name, err)
			}
		} else {
			log.Println("created post", item.Title)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
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
