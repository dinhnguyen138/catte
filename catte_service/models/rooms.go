package models

type Room struct {
	Id       string `json:"id" form:"-"`
	IsActive bool   `json:"isactive" form:"isactive"`
	NoPlayer int    `json:"noplayer" form:"noplayer"`
	Amount   int    `json:"amount"`
}
