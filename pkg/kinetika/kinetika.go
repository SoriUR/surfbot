package kinetika

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func FetchSessions() (*string, error) {
	sessions, err := requestSessions()
	if err != nil {
		return nil, err
	}

	text := formatSessions(*sessions)
	return &text, nil
}

var client = &http.Client{}

const dateFormat = "2006-01-02T15:04:05.000Z"
const timeFormat = "15:04"
const daysCount = 7

type Session struct {
	ID       int       `json:"id"`
	Name     string    `json:"sessionName"`
	Date     time.Time `json:"startsAt"`
	Link     string    `example:"https://momence.com/s/86299624"`
	Location string    `example:"Kedungu"`
	Teacher  string    `example:"Vladi Bagus"`
	Capacity int       `example:"5"`
	Booked   int       `json:"ticketsSold" example:"5"`
}

func requestSessions() (*[]Session, error) {
	fromDate := time.Now()
	toDate := fromDate.AddDate(0, 0, daysCount)

	fromDateFormatted := fromDate.UTC().Format(dateFormat)
	toDateFormatted := toDate.UTC().Format(dateFormat)

	baseUrl := "https://api.momence.com/host-plugins/host/31699/host-schedule/sessions"
	url := baseUrl + "?" + "fromDate=" + fromDateFormatted + "&" + "toDate=" + toDateFormatted

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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

func formatSessions(sessions []Session) string {
	var sessionsByDay = map[string][]Session{}
	var days = []string{}

	for _, session := range sessions {
		formatted := session.Date.Format("2006-01-02")

		if sessionsByDay[formatted] == nil {
			days = append(days, formatted)
		}
		sessionsByDay[formatted] = append(sessionsByDay[formatted], session)
	}

	var result = "Kinetika Schedule:"

	for _, day := range days {
		sessions := sessionsByDay[day]

		if sessions == nil || len(sessions) == 0 {
			continue
		}

		sort.Slice(sessions, func(i, j int) bool {
			return sessions[i].Date.Before(sessions[j].Date)
		})

		dayTime := sessions[0].Date
		result += "\n\n"
		result += dayTime.Weekday().String() + " " + dayTime.Format("01.02")

		var sessionsResult = ""

		for _, session := range sessions {
			sessionsResult += "\n"
			sessionsResult += formatSession(session)
		}

		result += sessionsResult
	}

	return result
}

func formatSession(session Session) string {
	elements := []string{
		"- " + session.Date.Format(timeFormat) + " " + session.Location,
		availability(session),
		session.Teacher,
	}

	return strings.Join(elements, " - ")
}

func availability(session Session) string {
	return strconv.Itoa(session.Booked) + "/" + strconv.Itoa(session.Capacity)
}

func localDate(date time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		return date
	}

	return date.In(loc)
}
