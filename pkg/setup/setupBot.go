package setup

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var Bot *tgbotapi.BotAPI

func SetupBot(apiToken string) error {
	envErr := godotenv.Load()
	if envErr != nil {
		return envErr
	}

	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return err
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	Bot = bot

	return nil
}
