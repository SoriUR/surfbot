package analytics

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

const amplitudeAPIURL = "https://api2.amplitude.com/2/httpapi"

var apiKey = os.Getenv("AMPLITUDE_API_KEY")

type AmplitudeEvent struct {
	UserID     string                 `json:"user_id"`
	EventType  string                 `json:"event_type"`
	EventProps map[string]interface{} `json:"event_properties,omitempty"`
}

type AmplitudePayload struct {
	APIKey string           `json:"api_key"`
	Events []AmplitudeEvent `json:"events"`
}

func LogEvent(eventType string, userID string, props map[string]interface{}) {
	var apiKey = os.Getenv("AMPLITUDE_API_KEY")

	event := AmplitudeEvent{
		UserID:     userID,
		EventType:  eventType,
		EventProps: props,
	}

	payload := AmplitudePayload{
		APIKey: apiKey,
		Events: []AmplitudeEvent{event},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Amplitude JSON marshal error: %v", err)
		return
	}

	log.Printf("Sending Amplitude Event: %v", payload)

	req, err := http.NewRequest("POST", amplitudeAPIURL, bytes.NewBuffer(body))

	if err != nil {
		log.Printf("Amplitude Create Req Error: %v", err)
		return
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Amplitude POST error: %v", err)
		return
	}

	log.Printf("Resp: %v", resp)
	defer resp.Body.Close()
}
