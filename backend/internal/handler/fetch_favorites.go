package handler

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"

	"top-news/backend/internal/service"
)

type FetchFavoritesHandler struct {
	FetchService *service.FetchService
}

func NewFetchFavoritesHandler(fetchService *service.FetchService) *FetchFavoritesHandler {
	return &FetchFavoritesHandler{
		FetchService: fetchService,
	}
}

func (h FetchFavoritesHandler) FetchFavoritesHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[POST] /api/favorites")

	ids := ctx.QueryArgs().PeekMulti("id")

	if len(ids) == 0 {
		log.Println("No ids found in request")
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}

	articles, err := h.FetchService.FetchFavorites(ids)
	if err != nil {
		log.Println("Error fetching favorites:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	log.Println("Favorites fetched successfully!")

	data, err := json.Marshal(articles)
	if err != nil {
		log.Println("Error marshalling articles:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetBody(data)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
