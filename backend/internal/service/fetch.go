package service

import (
	"fmt"
	"log"
	"top-news/backend/internal/models"
)

type FetchService struct {
	newsDB NewsDatabase
}

func NewFetchService(newsDB NewsDatabase) *FetchService {
	return &FetchService{
		newsDB: newsDB,
	}
}

func (s *FetchService) FetchFavorites(ids [][]byte) ([]*models.Article, error) {
	favoriteArticles, err := s.newsDB.GetFavorites(ids)
	if err != nil {
		return nil, fmt.Errorf("could not get favorites: %w", err)
	}
	log.Printf("fetched %d favorites", len(favoriteArticles))

	return favoriteArticles, nil
}
