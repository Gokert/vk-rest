package psx

import (
	"context"
	"vk-rest/pkg/models"
)

type IQuestionRepo interface {
	QuestionAdd(ctx context.Context, quest *models.Quest) (uint64, error)
	QuestionEvent(ctx context.Context, event *models.EventItem) error
	GetUserStat(ctx context.Context, userId uint64) (*models.UserStat, error)
	GetUserBalance(ctx context.Context, userId uint64) (uint64, error)
}
