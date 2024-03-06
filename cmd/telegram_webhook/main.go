package main

import (
	"cdnetwork/internal/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	config := util.GetConfig()
	// Create a new bot instance
	bot, err := tgbotapi.NewBotAPI(config.Telegram.TelegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	// Set the webhook URL
	webhookURL := fmt.Sprintf("https://%s/webhook/", config.Telegram.TelegramWebhook)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		log.Fatal(err)
	}

	// Start a simple HTTP server to handle incoming updates
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		update := tgbotapi.Update{}
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Println("Error decoding update:", err)
			return
		}

		// Handle the update (e.g., send a response)
		log.Printf("Received update: %+v\n", update)
	})

	// Start the HTTP server
	log.Println("Server listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
