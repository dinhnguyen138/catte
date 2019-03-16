package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dinhnguyen138/catte/catte_service/models"
	"github.com/dinhnguyen138/catte/catte_service/utilities"
)

var hosts []string

func RegisterHost(w http.ResponseWriter, r *http.Request) {
	request := new(models.RegisterHostMsg)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)
	log.Println(request.IpAddress)
	hosts = append(hosts, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func PickHost() string {
	for i, s := 0, len(hosts); i < s; i++ {
		if utilities.CheckPing(hosts[i]) == true {
			return hosts[i]
		} else {
			hosts = append(hosts[:i], hosts[i+1:]...)
			i--
			s--
		}
	}
	return ""
}
