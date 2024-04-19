package loader

import "time"

const (
	ParameterTypeSearch = "search"
	ParameterTypeUser   = "user"
)

type UlogParameters struct {
	Text    string    `json:"text"`
	Full    bool      `json:"full"`
	UserID  int       `json:"user_id"`
	IsVC    bool      `json:"is_vc"`
	IsAdmin bool      `json:"is_admin"`
	Date    time.Time `json:"date"`
	Listen  int       `json:"listen"`
	Type    string    `json:"type"`
}

type GamePanelParameters struct {
	Nickname         string            `json:"nickname"`
	Reason           string            `json:"reason"`
	Params           []string          `json:"params"`
	Date             time.Time         `json:"date"`
	Dates            []time.Time       `json:"dates"`
	Listen           int               `json:"listen"`
	Type             string            `json:"type"`
	AdditionalParams map[string]string `json:"additional_params"`
}
