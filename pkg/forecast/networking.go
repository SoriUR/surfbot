package forecast

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func requestForecast(spotName string) (*Forecast, error) {

	baseUrl := "https://api.surf-forecast.com/s1/breaks/"
	url := baseUrl + spotName + "/forecast"

	log.Println("url: ", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX3Rva2VuIjoiLUJtLXktU2YzSEFJTjBJOSIsImV4cCI6MjAyMTIxNTgxNn0.hvaU3rW8ja5vT_hA6gAAfApHBpoW2sPJVx8IsHOtn-o")
	req.Header.Set("x-surf-app-key", "8UMANwMz.aFM4z9nk*DGAd")

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

	data := struct {
		Forecast Forecast
	}{}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data.Forecast, nil
}
