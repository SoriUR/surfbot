package forecast

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
