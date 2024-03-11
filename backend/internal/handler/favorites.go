package handler

import (
	"html/template"
	"log"

	"github.com/valyala/fasthttp"
)

func (h *DisplayNewsHandler) FavoritesHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[GET] /favorites")

	tmpl, err := template.ParseFiles("frontend/html/favorites.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	ctx.SetContentType("text/html")
	err = tmpl.Execute(ctx, nil)
	if err != nil {
		log.Println("Error executing template:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
	}
}
