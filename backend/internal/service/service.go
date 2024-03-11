package service

import "top-news/backend/internal/models"

type NewsDatabase interface {
	GetNews() ([]*models.Article, error)
	GetTitles() ([]string, error)
	PutNews(articles []*models.Article) error
	GetTTL() int
	GetFavorites(ids [][]byte) ([]*models.Article, error)
}
