package utilities

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dinhnguyen138/catte/catte_backend/models"
	"github.com/dinhnguyen138/catte/catte_backend/settings"
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
	url := "http://" + settings.Get().ServiceIp + ":8080/register-host"
	msg := models.RegisterMsg{GetPublicIp()}
	msgData, _ := json.Marshal(msg)
	for {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(msgData))
		if err != nil {
			continue
		}
		if resp != nil && resp.StatusCode == http.StatusOK {
			break
		}
	}
}
