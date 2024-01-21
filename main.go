package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"u40apps.com/surfbot/pkg/forecast"
	"u40apps.com/surfbot/pkg/kinetika"
)

func main() {
	log.Println("Start fetching Kinetika sessions")
	sessionsMsg, err := kinetika.FetchSessions()
	if err != nil {
		log.Panic("Panic fetching Kinetika sessions", err)
	}

	log.Println("Sending sessions", *sessionsMsg)
	sendMessage(*sessionsMsg)

	log.Println("Start fetching forecast")
	forecastMsg, err := forecast.FetchForecast()
	if err != nil {
		log.Panic("Panic fetching forecast", err)
	}

	log.Println("Sending forecast", *forecastMsg)
	sendMessage(*forecastMsg)
}

var client = &http.Client{}

func sendMessage(text string) error {

	data := struct {
		ChatID string `json:"chat_id" example:"5"`
		Text   string `json:"text"`
	}{
		ChatID: "456464682",
		Text:   text,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	baseURL := "https://api.telegram.org/bot6720343526:AAHiT0Nlh5ZYcfkq8tv5MO53GJkA-IroQmQ"
	path := "/sendMessage"
	url := baseURL + path
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		return err1
	}

	return nil
}
