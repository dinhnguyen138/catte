package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dinhnguyen138/catte/catte_backend/models"
	"github.com/dinhnguyen138/catte/catte_backend/settings"
	_ "github.com/lib/pq"
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

func GetRoom(roomid string) *models.Room {
	stmt := "SELECT roomid, numplayer, amount, maxplayer FROM public.rooms WHERE roomid = $1"
	rows, err := db.Query(stmt, roomid)
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.Id, &room.NoPlayer, &room.Amount, &room.MaxPlayer)
		if err != nil {
			log.Println(err)
		}
		return &room
	}
	return nil
}

func JoinRoom(roomid string) {
	stmt, err := db.Prepare("SELECT numplayer FROM public.rooms WHERE roomid = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var numplayer int
	err = db.QueryRow(roomid).Scan(&numplayer)
	if err != nil && err == sql.ErrNoRows {
		log.Fatal(err)
	}

	statement := "UPDATE public.rooms SET numplayer = $1 WHERE roomid = $2"
	_, err = db.Exec(statement, numplayer+1, roomid)
	if err != nil {
		panic(err)
	}
}

func UpdateRoom(roomid string, numplayer int) {
	statement := "UPDATE public.rooms SET numplayer = $1 WHERE roomid = $2"
	_, err := db.Exec(statement, numplayer, roomid)
	if err != nil {
		log.Println(err)
	}
}

func UpdatePlayer(playerid string, amount int64) {
	statement := "UPDATE public.users SET amount = $1 WHERE userid = $2"
	_, err := db.Exec(statement, amount, playerid)
	if err != nil {
		log.Println(err)
	}
}
