package models

type Command struct {
	Action string `json:"action"`
	Room   string `json:"room"`
	Id     string `json:"id"`
	Card   string `json:"card"`
}

type Deck []string
