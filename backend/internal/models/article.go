package models

import (
	"github.com/mmcdole/gofeed"
	"top-news/backend/internal/utils"
)

type Article struct {
	ID          int
	Website     string
	CopyRight   string
	Title       string
	Description string
	Summary     string
	PublishDate string
	SourceURL   string
	Content     string
	Thumbnail   Thumbnail
}

func NewArticle(item *gofeed.Item, date, content string, thumbnail Thumbnail) *Article {
	return &Article{
		ID:          utils.GenerateUniqueID(),
		Title:       item.Title,
		Description: item.Description,
		PublishDate: date,
		SourceURL:   item.Link,
		Content:     content,
		Thumbnail:   thumbnail,
	}
}

type Thumbnail struct {
	URL    string
	Width  string
	Height string
}
