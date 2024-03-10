package rss

import (
	"sync"

	"github.com/mmcdole/gofeed"

	"top-news/backend/internal/models"
)

// FetchNews FetchRSS fetches the RSS feed from the given URL
func FetchNews(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

// ConvertRssResponseToFeedItems converts the RSS response to models.FeedItems
func ConvertRssResponseToFeedItems(feed *gofeed.Feed) []*models.FeedItem {
	var wg sync.WaitGroup

	feedItems := make([]*models.FeedItem, 0, len(feed.Items))
	feedItemsChan := make(chan *models.FeedItem, len(feed.Items))

	title, copyright := feed.Title, feed.Copyright

	for _, item := range feed.Items {
		wg.Add(1)
		go func(item *gofeed.Item) {
			defer wg.Done()
			newFeed := models.NewFeedItem(item, title, copyright)
			feedItemsChan <- newFeed
		}(item)
	}

	go func() {
		wg.Wait()
		close(feedItemsChan)
	}()

	for feedItem := range feedItemsChan {
		feedItems = append(feedItems, feedItem)
	}

	return feedItems
}
