package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"../db"
	"../models"
)

var rooms []models.Room

func GetRooms(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	db := db.InitDB()
	data, _ := json.Marshal(db.GetRooms())
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func JoinRoom(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	roomid := mux.Vars(r)["id"]
	validRoom := false
	for i := 0; i < len(rooms); i++ {
		x := &rooms[i]
		if x.Id == roomid {
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
		if x.Id == roomid {
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
