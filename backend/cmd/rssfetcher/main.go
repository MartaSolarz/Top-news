package main

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/valyala/fasthttp"
	"log"
	"time"
	"top-news/backend/internal/configuration"
	"top-news/backend/pkg/rss"
)

var configPath = "backend/internal/configuration/configuration.toml"

func sendPostRequest(url string, jsonData string) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.SetBodyString(jsonData)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	log.Printf("sending", jsonData)

	if err := fasthttp.Do(req, resp); err != nil {
		log.Printf("Error sending request: %v\n", err)
		return
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		log.Printf("Error response from server, status code: %d\n", resp.StatusCode())
		return
	}
}

func main() {
	configs := configuration.NewConfiguration(configPath)
	rssURL := configs.News.BBCNewsRSSURL
	duration := configs.News.FetchInterval
	maxRetries := configs.News.MaxRetries
	apiURL := fmt.Sprintf("https://%s:%d/api/process", configs.Server.Host, configs.Server.Port)

	for {
		fmt.Println("Fetching RSS feed...")
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

		data, err := json.Marshal(feeds)
		if err != nil {
			log.Printf("Error in marshaling articles: %v\n", err)
			time.Sleep(duration)
			continue
		}

		sendPostRequest(apiURL, string(data))
		fmt.Println("Articles sent to server, waiting for the next fetch...")

		time.Sleep(duration)
	}
}

//func parseFeedToArticle(feed *gofeed.Feed) []*models.Article {
//	articles := []*models.Article{}
//	for _, item := range feed.Items {
//		date := parser.ParseDate(item.Published)
//		content := parser.ExtractContent(item.Link)
//		thumbnail := parser.ParseThumbnail(item)
//		articles = append(articles, models.NewArticle(item, date, content, thumbnail))
//	}
//	return articles
//}
//
//type Summary struct {
//	Summary string `json:"summary"`
//}
//
//func processArticles(articlesChan <-chan *models.Article, doneChan chan<- *models.Article) {
//	for article := range articlesChan {
//		cmd := exec.Command("venv/bin/python", "python_scripts/ai/summarize.py")
//		cmd.Stdin = strings.NewReader(article.Content)
//		output, err := cmd.CombinedOutput()
//		if err != nil {
//			log.Printf("Failed to summarize article: %v", err)
//			continue
//		}
//
//		var summary Summary
//		if err = json.Unmarshal(output, &summary); err != nil {
//			log.Printf("Failed to unmarshal summary: %v", err)
//			continue
//		}
//		article.Summary = summary.Summary
//		doneChan <- article
//	}
//}
//
//func sendArticlesToChan(articlesChan chan<- *models.Article, items []*gofeed.Item, workerCount int) {
//	var wg sync.WaitGroup
//	itemChan := make(chan *gofeed.Item, len(items))
//
//	for i := 0; i < workerCount; i++ {
//		wg.Add(1)
//		go func() {
//			for item := range itemChan {
//				date := parser.ParseDate(item.Published)
//				content := parser.ExtractContent(item.Link)
//				thumbnail := parser.ParseThumbnail(item)
//				articlesChan <- models.NewArticle(item, date, content, thumbnail)
//			}
//			wg.Done()
//		}()
//	}
//
//	for _, item := range items {
//		itemChan <- item
//	}
//	close(itemChan)
//
//	wg.Wait()
//	close(articlesChan)
//}
