package handler

import (
	"github.com/valyala/fasthttp"
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
