package models

type Quest struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Cost uint64 `json:"cost"`
}
