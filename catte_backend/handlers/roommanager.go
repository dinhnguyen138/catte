package handlers

import (
	"fmt"

	"github.com/dinhnguyen138/catte/catte_backend/db"
	"github.com/firstrow/tcp_server"
)

type RoomManager struct {
	rooms map[string]*Room
}

// Find a room in room manager, if not, create new one
func (roomManager *RoomManager) FindRoom(id string) (*Room, bool) {
	if roomManager.rooms == nil {
		roomManager.rooms = map[string]*Room{}
	}
	if _, ok := roomManager.rooms[id]; ok {
		return roomManager.rooms[id], false
	}
	// TODO Load room from DB
	room := db.GetRoom(id)
	fmt.Println(room)
	if room == nil {
		return nil, false
	}
	roomManager.rooms[id] = NewRoom(room.Id, room.MaxPlayer, room.Amount)
	return roomManager.rooms[id], true
}

// Find a client in case of disconnected. In this case we don't know who are disconnect
func (roomManager *RoomManager) FindClient(c *tcp_server.Client) (room *Room, index int) {
	for _, v := range roomManager.rooms {
		index := v.FindClient(c)
		if index != -1 {
			return v, index
		}
	}
	return nil, -1
}

func (roomManager *RoomManager) RemoveRoom(id string) {
	fmt.Println("Remove room")
	delete(roomManager.rooms, id)
}
