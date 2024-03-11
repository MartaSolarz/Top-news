package main

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"top-news/backend/internal/adapter"
	"top-news/backend/internal/configuration"
	"top-news/backend/internal/handler"
	"top-news/backend/internal/repository/newsDB"
	"top-news/backend/internal/service"
)

var configPath = "backend/internal/configuration/configuration.toml"

func main() {
	configs := configuration.NewConfiguration(configPath)

	dbConn, err := adapter.NewDBConnection(
		configs.Database.User,
		configs.Database.Password,
		configs.Database.Host,
		configs.Database.DBName,
		configs.Database.Port,
	)
	if err != nil {
		log.Fatalf("Error in connecting to database: %v", err)
	}
	log.Printf("Connected to database on %s:%d", configs.Database.Host, configs.Database.Port)

	serverAddress := fmt.Sprintf(":%d", configs.Server.Port)

	displayNewsHandler := createDisplayHandler(dbConn, configs)
	processNewsHandler := createProcessHandler(dbConn, configs)
	fetchFavoritesHandler := createFetchFavoritesHandler(dbConn, configs)
	contactHandler := createContactHandler(configs)

	r := router.New()
	r.GET("/news", displayNewsHandler.DisplayNewsHandler)
	r.POST("/api/process", processNewsHandler.ProcessNewsHandler)
	r.GET("/favorites", displayNewsHandler.FavoritesHandler)
	r.POST("/api/favorites", fetchFavoritesHandler.FetchFavoritesHandler)
	r.GET("/contact", contactHandler.ContactHandler)

	fs := &fasthttp.FS{
		Root:               "frontend",
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: false,
		AcceptByteRange:    false,
	}
	fsHandler := fs.NewRequestHandler()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		switch {
		case path == "/news", path == "/favorites", path == "/api/process", path == "/contact", path == "/api/favorites":
			r.Handler(ctx)
		default:
			fsHandler(ctx)
		}
	}

	server := &fasthttp.Server{
		Handler: requestHandler,
	}

	log.Printf("Server is starting on %s...", serverAddress)
	if err = server.ListenAndServe(serverAddress); err != nil {
		log.Fatalf("Error in starting server: %v", err)
	}
}

func createDisplayHandler(dbConn *adapter.DBConnection, configs *configuration.Configuration) *handler.DisplayNewsHandler {
	dbRepo := newsDB.NewDBOperations(dbConn, configs.Database.DBTable, configs.Database.TTL)
	newsService := service.NewDisplayNewsService(dbRepo)

	return handler.NewDisplayNewsHandler(newsService)
}

func createProcessHandler(dbConn *adapter.DBConnection, configs *configuration.Configuration) *handler.ProcessNewsHandler {
	dbRepo := newsDB.NewDBOperations(dbConn, configs.Database.DBTable, configs.Database.TTL)
	newsService := service.NewProcessNewsService(dbRepo, configs.Workers.NumWorkers)

	return handler.NewProcessNewsHandler(newsService, configs.Workers.NumWorkers)
}

func createContactHandler(configs *configuration.Configuration) *handler.ContactHandler {
	return handler.NewContactHandler()
}

func createFetchFavoritesHandler(dbConn *adapter.DBConnection, configs *configuration.Configuration) *handler.FetchFavoritesHandler {
	dbRepo := newsDB.NewDBOperations(dbConn, configs.Database.DBTable, configs.Database.TTL)
	fetchService := service.NewFetchService(dbRepo)

	return handler.NewFetchFavoritesHandler(fetchService)
}
