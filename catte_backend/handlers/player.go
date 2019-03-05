package handlers

import (
	"encoding/json"
	"fmt"

	"../models"
	"github.com/firstrow/tcp_server"
)

type Player struct {
	Info         models.PlayerInfo `json:"info"`
	NumCard      int               `json:"numcard"`
	Index        int               `json:"index"`
	InGame       bool              `json:"ingame"`
	Finalist     bool              `json:"finalist"`
	IsHost       bool              `json:"isHost"`
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
