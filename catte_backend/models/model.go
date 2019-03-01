package models

type Command struct {
	Action string `json:"action"`
	Room   string `json:"room"`
	Id     string `json:"id"`
	Data   string `json:"data"`
}

type UserInfo struct {
	UUID     string  `json:"uuid" form:"-"`
	Username string  `json:"username" form:"username"`
	Amount   float32 `json:"amount" form:"amount"`
}

type Deck []string
