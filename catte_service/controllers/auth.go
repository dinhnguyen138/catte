package controllers

import (
	"encoding/json"
	"net/http"

	"../db"
	"../models"
	"../services"
	"github.com/dgrijalva/jwt-go"
)

func Login(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.LoginMsg)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)
	responseStatus, token := services.Login(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func GetInfo(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := r.Context().Value("user")
	claim := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userId, _ := claim["sub"].(string)
	foundUser := db.GetUser(userId)
	data, _ := json.Marshal(foundUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := r.Context().Value("user")
	claim := user.(*jwt.Token).Claims.(jwt.MapClaims)
	userId, _ := claim["sub"].(string)
	foundUser := db.GetUser(userId)
	token := services.RefreshToken(userId)
	w.Header().Set("Content-Type", "application/json")
	if foundUser != nil {
		w.WriteHeader(http.StatusOK)
		w.Write(token)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(""))
	}
}

func Logout(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// err := services.Logout(r)
	w.Header().Set("Content-Type", "application/json")
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// } else {
	w.WriteHeader(http.StatusOK)
	// }
}

func Register(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.RegisterMsg)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)
	db.CreateAppUser(requestUser.UserName, requestUser.Password)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func Login3rd(w http.ResponseWriter, r *http.Request) {
	requestUser := new(models.Login3rdMsg)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)
	// TODO: Verify access token
	user := db.Get3rdUser(requestUser.User3rdId, requestUser.Source)
	var userid string
	if user == nil {
		userid = db.Create3rdUser(requestUser.UserName, requestUser.User3rdId, requestUser.Source)
	} else {
		userid = user.UserId
	}
	token := services.RefreshToken(userid)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(token)
}
