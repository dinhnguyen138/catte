package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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
	var room *models.Room
	if host != "" {
		room = db.CreateRoom(request.Amount, host)
	} else {
		room = nil
	}

	if room == nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		data, _ := json.Marshal(room)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}

}

func QuickFind(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := r.Context().Value("user")
	claim := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userId, _ := claim["sub"].(string)
	foundUser := db.GetUser(userId)
	room := db.FindRoom(foundUser.Amount)
	data, _ := json.Marshal(room)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
