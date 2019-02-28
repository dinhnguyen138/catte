package controllers

import (
	"../models"
	"../rooms"
	"github.com/firstrow/tcp_server"
)

var roomManager rooms.RoomManager

func Init() {
	roomManager = rooms.RoomManager{}
}

func HandleCommand(command models.Command) {
	room := roomManager.FindRoom(command.Room)
	room.HandleCommand(command)
}

func JoinRoom(command models.Command, c *tcp_server.Client) {
	room := roomManager.FindRoom(command.Room)
	room.JoinRoom(command.Id, c)
}

func HandleDisconnect(c *tcp_server.Client) {
	room, id := roomManager.FindClient(c)
	room.Disconnect(id)
}
