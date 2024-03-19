package models

type UserStat struct {
	Balance uint64  `json:"balance"`
	Quests  []Quest `json:"quests"`
}
