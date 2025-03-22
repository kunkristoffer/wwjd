package models

type Prompt struct {
	ID        int
	DateAsked string
	Question  string
	Reply     string
	Votes     int
}
