package models

type Room struct {
	Id        string `json:"id" form:"-"`
	NoPlayer  int    `json:"noplayer"`
	MaxPlayer int    `json:"maxplayer"`
	Amount    int64  `json:"amount"`
	Host      string `json:"host"`
}

type CreateRoomMsg struct {
	Amount int64 `json:"amount"`
}

type RegisterHostMsg struct {
	IpAddress string `json:"ip"`
}
