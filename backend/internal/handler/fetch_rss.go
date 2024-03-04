package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"os/exec"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/valyala/fasthttp"

	"top-news/backend/internal/parser"
)

type Summary struct {
	Summary string `json:"summary"`
}

type Article struct {
	Title       string
	Description string
	Summary     string
	//Authors     []*gofeed.Person
	//Image  *gofeed.Image
}

func (h *Handler) FetchRSSHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[GET] /fetchrss")
	feed, err := parser.FetchRSS(h.NewsURL)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}

	tmpl, err := template.ParseFiles("frontend/articles.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	articlesChan := make(chan Article)
	doneChan := make(chan Article)

	go doSummaries(articlesChan, doneChan)

	go sentArticlesToChan(articlesChan, feed.Items)

	var articles []Article
	for i := 0; i < len(feed.Items); i++ {
		summary := <-doneChan
		articles = append(articles, summary)
	}

	ctx.SetContentType("text/html")
	err = tmpl.Execute(ctx, articles)
	if err != nil {
		log.Println("Error executing template:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
	}
}

func sentArticlesToChan(articlesChan chan Article, items []*gofeed.Item) {
	for _, item := range items {
		article := Article{
			Title:       item.Title,
			Description: item.Description,
			//Authors:     item.Authors,
			//Image: item.Image,
		}
		articlesChan <- article
	}
	close(articlesChan)
}

func doSummaries(articlesChan, doneChan chan Article) {
	for article := range articlesChan {
		cmd := exec.Command("venv/bin/python", "ai/summarize.py")
		cmd.Stdin = strings.NewReader(article.Description)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Failed to summarize article: %v", err)
			continue
		}

		var summary Summary
		if err = json.Unmarshal(output, &summary); err != nil {
			log.Printf("Failed to unmarshal summary: %v", err)
			continue
		}
		article.Summary = summary.Summary
		doneChan <- article
	}
}
