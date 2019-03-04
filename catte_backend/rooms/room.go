package rooms

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"../constants"
	"../models"
	"github.com/firstrow/tcp_server"
)

type Room struct {
	players        [6]*Player
	numPlayer      int
	finalRowPlayer int
	currentRow     int
	rowCount       int
	topCardIndex   int
	topCard        string
	hostIndex      int
}

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

var cardOrder map[string]int = map[string]int{
	"2":  1,
	"3":  2,
	"4":  3,
	"5":  4,
	"6":  5,
	"7":  6,
	"8":  7,
	"9":  8,
	"10": 9,
	"J":  10,
	"Q":  11,
	"K":  12,
	"A":  13,
}

type RoomManager struct {
	rooms map[string]*Room
}

func (player *Player) sendCommand(cmd models.ResponseCommand) {
	data, _ := json.Marshal(cmd)
	fmt.Println("sendCommand" + string(data))
	player.client.Send(string(data))
}

func (room *Room) JoinRoom(userId string, c *tcp_server.Client) {
	index := -1
	joint := false

	for i := 0; i < len(room.players); i++ {
		if room.players[i] == nil {
			index = i
		} else {
			fmt.Println(room.players[i].Id)
			if room.players[i].Id == userId {
				room.players[i].Disconnected = false
				room.players[i].client = c
				room.players[i].Index = i
				index = i
				joint = true
				break
			}
		}
	}
	if joint == false && index != -1 {
		player := &Player{}
		player.Id = userId
		player.Index = index
		player.client = c
		player.InGame = false
		player.Disconnected = false
		room.players[index] = player
		room.numPlayer++
		if room.numPlayer == 1 {
			room.hostIndex = index
		}
	}

	// send list of current players to all player to update
	players, _ := json.Marshal(room.players)
	fmt.Printf("%s\n", players)
	command := models.ResponseCommand{constants.PLAYERS, string(players)}
	for i := 0; i < len(room.players); i++ {
		if room.players[i] != nil {
			// Construct list of player
			room.players[i].sendCommand(command)
		}
	}

	fmt.Println(room)
}

func (room *Room) LeaveRoom(id string) (empty bool) {
	// changeHost := false
	// for i := 0; i < len(room.players); i++ {
	// 	if room.players[i].UserId == id {
	// 		if room.players[i].Index == room.hostIndex {
	// 			changeHost = true
	// 		}
	// 		room.players[i] = nil
	// 		room.numPlayer--
	// 	}
	// }
	// if room.numPlayer == 0 {
	// 	return true
	// } else if changeHost {
	// 	for i := 0; i < len(room.players); i++ {
	// 		if room.players[i] != nil {
	// 			room.hostIndex = i
	// 		}
	// 	}
	// }

	// // send list of current players to all player to update
	// players, _ := json.Marshal(room.players)
	// fmt.Printf("%s\n", players)
	// command := models.ResponseCommand{constants.PLAYERS, string(players)}
	// for i := 0; i < len(room.players); i++ {
	// 	if room.players[i] != nil {
	// 		room.players[i].sendCommand(command)
	// 	}
	// }
	return false
}

