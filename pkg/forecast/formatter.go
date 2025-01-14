package forecast

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func formatForecast(forecast Forecast, daysLimit int) string {
	grouped := make(map[string][]Period)
	order := []string{}
	var result = ""

	for _, period := range forecast.Periods {

		date := strings.Split(period.Localtime, " ")[0]

		if _, exists := grouped[date]; !exists {
			order = append(order, date)
		}

		grouped[date] = append(grouped[date], period)
	}

	for i, day := range order {
		if i >= daysLimit {
			break
		}

		periods := grouped[day]
		if len(periods) == 0 {
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

		for _, period := range periods {

			periodsResult += "\n"
			periodsResult += formatPeriod(period)
		}

		result += periodsResult
	}

	return result
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
