package handler

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"

	"top-news/backend/internal/models"
	"top-news/backend/internal/service"
)

type ProcessNewsHandler struct {
	NewsService *service.ProcessNewsService
	NumWorkers  int
}

func NewProcessNewsHandler(newsService *service.ProcessNewsService, numWorkers int) *ProcessNewsHandler {
	return &ProcessNewsHandler{
		NewsService: newsService,
		NumWorkers:  numWorkers,
	}
}

func (h *ProcessNewsHandler) ProcessNewsHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[POST] /api/process")
	ctx.SetBodyString("Processing news...")

	var data []*models.FeedItem
	err := json.Unmarshal(ctx.PostBody(), &data)
	if err != nil {
		log.Printf("Error unmarshalling data: %v", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Bad Request")
		return
	}

	log.Printf("Processing %d news items", len(data))

	err = h.NewsService.ProcessNews(data)
	if err != nil {
		log.Printf("Error processing news: %v", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	log.Printf("Processed news successfully!")

	ctx.SetStatusCode(fasthttp.StatusOK)
}
