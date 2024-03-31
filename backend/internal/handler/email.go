package handler

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/valyala/fasthttp"

	"top-news/backend/internal/service"
)

type EmailHandler struct {
	EmailService *service.EmailService
}

func NewEmailHandler(EmailService *service.EmailService) *EmailHandler {
	return &EmailHandler{
		EmailService: EmailService,
	}
}

func (h *EmailHandler) SaveEmailHandler(ctx *fasthttp.RequestCtx) {
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

	err := h.EmailService.SaveEmail(email)
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

func (h *EmailHandler) MailingHandler(ctx *fasthttp.RequestCtx) {
	log.Printf("[POST] /api/mailing")

	err := h.EmailService.SendEmails()
	if err != nil {
		log.Println("Error sending emails:", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Internal Server Error")
		return
	}

	log.Println("Emails sent successfully!")

	ctx.SetStatusCode(fasthttp.StatusOK)
}
