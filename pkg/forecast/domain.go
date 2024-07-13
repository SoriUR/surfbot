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
	Swell         Swell  `json:"foreground_swell"`
}

type Wind struct {
	Speed float64
}

type Swell struct {
	Height float64
}
