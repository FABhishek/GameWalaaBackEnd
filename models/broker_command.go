package models

type Command struct {
	Action string      `json:"action"`
	Data   GameDetails `json:"data"`
}
