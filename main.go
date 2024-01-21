package main

import (
	"log"

	"u40apps.com/surfbot/pkg/api"
	"u40apps.com/surfbot/pkg/forecast"
	"u40apps.com/surfbot/pkg/kinetika"
)

func main() {
	log.Println("Start fetching Kinetika sessions")
	sessionsMsg, err := kinetika.FetchSessions()
	if err != nil {
		log.Panic("Panic fetching Kinetika sessions: ", err)
	}

	log.Println("Sending sessions", *sessionsMsg)
	api.SendMessage(*sessionsMsg)

	log.Println("Start fetching forecast")
	forecastMsg, err := forecast.FetchForecast()
	if err != nil {
		log.Panic("Panic fetching forecast: ", err)
	}

	log.Println("Sending forecast", *forecastMsg)
	api.SendMessage(*forecastMsg)
}
