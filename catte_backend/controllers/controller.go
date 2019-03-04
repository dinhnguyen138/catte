package controllers

import (
	"fmt"

	"../handlers"
	"../models"
	"github.com/firstrow/tcp_server"
)

var roomManager handlers.RoomManager

func Init() {
	roomManager = handlers.RoomManager{}
}

func HandleCommand(command models.Command) {
	room := roomManager.FindRoom(command.Room)
	room.HandleCommand(command)
}

func JoinRoom(command models.Command, c *tcp_server.Client) {
	room := roomManager.FindRoom(command.Room)
	room.JoinRoom(command.Id, c)
}

func LeaveRoom(command models.Command) {
	room := roomManager.FindRoom(command.Room)
	room.LeaveRoom(command.Id)
	if len(room.Players) == 0 {
		roomManager.RemoveRoom(command.Id)
	}
}

func HandleDisconnect(c *tcp_server.Client) {
	room, id := roomManager.FindClient(c)
	kickUser := func(roomId string, userId string) {
		KickUser(roomId, userId)
	}
	if room != nil {
		room.Disconnect(id, kickUser)
	}
}

func KickUser(roomId string, userId string) {
	fmt.Println("Kick user " + userId)
	room := roomManager.FindRoom(roomId)
	room.LeaveRoom(userId)
	if len(room.Players) == 0 {
		roomManager.RemoveRoom(roomId)
	}
}
