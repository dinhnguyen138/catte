package rooms

import (
	"fmt"
	"math/rand"

	"../models"
	"github.com/firstrow/tcp_server"
)

type Room struct {
	players        [6]*Player
	numPlayer      int
	currRow        int
	normalRowCount int
	finalRowCount  int
	topCardIndex   int
	topCard        string
	hostIndex      int
}

type Player struct {
	id           string
	index        int
	ready        bool
	finalist     bool
	disconnected bool
	client       *tcp_server.Client
}

type RoomManager struct {
	rooms map[string]*Room
}

func (room *Room) JoinRoom(id string, c *tcp_server.Client) {
	nilIndex := -1
	joint := false
	for i := 0; i < len(room.players); i++ {
		if room.players[i] == nil {
			nilIndex = i
		} else if room.players[i].id == id {
			room.players[i].disconnected = true
			room.players[i].client = c
			joint = true
		}
	}
	if joint == false && nilIndex != -1 {
		player := &Player{}
		player.index = nilIndex
		player.client = c
		player.ready = false
		room.players[nilIndex] = player
		room.numPlayer++
		if room.numPlayer == 1 {
			room.hostIndex = nilIndex
		}
	}
}

func (room *Room) LeaveRoom(id string) (empty bool) {
	changeHost := false
	for i := 0; i < len(room.players); i++ {
		if room.players[i].id == id {
			if room.players[i].index == room.hostIndex {
				changeHost = true
			}
			room.players[i] = nil
			room.numPlayer--
		}
	}
	if room.numPlayer == 0 {
		return true
	} else {
		for i := 0; i < len(room.players); i++ {
			if room.players[i] != nil {
				room.hostIndex = i
				// Inform to update host
			}
		}
	}
	return false
}

func (room *Room) HandleCommand(command models.Command) {
	switch command.Action {
	case "LEAVE":
		roomEmpty := room.LeaveRoom(command.Id)
		if roomEmpty {
			// Remove room
		}
		break
	case "DEAL":
		room.newGame()
		break
	case "PLAY":
		room.play(command.Id, command.Card)
		break
	case "FOLD":
		room.fold(command.Id, command.Card)
		break
	case "SHOWFRONT":
		room.showFront(command.Id, command.Card)
		break
	case "SHOWBACK":
		break
	}
}

func (roomManager *RoomManager) FindRoom(id string) (room *Room) {
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
				return v, v.players[i].index
			}
		}
	}
	return nil, 0
}

func (room *Room) Disconnect(index int) {
	room.players[index].disconnected = true
	// Start reconnect timer
}

func (room *Room) newGame() {
	deck := deal()
	shuffle(deck)
	pos := 0
	for i := 0; i < len(room.players); i++ {
		if room.players[i] != nil {
			room.players[i].ready = true
			slice := deck[pos : pos+6]
			pos += 7
			fmt.Println(slice)
			// send slice to client
			// room.players[i].client.send
		}
	}
}

func (room *Room) play(id string, card string) {
	for i := 0; i < len(room.players); i++ {
		if room.players[i].id == id {
			room.topCard = card
			room.nextRowIndex = i
			room.rowCount++
		}
	}

	// Send card play to all player
	if room.rowCount == room.num {
		room.rowCount = 0
		room.row++
		// Note that player is allow to go final row
		// Inform lastRow top player to play
	}
}

func (room *Room) fold(id string, card string) {
	room.rowCount++
	// Send card fold to all player
	if room.rowCount == room.num {
		room.rowCount = 0
		room.row++
		// Note that player is allow to go final row
		if room.row == 5 {
			// Inform player that is out
		}
		// Inform lastRow top player to all player
	}
}

func (room *Room) showFront(id string, frontCard string) {
	if index == room.nextRowIndex {
		room.topCard = frontCard
		room.nextRowIndex = index
	} else if larger(room.topCard, frontCard) == true {
		room.topCard = frontCard
		room.nextRowIndex = index
	}
	// Inform all player the current player's front card
	if room.rowCount == room.num {
		room.rowCount = 0
		room.row++
	}
}

func (room *Room) showBack(index int, backCard string) {
	room.rowCount++
	if index == room.nextRowIndex {
		room.topCard = backCard
		room.nextRowIndex = index
	} else if larger(room.topCard, backCard) == true {
		room.topCard = backCard
		room.nextRowIndex = index
	}
	// Inform all player the current player's back card
	if room.rowCount == room.num {
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
