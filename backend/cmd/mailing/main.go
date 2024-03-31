package main

import (
	"fmt"
	"log"
	"time"

	"top-news/backend/internal/configuration"
	"top-news/backend/internal/utils"
)

var configPath = "backend/internal/configuration/configuration.toml"

func main() {
	configs := configuration.NewConfiguration(configPath)
	apiURL := fmt.Sprintf("http://%s:%d/api/mailing", configs.Server.Host, configs.Server.Port)

	var jsonData []byte
	utils.SendPostRequest(apiURL, jsonData)
	log.Println("Request to send mails sent to server, waiting for the next mailing...")

	time.Sleep(time.Hour * 24)
}
