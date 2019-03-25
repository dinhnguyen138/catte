package main

import (
	"encoding/json"
	"os"

	"github.com/kataras/golog"
	"github.com/natefinch/lumberjack"

	"github.com/dinhnguyen138/catte/catte_backend/constants"
	"github.com/dinhnguyen138/catte/catte_backend/controllers"
	"github.com/dinhnguyen138/catte/catte_backend/db"
	"github.com/dinhnguyen138/catte/catte_backend/models"
	"github.com/dinhnguyen138/catte/catte_backend/settings"
	"github.com/dinhnguyen138/catte/catte_backend/utilities"
	"github.com/dinhnguyen138/tcp_server"
)

func main() {
	golog.SetOutput(&lumberjack.Logger{
		Filename:   "log/backend/daily.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})

	golog.Info(os.Getenv("ENV"))
	settings.Init()
	db.InitDB()
	controllers.Init()
	go utilities.RegisterToService()
	server := tcp_server.NewWithTLS(":9999", settings.Get().ServerCertPath, settings.Get().ServerKeyPath)

	server.OnNewClient(func(c *tcp_server.Client) {
		golog.Info("Client connect")
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		defer func() {
			if r := recover(); r != nil {
				golog.Error("An error has been recovered")
				golog.Error(r)
			}
		}()
		golog.Info("Receive client message")
		golog.Info(message)
		// new message received
		var cmd models.Command
		err := json.Unmarshal([]byte(message), &cmd)
		if err != nil {
			golog.Error("Error parsing message")
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
		golog.Info("Client disconnected")
		controllers.HandleDisconnect(c)
	})

	server.Listen()
}
