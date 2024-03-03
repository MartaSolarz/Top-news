package handler

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"html/template"
	"log"
	"os/exec"
	"strings"
	"top-news/backend/internal/parser"
)

type Handler struct {
	NewsURL string
}

func NewHandler(newsURL string) *Handler {
	return &Handler{
		NewsURL: newsURL,
	}
}

func (h *Handler) HomeHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetBodyString("Welcome to the home page!")
}

type Summary struct {
	Summary string `json:"summary"`
}

type Article struct {
	Title       string
	Description string
	Summary     string
}

func (h *Handler) FetchRSSHandler(ctx *fasthttp.RequestCtx) {
	feed, err := parser.FetchRSS(h.NewsURL)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}

	tmpl, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	articlesChan := make(chan Article)
	summariesChan := make(chan Article)

	go func() {
		for article := range articlesChan {
			cmd := exec.Command("python3", "/Users/martasolarz/Github/Top-news/ai/summarize.py")
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
			summariesChan <- article
		}
	}()

	go func() {
		for _, item := range feed.Items {
			article := Article{
				Title:       item.Title,
				Description: item.Description,
			}
			articlesChan <- article
		}
		close(articlesChan)
	}()

	var articles []Article
	for i := 0; i < len(feed.Items); i++ {
		summary := <-summariesChan
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
