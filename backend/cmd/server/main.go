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
	subscribeHandler := createSubscribeHandler()
	emailHandler := createEmailHandler(dbConn, configs)

	r := router.New()
	r.GET("/", displayNewsHandler.DisplayNewsHandler)
	r.POST("/api/process", processNewsHandler.ProcessNewsHandler)
	r.GET("/favorites", displayNewsHandler.FavoritesHandler)
	r.POST("/api/favorites", fetchFavoritesHandler.FetchFavoritesHandler)
	r.GET("/subscription", subscribeHandler.SubscriptionHandler)
	r.POST("/api/save_email", emailHandler.SaveEmailHandler)
	r.POST("/api/mailing", emailHandler.MailingHandler)

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
		case path == "/", path == "/favorites", path == "/api/process", path == "/subscription",
			path == "/api/favorites", path == "/api/save_email", path == "/api/mailing":
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
	dbRepo := newsDB.NewDBOperations(dbConn, configs.Database.NewsTable, configs.Database.EmailTable, configs.Server.MainPath, configs.Database.TTL)
	newsService := service.NewDisplayNewsService(dbRepo)

	return handler.NewDisplayNewsHandler(newsService)
}

func createProcessHandler(dbConn *adapter.DBConnection, configs *configuration.Configuration) *handler.ProcessNewsHandler {
	dbRepo := newsDB.NewDBOperations(dbConn, configs.Database.NewsTable, configs.Database.EmailTable, configs.Server.MainPath, configs.Database.TTL)
	newsService := service.NewProcessNewsService(dbRepo, configs.Workers.NumWorkers, configs.OpenAPI.Url, configs.OpenAPI.Key, configs.OpenAPI.Disable)

	return handler.NewProcessNewsHandler(newsService, configs.Workers.NumWorkers)
}

func createSubscribeHandler() *handler.SubscriptionHandler {
	return handler.NewSubscriptionHandler()
}

func createFetchFavoritesHandler(dbConn *adapter.DBConnection, configs *configuration.Configuration) *handler.FetchFavoritesHandler {
	dbRepo := newsDB.NewDBOperations(dbConn, configs.Database.NewsTable, configs.Database.EmailTable, configs.Server.MainPath, configs.Database.TTL)
	fetchService := service.NewFetchService(dbRepo)

	return handler.NewFetchFavoritesHandler(fetchService)
}

func createEmailHandler(dbConn *adapter.DBConnection, configs *configuration.Configuration) *handler.EmailHandler {
	dbRepo := newsDB.NewDBOperations(dbConn, configs.Database.NewsTable, configs.Database.EmailTable, configs.Server.MainPath, -1)
	saveEmailService := service.NewSaveEmailService(dbRepo, configs.Email.HostEmail, configs.Email.HostPass, configs.Email.EmailTopic)
	return handler.NewEmailHandler(saveEmailService)
}
