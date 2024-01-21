package tides

type Tides struct {
	Days []TidesDay
}

type TidesDay struct {
	Date     string `json:"tide_date"`
	Extremes []TideExtrema
}

type TideExtrema struct {
	High float64
	Low  float64
	Time string
}

func (extrema TideExtrema) Extrema() float64 {
	if extrema.High != 0 {
		return extrema.High
	}
	return extrema.Low
}
