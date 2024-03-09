package service

import "top-news/backend/internal/models"

type NewsDatabase interface {
	GetNews() ([]*models.Article, error)
	PutNews(articles []*models.Article) error
}

type DisplayNewsService struct {
	newsDB NewsDatabase
}

func NewDisplayNewsService(newsDB NewsDatabase) *DisplayNewsService {
	return &DisplayNewsService{
		newsDB: newsDB,
	}
}
