package rssapi

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

func (c *Client) FetchRSS(url string) (Rss, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Rss{}, err
	}

	resp, err := c.httpClient.Do(req)
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
