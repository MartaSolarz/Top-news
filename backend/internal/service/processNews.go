package service

import "fmt"

type ProcessNewsService struct {
	newsDB NewsDatabase
}

func NewProcessNewsService(newsDB NewsDatabase) *ProcessNewsService {
	return &ProcessNewsService{
		newsDB: newsDB,
	}
}

func (s *ProcessNewsService) ProcessNews(news interface{}) {
	// TODO: check if news already exists in the database
	// TODO: if not, process the news - sent it to AI model to generate summarize and store it in the database
	for _, n := range news.([]interface{}) {
		fmt.Println(n)
	}
}
