package kinetika

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func requestSessions() (*[]Session, error) {
	fromDate := time.Now()
	toDate := fromDate.AddDate(0, 0, daysCount)

	fromDateFormatted := fromDate.UTC().Format(dateFormat)
	toDateFormatted := toDate.UTC().Format(dateFormat)

	url := "https://api.momence.com/host-plugins/host/31699/host-schedule/sessions"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.URL.Query().Set("fromDate", fromDateFormatted)
	req.URL.Query().Set("toDate", toDateFormatted)

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
		Sessions []Session `json:"payload"`
	}{}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	sessions := data.Sessions

	for i := 0; i < len(sessions); i++ {
		sessions[i].Date = localDate(sessions[i].Date)
	}

	return &sessions, nil
}
