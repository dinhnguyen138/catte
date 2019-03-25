package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dinhnguyen138/catte/catte_backend/constants"
	"github.com/dinhnguyen138/catte/catte_backend/db"
	"github.com/dinhnguyen138/catte/catte_backend/models"
	"github.com/dinhnguyen138/catte/catte_backend/utilities"
	"github.com/dinhnguyen138/tcp_server"
	"github.com/kataras/golog"
)

type Room struct {
	id             string
	players        []*Player
	amount         int64
	indexUsed      []bool
	finalRowPlayer int
	currentRow     int
	rowCount       int
	topCardIndex   int
	topCard        string
	showCard       string
	timer          *time.Timer
	turn           int
	inGame         bool
	cleanUp        bool
	maxPlayer      int
	kickUser       func(roomId string, index int)
}

var turnTimeout = 12
var cleanUpTimeout = 5
var newGameTimeout = 15

func NewRoom(id string, maxPlayer int, amount int64) *Room {
	room := new(Room)
	room.id = id
	room.amount = amount
	room.maxPlayer = maxPlayer
	room.inGame = false
	room.cleanUp = false
	for i := 0; i < room.maxPlayer; i++ {
		room.indexUsed = append(room.indexUsed, false)
	}
	return room
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
		room.play(command.Action, command.Index, command.Data, false)
		break
	}
}

func (room *Room) KickUserCallback(callback func(roomId string, index int)) {
	room.kickUser = callback
}

func (room *Room) sendBroadcast(command string, data interface{}) {
	var stringData string
	switch data.(type) {
	case string:
		stringData = data.(string)
		break
	default:
		temp, _ := json.Marshal(data)
		stringData = string(temp)
	}
	golog.Info("Broadcast message" + string(stringData))
	message := models.ResponseCommand{command, string(stringData)}
	for i := 0; i < len(room.players); i++ {
		room.players[i].sendCommand(message)
	}
}

func (room *Room) sendUnicast(index int, command string, data interface{}) {
	var stringData string
	switch data.(type) {
	case string:
		stringData = data.(string)
		break
	default:
		temp, _ := json.Marshal(data)
		stringData = string(temp)
	}
	golog.Infof("Unicast message to %v", index)
	golog.Info("Message content: " + string(stringData))
	message := models.ResponseCommand{command, string(stringData)}
	room.players[index].sendCommand(message)
}

func (room *Room) JoinRoom(userId string, data string, c *tcp_server.Client) {
	index := -1
	var playerInfo models.PlayerInfo
	json.Unmarshal([]byte(data), &playerInfo)

	for i := 0; i < len(room.players); i++ {
		// Check whether this is a new user or a disconnected user
		if room.players[i].Info.Id == userId {
			golog.Info("Player reconnect, update info")
			room.players[i].Disconnected = false
			room.players[i].client = c
			index = i
			break
		}
	}
	// New player
	if index == -1 {
		golog.Info("New player join room")
		if len(room.players) == room.maxPlayer {
			// Room is full, return error
			golog.Error("Maximum player amount is reach, send back error")
			utilities.SendClient(c, constants.ERROR, strconv.Itoa(constants.ERR_ROOM_FULL))
			return
		}
		var index int
		for i := 0; i < len(room.indexUsed); i++ {
			// Find an empty seat for new player
			if room.indexUsed[i] == false {
				index = i
				break
			}
		}
		// Create player
		room.indexUsed[index] = true
		player := newPlayer(playerInfo, index, c)
		room.players = append(room.players, player)
		if len(room.players) == 1 {
			// Update room's host
			room.players[0].IsHost = true
		}

		// Send list of current players to new player to update
		room.sendUnicast(len(room.players)-1, constants.PLAYERS, room.players)

		// Send new player to current players to update
		for i := 0; i < len(room.players)-1; i++ {
			room.sendUnicast(i, constants.NEWPLAYER, player)
		}

		// Update number player in room
		db.UpdateRoom(room.id, len(room.players))
		if len(room.players) == 2 {
			golog.Info("Start game main loop")
			room.timer = time.NewTimer(time.Second)
			go room.mainloop()
		}
	} else {
		// Send back list of players to current players
		room.sendUnicast(index, constants.PLAYERS, room.players)
		if room.inGame == true {
			// Send their current cards
			room.sendUnicast(index, constants.CARDS, room.players[index].cards)
		}
	}
}

