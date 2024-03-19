package psx

import "vk-rest/pkg/models"

type IQuestionRepo interface {
	QuestionAdd(quest *models.Quest) (uint64, error)
	QuestionEvent(event *models.EventItem) error
}
