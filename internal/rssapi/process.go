package rssapi

import (
	"fmt"
)

func (c *Client) ProcessRSS(data Rss) {
	for _, item := range data.Channel.Item {
		fmt.Printf("got item: %v\n", item.Title)
	}
}
