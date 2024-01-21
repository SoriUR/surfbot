package kinetika

import "time"

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
