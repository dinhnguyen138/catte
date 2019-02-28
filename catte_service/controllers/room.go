package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"../models"
)

var rooms []models.Room

func GetRooms(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if len(rooms) == 0 {
		rooms = []models.Room{models.Room{"1", true, 2}, models.Room{"2", true, 3}, models.Room{"3", false, 0}}
	}

	data, _ := json.Marshal(rooms)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func JoinRoom(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	roomid := mux.Vars(r)["id"]
	validRoom := false
	for i := 0; i < len(rooms); i++ {
		x := &rooms[i]
		if x.UUID == roomid {
			x.IsActive = true
			x.NoPlayer += 1
			validRoom = true
			break
		}
	}
	if validRoom == true {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func LeaveRoom(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	roomid := mux.Vars(r)["id"]
	validRoom := false
	for i := 0; i < len(rooms); i++ {
		x := &rooms[i]
		if x.UUID == roomid {
			x.IsActive = true
			x.NoPlayer += 1
			validRoom = true
			break
		}
	}
	if validRoom == true {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
