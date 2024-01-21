package tides

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func RequestTides() (*Tides, error) {
	url := "https://dao2000.pythonanywhere.com/tides2/5698456-Benoa"
	log.Println("url: ", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.URL.Query().Set("format", "json")

	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := Tides{}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
