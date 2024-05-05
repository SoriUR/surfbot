package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"u40apps.com/surfbot/pkg/forecast"
	"u40apps.com/surfbot/pkg/kinetika"
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

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			msg.Text = handleCommand(update.Message.Command())
		} else {
			msg.Text = "Use commands for interactions"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func handleCommand(command string) string {
	switch command {
	case "forecast":
		return *makeForecastMessage()

	case "kinetika":
		return *makeKinetikaMessage()

	default:
		return "I don't know that command"
	}
}

func makeForecastMessage() *string {
	log.Println("Start fetching forecast")
	forecastMsg, err := forecast.FetchForecast()
	if err != nil {
		log.Panic("Panic fetching forecast: ", err)
	}

	log.Println("Sending forecast", *forecastMsg)
	return forecastMsg
}

func makeKinetikaMessage() *string {
	log.Println("Start fetching Kinetika sessions")
	sessionsMsg, err := kinetika.FetchSessions()
	if err != nil {
		log.Panic("Panic fetching Kinetika sessions: ", err)
	}
	log.Println("Sending sessions", *sessionsMsg)
	return sessionsMsg
}
