package forecast

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func FetchForecast() (*string, error) {
	spotName := "Canggu"
	log.Println("Requesting forecast for: ", spotName)
	response, err := requestForecast(spotName)
	if err != nil {
		return nil, err
	}

	log.Println("Formatting forecast: ", *response)
	text := formatForecast(*response, spotName)
	return &text, nil
}

var client = &http.Client{}

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

func formatForecast(forecast Forecast, spotName string) string {

	var periodsByDay = map[string][]Period{}
	var days = []string{}
	var endDate = time.Now().AddDate(0, 0, 5).Unix()

	for _, period := range forecast.Periods {
		if period.Timestamp > endDate {
			continue
		}
		time := time.Unix(period.Timestamp, 0)
		formatted := time.Format("2006-01-02")

		if periodsByDay[formatted] == nil {
			days = append(days, formatted)
		}
		periodsByDay[formatted] = append(periodsByDay[formatted], period)
	}

	var result = spotName + " Forecast:"

	for _, day := range days {
		periods := periodsByDay[day]

		if periods == nil || len(periods) == 0 {
			continue
		}

		sort.Slice(periods, func(i, j int) bool {
			ti := time.Unix(periods[i].Timestamp, 0)
			tj := time.Unix(periods[j].Timestamp, 0)
			return ti.Before(tj)
		})

		dayTime := time.Unix(periods[0].Timestamp, 0)
		result += "\n\n"
		result += dayTime.Weekday().String() + " " + dayTime.Format("01.02")

		var periodsResult = ""

		for _, session := range periods {
			periodsResult += "\n"
			periodsResult += formatPeriod(session)
		}

		result += periodsResult
	}

	return result
}

func formatPeriod(period Period) string {
	time := time.Unix(period.Timestamp, 0)
	formatted := "- " + time.Format("15:04") + ": "

	stars := ""
	number, err := strconv.Atoi(period.Stars)
	if err == nil {
		for i := 0; i < number; i++ {
			stars += "⭐️"
		}
	}

	elements := []string{
		period.Energy + "⚡️",
		strconv.Itoa(int(period.Wind.Speed)) + " " + period.WindDirection + "💨",
		stars,
	}

	return formatted + strings.Join(elements, " / ")
}

type Forecast struct {
	Periods []Period `json:"t_periods_all_days"`
}

type Period struct {
	Stars         string
	Energy        string `json:"maxenergy"`
	Timestamp     int64  `json:"tstampstart_utc"`
	Wind          Wind
	WindDirection string `json:"offshoreness"`
}

type Wind struct {
	Speed float64
}
