package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dinhnguyen138/catte/catte_service/db"
	"github.com/dinhnguyen138/catte/catte_service/models"
)

var rooms []models.Room

func GetRooms(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data, _ := json.Marshal(db.GetRooms())
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func CreateRoom(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	request := new(models.CreateRoomMsg)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)
	host := PickHost()
	roomid := ""
	if host != "" {
		roomid = db.CreateRoom(request.Amount, host)
	}

	if roomid == "" {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Write([]byte(roomid))
}
