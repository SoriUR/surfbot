package forecast

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Forecast struct {
	Periods []Period `json:"t_periods_all_days"`
}

type Period struct {
	Stars         string
	Energy        string `json:"maxenergy"`
	Timestamp     int64  `json:"ts"`
	Localtime     string `json:"localtime"`
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

type DBForecast struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Spot     string             `bson:"spot"`
	Forecast Forecast           `bson:"forecast"`
}
