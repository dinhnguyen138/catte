package utilities

import (
	"io/ioutil"
	"net/http"

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
	for {
		resp, _ := http.Get(url)
		if resp.StatusCode == http.StatusOK {
			break
		}
	}
}
