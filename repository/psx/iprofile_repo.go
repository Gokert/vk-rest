package psx

import (
	"context"
	"vk-rest/pkg/models"
)

type IProfileRepo interface {
	GetUser(ctx context.Context, login string, password []byte) (*models.UserItem, bool, error)
	FindUser(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, login string, password []byte) error
	GetUserId(ctx context.Context, login string) (uint64, error)
}
