package handler

import (
	"log"

	"github.com/valyala/fasthttp"
)

func (h *DisplayNewsHandler) FavoritesHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[GET] /favorites")
	ctx.SetBodyString("Welcome to the favorites page!")
}
