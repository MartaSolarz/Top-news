package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mmcdole/gofeed"

	"top-news/backend/internal/configuration"
	"top-news/backend/internal/utils"
	"top-news/backend/pkg/rss"
)

var configPath = "backend/internal/configuration/configuration.toml"

func main() {
	configs := configuration.NewConfiguration(configPath)
	rssURL := configs.News.BBCNewsRSSURL
	duration := configs.News.FetchInterval
	maxRetries := configs.News.MaxRetries
	apiURL := fmt.Sprintf("http://%s:%d/api/process", configs.Server.Host, configs.Server.Port)

	for {
		log.Println("Fetching RSS feed...")
		var feeds *gofeed.Feed
		retries := 0
		for retries < maxRetries {
			var err error
			feeds, err = rss.FetchNews(rssURL)
			if err == nil {
				break
			}
			log.Printf("Attempt %d: Error in fetching RSS feed: %v\n", retries+1, err)
			time.Sleep(time.Duration(retries+1) * time.Minute)
			retries++
		}
		if retries == maxRetries {
			log.Println("Max retries reached. Skipping to next fetch cycle.")
			time.Sleep(duration)
			continue
		}

		rssResponses := rss.ConvertRssResponseToFeedItems(feeds)

		jsonData, err := json.Marshal(rssResponses)
		if err != nil {
			log.Printf("Error in marshaling articles: %v\n", err)
			time.Sleep(duration)
			continue
		}

		utils.SendPostRequest(apiURL, jsonData)
		log.Println("Articles sent to server, waiting for the next fetch...")

		time.Sleep(duration)
	}
}
