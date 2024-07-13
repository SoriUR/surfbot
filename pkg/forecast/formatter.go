package forecast

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cnkei/gospline"
	"u40apps.com/surfbot/pkg/tides"
)

var location, _ = time.LoadLocation("Asia/Singapore")

func FetchForecast() (*string, error) {

	log.Println("Requesting tides")
	tides, err := tides.RequestTides()
	if err != nil {
		return nil, err
	}

	spotName := "Canggu"

	log.Println("Requesting forecast for: ", spotName)
	response, err := requestForecast(spotName)
	if err != nil {
		return nil, err
	}

	log.Println("Formatting forecast: ")
	text := formatForecast(*response, spotName, *tides)
	return &text, nil
}

func groupPeriods(forecast Forecast) (map[string][]Period, []string) {
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

	return periodsByDay, days
}

func formatForecast(forecast Forecast, spotName string, tides tides.Tides) string {

	periodsByDay, days := groupPeriods(forecast)
	spline := makeSpline(tides)

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

		for _, period := range periods {

			periodsResult += "\n"
			periodsResult += formatPeriod(period, spline)
		}

		result += periodsResult
	}

	return result
}

func formatPeriod(period Period, spline gospline.Spline) string {
	time := time.Unix(period.Timestamp, 0)
	time = time.In(location)
	formatted := "- " + time.Format("15:04") + ": \n"

	starsStr := "⭐️ " + period.Stars
	energyStr := "⚡️" + period.Energy
	tideStr := "📈 " + fmt.Sprintf("%.2f", spline.At(float64(period.Timestamp)))
	windStr := "💨 " + strconv.Itoa(int(period.Wind.Speed))
	swellStr := "🌊 " + fmt.Sprintf("%.1f", period.Swell.Height)

	elements := []string{
		energyStr,
		tideStr,
		swellStr,
		windStr,
		starsStr,
	}

	return formatted + strings.Join(elements, "  ")
}

func makeSpline(tidesData tides.Tides) gospline.Spline {

	xs := []float64{}
	ys := []float64{}

	for _, day := range tidesData.Days {
		for _, extrema := range day.Extremes {
			t, err := time.ParseInLocation("2006-01-02T15:04:05", extrema.Time, location)
			if err == nil {
				xs = append(xs, float64(t.Unix()))
				ys = append(ys, extrema.Extrema())
			}
		}
	}

	return gospline.NewCubicSpline(xs, ys)
}
