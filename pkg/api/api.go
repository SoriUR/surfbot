package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Request struct {
	ChatID string `json:"chat_id" example:"5"`
	Text   string `json:"text"`
}

func SendMessage(text string) error {

	chatIDs := [3]string{
		"456464682",
		"813729",
		"124570373",
	}

	for _, chatID := range chatIDs {
		data := Request{
			ChatID: chatID,
			Text:   text,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		envErr := godotenv.Load()
		if envErr != nil {
			log.Fatal("Error loading .env file")
		}
		token := os.Getenv("API_TOKEN")

		baseURL := "https://api.telegram.org/bot" + token
		path := "/sendMessage"
		url := baseURL + path
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

		req.Header.Set("Content-Type", "application/json")

		var client = &http.Client{}
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
