package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"../constants"
	"../models"
	"github.com/firstrow/tcp_server"
)

type Room struct {
	Id             string `json:"id"`
	Players        []*Player
	Amount         int `json:"amount"`
	IndexUsed      []bool
	finalRowPlayer int
	currentRow     int
	rowCount       int
	topCardIndex   int
	topCard        string
	timer          *time.Timer
	turn           int
	inGame         bool
	kickUser       func(roomId string, index int)
}

func (room *Room) HandleCommand(command models.Command) {
	switch command.Action {
	case constants.LEAVE:
		room.LeaveRoom(command.Index)
		break
	case constants.DEAL:
		room.newGame()
		break
	case constants.PLAY:
		room.play(command.Index, command.Data)
		break
	case constants.FOLD:
		room.fold(command.Index, command.Data)
		break
	}
}

func (room *Room) KickUserCallback(callback func(roomId string, index int)) {
	room.kickUser = callback
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
		for i := 0; i < len(room.IndexUsed); i++ {
			if room.IndexUsed[i] == false {
				index = i
				break
			}
		}
		room.IndexUsed[index] = true
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

func (room *Room) LeaveRoom(index int) {
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].Index == index {
			changeHost := false
			if room.Players[i].IsHost == true {
				changeHost = true
			}

			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			room.IndexUsed[room.Players[i].Index] = false
			if changeHost == true {
				room.Players[0].IsHost = true
			}

			command := models.ResponseCommand{constants.LEAVE, strconv.Itoa(index)}
			fmt.Println(len(room.Players))
			for i := 0; i < len(room.Players); i++ {
				room.Players[i].sendCommand(command)
			}
			break
		}
	}
}

func (room *Room) Disconnect(index int) {
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].Index == index {
			fmt.Println("Find disconnected user " + room.Players[i].Info.Id)
			room.Players[i].Disconnected = true
		}
	}
}

func (room *Room) KickDisconnectedUser() {
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].Disconnected == true {
			room.kickUser(room.Id, room.Players[i].Index)
		}
	}
}

func (room *Room) mainloop() {
	for room.inGame == true {
		select {
		case <-room.timer.C:
			fmt.Println(room.turn)
		}
	}
}

func (room *Room) newGame() {
	room.turn = room.Players[room.topCardIndex].Index
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
		room.Players[i].Cards = slice
		pos += 7
		fmt.Println(slice)
		// send slice to client
		cards, _ := json.Marshal(slice)
		command := models.ResponseCommand{constants.CARDS, string(cards)}
		room.Players[i].sendCommand(command)

		command = models.ResponseCommand{constants.STARTROW, strconv.Itoa(room.turn)}
		room.Players[i].sendCommand(command)
	}
	room.inGame = true
	room.timer = time.NewTimer(time.Second * 30)
	go room.mainloop()
}

func (room *Room) play(index int, card string) {
	room.timer.Stop()
	room.rowCount++

	playData := models.PlayData{index, card}
	play, _ := json.Marshal(playData)
	command := models.ResponseCommand{constants.PLAY, string(play)}
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].Index == index {
			room.Players[i].NumCard--
			room.topCard = card
			room.topCardIndex = i
			for j := 0; j < len(room.Players[i].Cards); j++ {
				if room.Players[i].Cards[j] == card {
					room.Players[i].Cards = append(room.Players[i].Cards[:j], room.Players[i].Cards[j+1:]...)
				}
			}
		}
		room.Players[i].sendCommand(command)
	}

	// Send card play to all player
	if room.rowCount == len(room.Players) {
		if room.currentRow < 4 {
			room.rowCount = 0
			room.currentRow++
			// Note that player is allow to go final row
			room.Players[room.topCardIndex].Finalist = true
			if room.currentRow == 4 {
				room.finalRowPlayer = 0
				// Inform player that is out
				eliminatedPlayers := []int{}
				for i := 0; i < len(room.Players); i++ {
					if room.Players[i].Finalist == false {
						eliminatedPlayers = append(eliminatedPlayers, room.Players[i].Index)
					} else {
						room.finalRowPlayer++
					}
				}
				if len(eliminatedPlayers) > 0 {
					data, _ := json.Marshal(eliminatedPlayers)
					command := models.ResponseCommand{constants.ELIMINATED, string(data)}
					for i := 0; i < len(room.Players); i++ {
						room.Players[i].sendCommand(command)
					}
				}

				fmt.Println("Final row player %v", room.finalRowPlayer)
			}
			// Inform lastRow top player to play
			fmt.Println("SSSSS")
			command := models.ResponseCommand{constants.STARTROW, strconv.Itoa(room.Players[room.topCardIndex].Index)}
			for i := 0; i < len(room.Players); i++ {
				fmt.Println("Send STARTROW to ", i)
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
			command := models.ResponseCommand{constants.STARTROW, ""}
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
			command := models.ResponseCommand{constants.WINNER, strconv.Itoa(room.Players[room.topCardIndex].Index)}
			for i := 0; i < len(room.Players); i++ {
				if room.Players[i] != nil {
					room.Players[i].sendCommand(command)
				}
			}
		}
	}

	for {
		i := 1
		if room.IndexUsed[(room.turn+i)%6] == true {
			room.turn = (room.turn + i) % 6
			room.timer.Reset(time.Second * 30)
		}
		i++
	}
}

func (room *Room) fold(index int, card string) {
	room.timer.Stop()
	room.rowCount++
	playData := models.PlayData{index, card}
	play, _ := json.Marshal(playData)
	command := models.ResponseCommand{constants.FOLD, string(play)}
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i] != nil {
			room.Players[i].sendCommand(command)
		}
	}

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
			command := models.ResponseCommand{constants.STARTROW, strconv.Itoa(room.Players[room.topCardIndex].Index)}
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
			command := models.ResponseCommand{constants.STARTROW, ""}
			for i := 0; i < len(room.Players); i++ {
				room.Players[i].sendCommand(command)
			}
		}
	}

	for {
		i := 1
		if room.IndexUsed[(room.turn+i)%6] == true {
			room.turn = (room.turn + i) % 6
			room.timer.Reset(time.Second * 30)
		}
		i++
	}
}
