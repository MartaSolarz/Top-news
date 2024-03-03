package main

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"top-news/backend/internal/configuration"
	"top-news/backend/internal/handler"
)

var configPath = "backend/internal/configuration/configuration.toml"

func main() {
	configs := configuration.NewConfiguration(configPath)

	serverAddress := fmt.Sprintf(":%d", configs.Server.Port)

	h := handler.NewHandler(configs.News.BBCNewsRSSURL)

	r := router.New()
	r.GET("/", indexPageHandler)
	r.GET("/home", h.HomeHandler)
	r.GET("/fetchrss", h.FetchRSSHandler)

	server := &fasthttp.Server{
		Handler: r.Handler,
	}

	fmt.Println("Server is starting on port 8080...")
	if err := server.ListenAndServe(serverAddress); err != nil {
		fmt.Printf("Error in starting server: %v\n", err)
		return
	}
}

func indexPageHandler(ctx *fasthttp.RequestCtx) {
	fasthttp.ServeFile(ctx, "frontend/index.html")
}
