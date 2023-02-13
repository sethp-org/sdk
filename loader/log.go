package loader

import "time"

type Log struct {
	Date   time.Time `json:"date"`
	Text   string    `json:"text"`
	Raw    string    `json:"raw"`
	Sender BaseUser  `json:"sender"`
	Target BaseUser  `json:"target"`
}

type BaseUser struct {
	IP *IP `json:"ip"`
}

type IP struct {
	Reg  string `json:"reg"`
	Last string `json:"last"`
}

type User struct {
	Server    string    `json:"server"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	BanDays   int       `json:"ban_days"`
	LastLogin time.Time `json:"last_login"`
}
