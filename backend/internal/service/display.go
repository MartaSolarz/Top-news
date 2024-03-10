package service

import "top-news/backend/internal/models"

type DisplayNewsService struct {
	newsDB NewsDatabase
}

func NewDisplayNewsService(newsDB NewsDatabase) *DisplayNewsService {
	return &DisplayNewsService{
		newsDB: newsDB,
	}
}

func (s *DisplayNewsService) GetNews() ([]*models.Article, error) {
	return s.newsDB.GetNews()
}
