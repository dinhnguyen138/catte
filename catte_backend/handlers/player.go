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
	Cards        []string
	Commands     []models.ResponseCommand
	client       *tcp_server.Client
}

func (player *Player) sendCommand(cmd models.ResponseCommand) {
	player.Commands = append(player.Commands, cmd)
	if player.Disconnected == false {
		for i := 0; i < len(player.Commands); i++ {
			data, _ := json.Marshal(player.Commands[i])
			fmt.Println("sendCommand" + string(data))
			player.client.Send(string(data) + "\n")
		}
		player.Commands = player.Commands[:0]
	}
}
