package main

import (
	"encoding/json"
	"fmt"

	"./controllers"
	"./models"
	"github.com/firstrow/tcp_server"
)

func main() {
	controllers.Init()
	server := tcp_server.New("localhost:9999")

	server.OnNewClient(func(c *tcp_server.Client) {
		fmt.Println("Client connect")
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		// new message received
		var cmd models.Command
		json.Unmarshal([]byte(message), &cmd)
		if cmd.Action == "JOIN" {
			controllers.JoinRoom(cmd, c)
			return
		}
		controllers.HandleCommand(cmd)
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
	})

	server.Listen()
}
