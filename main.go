package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"go.mongodb.org/mongo-driver/mongo"

	"u40apps.com/surfbot/pkg/forecast"
	"u40apps.com/surfbot/pkg/setup"
)

var chatsCollection *mongo.Collection

func main() {

	apiToken := setup.ReadEnv("API_TOKEN")
	setup.SetupBot(apiToken)

	setup.SetupDB("surf_bot")
	collection, _ := setup.SetupCollection("chats")

	chatsCollection = collection

	handleUpdates()
}

func handleUpdates() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	bot := setup.Bot
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			callback.ShowAlert = false
			bot.Request(callback)

			handleCallback(update)
			continue
		}

		if update.Message.IsCommand() {
			handleCommand(update)
			continue
		}

		chatID := update.Message.Chat.ID
		sendMsg(chatID, "Only commands are supported. See available commands by typing /help")
	}
}

func handleCommand(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	args := strings.Fields(update.Message.Text)
	log.Printf("Handling command: %s", update.Message.Text)

	command := args[0]

	switch command {
	case "/start":
		buttons := map[string]string{
			"Try for Uluwatu": "try_uluwatu",
		}
		sendMsgButtons(chatID, startMessage, buttons)

	case "/help":
		sendMsg(chatID, helpMessage)

	case "/forecast":
		if len(args) < 3 {
			sendMsg(chatID, "Invalid arguments. Provide spot name and days limit. Example: /forecast Uluwatu 2")
			return
		}
		spotName := args[1]
		daysLimit, err := strconv.Atoi(args[2])
		if err != nil || daysLimit < 1 || daysLimit > 7 {
			sendMsg(chatID, "Invalid days limit. Should be between 1 and 7")
			return
		}

		handleForecastCommand(chatID, spotName, daysLimit)

	case "/ping":
		sendMsg(chatID, "pong")

	default:
		spotName, daysLimit, err := splitCommand(command)
		if err != nil {
			sendMsg(chatID, "Unknown command. See available commands by typing /help")
			return
		}
		log.Printf("Successfullt parsed command. Spot: %v. Days %v", spotName, daysLimit)
		handleForecastCommand(chatID, spotName, daysLimit)
	}
}

func handleForecastCommand(chatID int64, spotName string, daysLimit int) {
	forecast, err := makeForecastMessage(spotName, daysLimit)
	if err != nil {
		sendMsg(chatID, "Sorry. Unable to retrive forecase. Try later")
	} else {
		sendMsg(chatID, *forecast)
	}
}

func makeForecastMessage(spotName string, daysLimit int) (*string, error) {
	log.Printf("Start fetching forecast. Spot: %v. Days %v", spotName, daysLimit)
	forecastMsg, err := forecast.FetchForecast(spotName, daysLimit)
	if err != nil {
		return nil, err
	}

	log.Println("Forecast created")
	return forecastMsg, nil
}

const startMessage = `
Hi, I am Surf Forecast Bot. I will help you to get the surf forecast at your favourite spot!

How to read:
⚡️ - Energy (kilo Joules)
📈 - Tide level (meters)
🌊 - Wave heght (meters)
💨 - Wind (km/h)
⭐️ - Rating (10 - max)

Available commands:
1) /forecast - Forecast at any spot for number of days. 
Example: "/forecast Uluwatu 1"

2) /<spot_name> - Forecast at the spot for 3 days.
Example: "/uluwatu", "/airport_lefts"

3) /<spot_name>_<days_limit> - Forecast at the spot for number of days.
Example: "/uluwatu_1", "/airport_lefts_5"

You can find spot names at:
https://www.surf-forecast.com/countries
`

const helpMessage = `
1) /forecast - Forecast at any spot for number of days. 
Example: /forecast Uluwatu 1

2) /<spot_name> - Forecast at the spot for 3 days.
Example: /uluwatu, /airport_lefts

3) /<spot_name>_<days_limit> - Forecast at the spot for number of days.
Example: /uluwatu_1, /airport_lefts_5

You can find spot names at:
https://www.surf-forecast.com/countries
`

func sendMsg(chatID int64, message string) {
	sendMsgButtons(chatID, message, map[string]string{})
}

func sendMsgButtons(chatID int64, message string, buttonsMap map[string]string) {
	msg := tgbotapi.NewMessage(chatID, "")

	msg.Text = message
	if len(buttonsMap) != 0 {
		var buttons [][]tgbotapi.InlineKeyboardButton
		for text, callbackData := range buttonsMap {

			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(text, callbackData),
			))
		}
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	}

	setup.Bot.Send(msg)
}

func handleCallback(update tgbotapi.Update) {
	callback := update.CallbackQuery
	message := callback.Message
	chatID := message.Chat.ID

	switch callback.Data {
	case "try_uluwatu":
		sendMsg(chatID, "/uluwatu")
		handleForecastCommand(chatID, "Uluwatu", 3)

	default:
		return
	}
}

func splitCommand(command string) (string, int, error) {
	log.Printf("1 command %v", command)
	if len(command) > 0 && command[0] == '/' {
		command = command[1:]
	}

	log.Printf("2 command %v", command)

	parts := strings.Split(command, "_")

	log.Printf("3 parts %v", parts)

	if len(parts) == 0 {
		return "", 0, fmt.Errorf("input string is invalid")
	}

	lastPart := parts[len(parts)-1]
	number, err := strconv.Atoi(lastPart)
	if err != nil {
		number = 3
	} else {
		parts = parts[:len(parts)-1]
	}

	log.Printf("4 number %v", number)

	log.Printf("5 restParts %v", parts)

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.Title(part)
		}
	}
	restString := strings.Join(parts, "-")

	return restString, number, nil
}
