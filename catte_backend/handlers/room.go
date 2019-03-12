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
	Amount         int
	IndexUsed      []bool
	finalRowPlayer int
	currentRow     int
	rowCount       int
	topCardIndex   int
	topCard        string
	showCard       string
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
	case constants.PLAY, constants.FOLD:
		fmt.Println("XXXX")
		room.play(command.Action, command.Index, command.Data)
		break
	}
}

func (room *Room) KickUserCallback(callback func(roomId string, index int)) {
	room.kickUser = callback
}

func (room *Room) SendBroadcast(command string, data interface{}) {
	var stringData string
	switch data.(type) {
	case string:
		stringData = data.(string)
		break
	default:
		temp, _ := json.Marshal(data)
		stringData = string(temp)
	}
	message := models.ResponseCommand{command, string(stringData)}
	for i := 0; i < len(room.Players); i++ {
		room.Players[i].sendCommand(message)
	}
}

func (room *Room) SendUnicast(index int, command string, data interface{}) {
	var stringData string
	switch data.(type) {
	case string:
		stringData = data.(string)
		break
	default:
		temp, _ := json.Marshal(data)
		stringData = string(temp)
	}
	message := models.ResponseCommand{command, string(stringData)}
	room.Players[index].sendCommand(message)
}

func (room *Room) JoinRoom(userId string, data string, c *tcp_server.Client) {
	index := -1
	var playerInfo models.PlayerInfo
	json.Unmarshal([]byte(data), &playerInfo)

	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].Info.Id == userId {
			fmt.Println("Player reconnected " + userId)
			room.Players[i].Disconnected = false
			room.Players[i].client = c
			index = i
			break
		}
	}
	if index == -1 {
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
		room.SendUnicast(len(room.Players)-1, constants.PLAYERS, room.Players)
		for i := 0; i < len(room.Players)-1; i++ {
			room.SendUnicast(i, constants.NEWPLAYER, player)
		}
	} else {
		room.SendUnicast(index, constants.PLAYERS, room.Players)
		if room.inGame == true {
			room.SendUnicast(index, constants.CARDS, room.Players[index].cards)
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

			room.SendBroadcast(constants.LEAVE, strconv.Itoa(index))
			// command := models.ResponseCommand{constants.LEAVE, strconv.Itoa(index)}
			// fmt.Println(len(room.Players))
			// for i := 0; i < len(room.Players); i++ {
			// 	room.Players[i].sendCommand(command)
			// }
			break
		}
	}
	if len(room.Players) == 0 {
		room.timer.Stop()
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
	for len(room.Players) != 0 {
		select {
		case <-room.timer.C:
			fmt.Print("Timeout at turn ")
			fmt.Println(room.rowCount)
			if room.inGame == true {
				if room.rowCount == 0 {
					for i := 0; i < len(room.Players); i++ {
						if room.Players[i].Index == room.turn {
							room.play(constants.PLAY, room.turn, room.Players[i].cards[0])
							break
						}
					}
				} else {
					for i := 0; i < len(room.Players); i++ {
						if room.Players[i].Index == room.turn {
							room.play(constants.FOLD, room.turn, room.Players[i].cards[0])
							break
						}
					}
				}
			} else {
				room.newGame()
			}

		}
	}

}

func (room *Room) newGame() {
	if room.inGame == true {
		return
	}
	room.turn = room.Players[room.topCardIndex].Index
	room.finalRowPlayer = 0
	room.currentRow = 0
	room.rowCount = 0
	room.topCardIndex = -1
	room.topCard = ""
	for i := 0; i < len(room.Players); i++ {
		room.Players[i].InGame = true
		room.Players[i].Finalist = false
	}
	deck := deal()
	deck = shuffle(deck)
	pos := 0
	for i := 0; i < len(room.Players); i++ {
		slice := deck[pos : pos+6]
		room.Players[i].NumCard = 6
		room.Players[i].cards = slice
		pos += 7
		fmt.Println(slice)
		// send slice to client
		room.SendUnicast(i, constants.CARDS, slice)
		room.SendUnicast(i, constants.START, strconv.Itoa(room.turn))
	}
	room.inGame = true
	room.timer = time.NewTimer(time.Second * 20)
	go room.mainloop()
}

