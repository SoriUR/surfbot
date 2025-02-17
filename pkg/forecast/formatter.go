package forecast

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func formatForecast(forecast Forecast, daysLimit int) string {

	periodsByDate := make(map[string][]Period)
	datesOrder := []string{}

	location := calculateLocation(forecast.Periods[0])
	if location == nil {
		fmt.Println("Forecast Location is unknown")
		location = time.Now().Location()
	}
	fmt.Printf("Forecast Location %v", location)

	var result = ""

	for _, period := range forecast.Periods {

		localDate := strings.Split(period.Localtime, " ")[0]

		if _, exists := periodsByDate[localDate]; !exists {
			datesOrder = append(datesOrder, localDate)
		}

		periodsByDate[localDate] = append(periodsByDate[localDate], period)
	}

	for i, day := range datesOrder {
		if i >= daysLimit {
			break
		}

		periods := periodsByDate[day]
		if len(periods) == 0 {
			continue
		}
		sort.Slice(periods, func(i, j int) bool {
			ti := time.Unix(periods[i].Timestamp, 0)
			tj := time.Unix(periods[j].Timestamp, 0)
			return ti.Before(tj)
		})

		dayTime := time.Unix(periods[0].Timestamp, 0).In(location)
		result += "\n\n"
		result += dayTime.Weekday().String() + " " + dayTime.Format("02.01")

		var periodsResult = ""

		for _, period := range periods {

			periodsResult += "\n"
			periodsResult += formatPeriod(period)
		}

		result += periodsResult
	}

	return result
}

func calculateLocation(period Period) *time.Location {

	localTime, err := time.Parse("2006-01-02 15:04:05", period.Localtime)
	if err != nil {
		return nil
	}

	utcTime := time.Unix(period.Timestamp, 0).UTC()
	offset := int(localTime.Sub(utcTime).Hours())

	return time.FixedZone("", offset*60*60)
}

func formatPeriod(period Period) string {

	parsedTime, _ := time.Parse("2006-01-02 15:04:05", period.Localtime)
	formatted := "- " + parsedTime.Format("15:04") + ": \n"

	energyStr := "⚡️" + period.Energy
	swellStr := "🌊 " + fmt.Sprintf("%.1f", period.Swell.Height)
	windStr := "💨 " + strconv.Itoa(int(period.Wind.Speed))
	starsStr := "⭐️ " + period.Stars

	elements := []string{
		energyStr,
		swellStr,
		windStr,
		starsStr,
	}

	return formatted + strings.Join(elements, "  ")
}
