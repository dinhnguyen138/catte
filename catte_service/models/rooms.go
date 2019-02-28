package models

type Room struct {
	UUID     string `json:"uuid" form:"-"`
	IsActive bool   `json:"isactive" form:"isactive"`
	NoPlayer int    `json:"noplayer" form:"noplayer"`
}
