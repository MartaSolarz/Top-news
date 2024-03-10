package models

import (
	"github.com/mmcdole/gofeed"

	"top-news/backend/internal/parser"
)

type FeedItem struct {
	Website     string           `json:"website"`
	CopyRight   string           `json:"copyRight"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	PublishDate string           `json:"publishDate"`
	SourceURL   string           `json:"sourceURL"`
	Authors     []*gofeed.Person `json:"authors"`
	Thumbnail   Thumbnail        `json:"thumbnail"`
}

func NewFeedItem(item *gofeed.Item, title, copyright string) *FeedItem {
	return &FeedItem{
		Website:     title,
		CopyRight:   copyright,
		Title:       item.Title,
		Description: item.Description,
		PublishDate: parser.ParseDate(item.Published),
		SourceURL:   item.Link,
		Authors:     item.Authors,
		Thumbnail:   getThumbnail(item),
	}
}
