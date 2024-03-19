package handler

import (
	"github.com/valyala/fasthttp"
	"html/template"
	"log"
)

type SubscriptionHandler struct {
}

func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{}
}

func (h *SubscriptionHandler) SubscriptionHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[GET] /subscription")

	tmpl, err := template.ParseFiles("frontend/html/subscription.html")
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
