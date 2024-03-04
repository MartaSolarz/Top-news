package main

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"top-news/backend/internal/configuration"
	"top-news/backend/internal/handler"
)

var configPath = "backend/internal/configuration/configuration.toml"

func main() {
	configs := configuration.NewConfiguration(configPath)

	serverAddress := fmt.Sprintf(":%d", configs.Server.Port)

	h := handler.NewHandler(configs.News.BBCNewsRSSURL, configs.Workers.NumWorkers)

	r := router.New()
	r.GET("/home", h.HomeHandler)
	r.GET("/news", h.NewsHandler)
	r.GET("/favorites", h.FavoritesHandler)

	fs := &fasthttp.FS{
		Root:               "frontend",             // Katalog z plikami frontendowymi
		IndexNames:         []string{"index.html"}, // Nazwy plików indeksowych
		GenerateIndexPages: false,                  // Nie generuj stron indeksowych
		AcceptByteRange:    false,                  // Obsługa zakresów bajtów
	}
	fsHandler := fs.NewRequestHandler()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		switch {
		case path == "/", path == "/home", path == "/news", path == "/favorites":
			r.Handler(ctx)
		default:
			fsHandler(ctx)
		}
	}

	server := &fasthttp.Server{
		Handler: requestHandler,
	}

	log.Printf("Server is starting on %s...", serverAddress)
	if err := server.ListenAndServe(serverAddress); err != nil {
		log.Fatalf("Error in starting server: %v", err)
	}
}
