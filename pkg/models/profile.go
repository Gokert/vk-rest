package models

type UserItem struct {
	Id       uint64 `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
}
