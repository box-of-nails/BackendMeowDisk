package models

type UserData struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
