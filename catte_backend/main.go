package main

import (
	"encoding/json"
	"fmt"

	"./constants"
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
		fmt.Println(message)
		// new message received
		var cmd models.Command
		err := json.Unmarshal([]byte(message), &cmd)
		if err != nil {
			fmt.Println(message)
			c.Send(message)
		} else {
			if cmd.Action == constants.JOIN {
				controllers.JoinRoom(cmd, c)
				fmt.Println(cmd)
				return
			}
			controllers.HandleCommand(cmd)
		}
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
		fmt.Println("Client close")
	})

	server.Listen()
}
