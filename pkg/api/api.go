package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var client = &http.Client{}

type Request struct {
	ChatID string `json:"chat_id" example:"5"`
	Text   string `json:"text"`
}

func SendMessage(text string) error {

	chatIDs := [2]string{"456464682", "813729"}

	for _, chatID := range chatIDs {
		data := Request{
			ChatID: chatID,
			Text:   text,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		baseURL := "https://api.telegram.org/bot6407513466:AAE17rRvxcUIjnj9KDh-AAa8q-kOCFBkCRQ"
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
	}

	return nil
}
