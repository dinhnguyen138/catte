package models

type Room struct {
	Id       string `json:"id" form:"-"`
	NoPlayer int    `json:"noplayer" form:"noplayer"`
	Amount   int    `json:"amount"`
}

type CreateRoomMsg struct {
	Amount int `json:"amount"`
}
