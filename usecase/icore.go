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

	CreateUserAccount(login string, password string) error
	FindUserAccount(login string, password string) (*models.UserItem, bool, error)
	FindUserByLogin(login string) (bool, error)

	QuestionAdd(quest *models.Quest) (uint64, error)
	QuestionEvent(event *models.EventItem) error
}
