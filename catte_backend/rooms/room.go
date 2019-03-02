package rooms

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"../constants"
	"../models"
	"github.com/firstrow/tcp_server"
)

type Room struct {
	players        [6]*Player
	numPlayer      int
	finalRowPlayer int
	currentRow     int
	rowCount int
	normalRowCount int
	finalRowCount  int
	topCardIndex   int
	topCard        string
	hostIndex      int
}

type Player struct {
	UserId       string `json:"uuid"`
	NumCard      int    `json:"numcard"`
	Index        int    `json:"index"`
	InGame       bool   `json:"ingame"`
	Finalist     bool   `json:"finalist"`
	Disconnected bool   `json:"disconnected"`
	client       *tcp_server.Client
}

type RoomManager struct {
	rooms map[string]*Room
}

func (player *Player) sendCommand(cmd models.ResponseCommand) {
	data, _ := json.Marshal(cmd)
	player.client.Send(string(data) + "\n")
}

func (room *Room) JoinRoom(userId string, c *tcp_server.Client) {
	index := -1
	joint := false

	for i := 0; i < len(room.players); i++ {
		if room.players[i] == nil {
			index = i
		} else {
			fmt.Println(room.players[i].UserId)
			if room.players[i].UserId == userId {
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
		player.UserId = userId
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
	changeHost := false
	for i := 0; i < len(room.players); i++ {
		if room.players[i].UserId == id {
			if room.players[i].Index == room.hostIndex {
				changeHost = true
			}
			room.players[i] = nil
			room.numPlayer--
		}
	}
	if room.numPlayer == 0 {
		return true
	} else if changeHost {
		for i := 0; i < len(room.players); i++ {
			if room.players[i] != nil {
				room.hostIndex = i
			}
		}
	}

	// send list of current players to all player to update
	players, _ := json.Marshal(room.players)
	fmt.Printf("%s\n", players)
	command := models.ResponseCommand{constants.PLAYERS, string(players)}
	for i := 0; i < len(room.players); i++ {
		if room.players[i] != nil {
			room.players[i].sendCommand(command)
		}
	}
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
	case constants.BACK:
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
	deck := deal()
	shuffle(deck)
	pos := 0
	for i := 0; i < len(room.players); i++ {
		if room.players[i] != nil {
			room.players[i].InGame = true
			slice := deck[pos : pos+6]
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
	fmt.Println("Top card" + card)
	playData := models.PlayData{id, card}
	play, _ := json.Marshal(playData)
	command := models.ResponseCommand{constants.PLAY, string(play)}
	for i := 0; i < len(room.players); i++ {
		if room.players[i].UserId == id {
			room.topCard = card
			room.topCardIndex = i
			room.rowCount++
		}
		room.players[i].sendCommand(command)
	}

	// Send card play to all player
	if room.rowCount == room.numPlayer {
		room.rowCount = 0
		if room.rowCount < 5 {
			room.currentRow++
			// Note that player is allow to go final row
			room.players[room.topCardIndex].Finalist = true
			if room.currentRow == 5 {
				room.finalRowPlayer = 0
				// Inform player that is out
				for i := 0; i < len(room.players); i++ {
					command := models.ResponseCommand{Action: constants.ELIMINATED}
					if room.players[i].Finalist == false {
						room.players[i].sendCommand(command)
					} else {
						room.finalRowPlayer ++
					}
				}
			}
			// Inform lastRow top player to play
			command := models.ResponseCommand{constants.STARTROW, room.players[room.topCardIndex].UserId}
			for i := 0; i < len(room.players); i++ {
				room.players[i].sendCommand(command)
			}
		}
	}
	if room.rowCount == room.finalRowPlayer {

	}
}

func (room *Room) fold(id string, card string) {
	playData := models.PlayData{id, card}
	play, _ := json.Marshal(playData)
	command := models.ResponseCommand{constants.FOLD, string(play)}
	for i := 0; i < len(room.players); i++ {
		room.players[i].sendCommand(command)
	}
	room.normalRowCount++
	// Send card fold to all player
	if room.normalRowCount == room.numPlayer {
		room.normalRowCount = 0
		if room.currentRow < 5 {
			room.currentRow++
			// Note that player is allow to go final row
			if room.currentRow == 5 {
				// Inform player that is out
				room.finalRowPlayer
				for i := 0; i < len(room.players); i++ {
					command := models.ResponseCommand{Action: constants.ELIMINATED}
					if room.players[i].Finalist == false {
						room.players[i].sendCommand(command)
					}
				}
			}
			// Inform lastRow top player to play
			command := models.ResponseCommand{constants.STARTROW, room.players[room.topCardIndex].UserId}
			for i := 0; i < len(room.players); i++ {
				room.players[i].sendCommand(command)
			}
		}
	}

	if room.
}

func (room *Room) showBack(index int, backCard string) {
	room.finalRowCount++
	if index == room.topCardIndex {
		room.topCard = backCard
		room.topCardIndex = index
	} else if larger(room.topCard, backCard) == true {
		room.topCard = backCard
		room.topCardIndex = index
	}
	// Inform all player the current player's back card
	if room.finalRowCount == room.finalRowPlayer {
		//Calculate result
		//Adjust money
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
	for i := 1; i < len(d); i++ {
		// Create a random int up to the number of cards
		r := rand.Intn(i + 1)

		// If the the current card doesn't match the random
		// int we generated then we'll switch them out
		if i != r {
			d[r], d[i] = d[i], d[r]
		}
	}
	return d
}

func larger(leftCard string, rightCard string) bool {
	return true
}
