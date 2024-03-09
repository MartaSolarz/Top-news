package rss

import "github.com/mmcdole/gofeed"

// FetchRSS fetches the RSS feed from the given URL
func FetchNews(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return feed, nil
}
