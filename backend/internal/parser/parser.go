package parser

import (
	"log"
	"os/exec"
	"time"

	"github.com/mmcdole/gofeed"

	"top-news/backend/internal/models"
)

func ParseDate(date string) string {
	t, err := time.Parse(time.RFC1123, date)
	if err != nil {
		log.Println(err)
	}
	return t.Format("2006-01-02 15:04:05")
}

func ParseThumbnail(item *gofeed.Item) models.Thumbnail {
	media, ok := item.Extensions["media"]
	if !ok {
		return models.Thumbnail{}
	}

	if thumbnails, ok := media["thumbnail"]; ok && len(thumbnails) > 0 {
		th := thumbnails[0]

		url, urlOk := th.Attrs["url"]
		width, widthOk := th.Attrs["width"]
		height, heightOk := th.Attrs["height"]

		if urlOk && widthOk && heightOk {
			thb := models.Thumbnail{
				URL:    url,
				Width:  width,
				Height: height,
			}
			return thb
		}
	}
	return models.Thumbnail{}
}

func ExtractContent(url string) string {
	cmd := exec.Command("venv/bin/python", "python_scripts/extract.py", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to extract content: %v", err)
		return ""
	}
	return string(output)
}
