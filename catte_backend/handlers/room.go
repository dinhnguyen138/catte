package handlers

import (
	"encoding/json"
	"fmt"

	"../constants"
	"../models"
	"github.com/firstrow/tcp_server"
)

type Room struct {
	Id             string `json:"id"`
	Players        []*Player
	finalRowPlayer int
	currentRow     int
	rowCount       int
	topCardIndex   int
	topCard        string
	Amount         int `json:"amount"`
}

func (room *Room) JoinRoom(userId string, data string, c *tcp_server.Client) {
	joint := false
	var playerInfo models.PlayerInfo
	json.Unmarshal([]byte(data), &playerInfo)

	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].Info.Id == userId {
			room.Players[i].Disconnected = false
			room.Players[i].client = c
			joint = true
			break
		}
	}
	if joint == false {
		var index int
		for i := 0; i < constants.MAXPLAYER; i++ {
			found := false
			for j := 0; j < len(room.Players); j++ {
				if room.Players[j].Index == i {
					found = true
					break
				}
			}
			if found == false {
				index = i
				break
			}
		}
		player := &Player{}
		player.Info = playerInfo
		player.Index = index
		player.client = c
		player.InGame = false
		player.Disconnected = false
		room.Players = append(room.Players, player)
		if len(room.Players) == 1 {
			room.Players[0].IsHost = true
		}

		// send list of current players to all player to update
		players, _ := json.Marshal(room.Players)
		fmt.Printf("%s\n", players)
		command := models.ResponseCommand{constants.PLAYERS, string(players)}
		room.Players[len(room.Players)-1].sendCommand(command)

		newPlayer, _ := json.Marshal(player)
		command = models.ResponseCommand{constants.NEWPLAYER, string(newPlayer)}
		for i := 0; i < len(room.Players)-1; i++ {
			room.Players[i].sendCommand(command)
		}
	}
}

func (room *Room) LeaveRoom(id string) {
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].Info.Id == id {
			changeHost := false
			if room.Players[i].IsHost == true {
				changeHost = true
			}

			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			if changeHost == true {
				room.Players[0].IsHost = true
			}

			command := models.ResponseCommand{constants.LEAVE, room.Players[i].Info.Id}
			fmt.Println(len(room.Players))
			for i := 0; i < len(room.Players); i++ {
				room.Players[i].sendCommand(command)
			}
			break
		}
	}
}

func (room *Room) HandleCommand(command models.Command) {
	switch command.Action {
	case constants.LEAVE:
		room.LeaveRoom(command.Id)
		break
	case constants.DEAL:
		room.newGame()
		break
	case constants.PLAY:
		room.play(command.Id, command.Data)
		break
	case constants.FOLD:
		room.fold(command.Id, command.Data)
		break
	}
}

func (room *Room) Disconnect(id string, kickUser func(roomId string, userId string)) {
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].Info.Id == id {
			fmt.Println("Find disconnected user " + id)
			room.Players[i].Disconnected = true
			kickUser(room.Id, id)
		}
	}

	// Start reconnect timer
}

func (room *Room) newGame() {
	room.finalRowPlayer = 0
	room.currentRow = 0
	room.rowCount = 0
	room.topCardIndex = -1
	room.topCard = ""
	for i := 0; i < len(room.Players); i++ {
		room.Players[i].InGame = true
		room.Players[i].Finalist = false
		room.Players[i].finalCard = ""
	}
	deck := deal()
	deck = shuffle(deck)
	pos := 0
	for i := 0; i < len(room.Players); i++ {
		slice := deck[pos : pos+6]
		room.Players[i].NumCard = 6
		pos += 7
		fmt.Println(slice)
		// send slice to client
		cards, _ := json.Marshal(slice)
		command := models.ResponseCommand{constants.CARDS, string(cards)}
		room.Players[i].sendCommand(command)
	}
}

