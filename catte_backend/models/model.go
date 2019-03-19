package models

type Command struct {
	Action string `json:"action"`
	Room   string `json:"room"`
	Index  int    `json:"index"`
	Data   string `json:"data"`
}

type ResponseCommand struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

type PlayData struct {
	Index    int    `json:"index"`
	Row      int    `json:"row"`
	NextTurn int    `json:"nextturn"`
	NewRow   bool   `json:"newrow"`
	Data     string `json:"data"`
}

type WinnerCommand struct {
	Index int        `json:"index"`
	Data  []PlayData `json:"lastplays"`
}

type PlayerInfo struct {
	Id       string `json:"id" form:"-"`
	Username string `json:"username"`
	Amount   int64  `json:"amount"`
	Image    string `json:"image"`
}

type RegisterMsg struct {
	IpAddress string `json:"ip"`
}

type ResultMsg struct {
	Index  int   `json:"index"`
	Change int64 `json:"change"`
	Amount int64 `json:"amount"`
}

type LeaveMsg struct {
	Index int `json:"index"`
	Host  int `json:"host"`
}

type Deck []string
