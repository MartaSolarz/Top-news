package handler

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/valyala/fasthttp"

	"top-news/backend/internal/service"
)

type SaveEmailHandler struct {
	SaveEmailService *service.SaveEmailService
}

func NewSaveEmailHandler(saveEmailService *service.SaveEmailService) *SaveEmailHandler {
	return &SaveEmailHandler{
		SaveEmailService: saveEmailService,
	}
}

func (h *SaveEmailHandler) SaveEmailHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[POST] /api/save_email")

	var requestBody struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(ctx.PostBody(), &requestBody); err != nil {
		log.Println("Error decoding JSON body:", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	email := requestBody.Email
	if email == "" {
		log.Println("No email found in request")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err := h.SaveEmailService.SaveEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidEmail):
			log.Println("Invalid email:", err)
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBodyString("Invalid email")
			return
		default:
		}
		log.Println("Error saving email:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	log.Println("Email saved successfully!")

	ctx.SetStatusCode(fasthttp.StatusOK)

}