func (room *Room) LeaveRoom(index int) {
	for i := 0; i < len(room.players); i++ {
		if room.players[i].Index == index {
			// Receive leave command when in game, rejected
			if room.inGame == true {
				golog.Error("Received LEAVE while in game, rejected")
				room.sendUnicast(index, constants.ERROR, strconv.Itoa(constants.ERR_LEAVE_IN_GAME))
				break
			}

			// Update host if leave player is host
			msg := models.LeaveMsg{Index: index}
			if room.players[i].IsHost == true {
				room.players[i].IsHost = false
				for j := 0; j < len(room.players); j++ {
					if i == j {
						continue
					}
					room.players[j].IsHost = true
					msg.Host = room.players[j].Index
					break
				}
			} else {
				msg.Host = -1
			}

			// Send leave message to all players
			room.sendBroadcast(constants.LEAVE, msg)

			// Update players in DB
			room.indexUsed[room.players[i].Index] = false
			room.players = append(room.players[:i], room.players[i+1:]...)
			db.UpdateRoom(room.id, len(room.players))
			break
		}
	}
	if len(room.players) == 1 {
		room.inGame = false
		room.cleanUp = false
		room.timer.Stop()
	}
}

func (room *Room) FindClient(c *tcp_server.Client) int {
	for i := 0; i < len(room.players); i++ {
		if room.players[i].client == c {
			return room.players[i].Index
		}
	}
	return -1
}

func (room *Room) IsEmpty() bool {
	return len(room.players) == 0
}

func (room *Room) Disconnect(index int) {
	for i := 0; i < len(room.players); i++ {
		if room.players[i].Index == index {
			golog.Info("Find disconnected user " + room.players[i].Info.Id)
			room.players[i].Disconnected = true
			if room.inGame == false {
				room.kickUser(room.id, room.players[i].Index)
				return
			}
		}
	}
}

func (room *Room) KickUsers() {
	for i := 0; i < len(room.players); i++ {
		if room.players[i].Disconnected == true || room.players[i].isInactive == true {
			room.kickUser(room.id, room.players[i].Index)
		}
	}
	room.cleanUp = true
}

func (room *Room) mainloop() {
	for len(room.players) > 1 {
		select {
		case <-room.timer.C:
			golog.Infof("Timeout at turn %v", room.turn)
			if room.inGame == true {
				if room.rowCount == 0 {
					for i := 0; i < len(room.players); i++ {
						if room.players[i].Index == room.turn {
							room.play(constants.PLAY, room.turn, room.players[i].cards[0], true)
							break
						}
					}
				} else {
					for i := 0; i < len(room.players); i++ {
						if room.players[i].Index == room.turn {
							room.play(constants.FOLD, room.turn, room.players[i].cards[0], true)
							break
						}
					}
				}
			} else if room.cleanUp == false {
				golog.Info("Kick disconnect and inactive users")
				room.KickUsers()
				// TODO: Send broadcast to user to info automatic newgame
				if len(room.players) > 1 {
					golog.Info("Inform users prepare for new game")
					room.sendBroadcast(constants.INFORM, "")
					room.resetTimer(newGameTimeout)
				}
			} else {
				if len(room.players) > 1 {
					golog.Info("Start new game due to timeout")
					room.newGame()
				}
			}

		}
	}

}

func (room *Room) reset() {
	room.turn = room.players[room.topCardIndex].Index
	room.finalRowPlayer = 0
	room.currentRow = 0
	room.rowCount = 0
	room.topCardIndex = -1
	room.topCard = ""
	for i := 0; i < len(room.players); i++ {
		room.players[i].InGame = true
		room.players[i].Finalist = false
	}
}

