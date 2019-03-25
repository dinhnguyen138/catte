package db

import (
	"database/sql"
	"fmt"

	"github.com/dinhnguyen138/catte/catte_service/models"
	"github.com/dinhnguyen138/catte/catte_service/settings"
	"github.com/google/uuid"
	"github.com/kataras/golog"
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
	fmt.Println(username + " " + password)
	stmt := "SELECT userid, password FROM public.users WHERE username = $1"

	rows, err := db.Query(stmt, username)
	if err != nil {
		golog.Println(err)
		return ""
	}
	for rows.Next() {
		var userid string
		var pass string
		err := rows.Scan(&userid, &pass)
		fmt.Println(userid + " " + pass)
		if err != nil {
			golog.Fatal(err)
		}
		if bcrypt.CompareHashAndPassword([]byte(pass), []byte(password)) == nil {
			return userid
		}
	}
	return ""
}

func CheckIn(userid string) int64 {
	stmt, err := db.Prepare("SELECT lastcheckin = current_date, amount FROM public.users WHERE userid = $1")
	if err != nil {
		golog.Println(err)
		return 0
	}
	defer stmt.Close()
	var result bool
	var amount int64
	update := int64(0)
	err = stmt.QueryRow(userid).Scan(&result, &amount)
	if result == false {
		if amount < 10000 {
			amount += 30000
			update = int64(30000)
		}
		statement := "UPDATE public.users SET lastcheckin = current_date, amount = $2 WHERE userid = $1"
		_, _ = db.Exec(statement, userid, amount)
	}

	return update
}

func GetUser(userid string) *models.UserInfo {
	fmt.Println(userid)
	stmt := "SELECT userid, username, user3rdid, amount, source, image FROM public.users WHERE userid = $1"
	rows, err := db.Query(stmt, userid)
	if err != nil {
		golog.Println(err)
		return nil
	}
	for rows.Next() {
		var user models.UserInfo
		err = rows.Scan(&user.UserId, &user.UserName, &user.User3rdId, &user.Amount, &user.Source, &user.Image)
		if err != nil {
			golog.Println(err)
		}
		return &user
	}
	return nil
}

func Get3rdUser(user3rdid string, source string) *models.UserInfo {
	stmt := "SELECT userid, username, user3rdid, amount, source, image FROM public.users WHERE user3rdid = $1 AND source = $2"
	rows, err := db.Query(stmt, user3rdid, source)
	if err != nil {
		golog.Println(err)
		return nil
	}
	for rows.Next() {
		var user models.UserInfo
		err = rows.Scan(&user.UserId, &user.UserName, &user.User3rdId, &user.Amount, &user.Source, &user.Image)
		if err != nil {
			golog.Println(err)
		}
		return &user
	}
	return nil
}

func CreateAppUser(username string, password string) {
	fmt.Println(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	stmt, err := db.Prepare("INSERT INTO public.users (userid, username, source, password, amount, user3rdid, image) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		golog.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid.New().String(), username, "App", string(hashedPassword), 50000, "", "")
	if err != nil {
		golog.Fatal(err)
	}
}

func Create3rdUser(username string, user3rdid string, source string, image string) string {
	stmt, err := db.Prepare("INSERT INTO public.users (userid, username, source, password, amount, user3rdid, image) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		golog.Println(err)
		return ""
	}
	defer stmt.Close()
	userid := uuid.New().String()
	_, err = stmt.Exec(userid, username, source, "", 50000, user3rdid, image)
	if err != nil {
		golog.Fatal(err)
		return ""
	}
	return userid
}

func GetRooms() []models.Room {
	stmt := "SELECT roomid, numplayer, amount, host, maxplayer FROM public.rooms WHERE numplayer > 0"
	var rooms = []models.Room{}
	rows, err := db.Query(stmt)
	if err != nil {
		golog.Println(err)
		return rooms
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.Id, &room.NoPlayer, &room.Amount, &room.Host, &room.MaxPlayer)
		if err != nil {
			golog.Println(err)
		}
		rooms = append(rooms, room)
	}
	golog.Println(rooms)
	return rooms
}

func FindRoom(amount int64) *models.Room {
	stmt := "SELECT roomid, numplayer, amount, host, maxplayer FROM public.rooms WHERE amount < $1 and numplayer > 0"

	rows, err := db.Query(stmt, amount/2)
	if err != nil {
		golog.Println(err)
		return nil
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.Id, &room.NoPlayer, &room.Amount, &room.Host, &room.MaxPlayer)
		if err != nil {
			golog.Println(err)
		} else {
			return &room
		}
	}
	return nil
}

func CreateRoom(amount int64, maxplayer int, host string) *models.Room {
	stmt := "SELECT roomid FROM public.rooms WHERE numplayer = 0 LIMIT 1"
	var roomid string
	rows, err := db.Query(stmt)

	if err != nil {
		golog.Println(err)
		return nil
	}
	for rows.Next() {
		err := rows.Scan(&roomid)
		if err != nil {
			golog.Fatal(err)
		} else {
			break
		}
	}
	stmt = "UPDATE public.rooms SET amount = $1, maxplayer = $2, host = $3 WHERE roomid = $4"
	_, err = db.Exec(stmt, amount, maxplayer, host, roomid)
	if err != nil {
		golog.Println(err)
		return nil
	}
	return &models.Room{Id: roomid, Amount: amount, Host: host}
}
