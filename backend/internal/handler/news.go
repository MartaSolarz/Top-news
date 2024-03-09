package handler

import (
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"github.com/valyala/fasthttp"
	"html/template"
	"log"
	"os/exec"
	"strings"
	"sync"
	"top-news/backend/internal/service"

	"top-news/backend/internal/models"
	"top-news/backend/internal/parser"
)

type DisplayNewsHandler struct {
	NewsService *service.DisplayNewsService
	NewsURL     string
	NumWorkers  int
}

func NewDisplayNewsHandler(newsService *service.DisplayNewsService, numWorkers int) *DisplayNewsHandler {
	return &DisplayNewsHandler{
		NewsService: newsService,
		NumWorkers:  numWorkers,
	}
}

func (h *DisplayNewsHandler) DisplayNewsHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[GET] /news")
	feed, err := parser.FetchRSS(h.NewsURL)
	if err != nil {
		log.Printf("Error fetching RSS feed: %v", err)
		ctx.SetBodyString(err.Error())
		return
	}

	tmpl, err := template.ParseFiles("frontend/html/news.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	waitChan, doneChan := make(chan *models.Article), make(chan *models.Article)

	var wg sync.WaitGroup
	for i := 0; i < h.NumWorkers; i++ {
		wg.Add(1)
		go summarizeArticles(waitChan, doneChan, &wg)
	}

	go h.sendArticlesToChan(waitChan, feed.Items)

	var articles []*models.Article
	for i := 0; i < len(feed.Items); i++ {
		article := <-doneChan
		article.Website = feed.Title
		article.CopyRight = feed.Copyright
		articles = append(articles, article)
	}

	wg.Wait()

	ctx.SetContentType("text/html")
	err = tmpl.Execute(ctx, articles)
	if err != nil {
		log.Println("Error executing template:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
	}
}

func (h *DisplayNewsHandler) sendArticlesToChan(articlesChan chan<- *models.Article, items []*gofeed.Item) {
	var wg sync.WaitGroup
	itemChan := make(chan *gofeed.Item, len(items))

	for i := 0; i < h.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			for item := range itemChan {
				date := parser.ParseDate(item.Published)
				content := parser.ExtractContent(item.Link)
				thumbnail := parser.ParseThumbnail(item)
				articlesChan <- models.NewArticle(item, date, content, thumbnail)
			}
			wg.Done()
		}()
	}

	for _, item := range items {
		itemChan <- item
	}
	close(itemChan)

	wg.Wait()
	close(articlesChan)
}

type Summary struct {
	Summary string `json:"summary"`
}

func summarizeArticles(articlesChan, doneChan chan *models.Article, wg *sync.WaitGroup) {
	for article := range articlesChan {
		cmd := exec.Command("venv/bin/python", "python_scripts/ai/summarize.py")
		cmd.Stdin = strings.NewReader(article.Content)
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
	wg.Done()
}
