package handler

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
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

	var body interface{}
	err := json.Unmarshal(ctx.PostBody(), &body)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(fmt.Sprintf("Error processing request: %v", err))
		return
	}

	h.NewsService.ProcessNews(body)
}
