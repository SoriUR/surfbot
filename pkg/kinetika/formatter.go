package kinetika

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

func FetchSessions() (*string, error) {
	log.Println("Requesting sessions")
	sessions, err := requestSessions()
	if err != nil {
		return nil, err
	}

	log.Println("Formatting sessions")
	text := formatSessions(*sessions)
	return &text, nil
}

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

	formatted := "- " + session.Date.Format(timeFormat) + ": "

	locationStr := "📍" + session.Location
	availabilityStr := "🏄 " + availability(session)
	teacherStr := "🥸 " + session.Teacher

	elements := []string{
		fmt.Sprintf("%-12s", locationStr),
		fmt.Sprintf("%-9s", availabilityStr),
		fmt.Sprintf("%-10s", teacherStr),
	}

	return formatted + strings.Join(elements, " ")
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
