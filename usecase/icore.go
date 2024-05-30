package usecase

import (
	"context"
	"vk-rest/pkg/models"
)

type ICore interface {
	GetUserName(ctx context.Context, sid string) (string, error)
	CreateSession(ctx context.Context, login string) (models.Session, error)
	FindActiveSession(ctx context.Context, sid string) (bool, error)
	KillSession(ctx context.Context, sid string) error
	GetUserId(ctx context.Context, sid string) (uint64, error)

	CreateUserAccount(ctx context.Context, login string, password string) error
	FindUserAccount(ctx context.Context, login string, password string) (*models.UserItem, bool, error)
	FindUserByLogin(ctx context.Context, login string) (bool, error)
	GetUserStat(ctx context.Context, userId uint64) (*models.UserStat, error)

	QuestionAdd(ctx context.Context, quest *models.Quest) (uint64, error)
	QuestionEvent(ctx context.Context, event *models.EventItem) error
	GetUserBalance(ctx context.Context, userId uint64) (uint64, error)
}
