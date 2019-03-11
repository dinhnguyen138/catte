package constants

const JOIN string = "JOIN"
const LEAVE string = "LEAVE"
const PLAY string = "PLAY"
const DEAL string = "DEAL"
const FOLD string = "FOLD"
const BACK string = "BACK"

const PLAYERS string = "PLAYERS"
const NEWPLAYER string = "NEWPLAYER"
const CARDS string = "CARDS"
const ELIMINATED string = "ELIMINATED"
const START string = "START"
const SHOWBACK string = "SHOWBACK"
const WINNER string = "WINNER"

const MAXPLAYER int = 6

var CardOrder map[string]int = map[string]int{
	"2":  1,
	"3":  2,
	"4":  3,
	"5":  4,
	"6":  5,
	"7":  6,
	"8":  7,
	"9":  8,
	"10": 9,
	"J":  10,
	"Q":  11,
	"K":  12,
	"A":  13,
}

var Types = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

var Suits = []string{"H", "D", "C", "S"}
