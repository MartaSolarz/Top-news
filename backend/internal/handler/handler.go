package handler

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"

	"top-news/backend/internal/parser"
)

func (h *DisplayNewsHandler) HomeHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[GET] /home")
	feed, err := parser.FetchRSS(h.NewsURL)
	if err != nil {
		ctx.SetBodyString(err.Error())
		return
	}
	fmt.Fprint(ctx, feed)

	//ctx.SetBodyString("Welcome to the home page!")
}
