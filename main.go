package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("API_TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}

	// log.Println("Start fetching Kinetika sessions")
	// sessionsMsg, err := kinetika.FetchSessions()
	// if err != nil {
	// 	log.Panic("Panic fetching Kinetika sessions: ", err)
	// }

	// log.Println("Sending sessions", *sessionsMsg)
	// api.SendMessage(*sessionsMsg)

	// log.Println("Start fetching forecast")
	// forecastMsg, err := forecast.FetchForecast()
	// if err != nil {
	// 	log.Panic("Panic fetching forecast: ", err)
	// }

	// log.Println("Sending forecast", *forecastMsg)
	// api.SendMessage(*forecastMsg)

	// envErr := godotenv.Load()
	// if envErr != nil {
	// 	log.Fatal("Error loading .env file")
	// }
}
