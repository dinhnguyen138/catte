package gamecontroller

import "github.com/gorilla/websocket"

type Player struct {
	id   string
	room *RoomController
	conn *websocket.Conn
	send chan []byte
}

type RoomController struct {
	id         string
	player     map[*Player]bool
	broadcast  chan []byte
	register   chan *Player
	unregister chan *Player
}

var roomControllers []*RoomController

func findRoom(id string) *RoomController {
	for i := 0; i < len(roomControllers); i++ {
		if roomControllers[i].id == id {
			return roomControllers[i]
		}
	}
	newRoom := &RoomController{
		id:         id,
		player:     make(map[*Player]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Player),
		unregister: make(chan *Player),
	}

	roomControllers = append(roomControllers, newRoom)
	return newRoom
}
