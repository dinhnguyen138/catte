package handlers

import (
	"fmt"

	"github.com/firstrow/tcp_server"
)

type RoomManager struct {
	rooms map[string]*Room
}

// Find a room in room manager, if not, create new one
func (roomManager *RoomManager) FindRoom(id string) (room *Room, isNew bool) {
	if roomManager.rooms == nil {
		roomManager.rooms = map[string]*Room{}
	}
	if _, ok := roomManager.rooms[id]; ok {
		return roomManager.rooms[id], false
	}
	// TODO Load room from DB
	roomManager.rooms[id] = &Room{Id: id, IndexUsed: []bool{false, false, false, false, false, false}}
	return roomManager.rooms[id], true
}

// Find a client in case of disconnected. In this case we don't know who are disconnect
func (roomManager *RoomManager) FindClient(c *tcp_server.Client) (room *Room, index int) {
	for _, v := range roomManager.rooms {
		for i := 0; i < len(v.Players); i++ {
			if v.Players[i] != nil && v.Players[i].client == c {
				return v, v.Players[i].Index
			}
		}
	}
	return nil, -1
}

func (roomManager *RoomManager) RemoveRoom(id string) {
	fmt.Println("Remove room")
	delete(roomManager.rooms, id)
}
