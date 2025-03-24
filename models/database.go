package models

import "time"

type Prompt struct {
	ID        int
	DateAsked time.Time
	Question  string
	Reply     string
	Votes     int
}
