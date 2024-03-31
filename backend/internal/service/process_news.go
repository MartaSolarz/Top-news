package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"top-news/backend/internal/models"
	"top-news/backend/internal/parser"
)

type ProcessNewsService struct {
	newsDB         NewsDatabase
	NumWorkers     int
	openAPIKey     string
	openAPIUrl     string
	disableOpenAPI bool
}

func NewProcessNewsService(newsDB NewsDatabase, workers int, url, key string, disableOpenAPI bool) *ProcessNewsService {
	return &ProcessNewsService{
		newsDB:         newsDB,
		NumWorkers:     workers,
		openAPIKey:     key,
		openAPIUrl:     url,
		disableOpenAPI: disableOpenAPI,
	}
}

func (s *ProcessNewsService) ProcessNews(rssResponses []*models.FeedItem) error {
	titleToItem := getTitleToItem(rssResponses)

	newFeeds, err := s.getNewArticles(titleToItem)
	if err != nil {
		log.Printf("Error getting new articles: %v", err)
		newFeeds = rssResponses
	}

	if len(newFeeds) == 0 {
		log.Printf("No new articles found")
		return nil
	}

	log.Printf("Found %d new articles", len(newFeeds))

	articlesChan := make(chan *models.Article, len(newFeeds))
	for _, feed := range newFeeds {
		articlesChan <- models.ConvertFeedToArticle(feed, s.newsDB.GetTTL())
	}

	close(articlesChan)

	articlesWithContentChan := make(chan *models.Article, len(newFeeds))
	s.addContentToArticles(articlesChan, articlesWithContentChan)

	close(articlesWithContentChan)

	doneChan := make(chan *models.Article, len(newFeeds))
	s.summarizeArticles(articlesWithContentChan, doneChan)

	close(doneChan)

	articlesToSaveInDB := make([]*models.Article, 0, len(doneChan))
	for article := range doneChan {
		articlesToSaveInDB = append(articlesToSaveInDB, article)
	}

	err = s.newsDB.PutNews(articlesToSaveInDB)
	if err != nil {
		return fmt.Errorf("error putting news in DB: %w", err)
	}

	return nil
}

func getTitleToItem(rssResponses []*models.FeedItem) map[string]*models.FeedItem {
	titleToItem := make(map[string]*models.FeedItem, len(rssResponses))
	for _, n := range rssResponses {
		titleToItem[n.Title] = n
	}
	return titleToItem
}

func (s *ProcessNewsService) getNewArticles(titleToItem map[string]*models.FeedItem) ([]*models.FeedItem, error) {
	titlesFromDB, err := s.newsDB.GetTitles()
	if err != nil {
		return nil, fmt.Errorf("error getting titles from DB: %w", err)
	}

	newArticles := make([]*models.FeedItem, 0)
	for title, item := range titleToItem {
		if !contains(titlesFromDB, title) {
			newArticles = append(newArticles, item)
		}
	}

	return newArticles, nil
}

func contains(titles []string, title string) bool {
	for _, t := range titles {
		if t == title {
			return true
		}
	}
	return false
}

func (s *ProcessNewsService) addContentToArticles(articlesChan, articlesWithContentChan chan *models.Article) {
	var wg sync.WaitGroup

	for i := 0; i < s.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			for article := range articlesChan {
				article.Content = parser.ExtractContent(article.SourceURL)
				articlesWithContentChan <- article
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

type Summary struct {
	Summary string `json:"summary"`
}

func (s *ProcessNewsService) summarizeArticles(articlesWithContentChan, doneChan chan *models.Article) {
	var wg sync.WaitGroup
	for i := 0; i < s.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			for article := range articlesWithContentChan {
				s.doSummarize(article)
				doneChan <- article
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func (s *ProcessNewsService) doSummarize(article *models.Article) {
	if s.disableOpenAPI {
		article.Summary = article.Content
		return
	}

	cmd := exec.Command("venv/bin/python", "backend/python/ai/main.py",
		article.Content, s.openAPIUrl, s.openAPIKey)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to summarize article: %v", err)
		return
	}

	var summary Summary
	if err = json.Unmarshal(output, &summary); err != nil {
		log.Printf("Failed to unmarshal summary: %v", err)
		return
	}
	article.Summary = summary.Summary
}
