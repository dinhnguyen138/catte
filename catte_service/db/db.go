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
	stmt, err := db.Prepare("SELECT UserId FROM User WHERE Username = $1 AND Password = $2)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var userid string
	err = db.QueryRow(username, string(hashedPassword)).Scan(&userid)
	if err != nil && err == sql.ErrNoRows {
		return ""
	}
	return userid
}

func GetUser(userid string) *models.UserInfo {
	stmt, err := db.Prepare("SELECT UserId, UserName, User3rdId, Amount, Source FROM User WHERE UserId = $1")
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
	stmt, err := db.Prepare("SELECT UserId, UserName, User3rdId, Amount, Source FROM User WHERE User3rdId = $1")
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
	stmt, err := db.Prepare("INSERT INTO User(UserId, Username, Password, Amount, Source) ($1, $2, $3, $4, $5)")
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
	stmt, err := db.Prepare("INSERT INTO User(UserId, Username, User3rdId, Amount, Source) ($1, $2, $3, $4, $5)")
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
	stmt := "SELECT Id, NumPlayer, Amount FROM Room WHERE IsActive = false"
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
	return rooms
}
