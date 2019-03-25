package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dinhnguyen138/catte/catte_backend/models"
	"github.com/dinhnguyen138/catte/catte_backend/settings"
	"github.com/dinhnguyen138/tcp_server"
)

func GetPublicIp() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	return string(ip)
}

func RegisterToService() {
	url := settings.Get().ServiceIp + "/register-host"
	msg := models.RegisterMsg{GetPublicIp()}
	msgData, _ := json.Marshal(msg)
	for {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgData))
		if err != nil {
			fmt.Println("Failed to register-host")
			continue
		}
		if resp != nil && resp.StatusCode == http.StatusOK {
			break
		}
	}
}

func SendClient(c *tcp_server.Client, command string, data interface{}) {
	var stringData string
	switch data.(type) {
	case string:
		stringData = data.(string)
		break
	default:
		temp, _ := json.Marshal(data)
		stringData = string(temp)
	}
	message := models.ResponseCommand{command, string(stringData)}
	resp, _ := json.Marshal(message)
	c.Send(string(resp) + "\n")
}
