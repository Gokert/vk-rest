package psx

import "vk-rest/pkg/models"

type IQuestionRepo interface {
	QuestionAdd(quest *models.Quest) (uint64, error)
	QuestionEvent(event *models.EventItem) error
	GetUserStat(userId uint64) (*models.UserStat, error)
	GetUserBalance(userId uint64) (uint64, error)
}
