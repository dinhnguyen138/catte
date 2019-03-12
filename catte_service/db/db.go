package db

import (
	"database/sql"
	"fmt"
	"log"

	"../models"
	"../settings"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func InitDB() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		settings.Get().DBHost, settings.Get().DBPort,
		settings.Get().DBUser, settings.Get().DBPassword, settings.Get().DBName)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}

func CloseDB() {
	db.Close()
}

func AuthUser(username string, password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	stmt := "SELECT userid, password FROM public.users WHERE username = $1"

	rows, err := db.Query(stmt, username)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var userid string
		var pass string
		err := rows.Scan(&userid, &pass)
		if err != nil {
			log.Fatal(err)
		}
		if bcrypt.CompareHashAndPassword(hashedPassword, []byte(pass)) != nil {
			return userid
		}
	}
	return ""
}

func GetUser(userid string) *models.UserInfo {
	stmt, err := db.Prepare("SELECT userid, username, user3rdid, amount, source FROM public.users WHERE userid = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var user = &models.UserInfo{}
	err = db.QueryRow(userid).Scan(&user.UserId, &user.UserName, &user.User3rdId, &user.Amount, &user.Source)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return user
}

func Get3rdUser(user3rdid string, source string) *models.UserInfo {
	stmt, err := db.Prepare("SELECT userid, username, user3rdid, amount, source FROM public.users WHERE user3rdid = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var user = &models.UserInfo{}
	err = db.QueryRow(user3rdid).Scan(&user.UserId, &user.UserName, &user.User3rdId, &user.Amount, &user.Source)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return user
}

func CreateAppUser(username string, password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	stmt, err := db.Prepare("INSERT INTO public.users (userid, username, password, amount, source) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid.New().String(), username, string(hashedPassword), 50000, "App")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

func Create3rdUser(username string, user3rdid string, source string) string {
	stmt, err := db.Prepare("INSERT INTO public.users (userid, username, user3rdid, amount, source) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	userid := uuid.New().String()
	_, err = stmt.Exec(userid, username, user3rdid, 50000, source)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return userid
}

func GetRooms() []models.Room {
	stmt := "SELECT roomid, numplayer, amount FROM public.rooms WHERE isactive = true"
	var rooms = []models.Room{}
	rows, err := db.Query(stmt)
	if err != nil && err == sql.ErrNoRows {
		return rooms
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.Id, &room.NoPlayer, &room.Amount)
		if err != nil {
			log.Fatal(err)
		}
		rooms = append(rooms, room)
	}
	log.Println(rooms)
	return rooms
}

func CreateRoom(amount int) string {
	stmt := "SELECT roomid FROM public.rooms WHERE isactive = false LIMIT 1"
	var roomid string
	rows, err := db.Query(stmt)
	if err != nil && err == sql.ErrNoRows {
		return ""
	}
	for rows.Next() {
		err := rows.Scan(&roomid)
		if err != nil {
			log.Fatal(err)
		} else {
			break
		}
	}
	stmt = "UPDATE public.rooms SET amount = $1 WHERE roomid = $2"
	_, err = db.Exec(stmt, amount, roomid)
	if err != nil {
		return ""
	}
	return roomid
}