func (room *Room) play(action string, index int, card string) {
	room.timer.Stop()
	room.rowCount++
	fmt.Println(room.rowCount)

	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].Index == index {
			found := room.Players[i].playCard(card)
			if found == true {
				if action == constants.PLAY && larger(card, room.topCard) {
					room.topCard = card
					room.topCardIndex = i
				} else {
					action = constants.FOLD
				}
			} else {
				card = room.Players[i].cards[0]
				action = constants.FOLD
				room.Players[i].playCard(card)
			}
		}
	}
	row := room.currentRow
	eliminatedPlayers := []int{}
	room.turn = room.getNext(room.turn)

	if room.currentRow == 4 && room.rowCount == room.finalRowPlayer {
		room.rowCount = 0
		room.currentRow++
		room.turn = room.topCardIndex
		// Reset top card to empty after row
		room.showCard = room.topCard
		room.topCard = ""
	}
	// Send card play to all player
	if room.currentRow < 4 && room.rowCount == len(room.Players) {
		room.rowCount = 0
		room.currentRow++
		// Reset top card to empty after row
		room.topCard = ""
		// Note that player is allow to go final row
		room.Players[room.topCardIndex].Finalist = true
		if room.currentRow == 4 {
			room.finalRowPlayer = 0
			// Inform player that is out
			for i := 0; i < len(room.Players); i++ {
				if room.Players[i].Finalist == false {
					eliminatedPlayers = append(eliminatedPlayers, room.Players[i].Index)
				} else {
					room.finalRowPlayer++
				}
			}
		}
		// Inform lastRow top player to play
		room.turn = room.topCardIndex
	}

	if room.inGame == true {
		playData := models.PlayData{index, row, room.turn, (row != room.currentRow), card}
		room.SendBroadcast(action, playData)

		if len(eliminatedPlayers) > 0 {
			room.SendBroadcast(constants.ELIMINATED, eliminatedPlayers)
		}
	}

	if room.finalRowPlayer == 1 {
		room.SendBroadcast(constants.WINNER, strconv.Itoa(room.Players[room.topCardIndex].Index))
		room.calculateWinnerAmount(true)
		room.inGame = false
	}

	if room.currentRow == 5 && room.rowCount == room.finalRowPlayer {
		room.rowCount = 0
		room.currentRow = 0
		room.turn = room.topCardIndex
		// Reset top card to empty after row
		room.topCard = ""
		room.SendBroadcast(constants.WINNER, strconv.Itoa(room.Players[room.topCardIndex].Index))
		room.calculateWinnerAmount(room.isDouble())
		room.inGame = false
	}

	room.resetTimer()
}

func (room *Room) getNext(index int) int {
	for i := 1; i < len(room.Players); i++ {
		var j = (index + i) % len(room.Players)
		fmt.Println()
		if room.IndexUsed[j] == true {
			return j
		}
	}
	return -1
}

func (room *Room) calculateWinnerAmount(isDouble bool) {
	amount := float32(room.Amount)
	if isDouble == true {
		amount *= 2
	}
	totalWin := float32(0)
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i].Index != room.topCardIndex {
			if room.Players[i].Info.Amount < amount {
				totalWin += room.Players[i].Info.Amount
				room.Players[i].Info.Amount = 0
			} else {
				totalWin += amount
				room.Players[i].Info.Amount -= amount
			}
		}
	}
	room.Players[room.topCardIndex].Info.Amount += totalWin * 0.8
	// Update db
	// Send result
}

func (room *Room) resetTimer() {
	var duration time.Duration
	if room.inGame == true {
		duration = time.Duration(20)
	} else {
		duration = time.Duration(30)
	}
	room.timer.Reset(time.Second * time.Duration(duration))
}

func (room *Room) isDouble() bool {
	leftValue := room.showCard[:len(room.showCard)-1]
	rightValue := room.topCard[:len(room.topCard)-1]
	if constants.CardOrder[leftValue]+constants.CardOrder[rightValue] <= 3 {
		return true
	} else {
		return false
	}
}
