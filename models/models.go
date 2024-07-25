package models

type CronResponse struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Year       string
	Command    string
}

type Config struct {
	Min, Max int
	Name     string
}

const (
	Minute = iota
	Hour
	DayOfMonth
	Month     = 3
	DayOfWeek = 4
	//Command   = 5
)

var CronConfig = map[int]Config{
	0: {Min: 0, Max: 59, Name: "Minute"},
	1: {Min: 0, Max: 23, Name: "Hour"},
	2: {Min: 1, Max: 31, Name: "DayOfMonth"},
	3: {Min: 1, Max: 12, Name: "Month"},
	4: {Min: 0, Max: 6, Name: "DayOfWeek"},
	5: {Min: 1, Max: 3000, Name: "Year"},
}
