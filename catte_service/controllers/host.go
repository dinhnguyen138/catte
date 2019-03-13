package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dinhnguyen138/catte/catte_service/models"
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
