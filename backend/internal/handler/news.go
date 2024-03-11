package handler

import (
	"html/template"
	"log"

	"github.com/valyala/fasthttp"

	"top-news/backend/internal/service"
)

type DisplayNewsHandler struct {
	NewsService *service.DisplayNewsService
}

func NewDisplayNewsHandler(newsService *service.DisplayNewsService) *DisplayNewsHandler {
	return &DisplayNewsHandler{
		NewsService: newsService,
	}
}

func (h *DisplayNewsHandler) DisplayNewsHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[GET] /news")

	tmpl, err := template.ParseFiles("frontend/html/news.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	articles, err := h.NewsService.GetNews()
	if err != nil {
		log.Println("Error getting news:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	log.Println("Articles retrived successfully!")

	ctx.SetContentType("text/html")
	err = tmpl.Execute(ctx, articles)
	if err != nil {
		log.Println("Error executing template:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
	}
}
