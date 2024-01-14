package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

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

func main() {
	sessions, err := requestSessions()
	if err != nil {
		panic(err)
	}

	text := formatSessions(*sessions)
	// fmt.Println(text)

	sendSessions(text)
}

var client = &http.Client{}

const dateFormat = "2006-01-02T15:04:05.000Z"
const timeFormat = "15:04"
const daysCount = 7

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

	var result = ""

	for _, day := range days {
		sessions := sessionsByDay[day]

		if sessions == nil || len(sessions) == 0 {
			continue
		}

		sort.Slice(sessions, func(i, j int) bool {
			return sessions[i].Date.Before(sessions[j].Date)
		})

		if result != "" {
			result += "\n\n"
		}
		result += sessions[0].Date.Weekday().String()

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

func requestSessions() (*[]Session, error) {
	fromDate := time.Now()
	toDate := fromDate.AddDate(0, 0, daysCount)

	fromDateFormatted := fromDate.UTC().Format(dateFormat)
	toDateFormatted := toDate.UTC().Format(dateFormat)

	baseUrl := "https://api.momence.com/host-plugins/host/31699/host-schedule/sessions"
	url := baseUrl + "?" + "fromDate=" + fromDateFormatted + "&" + "toDate=" + toDateFormatted

	// fmt.Println(url)

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

	_, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return err1
	}

	return nil
}
