package models

type Room struct {
	Id       string `json:"id" form:"-"`
	NoPlayer int    `json:"noplayer" form:"noplayer"`
	Amount   int    `json:"amount"`
	Host     string `json:"host"`
}

type CreateRoomMsg struct {
	Amount int `json:"amount"`
}

type RegisterHostMsg struct {
	IpAddress string `json:"ip"`
}
