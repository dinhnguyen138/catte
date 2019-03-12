package controllers

import (
	"encoding/json"
	"net/http"

	"../db"
	"../models"
)

var rooms []models.Room

func GetRooms(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	data, _ := json.Marshal(db.GetRooms())
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
