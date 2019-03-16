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
	Id        string  `json:"id" form:"-"`
	Username  string  `json:"username"`
	User3rdId string  `json:"user3rdid"`
	Source    string  `json:"source"`
	Amount    float32 `json:"amount"`
}

type RegisterMsg struct {
	IpAddress string `json:"ip"`
}

type Deck []string
