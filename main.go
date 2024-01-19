package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"u40apps.com/surfbot/pkg/kinetika"
)

func main() {
	sessions, err := kinetika.FetchSessions()
	if err != nil {
		panic(err)
	}

	sendSessions(*sessions)
}

var client = &http.Client{}

func sendSessions(text string) error {

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