func (room *Room) play(id string, card string) {
	room.rowCount++
	if room.currentRow < 5 {
		playData := models.PlayData{id, card}
		play, _ := json.Marshal(playData)
		command := models.ResponseCommand{constants.PLAY, string(play)}
		for i := 0; i < len(room.Players); i++ {
			if room.Players[i].Info.Id == id {
				room.Players[i].NumCard--
				room.topCard = card
				room.topCardIndex = i
			}
			room.Players[i].sendCommand(command)
		}
	} else {
		playData := models.PlayData{id, card}
		play, _ := json.Marshal(playData)
		command := models.ResponseCommand{constants.BACK, string(play)}
		for i := 0; i < len(room.Players); i++ {
			if room.Players[i].Info.Id == id {
				room.Players[i].finalCard = card
			}
			room.Players[i].sendCommand(command)
		}
	}

	fmt.Println(room.rowCount)

	// Send card play to all player
	if room.rowCount == len(room.Players) {
		if room.currentRow < 4 {
			room.rowCount = 0
			fmt.Println("XXXXXXX")
			room.currentRow++
			// Note that player is allow to go final row
			room.Players[room.topCardIndex].Finalist = true
			if room.currentRow == 4 {
				room.finalRowPlayer = 0
				// Inform player that is out
				for i := 0; i < len(room.Players); i++ {
					if room.Players[i] != nil {
						command := models.ResponseCommand{Action: constants.ELIMINATED}
						if room.Players[i].Finalist == false {
							room.Players[i].sendCommand(command)
						} else {
							room.finalRowPlayer++
						}
					}
				}
				fmt.Println("Final row player %v", room.finalRowPlayer)
			}
			// Inform lastRow top player to play
			fmt.Println("SSSSS")
			command := models.ResponseCommand{constants.STARTROW, room.Players[room.topCardIndex].Info.Id}
			for i := 0; i < len(room.Players); i++ {
				if room.Players[i] != nil {
					fmt.Println("Send STARTROW to ", i)
					room.Players[i].sendCommand(command)
				}
			}
		}
	}
	if room.rowCount == room.finalRowPlayer {
		if room.currentRow < 4 {
			return
		}
		if room.currentRow == 4 {
			room.rowCount = 0
			room.currentRow++
			// Send to player
			command := models.ResponseCommand{Action: constants.SHOWBACK}
			for i := 0; i < len(room.Players); i++ {
				if room.Players[i] != nil {
					room.Players[i].sendCommand(command)
				}
			}
		}
		if room.currentRow == 5 {
			room.topCard = room.Players[room.topCardIndex].finalCard
			// Calculate winner
			for i := 0; i < len(room.Players); i++ {
				if room.Players[i] != nil && room.Players[i].Finalist == true && larger(room.Players[i].finalCard, room.topCard) == true {
					room.topCardIndex = i
					room.topCard = room.Players[i].finalCard
				}
			}
			command := models.ResponseCommand{constants.WINNER, room.Players[room.topCardIndex].Info.Id}
			for i := 0; i < len(room.Players); i++ {
				if room.Players[i] != nil {
					room.Players[i].sendCommand(command)
				}
			}
		}
	}
}

func (room *Room) fold(id string, card string) {
	playData := models.PlayData{id, card}
	play, _ := json.Marshal(playData)
	command := models.ResponseCommand{constants.FOLD, string(play)}
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i] != nil {
			room.Players[i].sendCommand(command)
		}
	}
	room.rowCount++
	fmt.Println(room.rowCount)
	// Send card fold to all player
	if room.rowCount == len(room.Players) {
		if room.currentRow < 4 {
			room.rowCount = 0
			room.currentRow++
			room.Players[room.topCardIndex].Finalist = true
			// Note that player is allow to go final row
			if room.currentRow == 4 {
				// Inform player that is out
				room.finalRowPlayer = 0
				for i := 0; i < len(room.Players); i++ {
					command := models.ResponseCommand{Action: constants.ELIMINATED}
					if room.Players[i].Finalist == false {
						room.Players[i].sendCommand(command)
					} else {
						room.finalRowPlayer++
					}
				}
				fmt.Println("Final row player %v", room.finalRowPlayer)
			}
			// Inform lastRow top player to play
			command := models.ResponseCommand{constants.STARTROW, room.Players[room.topCardIndex].Info.Id}
			for i := 0; i < len(room.Players); i++ {
				room.Players[i].sendCommand(command)
			}
		}
	}

	if room.rowCount == room.finalRowPlayer {
		if room.currentRow < 4 {
			return
		}
		if room.currentRow == 4 {
			room.rowCount = 0
			room.currentRow++
			// Send to player
			command := models.ResponseCommand{Action: constants.SHOWBACK}
			for i := 0; i < len(room.Players); i++ {
				room.Players[i].sendCommand(command)
			}
		}
	}
}
