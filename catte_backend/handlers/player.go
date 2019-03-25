package handlers

import (
	"encoding/json"

	"github.com/dinhnguyen138/catte/catte_backend/models"
	"github.com/dinhnguyen138/tcp_server"
	"github.com/kataras/golog"
)

type Player struct {
	Info         models.PlayerInfo `json:"info"`
	NumCard      int               `json:"numcard"`
	Index        int               `json:"index"`
	InGame       bool              `json:"ingame"`
	Finalist     bool              `json:"finalist"`
	IsHost       bool              `json:"ishost"`
	Disconnected bool              `json:"disconnected"`
	isInactive   bool
	cards        []string
	client       *tcp_server.Client
}

func (player *Player) sendCommand(cmd models.ResponseCommand) {
	data, _ := json.Marshal(cmd)
	golog.Info("sendCommand" + string(data))
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

func newPlayer(info models.PlayerInfo, index int, c *tcp_server.Client) *Player {
	player := new(Player)
	player.Info = info
	player.Index = index
	player.InGame = false
	player.Disconnected = false
	player.isInactive = false
	player.client = c
	return player
}
