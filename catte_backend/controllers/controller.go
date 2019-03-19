package controllers

import (
	"encoding/json"

	"github.com/dinhnguyen138/catte/catte_backend/handlers"
	"github.com/dinhnguyen138/catte/catte_backend/models"

	"github.com/dinhnguyen138/tcp_server"
)

var roomManager handlers.RoomManager

func Init() {
	roomManager = handlers.RoomManager{}
}

func HandleCommand(command models.Command) {
	room, _ := roomManager.FindRoom(command.Room)
	room.HandleCommand(command)
}

func JoinRoom(command models.Command, c *tcp_server.Client) {
	room, isNew := roomManager.FindRoom(command.Room)
	if room == nil {
		// should return failure here
		return
	}
	if isNew == true {
		room.KickUserCallback(func(roomId string, index int) {
			KickUser(roomId, index)
		})
	}
	var playerInfo models.PlayerInfo
	json.Unmarshal([]byte(command.Data), &playerInfo)
	room.JoinRoom(playerInfo.Id, command.Data, c)
}

func LeaveRoom(command models.Command) {
	room, _ := roomManager.FindRoom(command.Room)
	room.LeaveRoom(command.Index)
	if room.IsEmpty() {
		roomManager.RemoveRoom(command.Room)
	}
}

func HandleDisconnect(c *tcp_server.Client) {
	room, index := roomManager.FindClient(c)
	if room != nil {
		room.Disconnect(index)
	}
}

func KickUser(roomId string, index int) {
	room, _ := roomManager.FindRoom(roomId)
	room.LeaveRoom(index)
	if room.IsEmpty() {
		roomManager.RemoveRoom(roomId)
	}
}
