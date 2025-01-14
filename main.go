package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"u40apps.com/surfbot/pkg/forecast"
	"u40apps.com/surfbot/pkg/setup"
)

func main() {

	setup.SetupBot(setup.ApiToken())
	setup.SetupDB("surf_bot")
	defer setup.DisconnectDB()

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
			"Check Uluwatu forecast": "try_uluwatu",
		}
		sendMsgButtons(chatID, startMessage, buttons)

	case "/help":
		sendMsg(chatID, helpMessage)

	case "/ping":
		sendMsg(chatID, "pong")

	default:
		spotName, daysLimit, err := splitCommand(command)
		if err != nil {
			sendMsg(chatID, "Unknown command. See available commands by typing /help")
			return
		}
		log.Printf("Successfully parsed command. Spot: %v. Days: %v", spotName, daysLimit)
		handleForecastCommand(chatID, spotName, daysLimit)
	}
}

func handleForecastCommand(chatID int64, spotName string, daysLimit int) {
	forecast, err := makeForecastMessage(spotName, daysLimit)
	if err != nil {
		sendMsg(chatID, "Sorry. Unable to get forecast. Try later")
	} else {
		sendMsg(chatID, *forecast)
	}
}

func makeForecastMessage(spotName string, daysLimit int) (*string, error) {
	log.Printf("Start fetching forecast. Spot: %v. Days: %v", spotName, daysLimit)
	forecastMsg, err := forecast.GetForecast(spotName, daysLimit)
	if err != nil {
		return nil, err
	}

	log.Println("Forecast created")
	return forecastMsg, nil
}

const startMessage = `
Hi, I am Surf Forecast Bot. I will help you to get the surf forecast at your favourite spot!

Example
Tuesday 01.14
- 08:00: 
⚡️908  📈 0.66  🌊 1.5  💨 10  ⭐️ 2

## How to read
⚡️ - Energy (kilo Joules)
📈 - Tide level (meters)
🌊 - Wave height (meters)
💨 - Wind (km/h)
⭐️ - Rating (10 - max)

` + helpMessage

const helpMessage = `
## How to use
1) Choose a spot. Available spots can be found at https://www.surf-forecast.com/countries.
2) Send me one of this commands:

- /<spot_name> - Forecast at the spot for 5 days.
Example: /uluwatu, /airport_lefts

- /<spot_name>_<days_limit> - Forecast at the spot for number of days.
Example: /uluwatu_1, /airport_lefts_5
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
	if len(command) > 0 && command[0] == '/' {
		command = command[1:]
	}

	parts := strings.Split(command, "_")

	if len(parts) == 0 {
		return "", 0, fmt.Errorf("input string is invalid")
	}

	lastPart := parts[len(parts)-1]
	number, err := strconv.Atoi(lastPart)
	if err != nil {
		number = 5
	} else {
		parts = parts[:len(parts)-1]
	}

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.Title(part)
		}
	}
	restString := strings.Join(parts, "-")

	return restString, number, nil
}
