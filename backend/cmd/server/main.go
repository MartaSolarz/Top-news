package main

import (
	"fmt"
	"log"
	"top-news/backend/internal/adapter"
	"top-news/backend/internal/repository/newsDB"
	"top-news/backend/internal/service"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"top-news/backend/internal/configuration"
	"top-news/backend/internal/handler"
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

	displayNewsHandler := createDisplayHandler(dbConn, configs.Database.DBTable, configs.Workers.NumWorkers)
	processNewsHandler := createProcessHandler(dbConn, configs.Database.DBTable, configs.Workers.NumWorkers)

	r := router.New()
	r.GET("/news", displayNewsHandler.DisplayNewsHandler)
	r.GET("/favorites", displayNewsHandler.FavoritesHandler)
	r.POST("/api/process", processNewsHandler.ProcessNewsHandler)

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
		case path == "/", path == "/home", path == "/news", path == "/favorites", path == "/api/process":
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

func createDisplayHandler(dbConn *adapter.DBConnection, tableName string, numOfWorkers int) *handler.DisplayNewsHandler {
	dbRepo := newsDB.NewDBOperations(dbConn, tableName)
	newsService := service.NewDisplayNewsService(dbRepo)

	return handler.NewDisplayNewsHandler(newsService, numOfWorkers)
}

func createProcessHandler(dbConn *adapter.DBConnection, tableName string, numOfWorkers int) *handler.ProcessNewsHandler {
	dbRepo := newsDB.NewDBOperations(dbConn, tableName)
	newsService := service.NewProcessNewsService(dbRepo)

	return handler.NewProcessNewsHandler(newsService, numOfWorkers)
}
