package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/dinhnguyen138/catte/catte_backend/models"
	"github.com/firstrow/tcp_server"
)

type Player struct {
	Info         models.PlayerInfo `json:"info"`
	NumCard      int               `json:"numcard"`
	Index        int               `json:"index"`
	InGame       bool              `json:"ingame"`
	Finalist     bool              `json:"finalist"`
	IsHost       bool              `json:"ishost"`
	Disconnected bool              `json:"disconnected"`
	cards        []string
	client       *tcp_server.Client
}

func (player *Player) sendCommand(cmd models.ResponseCommand) {
	data, _ := json.Marshal(cmd)
	fmt.Println("sendCommand" + string(data))
	if player.Disconnected == false {
		player.client.Send(string(data) + "\n")
	}
}

func (player *Player) playCard(card string) bool {
	found := false
	for i := 0; i < len(player.cards); i++ {
		if player.cards[i] == card {
			found = true
			player.cards = append(player.cards[:i], player.cards[i+1:]...)
			player.NumCard--
		}
	}
	return found
}
