package main

import (
	"fmt"
	"github.com/beloslav13/telegram-vk-posts/config"
	"github.com/beloslav13/telegram-vk-posts/logger"
	"github.com/beloslav13/telegram-vk-posts/pkg/telegram/models"
	"github.com/beloslav13/telegram-vk-posts/server"
	"github.com/joho/godotenv"
	"log"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := config.NewConfig()
	w := models.Webhook{}
	webhook, err := w.CheckWebhook(*conf)
	if err != nil && !webhook {
		fmt.Println("LOG FATAl>>>>>>>")
		log.Fatal(fmt.Sprintf("webhook is %t. error: %s", webhook, err))
	}

	_, err = server.StartServer()
	if err != nil {
		logger.LogFile.Fatalln(fmt.Sprintf("Server not started... Error: %s", err))
	}
}