func (room *Room) newGame() {
	if room.inGame == true || len(room.players) == 0 {
		return
	}
	golog.Info("Reset room for new game")
	room.reset()

	deck := deal()
	deck = shuffle(deck)
	pos := 0
	for i := 0; i < len(room.players); i++ {
		slice := deck[pos : pos+6]
		room.players[i].NumCard = 6
		room.players[i].cards = slice
		pos += 7
		fmt.Println(slice)
		// send slice to client
		room.sendUnicast(i, constants.CARDS, slice)
		room.sendUnicast(i, constants.START, strconv.Itoa(room.turn))
	}
	room.inGame = true
	room.cleanUp = false
	room.resetTimer(turnTimeout)
}

func (room *Room) play(action string, index int, card string, isAuto bool) {
	room.timer.Stop()
	room.rowCount++
	fmt.Println(room.rowCount)

	for i := 0; i < len(room.players); i++ {
		if room.players[i].Index == index {
			found := room.players[i].playCard(card)
			if found == true {
				if action == constants.PLAY && larger(card, room.topCard) {
					room.topCard = card
					room.topCardIndex = i
				} else {
					action = constants.FOLD
				}
			} else {
				card = room.players[i].cards[0]
				action = constants.FOLD
				room.players[i].playCard(card)
			}
			room.players[i].isInactive = isAuto
		}
	}
	row := room.currentRow
	eliminatedplayers := []int{}
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
	if room.currentRow < 4 && room.rowCount == len(room.players) {
		room.rowCount = 0
		room.currentRow++
		// Reset top card to empty after row
		room.topCard = ""
		// Note that player is allow to go final row
		room.players[room.topCardIndex].Finalist = true
		if room.currentRow == 4 {
			room.finalRowPlayer = 0
			// Inform player that is out
			for i := 0; i < len(room.players); i++ {
				if room.players[i].Finalist == false {
					eliminatedplayers = append(eliminatedplayers, room.players[i].Index)
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
		room.sendBroadcast(action, playData)

		if len(eliminatedplayers) > 0 {
			room.sendBroadcast(constants.ELIMINATED, eliminatedplayers)
		}
	}

	if room.finalRowPlayer == 1 {
		room.calculateWinnerAmount(true)
		room.inGame = false
	} else if room.currentRow == 5 && room.rowCount == room.finalRowPlayer {
		room.rowCount = 0
		room.currentRow = 0
		room.turn = room.topCardIndex
		// Reset top card to empty after row
		room.calculateWinnerAmount(room.isDouble())
		room.inGame = false
	}

	if room.inGame == true {
		room.resetTimer(turnTimeout)
	} else {
		room.resetTimer(cleanUpTimeout)
	}
}

func (room *Room) getNext(index int) int {
	for i := 1; i < len(room.players); i++ {
		var j = (index + i) % len(room.players)
		fmt.Println()
		if room.indexUsed[j] == true {
			return j
		}
	}
	return -1
}

func (room *Room) calculateWinnerAmount(isDouble bool) {
	var result []models.ResultMsg
	amount := room.amount
	if isDouble == true {
		amount *= 2
	}
	totalWin := int64(0)
	for i := 0; i < len(room.players); i++ {
		if room.players[i].Index != room.topCardIndex {
			if room.players[i].Info.Amount < amount {
				totalWin += room.players[i].Info.Amount
				result = append(result, models.ResultMsg{room.players[i].Index, -room.players[i].Info.Amount, room.players[i].Info.Amount})
				room.players[i].Info.Amount = 0
				db.UpdatePlayer(room.players[i].Info.Id, room.players[i].Info.Amount)
			} else {
				totalWin += amount
				room.players[i].Info.Amount -= amount
				result = append(result, models.ResultMsg{room.players[i].Index, -amount, room.players[i].Info.Amount})
				db.UpdatePlayer(room.players[i].Info.Id, room.players[i].Info.Amount)
			}
		}
	}
	room.players[room.topCardIndex].Info.Amount += (int64)((float64)(totalWin) * 0.8)
	db.UpdatePlayer(room.players[room.topCardIndex].Info.Id, room.players[room.topCardIndex].Info.Amount)
	result = append(result, models.ResultMsg{room.players[room.topCardIndex].Index, totalWin, room.players[room.topCardIndex].Info.Amount})

	room.sendBroadcast(constants.RESULT, result)
}

func (room *Room) resetTimer(timeout int) {
	duration := time.Duration(timeout)
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
