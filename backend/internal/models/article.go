package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
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
	Authors     string
	Thumbnail   Thumbnail
	ExpireAt    time.Time
	MyDomain    string
}

func ConvertFeedToArticle(feed *FeedItem, ttl int) *Article {
	return &Article{
		Website:     feed.Website,
		CopyRight:   feed.CopyRight,
		Title:       feed.Title,
		Description: feed.Description,
		PublishDate: feed.PublishDate,
		SourceURL:   feed.SourceURL,
		Authors:     parseAuthors(feed.Authors),
		Thumbnail:   feed.Thumbnail,
		ExpireAt:    time.Now().AddDate(0, ttl, 0),
	}
}

func parseAuthors(authors []*gofeed.Person) string {
	result := make([]string, 0, len(authors))
	for _, author := range authors {
		result = append(result, author.Name)
	}
	return strings.Join(result, ", ")
}

func NewArticleFromDB(
	id, thbWidth, thbHeight int,
	website,
	copyRight,
	title,
	description,
	summary,
	sourceUrl,
	authors,
	thbUrl string,
	publishDate, expireAt time.Time,
) *Article {
	return &Article{
		ID:          id,
		Website:     website,
		CopyRight:   copyRight,
		Title:       title,
		Description: description,
		Summary:     summary,
		PublishDate: publishDate.Format("2006-01-02 15:04:05"),
		SourceURL:   sourceUrl,
		Thumbnail: Thumbnail{
			URL:    thbUrl,
			Width:  thbWidth,
			Height: thbHeight,
		},
		Authors:  authors,
		ExpireAt: expireAt,
	}
}

type Thumbnail struct {
	URL    string
	Width  int
	Height int
}

func getThumbnail(item *gofeed.Item) Thumbnail {
	media, ok := item.Extensions["media"]
	if !ok {
		return Thumbnail{}
	}

	if thumbnails, ok := media["thumbnail"]; ok && len(thumbnails) > 0 {
		th := thumbnails[0]

		url, urlOk := th.Attrs["url"]
		width, widthOk := th.Attrs["width"]
		height, heightOk := th.Attrs["height"]

		if urlOk && widthOk && heightOk {
			widthInt, err := strconv.Atoi(width)
			if err != nil {
				widthInt = 0
			}
			heightInt, err := strconv.Atoi(height)
			if err != nil {
				heightInt = 0
			}

			thb := Thumbnail{
				URL:    url,
				Width:  widthInt,
				Height: heightInt,
			}
			return thb
		}
	}
	return Thumbnail{}
}
