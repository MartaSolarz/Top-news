package handler

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"

	"top-news/backend/internal/parser"
)

type Handler struct {
	NewsURL    string
	NumWorkers int
}

func NewHandler(newsURL string, numWorkers int) *Handler {
	return &Handler{
		NewsURL:    newsURL,
		NumWorkers: numWorkers,
	}
}

func (h *Handler) HomeHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[GET] /home")
	feed, err := parser.FetchRSS(h.NewsURL)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	fmt.Fprint(ctx, feed)

	//ctx.SetBodyString("Welcome to the home page!")
}
