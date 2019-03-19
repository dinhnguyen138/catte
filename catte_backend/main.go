package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dinhnguyen138/catte/catte_backend/constants"
	"github.com/dinhnguyen138/catte/catte_backend/controllers"
	"github.com/dinhnguyen138/catte/catte_backend/db"
	"github.com/dinhnguyen138/catte/catte_backend/models"
	"github.com/dinhnguyen138/catte/catte_backend/settings"
	"github.com/dinhnguyen138/catte/catte_backend/utilities"
	"github.com/dinhnguyen138/tcp_server"
)

func main() {
	fmt.Println(os.Getenv("ENV"))
	settings.Init()
	db.InitDB()
	controllers.Init()
	go utilities.RegisterToService()
	server := tcp_server.NewWithTLS(":9999", settings.Get().ServerCertPath, settings.Get().ServerKeyPath)

	server.OnNewClient(func(c *tcp_server.Client) {
		fmt.Println("Client connect")
		fmt.Println(c)
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
				return
			}
			if cmd.Action == constants.LEAVE {
				controllers.LeaveRoom(cmd)
			}
			controllers.HandleCommand(cmd)
		}
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
		fmt.Println("Client close")
		fmt.Println(c)
		controllers.HandleDisconnect(c)
	})

	server.Listen()
}
