package models

type Room struct {
	Id        string `json:"id"`
	NoPlayer  int    `json:"noplayer"`
	MaxPlayer int    `json:"maxplayer"`
	Amount    int64  `json:"amount"`
}
