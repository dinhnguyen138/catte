package db

import (
	"../models"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	users []models.User
	rooms []models.Room
}

var db *DB = nil

func InitDB() *DB {
	if db == nil {
		db = &DB{[]models.User{}, []models.Room{}}
		db.InitDB()
	}

	return db
}

func (db *DB) InitDB() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("pham"), 10)
	db.users = append(db.users, models.User{
		UUID:     "1",
		Username: "hoan",
		Amount:   1000,
		Password: string(hashedPassword),
	})

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("nguyen"), 10)
	db.users = append(db.users, models.User{
		UUID:     "2",
		Username: "dinh",
		Amount:   5000,
		Password: string(hashedPassword),
	})

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("huynh"), 10)
	db.users = append(db.users, models.User{
		UUID:     "3",
		Username: "thai",
		Amount:   2000,
		Password: string(hashedPassword),
	})

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("huynh"), 10)
	db.users = append(db.users, models.User{
		UUID:     "4",
		Username: "phuong",
		Amount:   2000,
		Password: string(hashedPassword),
	})

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("nguyen"), 10)
	db.users = append(db.users, models.User{
		UUID:     "3",
		Username: "minh",
		Amount:   1500,
		Password: string(hashedPassword),
	})

	db.rooms = append(db.rooms,
		models.Room{"1", true, 2, 1000},
	)

	db.rooms = append(db.rooms,
		models.Room{"2", true, 3, 2000},
	)

	db.rooms = append(db.rooms,
		models.Room{"3", false, 0, 1000},
	)
}

func (db *DB) FindUserByName(username string) *models.User {
	for i := 0; i < len(db.users); i++ {
		if username == db.users[i].Username {
			return &db.users[i]
		}
	}
	return nil
}

func (db *DB) FindUserById(id string) *models.User {
	for i := 0; i < len(db.users); i++ {
		if id == db.users[i].UUID {
			return &db.users[i]
		}
	}
	return nil
}

func (db *DB) GetRooms() []models.Room {
	return db.rooms
}
