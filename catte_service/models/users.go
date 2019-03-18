package models

type UserInfo struct {
	UserId    string `json:"id"`
	UserName  string `json:"username"`
	User3rdId string `json:"user3rdid"`
	Amount    int64  `json:"amount"`
	Source    string `json:"source"`
	Image     string `json:"image"`
}

type LoginMsg struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type RegisterMsg struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Login3rdMsg struct {
	UserName  string `json:"username"`
	User3rdId string `json:"user3rdid"`
	Token     string `json:"access_token"`
	Source    string `json:"source"`
	Image     string `json:"image"`
}
