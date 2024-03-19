package service

import (
	"sort"
	"top-news/backend/internal/models"
)

type DisplayNewsService struct {
	newsDB NewsDatabase
}

func NewDisplayNewsService(newsDB NewsDatabase) *DisplayNewsService {
	return &DisplayNewsService{
		newsDB: newsDB,
	}
}

func (s *DisplayNewsService) GetNews() ([]*models.Article, error) {
	articles, err := s.newsDB.GetNews()
	if err != nil {
		return nil, err
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishDate > articles[j].PublishDate
	})

	return articles, nil
}