func (room *Room) HandleCommand(command models.Command) {
	switch command.Action {
	case constants.LEAVE:
		roomEmpty := room.LeaveRoom(command.Id)
		if roomEmpty {
			// Remove room
		}
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

func (roomManager *RoomManager) FindRoom(id string) (room *Room) {
	if roomManager.rooms == nil {
		roomManager.rooms = map[string]*Room{}
	}
	if _, ok := roomManager.rooms[id]; ok {
		return roomManager.rooms[id]
	}
	roomManager.rooms[id] = &Room{}
	return roomManager.rooms[id]
}

func (roomManager *RoomManager) FindClient(c *tcp_server.Client) (room *Room, index int) {
	for _, v := range roomManager.rooms {
		for i := 0; i < len(v.players); i++ {
			if v.players[i].client == c {
				return v, v.players[i].Index
			}
		}
	}
	return nil, 0
}

func (room *Room) Disconnect(index int) {
	room.players[index].Disconnected = true
	// Start reconnect timer
}

func (room *Room) newGame() {
	room.finalRowPlayer = 0
	room.currentRow = 0
	room.rowCount = 0
	room.topCardIndex = -1
	room.topCard = ""
	for i := 0; i < len(room.players); i++ {
		if room.players[i] != nil {
			room.players[i].InGame = true
			room.players[i].Finalist = false
			room.players[i].finalCard = ""
		}
	}
	fmt.Println("Room Player %v", room.numPlayer)
	deck := deal()
	deck = shuffle(deck)
	pos := 0
	for i := 0; i < len(room.players); i++ {
		if room.players[i] != nil {
			slice := deck[pos : pos+6]
			room.players[i].NumCard = 6
			pos += 7
			fmt.Println(slice)
			// send slice to client
			cards, _ := json.Marshal(slice)
			command := models.ResponseCommand{constants.CARDS, string(cards)}
			room.players[i].sendCommand(command)
		}
	}
}

func (room *Room) play(id string, card string) {
	room.rowCount++
	if room.currentRow < 5 {
		playData := models.PlayData{id, card}
		play, _ := json.Marshal(playData)
		command := models.ResponseCommand{constants.PLAY, string(play)}
		for i := 0; i < len(room.players); i++ {
			if room.players[i] != nil && room.players[i].Id == id {
				room.players[i].NumCard--
				room.topCard = card
				room.topCardIndex = i
			}
			if room.players[i] != nil {
				room.players[i].sendCommand(command)
			}
		}
	} else {
		playData := models.PlayData{id, card}
		play, _ := json.Marshal(playData)
		command := models.ResponseCommand{constants.BACK, string(play)}
		for i := 0; i < len(room.players); i++ {
			if room.players[i] != nil {
				if room.players[i].Id == id {
					room.players[i].finalCard = card
				}
				room.players[i].sendCommand(command)
			}
		}
	}

	fmt.Println(room.rowCount)

	// Send card play to all player
	if room.rowCount == room.numPlayer {
		if room.currentRow < 4 {
			room.rowCount = 0
			fmt.Println("XXXXXXX")
			room.currentRow++
			// Note that player is allow to go final row
			room.players[room.topCardIndex].Finalist = true
			if room.currentRow == 4 {
				room.finalRowPlayer = 0
				// Inform player that is out
				for i := 0; i < len(room.players); i++ {
					if room.players[i] != nil {
						command := models.ResponseCommand{Action: constants.ELIMINATED}
						if room.players[i].Finalist == false {
							room.players[i].sendCommand(command)
						} else {
							room.finalRowPlayer++
						}
					}
				}
				fmt.Println("Final row player %v", room.finalRowPlayer)
			}
			// Inform lastRow top player to play
			fmt.Println("SSSSS")
			command := models.ResponseCommand{constants.STARTROW, room.players[room.topCardIndex].Id}
			for i := 0; i < len(room.players); i++ {
				if room.players[i] != nil {
					fmt.Println("Send STARTROW to ", i)
					room.players[i].sendCommand(command)
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
			for i := 0; i < len(room.players); i++ {
				if room.players[i] != nil {
					room.players[i].sendCommand(command)
				}
			}
		}
		if room.currentRow == 5 {
			room.topCard = room.players[room.topCardIndex].finalCard
			// Calculate winner
			for i := 0; i < len(room.players); i++ {
				if room.players[i] != nil && room.players[i].Finalist == true && larger(room.players[i].finalCard, room.topCard) == true {
					room.topCardIndex = i
					room.topCard = room.players[i].finalCard
				}
			}
			command := models.ResponseCommand{constants.WINNER, room.players[room.topCardIndex].Id}
			for i := 0; i < len(room.players); i++ {
				if room.players[i] != nil {
					room.players[i].sendCommand(command)
				}
			}
		}
	}
}

func (room *Room) fold(id string, card string) {
	playData := models.PlayData{id, card}
	play, _ := json.Marshal(playData)
	command := models.ResponseCommand{constants.FOLD, string(play)}
	for i := 0; i < len(room.players); i++ {
		if room.players[i] != nil {
			room.players[i].sendCommand(command)
		}
	}
	room.rowCount++
	fmt.Println(room.rowCount)
	// Send card fold to all player
	if room.rowCount == room.numPlayer {
		if room.currentRow < 4 {
			room.rowCount = 0
			room.currentRow++
			room.players[room.topCardIndex].Finalist = true
			fmt.Println("YYYYY")
			// Note that player is allow to go final row
			if room.currentRow == 4 {
				// Inform player that is out
				room.finalRowPlayer = 0
				for i := 0; i < len(room.players); i++ {
					command := models.ResponseCommand{Action: constants.ELIMINATED}
					if room.players[i] != nil {
						if room.players[i].Finalist == false {
							room.players[i].sendCommand(command)
						} else {
							room.finalRowPlayer++
						}
					}
				}
				fmt.Println("Final row player %v", room.finalRowPlayer)
			}
			fmt.Println("SSSS")
			// Inform lastRow top player to play
			command := models.ResponseCommand{constants.STARTROW, room.players[room.topCardIndex].Id}
			for i := 0; i < len(room.players); i++ {
				if room.players[i] != nil {
					fmt.Println("Send STARTROW to ", i)
					room.players[i].sendCommand(command)
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
			for i := 0; i < len(room.players); i++ {
				if room.players[i] != nil {
					room.players[i].sendCommand(command)
				}
			}
		}
	}
}

func deal() (deck models.Deck) {
	// Valid types include Two, Three, Four, Five, Six
	// Seven, Eight, Nine, Ten, Jack, Queen, King & Ace
	types := []string{"2", "3", "4", "5", "6", "7",
		"8", "9", "10", "J", "Q", "K", "A"}

	// Valid suits include Heart, Diamond, Club & Spade
	suits := []string{"H", "D", "C", "S"}

	// Loop over each type and suit appending to the deck
	for i := 0; i < len(types); i++ {
		for n := 0; n < len(suits); n++ {
			deck = append(deck, types[i]+suits[n])
		}
	}
	return
}

// Shuffle the deck
func shuffle(d models.Deck) models.Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) { d[i], d[j] = d[j], d[i] })
	return d
}

func larger(leftCard string, rightCard string) bool {
	leftValue := leftCard[:len(leftCard)-1]
	leftType := leftCard[len(leftCard)-1:]
	rightValue := rightCard[:len(rightCard)-1]
	rightType := rightCard[len(rightCard)-1:]
	if leftType == rightType {
		return cardOrder[leftValue] > cardOrder[rightValue]
	}
	return false
}
