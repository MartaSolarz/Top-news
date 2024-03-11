package handler

import "github.com/valyala/fasthttp"

type ContactHandler struct {
}

func NewContactHandler() *ContactHandler {
	return &ContactHandler{}
}

func (h *ContactHandler) ContactHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetBodyString("Welcome to the contact page!")
}
