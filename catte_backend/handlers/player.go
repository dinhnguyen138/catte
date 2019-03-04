package handlers

import (
	"encoding/json"
	"fmt"

	"../models"
	"github.com/firstrow/tcp_server"
)

type Player struct {
	Id           string `json:"id"`
	NumCard      int    `json:"numcard"`
	Index        int    `json:"index"`
	InGame       bool   `json:"ingame"`
	Finalist     bool   `json:"finalist"`
	finalCard    string
	Disconnected bool `json:"disconnected"`
	client       *tcp_server.Client
}

func (player *Player) sendCommand(cmd models.ResponseCommand) {
	data, _ := json.Marshal(cmd)
	fmt.Println("sendCommand" + string(data))
	if player.Disconnected == false {
		player.client.Send(string(data) + "\n")
	}
}
